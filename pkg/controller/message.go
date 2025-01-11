package controller

import (
	"log/slog"
	"net/http"

	"github.com/bdreece/inbox/pkg/email"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Message struct {
	client    email.Client
	validator *validator.Validate
	logger    *slog.Logger
	opts      MessageOptions
}

type MessageOptions struct {
	Destination string
}

type messageForm struct {
	From    string `form:"from" validate:"required,email"`
	Subject string `form:"subject" validate:"required"`
	Body    string `form:"body" validate:"required"`
}

var (
	_ http.Handler = (*Message)(nil)
)

func (c *Message) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var form messageForm

	c.logger.Debug("binding form...")
	if err := binding.FormPost.Bind(r, &form); err != nil {
		c.logger.
			With("error", err).
			Warn("failed to bind form")

		w.WriteHeader(http.StatusBadRequest)
		return
	}

	c.logger.Debug("validating form...")
	if err := c.validator.Struct(&form); err != nil {
		c.logger.
			With("error", err).
			Warn("failed to validate form")

		w.WriteHeader(http.StatusBadRequest)
		return
	}

	c.logger.
		With("from", form.From).
		With("subject", form.Subject).
		Info("sending email message...")

	err := c.client.Send(r.Context(), email.Message{
		To:      []string{c.opts.Destination},
		ReplyTo: []string{form.From},
		Subject: form.Subject,
		Text:    form.Body,
	})

	if err != nil {
		c.logger.
			With("error", err).
			Error("failed to validate form")

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	c.logger.Info("email message sent successfully!")

	w.WriteHeader(http.StatusCreated)
}

func NewMessage(
	client email.Client,
	validator *validator.Validate,
	logger *slog.Logger,
	opts MessageOptions,
) *Message {
	return &Message{
		client,
		validator,
		logger,
		opts,
	}
}
