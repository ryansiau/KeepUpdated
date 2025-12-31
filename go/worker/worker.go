package worker

import (
	"container/heap"
	"context"
	"slices"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/avast/retry-go/v5"

	"github.com/ryansiau/KeepUpdated/go/config"
	"github.com/ryansiau/KeepUpdated/go/model"
	"github.com/ryansiau/KeepUpdated/go/pkg/database"
	graceful_shutdown "github.com/ryansiau/KeepUpdated/go/pkg/graceful-shutdown"
	workflow_heap "github.com/ryansiau/KeepUpdated/go/worker/workflow-heap"
)

type Worker struct {
	executions workflow_heap.WorkflowHeap
	workflows  []config.Workflow

	db *gorm.DB

	gracefulShutdown graceful_shutdown.GracefulShutdown
}

func NewWorker(config *config.Config, gracefulShutdown graceful_shutdown.GracefulShutdown) (*Worker, error) {
	executions := make([]*workflow_heap.Execution, len(config.Workflows))
	for idx, workflow := range config.Workflows {
		// next execution automatically assigned with 0
		nextExecution := time.Now()

		executions[idx] = &workflow_heap.Execution{
			Workflow:      workflow,
			Interval:      workflow.Interval,
			NextExecution: nextExecution,
		}
	}

	// initiate DB
	db, err := database.NewDB(&config.Database)
	if err != nil {
		return nil, err
	}

	err = database.Migrate(db)
	if err != nil {
		return nil, err
	}

	return &Worker{
		executions:       executions,
		workflows:        config.Workflows,
		db:               db,
		gracefulShutdown: gracefulShutdown,
	}, nil
}

func (w *Worker) Run() error {
	logrus.Info("Setting up worker")

	// test db the ability to read and write
	err := database.CheckConnectionCapability(w.db)
	if err != nil {
		return err
	}

	retrier := retry.New(
		retry.Attempts(3),
		retry.Delay(100*time.Millisecond),
		retry.DelayType(retry.BackOffDelay))

	for {
		if w.gracefulShutdown.IsTerminated() {
			break
		}

		for {
			// find which workflows need to be processed
			top := w.executions.Peek()
			if top.NextExecution.After(time.Now()) {
				break
			} else {
				heap.Pop(&w.executions)
			}

			// init usable variables
			startTime := time.Now()
			workflow := top.Workflow

			var source model.Source
			var filters []model.Filter
			var notifiers []model.Notifier

			// build configs into its own implementor
			source, err = workflow.Source.Config.Build(workflow.Source.Name)
			if err != nil {
				return err
			}

			for _, filter := range workflow.Filters {
				f, err := filter.Config.Build()
				if err != nil {
					return err
				}
				filters = append(filters, f)
			}

			for _, notifier := range workflow.Notifiers {
				n, err := notifier.Config.Build()
				if err != nil {
					return err
				}
				notifiers = append(notifiers, n)
			}

			// get the latest PublishedAt recorded in the database
			// TODO utilize this for data filtering instead of using id
			var latestPublishedAt time.Time
			dbResult := w.db.Model(&model.Content{}).
				Select("published_at").
				Where("source_id = ?", source.SourceID()).
				Order("published_at DESC").
				Limit(1).
				Scan(&latestPublishedAt)
			if dbResult.Error != nil {
				return dbResult.Error
			}

			// check if this is a new source
			var isNewSource bool
			if latestPublishedAt.IsZero() {
				isNewSource = true
			}

			// call into the sources
			ctx := context.Background()

			contents, err := source.Fetch(ctx)
			if err != nil {
				return err
			}

			logrus.WithField("workflow", workflow.Name).Infof("Fetched %d contents from %s", len(contents), source.Name())

			// if the db is empty, fill the db with every update except the latest
			if isNewSource {
				dbResult = w.db.Create(contents[1:])
				if dbResult.Error != nil {
					return dbResult.Error
				}
			}

			// fetch into db and filter out old updates
			var contentIDs []string
			for _, content := range contents {
				contentIDs = append(contentIDs, content.ID)
			}

			// currently the most reliable way, by comparing the ids.
			// however, comparison by PublishedAt is a good choice to consider
			var trackedContents []model.Content
			dbResult = w.db.
				Where("source_id = ?", source.SourceID()).
				Find(&trackedContents, contentIDs)
			if dbResult.Error != nil {
				return dbResult.Error
			}

			trackedContentMaps := map[string]struct{}{}
			for _, trackedContent := range trackedContents {
				trackedContentMaps[trackedContent.ID] = struct{}{}
			}

			var newContents []model.Content
			for _, content := range contents {
				if _, ok := trackedContentMaps[content.ID]; !ok {
					newContents = append(newContents, content)
					logrus.Infof("Fetched new content from %s: %s", content.Platform, content.Title)
				}
			}

			// apply filters
			var filteredContents []model.Content

			for _, content := range newContents {
				valid := true

				for _, filter := range filters {
					valid = valid && filter.Apply(content)
				}

				if valid {
					filteredContents = append(filteredContents, content)
				} else {
					logrus.Infof("Filtered out content from %s: %s", content.Platform, content.Title)
				}
			}
			logrus.WithField("workflow", workflow.Name).Infof("Filtered %d contents", len(filteredContents))

			// sort contents by PublishedAt
			slices.SortFunc(filteredContents, func(a, b model.Content) int {
				if a.PublishedAt.Before(b.PublishedAt) {
					return -1
				} else if a.PublishedAt.After(b.PublishedAt) {
					return 1
				}
				return 0
			})

			// call notifier
			for _, notifier := range notifiers {
				for _, content := range filteredContents {
					err = retrier.Do(func() error {
						return notifier.Send(ctx, content)
					})
					if err != nil {
						return err
					}
				}
			}

			// TODO log request response history when sending notification. This serves as debugging log, but is it needed and is it secure?
			//      logging req response, including url methods and auth seems scary as it'll store the complete url and auth header too
			//      perhaps add an env or a new field in config to decide whether to log this?

			// if notifier succeeds, the new updates MUST BE updated to the db. failing to update means they will be resented
			//   in the next iteration. if the storing process errors or fails, find an alternative way to store this data.
			//   perhaps just throw an error and stop the program?
			if len(filteredContents) > 0 {
				err = retrier.Do(func() error {
					dbResult := w.db.Create(&filteredContents)
					if dbResult.Error != nil {
						return dbResult.Error
					}
					return nil
				})
				if err != nil {
					return dbResult.Error
				}
			}

			// log:
			//   - which workflow has been called
			//   - when
			//   - the duration until the process finished
			//   - data count information:
			//       - new updates
			//       - filtered out
			//       - notified
			//   - notification channels
			var notificationChannelNames []string
			for _, notifier := range notifiers {
				notificationChannelNames = append(notificationChannelNames, notifier.Name())
			}

			logrus.WithFields(logrus.Fields{
				"workflow":    workflow.Name,
				"started_at":  startTime,
				"finished_at": time.Now(),
				"duration_ms": time.Since(startTime).Milliseconds(),
				"summary": map[string]interface{}{
					"new_updates":  len(newContents),
					"filtered_out": len(newContents) - len(filteredContents),
					"notified":     len(filteredContents),
				},
				"channels": notificationChannelNames,
			}).Info("Finished processing workflow")

			// assign a new schedule and re-register to the heap
			top.NextExecution = time.Now().Add(workflow.Interval)
			heap.Push(&w.executions, top)

			if w.gracefulShutdown.IsTerminated() {
				break
			}
		}

		if w.gracefulShutdown.IsTerminated() {
			break
		}

		// how often the worker check if there's a workflow to run
		time.Sleep(5 * time.Second)
	}

	return nil
}
