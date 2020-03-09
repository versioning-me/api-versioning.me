package models

import (
	"api-versioning-me/config"
	"cloud.google.com/go/storage"
	"context"
	"google.golang.org/api/option"
	"io"
	"log"
	"os"
)

func StoreGCS(bucketName string, f *File, v *Version) (*storage.ObjectAttrs, error) {
	ctx := context.Background()
	client, err := SetGCSClient(ctx)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer client.Close()

	bucket := client.Bucket(bucketName)
	_, err = bucket.Attrs(ctx)
	if err != nil {
		log.Printf("Don't exist a bucket. Create %s.", bucketName)
		attrs := &storage.BucketAttrs{
			Location:                   "asia-northeast1",
			LocationType:               "region",
		}
		if err := bucket.Create(ctx, "backend-versioning-me-dev", attrs); err != nil {
			log.Fatalf("Can't create bucket. (Reason: %s)", err)
		}
		log.Printf("Success to Create a bucket %+v.", bucket)
	}

	// fileId/versionNum_versionName
	obj := bucket.Object("v0.0.1/" + f.ConvertFileIdToStoring() + "/" + v.ConvertVersionNumToString() + "_" + v.Name)

	w := obj.NewWriter(ctx)
	w.ContentType = v.Mime
	w.ChunkSize = 1024

	if _, err = io.Copy(w, v.Reader); err != nil {
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

func SetGCSClient(ctx context.Context) (*storage.Client, error){
	if config.Env == "localhost" {
		opt := option.WithCredentialsFile(os.Getenv("CREDENTIAL_FILE"))
		client, err := storage.NewClient(ctx, opt)
		if err != nil {
			return nil, err
		}
		return client, nil
	} else {
		client, err := storage.NewClient(ctx)
		if err != nil {
			return nil, err
		}
		return client, nil
	}
}
