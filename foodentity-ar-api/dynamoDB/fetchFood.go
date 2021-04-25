package dynamoDB

import (
	"fmt"
	"foodentity-ar-api/model"
	"os"

	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func NewFetchFoodRepositoryImpl() FetchFoodRepository {
	return fetchFoodRepositoryImpl{}
}

type FetchFoodRepository interface {
	Fetch(foodNames []string) (*model.Response, error)
}

type fetchFoodRepositoryImpl struct{}

func (impl fetchFoodRepositoryImpl) Fetch(foodNames []string) (*model.Response, error) {
	fmt.Println("Start fetching food from DynamoDB.")
	db := dynamodb.New(impl.newSession())

	result := dynamodb.GetItemOutput{}
	for _, foodName := range foodNames {
		itemInput := &dynamodb.GetItemInput{
			TableName: aws.String(os.Getenv("TABLE_NAME")),
			Key: map[string]*dynamodb.AttributeValue{
				"food_name": {
					S: aws.String(foodName),
				},
			},
		}

		tmpResult, err := db.GetItem(itemInput)
		if err != nil {
			return nil, err
		}

		if len(tmpResult.Item) > 0 {
			result = *tmpResult
			break
		}
	}
	fmt.Println("Finish fetching food from DynamoDB.")
	fmt.Printf("fetched food: %v\n", result.Item)

	response := &model.Response{}
	if err := dynamodbattribute.UnmarshalMap(result.Item, &response); err != nil {
		return nil, err
	}

	return response, nil
}

func (impl fetchFoodRepositoryImpl) newSession() *session.Session {
	return session.Must(session.NewSession(&aws.Config{
		Region: aws.String("ap-northeast-1"),
	}))
}
