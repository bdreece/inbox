package email

import "context"

type Message struct {
	To      []string
	ReplyTo []string
	CC      []string
	BCC     []string
	Subject string
	HTML    string
	Text    string
}

type Client interface {
    Send(ctx context.Context, msg Message) error
}
