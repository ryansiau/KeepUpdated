package workflow_heap

import (
	"container/heap"
	"time"

	"github.com/ryansiau/utilities/go/config"
)

type WorkflowHeap []*Execution

type Execution struct {
	Workflow      config.Workflow
	Interval      time.Duration
	NextExecution time.Time
}

var _ heap.Interface = (*WorkflowHeap)(nil)

func (h *WorkflowHeap) Len() int {
	return len(*h)
}

func (h *WorkflowHeap) Less(i, j int) bool {
	return (*h)[i].NextExecution.Before((*h)[j].NextExecution)
}

func (h *WorkflowHeap) Swap(i, j int) {
	(*h)[i], (*h)[j] = (*h)[j], (*h)[i]
}

func (h *WorkflowHeap) Push(x any) {
	*h = append(*h, x.(*Execution))
}

func (h *WorkflowHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func (h *WorkflowHeap) Peek() *Execution {
	return (*h)[0]
}
