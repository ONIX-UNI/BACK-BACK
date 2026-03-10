package storage

import (
	"context"
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/minio/minio-go/v7"
)

func EnsureBucket(client *minio.Client, bucketName string) error {
	ctx := context.Background()

	if bucketName == "" {
		return fmt.Errorf("bucket name is empty")
	}

	exists, err := client.BucketExists(ctx, bucketName)
	if err != nil {
		return fmt.Errorf("error checking bucket: %w", err)
	}

	if !exists {
		err = client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return fmt.Errorf("error creating bucket: %w", err)
		}
	}

	return nil
}

/*
minioClient := NewMinioClient()
EnsureBucket(minioClient)
*/

func UploadFile(client *minio.Client, objectName string, filePath string) error {
	ctx := context.Background()
	bucketName := os.Getenv("MINIO_BUCKET")

	_, err := client.FPutObject(ctx, bucketName, objectName, filePath, minio.PutObjectOptions{
		ContentType: "application/pdf",
	})

	return err
}

func UploadHandler(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(400).SendString("Archivo requerido")
	}

	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	client := NewMinioClient()
	bucket := os.Getenv("MINIO_BUCKET")

	_, err = client.PutObject(
		context.Background(),
		bucket,
		file.Filename,
		src,
		file.Size,
		minio.PutObjectOptions{
			ContentType: file.Header.Get("Content-Type"),
		},
	)

	if err != nil {
		return err
	}

	return c.SendString("Archivo subido correctamente")
}
