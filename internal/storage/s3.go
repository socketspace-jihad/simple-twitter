package storage

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/rs/zerolog/log"
)

var S3Client *S3Storage

type S3Storage struct {
	client *s3.Client
	bucket string
	region string
}

func init() {
	bucket := os.Getenv("S3_BUCKET")
	region := os.Getenv("S3_REGION")
	accessKey := os.Getenv("S3_ACCESS_KEY")
	secretKey := os.Getenv("S3_SECRET_KEY")
	endpoint := os.Getenv("S3_ENDPOINT")

	if bucket == "" || region == "" || accessKey == "" || secretKey == "" {
		log.Warn().Msg("S3 configuration not set, image upload will be disabled")
		return
	}

	creds := credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")

	var cfg aws.Config
	var err error

	if endpoint != "" {
		cfg, err = config.LoadDefaultConfig(context.Background(),
			config.WithRegion(region),
			config.WithCredentialsProvider(creds),
		)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to load AWS config")
		}
		S3Client = &S3Storage{
			client: s3.NewFromConfig(cfg, func(o *s3.Options) {
				o.BaseEndpoint = aws.String(endpoint)
				o.UsePathStyle = true
			}),
			bucket: bucket,
			region: region,
		}
	} else {
		cfg, err = config.LoadDefaultConfig(context.Background(),
			config.WithRegion(region),
			config.WithCredentialsProvider(creds),
		)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to load AWS config")
		}
		S3Client = &S3Storage{
			client: s3.NewFromConfig(cfg),
			bucket: bucket,
			region: region,
		}
	}

	log.Info().Str("bucket", bucket).Str("region", region).Msg("S3 storage initialized")
}

func (s *S3Storage) Upload(data []byte, key string, contentType string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(data),
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload to S3: %w", err)
	}

	endpoint := os.Getenv("S3_ENDPOINT")
	var url string
	if endpoint != "" {
		url = fmt.Sprintf("%s/%s/%s", endpoint, s.bucket, key)
	} else {
		url = fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", s.bucket, s.region, key)
	}

	return url, nil
}
