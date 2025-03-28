package aws

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
	"github.com/sparkfund/pkg/errors"
)

// Config represents AWS configuration
type Config struct {
	Region          string
	AccessKeyID     string
	SecretAccessKey string
	BucketName      string
}

// Client represents an AWS client
type Client struct {
	s3  *s3.Client
	ses *ses.Client
	cfg *Config
}

// NewClient creates a new AWS client
func NewClient(cfg *Config) (*Client, error) {
	creds := credentials.NewStaticCredentialsProvider(
		cfg.AccessKeyID,
		cfg.SecretAccessKey,
		"",
	)

	awsCfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(cfg.Region),
		config.WithCredentialsProvider(creds),
	)
	if err != nil {
		return nil, errors.ErrInternalServer(err)
	}

	return &Client{
		s3:  s3.NewFromConfig(awsCfg),
		ses: ses.NewFromConfig(awsCfg),
		cfg: cfg,
	}, nil
}

// UploadFile uploads a file to S3
func (c *Client) UploadFile(key string, reader io.Reader, contentType string) error {
	_, err := c.s3.PutObject(context.Background(), &s3.PutObjectInput{
		Bucket:      aws.String(c.cfg.BucketName),
		Key:         aws.String(key),
		Body:        reader,
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return errors.ErrInternalServer(err)
	}
	return nil
}

// DownloadFile downloads a file from S3
func (c *Client) DownloadFile(key string) (io.ReadCloser, error) {
	result, err := c.s3.GetObject(context.Background(), &s3.GetObjectInput{
		Bucket: aws.String(c.cfg.BucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, errors.ErrInternalServer(err)
	}
	return result.Body, nil
}

// DeleteFile deletes a file from S3
func (c *Client) DeleteFile(key string) error {
	_, err := c.s3.DeleteObject(context.Background(), &s3.DeleteObjectInput{
		Bucket: aws.String(c.cfg.BucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		return errors.ErrInternalServer(err)
	}
	return nil
}

// GetPresignedURL generates a presigned URL for S3 object
func (c *Client) GetPresignedURL(key string, expiration time.Duration) (string, error) {
	presignClient := s3.NewPresignClient(c.s3)
	presignResult, err := presignClient.PresignGetObject(context.Background(), &s3.GetObjectInput{
		Bucket: aws.String(c.cfg.BucketName),
		Key:    aws.String(key),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = expiration
	})
	if err != nil {
		return "", errors.ErrInternalServer(err)
	}
	return presignResult.URL, nil
}

// SendEmail sends an email using SES
func (c *Client) SendEmail(to, subject, body string) error {
	input := &ses.SendEmailInput{
		Destination: &types.Destination{
			ToAddresses: []string{to},
		},
		Message: &types.Message{
			Body: &types.Body{
				Text: &types.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String(body),
				},
			},
			Subject: &types.Content{
				Charset: aws.String("UTF-8"),
				Data:    aws.String(subject),
			},
		},
		Source: aws.String("noreply@sparkfund.com"),
	}

	_, err := c.ses.SendEmail(context.Background(), input)
	if err != nil {
		return errors.ErrInternalServer(err)
	}
	return nil
}

// SendTemplatedEmail sends a templated email using SES
func (c *Client) SendTemplatedEmail(to, templateName string, templateData map[string]string) error {
	input := &ses.SendTemplatedEmailInput{
		Destination: &types.Destination{
			ToAddresses: []string{to},
		},
		Source:       aws.String("noreply@sparkfund.com"),
		Template:     aws.String(templateName),
		TemplateData: aws.String(fmt.Sprintf("%v", templateData)),
	}

	_, err := c.ses.SendTemplatedEmail(context.Background(), input)
	if err != nil {
		return errors.ErrInternalServer(err)
	}
	return nil
} 