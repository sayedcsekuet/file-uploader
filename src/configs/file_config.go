package configs

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	logger "github.com/sirupsen/logrus"
	"gocloud.dev/blob"
	"gocloud.dev/blob/fileblob"
	"gocloud.dev/blob/s3blob"
	"golang.org/x/net/context"
)

type FileStorageConfig struct {
	Provider         string
	Profile          string
	Endpoint         string
	DisableSSL       bool
	Region           string
	AccessKeyID      string
	SecretAccessKey  string
	S3ForcePathStyle bool
}
type FileStorage struct {
	*FileStorageConfig
	Ctx context.Context
}

func NewFileStorage(config *FileStorageConfig) *FileStorage {
	return &FileStorage{config, context.Background()}
}

func (fs *FileStorage) Open(path string, provider string) *blob.Bucket {
	if provider == "" {
		provider = fs.Provider
	}
	var bucket *blob.Bucket
	var err error
	switch provider {
	case "file":
		bucket, err = fileblob.OpenBucket(path, &fileblob.Options{CreateDir: true})
	case "s3":
		config := aws.Config{
			Endpoint:                      aws.String(fs.Endpoint),
			DisableSSL:                    aws.Bool(fs.DisableSSL),
			Region:                        aws.String(fs.Region),
			S3ForcePathStyle:              aws.Bool(fs.S3ForcePathStyle),
			CredentialsChainVerboseErrors: aws.Bool(true),
			LogLevel:                      aws.LogLevel(4),
		}
		options := session.Options{Config: config}
		if fs.Profile != "" {
			options.Profile = fs.Profile
		}
		if fs.AccessKeyID != "" && fs.SecretAccessKey != "" {
			config.Credentials = credentials.NewStaticCredentials(fs.AccessKeyID, fs.SecretAccessKey, "")
		}
		s3Session, err := session.NewSessionWithOptions(options)
		if err != nil {
			logger.Fatal(err, fs.Ctx)
		}
		bucket, err = s3blob.OpenBucket(fs.Ctx, s3Session, path, nil)
	default:
		logger.Fatal("not supported file service", fs.Ctx)
	}
	if err != nil {
		logger.Fatal(err, fs.Ctx)
	}
	return bucket
}
