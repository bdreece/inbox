package email

import (
	"context"
	"fmt"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
	"go.uber.org/config"
)

type SESOptions struct {
	Region string
	From   string
	To     string
}

type SESClient struct {
	client *ses.Client
	source string
}

var (
	utf8        = "UTF-8"
	_    Client = (*SESClient)(nil)
)

// Send implements Client.
func (c *SESClient) Send(ctx context.Context, msg Message) error {
	body := new(types.Body)
	if msg.HTML != "" {
		body.Html = &types.Content{
			Data:    &msg.HTML,
			Charset: &utf8,
		}
	}
    if msg.Text != "" {
        body.Text = &types.Content{
			Data:    &msg.Text,
			Charset: &utf8,
        }
    }

	_, err := c.client.SendEmail(ctx, &ses.SendEmailInput{
		Source:           &c.source,
		ReplyToAddresses: msg.ReplyTo,
		Destination: &types.Destination{
			ToAddresses:  msg.To,
			CcAddresses:  msg.CC,
			BccAddresses: msg.BCC,
		},
		Message: &types.Message{
			Subject: &types.Content{
				Data:    &msg.Subject,
				Charset: &utf8,
			},
			Body: body,
		},
	})

	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

func NewSESClient(ctx context.Context, opts SESOptions) (*SESClient, error) {
	cfg, err := awsconfig.LoadDefaultConfig(ctx,
		awsconfig.WithRegion(opts.Region),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to initialize AWS config: %w", err)
	}

	return &SESClient{ses.NewFromConfig(cfg), opts.From}, nil
}

func ConfigureSES(cfg config.Provider) (*SESOptions, error) {
	const key string = "ses"
	var opts SESOptions

	if err := cfg.Get(key).Populate(&opts); err != nil {
		return nil, fmt.Errorf("failed to populate SES options: %w", err)
	}

	return &opts, nil
}
