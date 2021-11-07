package s3

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"path"

	"arhat.dev/rs"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type inputS3Sepc struct {
	rs.BaseField `yaml:"-"`

	Path string `yaml:"path"`

	Config rendererS3Config `yaml:"config"`
}

type rendererS3Config struct {
	rs.BaseField `yaml:"-"`

	EndpointURL string `yaml:"endpoint_url"`
	Region      string `yaml:"region"`

	Bucket   string `yaml:"bucket"`
	BasePath string `yaml:"base_path"`

	AccessKeyID     string `yaml:"access_key_id"`
	AccessKeySecret string `yaml:"access_key_secret"`
}

type s3Client struct {
	client *minio.Client

	bucket   string
	region   string
	basePath string
}

func (c *s3Client) download(ctx context.Context, objPath string) (io.ReadCloser, error) {
	obj, err := c.client.GetObject(
		ctx,
		c.bucket,
		path.Join(c.basePath, objPath),
		minio.GetObjectOptions{},
	)
	if err != nil {
		return nil, err
	}

	return obj, err
}

func (c *rendererS3Config) createClient() (*s3Client, error) {
	eURL, err := url.Parse(c.EndpointURL)
	if err != nil {
		return nil, fmt.Errorf("invalid endpoint url: %w", err)
	}

	client, err := minio.New(eURL.Host, &minio.Options{
		Creds:  credentials.NewStaticV4(c.AccessKeyID, c.AccessKeySecret, ""),
		Secure: eURL.Scheme == "https",
		Region: c.Region,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create s3 client: %w", err)
	}

	return &s3Client{
		client: client,

		bucket:   c.Bucket,
		region:   c.Region,
		basePath: c.BasePath,
	}, nil
}
