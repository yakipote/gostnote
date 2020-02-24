package firebase

import (
	"cloud.google.com/go/storage"
	"context"
	firebase "firebase.google.com/go"
	"fmt"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

type Config struct {
	KeyPath       string
	StorageBucket string
}

type Firebase struct {
	bucket *storage.BucketHandle
}

func NewStorage(config Config) *Firebase {
	// firebase 初期化
	firebaseConfig := &firebase.Config{
		StorageBucket: config.StorageBucket,
	}
	opt := option.WithCredentialsFile(config.KeyPath)
	app, err := firebase.NewApp(context.Background(), firebaseConfig, opt)
	if err != nil {
		log.Fatalln(err)
	}
	// get storage bucket
	client, err := app.Storage(context.Background())
	if err != nil {
		log.Fatalln(err)
	}

	bucket, err := client.DefaultBucket()
	if err != nil {
		log.Fatalln(err)
	}

	return &Firebase{
		bucket: bucket,
	}

}

func (f Firebase) Upload(file *os.File) {
	// upload file
	contentType := "text/plain"
	ctx := context.Background()

	remoteFilename := filepath.Base(file.Name())
	writer := f.bucket.Object(remoteFilename).NewWriter(ctx)
	writer.ObjectAttrs.ContentType = contentType
	writer.ObjectAttrs.CacheControl = "no-cache"
	writer.ObjectAttrs.ACL = []storage.ACLRule{
		{
			Entity: storage.AllUsers,
			Role:   storage.RoleReader,
		},
	}
	if _, err := io.Copy(writer, file); err != nil {
		log.Fatalln(err)
	}

	if err := writer.Close(); err != nil {
		log.Fatalln(err)
	}
}
func (f Firebase) List() error {
	ctx := context.Background()
	it := f.bucket.Objects(ctx, nil)
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		fmt.Println(attrs.Name)
	}
	return nil
}

func (f Firebase) GetFileList() ([]string, error) {
	list := []string{}
	ctx := context.Background()
	it := f.bucket.Objects(ctx, nil)
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		list = append(list, attrs.Name)
	}
	return list, nil
}

func (f Firebase) Download(fileName string) []byte {
	ctx := context.Background()
	rc, err := f.bucket.Object(fileName).NewReader(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	defer rc.Close()

	data, err := ioutil.ReadAll(rc)
	if err != nil {
		log.Fatalln(err)
	}
	return data
}
