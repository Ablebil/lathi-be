package minio

import (
	"fmt"

	"github.com/Ablebil/lathi-be/internal/config"
	m "github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioItf interface {
	GetObjectURL(object string) string
}

type minio struct {
	client        *m.Client
	endpoint      string
	publicBaseURL string
	bucket        string
}

func New(env *config.Env) (MinioItf, error) {
	client, err := m.New(env.StorageEndpoint, &m.Options{
		Creds:  credentials.NewStaticV4(env.StorageAccessKey, env.StorageSecretKey, ""),
		Secure: true,
	})
	if err != nil {
		return nil, err
	}

	return &minio{
		client:        client,
		endpoint:      env.StorageEndpoint,
		publicBaseURL: env.StoragePublicUrl,
		bucket:        env.StorageBucket,
	}, nil
}

func (m *minio) GetObjectURL(object string) string {
	if object == "" {
		return ""
	}
	return fmt.Sprintf("https://%s/%s/%s", m.publicBaseURL, m.bucket, object)
}
