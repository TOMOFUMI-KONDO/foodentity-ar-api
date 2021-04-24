package detectLabel

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rekognition"
)

func NewDetectLabelRepositoryImpl() DetectLabelRepository {
	return detectLabelRepositoryImpl{}
}

type DetectLabelRepository interface {
	Detect(imageName string) (*rekognition.DetectLabelsOutput, error)
}

type detectLabelRepositoryImpl struct{}

func (impl detectLabelRepositoryImpl) Detect(imageName string) (*rekognition.DetectLabelsOutput, error) {
	fmt.Println("Start detection label.")
	svc := rekognition.New(impl.newSession())

	s3object := &rekognition.S3Object{
		Bucket: aws.String(os.Getenv("BUCKET_NAME")),
		Name:   aws.String(imageName),
	}
	input := &rekognition.DetectLabelsInput{
		Image: &rekognition.Image{
			S3Object: s3object,
		},
	}

	res, err := svc.DetectLabels(input)
	if err != nil {
		return nil, err
	}
	fmt.Println("Finish detection label.")

	return res, nil
}

func (impl *detectLabelRepositoryImpl) newSession() *session.Session {
	return session.Must(session.NewSession(&aws.Config{
		Region: aws.String("ap-northeast-1"),
	}))
}
