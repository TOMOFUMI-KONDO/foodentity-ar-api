package main

import (
	"encoding/json"
	"fmt"
	"foodentity-ar-api/detectLabel"
	"foodentity-ar-api/dynamoDB"
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

	req, convertJsonError := convertJsonToStruct(request.Body)
	if convertJsonError != nil {
		return handleError("convertJsonError", convertJsonError)
	}

	fileRepository := localFile.NewLocalFileRepositoryImpl()
	s3Repository := s3.NewS3RepositoryImpl()
	detectLabelRepository := detectLabel.NewDetectLabelRepositoryImpl()
	fetchFoodRepository := dynamoDB.NewFetchFoodRepositoryImpl()
	jstCurrentUnixTime := int(time.Now().UTC().In(time.FixedZone("Asia/Tokyo", 9*60*60)).Unix())
	imageName := strconv.Itoa(jstCurrentUnixTime) + ".png"

	response, useCaseErr := uploadAndDetect.NewUploadAndDetectUseCase(
		req,
		fileRepository,
		s3Repository,
		detectLabelRepository,
		fetchFoodRepository,
	).Exec(imageName)
	if useCaseErr != nil {
		return handleError("useCaseErr", useCaseErr)
	}

	body, jsonEncodeErr := json.Marshal(response)
	if jsonEncodeErr != nil {
		return handleError("jsonEncode", jsonEncodeErr)
	}

	fmt.Println("Finish lambda process.")

	//body, _ := json.Marshal(model.Response{
	//	Food:       "hum",
	//	Identities: []string{"cochineal", "nitrous acid", "chemical protein"},
	//})
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

func handleError(key string, error error) (events.APIGatewayProxyResponse, error) {
	fmt.Printf(key+": %v\n", error)
	return events.APIGatewayProxyResponse{
		Body:       error.Error(),
		StatusCode: 500,
	}, error
}

func main() {
	lambda.Start(handler)
}
