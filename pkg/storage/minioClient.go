package storage

// var MinioClient *minio.Client
import (
	"log"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Client = minio.Client
type Object = minio.Object
type PutObjectOptions = minio.PutObjectOptions
type RemoveObjectOptions = minio.RemoveObjectOptions
type GetObjectOptions = minio.GetObjectOptions

type MinioStorage struct {
	client *minio.Client
	bucket string
}

func NewMinioClient() *minio.Client {
	endpoint := os.Getenv("MINIO_ENDPOINT")
	accessKey := os.Getenv("MINIO_ROOT_USER")
	secretKey := os.Getenv("MINIO_ROOT_PASSWORD")
	useSSL := os.Getenv("MINIO_USE_SSL") == "true"

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})

	if err != nil {
		log.Fatalf("Error conectando a MinIO: %v", err)
	}

	// MinioClient = client
	log.Printf("✅ Successfully connected to MinIO at %s\n", endpoint)

	return client
}
