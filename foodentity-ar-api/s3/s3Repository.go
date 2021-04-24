package s3

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func NewS3RepositoryImpl() S3Repository {
	return s3RepositoryImpl{}
}

type S3Repository interface {
	Add(imageName string) error
}

type s3RepositoryImpl struct{}

func (impl s3RepositoryImpl) Add(imageName string) error {
	file, openErr := os.Open("/tmp/" + imageName)
	if openErr != nil {
		return openErr
	}

	defer file.Close()
	if uploadErr := impl.upload(file, imageName); uploadErr != nil {
		return uploadErr
	}

	return nil
}

func (impl s3RepositoryImpl) upload(file *os.File, imageName string) error {
	uploader := s3manager.NewUploader(impl.newSession())

	fmt.Println("Start uploading to s3.")
	if _, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(os.Getenv("BUCKET_NAME")),
		Key:    aws.String(imageName),
		Body:   file,
	}); err != nil {
		return err
	}
	fmt.Printf("End uploading '%v' to s3.\n", imageName)

	return nil
}

func (impl s3RepositoryImpl) newSession() *session.Session {
	return session.Must(session.NewSession(&aws.Config{
		Region: aws.String("ap-northeast-1"),
	}))
}
