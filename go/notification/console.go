package notification

import (
	"context"
	"fmt"
	"time"

	"github.com/ryansiau/utilities/go/model"
)

// ConsoleNotifier is a simple notifier that prints to the console
type ConsoleNotifier struct {
	name string
}

// NewConsoleNotifier creates a new console notifier
func NewConsoleNotifier(name string) *ConsoleNotifier {
	if name == "" {
		name = "Console"
	}
	return &ConsoleNotifier{
		name: name,
	}
}

// Name returns the name of the notifier
func (n *ConsoleNotifier) Name() string {
	return n.name
}

// Type returns the notification type
func (n *ConsoleNotifier) Type() string {
	return "Console"
}

// Send sends a notification for the given content
func (n *ConsoleNotifier) Send(ctx context.Context, content model.Content) error {
	fmt.Printf("🔔 NEW CONTENT FROM %s 🔔\n", content.Platform)
	fmt.Printf("Title: %s\n", content.Title)
	fmt.Printf("Author: %s\n", content.Author)
	fmt.Printf("URL: %s\n", content.URL)
	fmt.Printf("Published: %s\n", content.PublishedAt.Format(time.RFC1123))
	fmt.Println("-----------------------------------")
	return nil
}
