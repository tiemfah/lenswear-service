package bucket

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"strconv"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
)

type GCPBucket struct {
	ctx       context.Context
	client    *storage.Client
	projectID string
}

type Bucket interface {
	UploadFile(bucketName, folderName string, file []byte) ([]string, error)
	DeleteFolder(bucketName, folderName string) error
}

func NewGCPBucket(ctx context.Context, client *storage.Client, projectID string) *GCPBucket {
	return &GCPBucket{
		ctx:       ctx,
		client:    client,
		projectID: projectID,
	}
}

func (g *GCPBucket) UploadFile(bucketName, folderName string, form *multipart.Form) ([]string, error) {
	resFilePaths := []string{}
	bucket := g.client.Bucket(g.projectID + "-" + bucketName)
	for i, mf := range form.File["images"] {
		filePath := folderName + "/" + strconv.Itoa(i) + ".jpg"
		writer := bucket.Object(filePath).NewWriter(g.ctx)
		file, err := mf.Open()
		if err != nil {
			return nil, fmt.Errorf("mf.Open: %v", err)
		}
		if _, err := io.Copy(writer, file); err != nil {
			return nil, fmt.Errorf("io.Copy: %v", err)
		}
		if err := writer.Close(); err != nil {
			return nil, fmt.Errorf("wc.Close: %v", err)
		}
		if err := file.Close(); err != nil {
			return nil, fmt.Errorf("file.Close: %v", err)
		}
		resFilePaths = append(resFilePaths, filePath)
	}
	trainBucket := g.client.Bucket(g.projectID + "-" + bucketName + "-train")
	for i, mf := range form.File["train_images"] {
		filePath := folderName + "/" + strconv.Itoa(i) + ".jpg"
		writer := trainBucket.Object(filePath).NewWriter(g.ctx)
		file, err := mf.Open()
		if err != nil {
			return nil, fmt.Errorf("mf.Open: %v", err)
		}
		if _, err := io.Copy(writer, file); err != nil {
			return nil, fmt.Errorf("io.Copy: %v", err)
		}
		if err := writer.Close(); err != nil {
			return nil, fmt.Errorf("wc.Close: %v", err)
		}
		if err := file.Close(); err != nil {
			return nil, fmt.Errorf("file.Close: %v", err)
		}
	}
	return resFilePaths, nil
}

func (g *GCPBucket) DeleteFolder(bucketName, folderName string) error {
	bucket := g.client.Bucket(g.projectID + "-" + bucketName)
	it := bucket.Objects(g.ctx, &storage.Query{Prefix: folderName})
	for {
		objAttrs, err := it.Next()
		if err != nil && err != iterator.Done {
			return err
		}
		if err == iterator.Done {
			break
		}
		if err := bucket.Object(objAttrs.Name).Delete(g.ctx); err != nil {
			return err
		}
	}
	trainBucket := g.client.Bucket(g.projectID + "-" + bucketName + "-train")
	trainIt := trainBucket.Objects(g.ctx, &storage.Query{Prefix: folderName})
	for {
		objAttrs, err := trainIt.Next()
		if err != nil && err != iterator.Done {
			return err
		}
		if err == iterator.Done {
			break
		}
		if err := trainBucket.Object(objAttrs.Name).Delete(g.ctx); err != nil {
			return err
		}
	}
	return nil
}
