package models

import (
	"cloud.google.com/go/storage"
	"context"
	"google.golang.org/api/option"
	"io"
	"log"
)

func StoreGCS(bucketName string, uploadedFile *UploadedFile) (*storage.ObjectAttrs, error) {
	ctx := context.Background()
	opt := option.WithCredentialsFile("./config/credentials/backend-versioning-me-dev-e743b229e23b.json")

	client, err := storage.NewClient(ctx, opt)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer client.Close()

	obj := client.Bucket(bucketName).Object(uploadedFile.VersionName + "/" + uploadedFile.UUID)

	w := obj.NewWriter(ctx)
	w.ContentType = uploadedFile.Mime
	w.ChunkSize = 1024

	if _, err = io.Copy(w, uploadedFile.Reader); err != nil {
		return nil, err
	}
	if err = w.Close(); err != nil {
		return nil, err
	}

	if err := obj.ACL().Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
		return nil, err
	}

	attrs, err := obj.Attrs(ctx)
	if err != nil {
		return nil, err
	}

	return attrs, nil
}
