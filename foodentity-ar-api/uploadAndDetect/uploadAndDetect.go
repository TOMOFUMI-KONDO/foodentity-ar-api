package uploadAndDetect

import (
	"fmt"
	"foodentity-ar-api/detectLabel"
	"foodentity-ar-api/dynamoDB"
	"foodentity-ar-api/localFile"
	"foodentity-ar-api/model"
	"foodentity-ar-api/s3"
)

func NewUploadAndDetectUseCase(
	request *model.Request,
	fileRepository localFile.LocalFileRepository,
	s3Repository s3.S3Repository,
	detectLabelRepository detectLabel.DetectLabelRepository,
	fetchFoodRepository dynamoDB.FetchFoodRepository,
) UploadAndDetectUseCase {
	return uploadAndDetectUseCaseImpl{
		Request:               request,
		LocalFileRepository:   fileRepository,
		S3Repository:          s3Repository,
		DetectLabelRepository: detectLabelRepository,
		FetchFoodRepository:   fetchFoodRepository,
	}
}

type UploadAndDetectUseCase interface {
	Exec(imageName string) (*model.Response, error)
}

type uploadAndDetectUseCaseImpl struct {
	Request               *model.Request
	LocalFileRepository   localFile.LocalFileRepository
	S3Repository          s3.S3Repository
	DetectLabelRepository detectLabel.DetectLabelRepository
	FetchFoodRepository   dynamoDB.FetchFoodRepository
}

func (impl uploadAndDetectUseCaseImpl) Exec(imageName string) (*model.Response, error) {
	if localFileErr := impl.LocalFileRepository.Add(impl.Request, imageName); localFileErr != nil {
		return nil, localFileErr
	}

	if s3Err := impl.S3Repository.Add(imageName); s3Err != nil {
		return nil, s3Err
	}

	result, detectErr := impl.DetectLabelRepository.Detect(imageName)
	if detectErr != nil {
		return nil, detectErr
	}

	var labels []string
	for _, label := range result.Labels {
		labels = append(labels, *label.Name)
	}
	fmt.Printf("labels: %v\n", labels)

	response, fetchErr := impl.FetchFoodRepository.Fetch(labels)
	//response, fetchErr := impl.FetchFoodRepository.Fetch([]string{"hum"})
	if fetchErr != nil {
		return nil, fetchErr
	}

	return response, nil
}
