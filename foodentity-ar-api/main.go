package main

import (
	"encoding/json"
	"fmt"
	"foodentity-ar-api/detectLabel"
	"foodentity-ar-api/localFile"
	"foodentity-ar-api/model"
	"foodentity-ar-api/s3"
	"foodentity-ar-api/uploadAndDetect"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	fmt.Println("Start lambda process.")

	req, err := convertJsonToStruct(request.Body)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return events.APIGatewayProxyResponse{
			Body:       err.Error(),
			StatusCode: 500,
		}, err
	}

	fileRepository := localFile.NewLocalFileRepositoryImpl()
	s3Repository := s3.NewS3RepositoryImpl()
	detectLabelRepository := detectLabel.NewDetectLabelRepositoryImpl()
	jstCurrentUnixTime := int(time.Now().UTC().In(time.FixedZone("Asia/Tokyo", 9*60*60)).Unix())
	imageName := strconv.Itoa(jstCurrentUnixTime) + ".png"

	detectResult, useCaseErr := uploadAndDetect.NewUploadAndDetectUseCase(
		req,
		fileRepository,
		s3Repository,
		detectLabelRepository,
	).Exec(imageName)
	if useCaseErr != nil {
		fmt.Printf("useCaseError: %v\n", useCaseErr)
		return events.APIGatewayProxyResponse{
			Body:       useCaseErr.Error(),
			StatusCode: 500,
		}, useCaseErr
	}

	body, jsonEncodeErr := json.Marshal(detectResult.Labels)
	if jsonEncodeErr != nil {
		fmt.Printf("jsonENcodeError: %v\n", jsonEncodeErr)
		return events.APIGatewayProxyResponse{
			Body:       jsonEncodeErr.Error(),
			StatusCode: 500,
		}, useCaseErr
	}

	fmt.Println("Finish lambda process.")

	return events.APIGatewayProxyResponse{
		Body:       string(body),
		StatusCode: 200,
		Headers:    map[string]string{"Content-Type": "application/json"},
	}, nil
}

func convertJsonToStruct(inputs string) (*model.Request, error) {
	var req model.Request
	err := json.Unmarshal([]byte(inputs), &req)
	if err != nil {
		return nil, err
	}
	return &req, nil
}

func main() {
	lambda.Start(handler)
}
