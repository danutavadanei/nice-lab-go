package aws

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"io"
	"log"
	"net/http"
)

type Client interface {
	ListBuckets(ctx context.Context) ([]byte, error)
	ListBucketObjects(ctx context.Context, bucket *string) ([]byte, error)
	SinkFileToWriter(ctx context.Context, bucket *string, key *string, w http.ResponseWriter) (int64, error)
	UploadFileToBucket(ctx context.Context, bucket *string, key *string, body io.Reader) ([]byte, error)
}

type client struct {
	sess       *session.Session
	s3Client   *s3.S3
	s3Uploader *s3manager.Uploader
}

func (c client) ListBuckets(ctx context.Context) ([]byte, error) {
	req, resp := c.s3Client.ListBucketsRequest(&s3.ListBucketsInput{})

	err := req.Send()

	if err != nil {
		log.Printf("error listing s3 buckets: %v", err)

		return nil, err
	}

	bytes, err := json.Marshal(resp)

	if err != nil {
		log.Printf("error encoding s3 response: %v", err)

		return nil, err
	}

	return bytes, nil
}

func (c client) ListBucketObjects(ctx context.Context, bucket *string) ([]byte, error) {
	resp, err := c.s3Client.ListObjects(&s3.ListObjectsInput{Bucket: bucket})

	if err != nil {
		log.Printf("error listing s3 bucket (%s): %v", *bucket, err)
		return nil, err
	}

	bytes, err := json.Marshal(resp)

	if err != nil {
		log.Printf("error encoding s3 response: %v", err)
		return nil, err
	}

	return bytes, nil
}

func (c client) SinkFileToWriter(ctx context.Context, bucket *string, key *string, w http.ResponseWriter) (int64, error) {
	o, err := c.s3Client.GetObject(&s3.GetObjectInput{
		Key:    aws.String(*key),
		Bucket: aws.String(*bucket),
	})

	if err != nil {
		log.Printf("error getting file with key: %s from s3 bucket (%s): %v", *key, *bucket, err)

		return 0, err
	}

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", *key))
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Content-Type", *o.ContentType)

	bytesWritten, err := io.Copy(w, o.Body)

	if err != nil {
		log.Printf("error copying file with key: %s to the http response from s3 bucket (%s): %v", *key, *bucket, err)

		return 0, err
	}

	return bytesWritten, err
}

func (c client) UploadFileToBucket(ctx context.Context, bucket *string, key *string, body io.Reader) ([]byte, error) {
	o, err := c.s3Uploader.Upload(&s3manager.UploadInput{
		Body:   body,
		Bucket: aws.String(*bucket),
		Key:    aws.String(*key),
	})

	if err != nil {
		log.Printf("error uploading file to s3 bucket(%s): %v", *bucket, err)

		return nil, err
	}

	bytes, err := json.Marshal(o)

	if err != nil {
		log.Printf("error encoding s3 response: %v", err)
		return nil, err
	}

	return bytes, nil
}

func NewClient(cfg *aws.Config) Client {
	sess := session.Must(session.NewSession(cfg))

	return client{
		sess:       sess,
		s3Client:   s3.New(sess),
		s3Uploader: s3manager.NewUploader(sess),
	}
}
