package uploadAndDetect

import (
	"foodentity-ar-api/detectLabel"
	"foodentity-ar-api/localFile"
	"foodentity-ar-api/model"
	"foodentity-ar-api/s3"

	"github.com/aws/aws-sdk-go/service/rekognition"
)

func NewUploadAndDetectUseCase(
	request *model.Request,
	fileRepository localFile.LocalFileRepository,
	s3Repository s3.S3Repository,
	detectLabelRepository detectLabel.DetectLabelRepository,
) UploadAndDetectUseCase {
	return uploadAndDetectUseCaseImpl{
		Request:               request,
		LocalFileRepository:   fileRepository,
		S3Repository:          s3Repository,
		DetectLabelRepository: detectLabelRepository,
	}
}

type UploadAndDetectUseCase interface {
	Exec(imageName string) (*rekognition.DetectLabelsOutput, error)
}

type uploadAndDetectUseCaseImpl struct {
	Request               *model.Request
	LocalFileRepository   localFile.LocalFileRepository
	S3Repository          s3.S3Repository
	DetectLabelRepository detectLabel.DetectLabelRepository
}

func (impl uploadAndDetectUseCaseImpl) Exec(imageName string) (*rekognition.DetectLabelsOutput, error) {
	if localFileErr := impl.LocalFileRepository.Add(impl.Request, imageName); localFileErr != nil {
		return nil, localFileErr
	}

	if s3Err := impl.S3Repository.Add(imageName); s3Err != nil {
		return nil, s3Err
	}

	res, detectErr := impl.DetectLabelRepository.Detect(imageName)
	if detectErr != nil {
		return nil, detectErr
	}

	return res, nil
}
