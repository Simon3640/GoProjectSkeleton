package aws

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"os"
	"strings"

	contractsProviders "github.com/simon3640/goprojectskeleton/src/application/contracts/providers"
	application_errors "github.com/simon3640/goprojectskeleton/src/application/shared/errors"
	"github.com/simon3640/goprojectskeleton/src/application/shared/locales/messages"
	email_models "github.com/simon3640/goprojectskeleton/src/application/shared/services/emails/models"
	"github.com/simon3640/goprojectskeleton/src/application/shared/status"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// S3RendererBase is a renderer that loads templates from S3
type S3RendererBase[T any] struct {
	s3Client *s3.Client
	bucket   string
	Data     T
}

var _ contractsProviders.IRendererProvider[any] = (*S3RendererBase[any])(nil)

// NewS3RendererBase creates a new S3 renderer instance
func NewS3RendererBase[T any](bucket string) (*S3RendererBase[T], error) {
	ctx := context.Background()
	awsCfg, err := awsconfig.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	return &S3RendererBase[T]{
		s3Client: s3.NewFromConfig(awsCfg),
		bucket:   bucket,
	}, nil
}

// Render renders a template from S3 or local filesystem
func (r *S3RendererBase[T]) Render(templatePath string, data T) (string, *application_errors.ApplicationError) {
	var templateContent string
	var err error

	// Check if template path is an S3 path (s3://bucket/key)
	if strings.HasPrefix(templatePath, "s3://") {
		templateContent, err = r.loadFromS3(templatePath)
		if err != nil {
			return "", application_errors.NewApplicationError(
				status.ProviderError,
				messages.MessageKeysInstance.SOMETHING_WENT_WRONG,
				fmt.Sprintf("failed to load template from S3: %v", err),
			)
		}
	} else if r.bucket != "" {
		// If bucket is configured, construct S3 path from local path
		// templatePath is like "templates/emails/new_user_en.gohtml"
		// Convert to "s3://bucket/templates/emails/new_user_en.gohtml"
		s3Path := fmt.Sprintf("s3://%s/%s", r.bucket, strings.TrimPrefix(templatePath, "/"))
		templateContent, err = r.loadFromS3(s3Path)
		if err != nil {
			return "", application_errors.NewApplicationError(
				status.ProviderError,
				messages.MessageKeysInstance.SOMETHING_WENT_WRONG,
				fmt.Sprintf("failed to load template from S3: %v", err),
			)
		}
	} else {
		// Fallback to local filesystem
		templateContent, err = r.loadFromLocal(templatePath)
		if err != nil {
			return "", application_errors.NewApplicationError(
				status.InternalError,
				messages.MessageKeysInstance.SOMETHING_WENT_WRONG,
				err.Error(),
			)
		}
	}

	// Parse and execute template
	tmpl, err := template.New("email").Parse(templateContent)
	if err != nil {
		return "", application_errors.NewApplicationError(
			status.InternalError,
			messages.MessageKeysInstance.SOMETHING_WENT_WRONG,
			err.Error(),
		)
	}

	var rendered bytes.Buffer
	if err := tmpl.Execute(&rendered, data); err != nil {
		return "", application_errors.NewApplicationError(
			status.InternalError,
			messages.MessageKeysInstance.SOMETHING_WENT_WRONG,
			err.Error(),
		)
	}

	return rendered.String(), nil
}

// loadFromS3 loads a template from S3
func (r *S3RendererBase[T]) loadFromS3(s3Path string) (string, error) {
	// Parse S3 path: s3://bucket/key
	path := strings.TrimPrefix(s3Path, "s3://")
	parts := strings.SplitN(path, "/", 2)
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid S3 path format: %s", s3Path)
	}

	bucket := parts[0]
	key := parts[1]

	ctx := context.Background()
	result, err := r.s3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return "", fmt.Errorf("failed to get object from S3: %w", err)
	}
	defer result.Body.Close()

	var buf bytes.Buffer
	_, err = buf.ReadFrom(result.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read S3 object body: %w", err)
	}

	return buf.String(), nil
}

// loadFromLocal loads a template from local filesystem (fallback)
func (r *S3RendererBase[T]) loadFromLocal(templatePath string) (string, error) {
	data, err := os.ReadFile(templatePath)
	if err != nil {
		return "", fmt.Errorf("failed to read template file: %w", err)
	}
	return string(data), nil
}

// S3RenderNewUserEmail renders new user emails from S3
type S3RenderNewUserEmail struct {
	*S3RendererBase[email_models.NewUserEmailData]
}

// S3RenderResetPasswordEmail renders reset password emails from S3
type S3RenderResetPasswordEmail struct {
	*S3RendererBase[email_models.ResetPasswordEmailData]
}

// S3RenderOTPEmail renders OTP emails from S3
type S3RenderOTPEmail struct {
	*S3RendererBase[email_models.OneTimePasswordEmailData]
}

// NewS3RenderProviders creates all S3 render providers
func NewS3RenderProviders(bucket string) (*S3RenderNewUserEmail, *S3RenderResetPasswordEmail, *S3RenderOTPEmail, error) {
	baseNewUser, err := NewS3RendererBase[email_models.NewUserEmailData](bucket)
	if err != nil {
		return nil, nil, nil, err
	}

	baseResetPassword, err := NewS3RendererBase[email_models.ResetPasswordEmailData](bucket)
	if err != nil {
		return nil, nil, nil, err
	}

	baseOTP, err := NewS3RendererBase[email_models.OneTimePasswordEmailData](bucket)
	if err != nil {
		return nil, nil, nil, err
	}

	return &S3RenderNewUserEmail{baseNewUser},
		&S3RenderResetPasswordEmail{baseResetPassword},
		&S3RenderOTPEmail{baseOTP},
		nil
}
