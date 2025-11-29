package model

import "context"

// Notifier defines the interface for notification channels
type Notifier interface {
	// Name returns the name of the notifier
	Name() string

	// Send sends a notification for the given content
	Send(ctx context.Context, content Content) error

	// Type returns the notification type (Console, Email, Discord, etc.)
	Type() string
}
