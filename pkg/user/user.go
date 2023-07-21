package user

import(
	"encoding/json"
	"errors"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/aws/aws-sdk-go/dynamodb/dynamodbattribute"
)

var(
	ErrorFailedToFetchRecord = "Failed To Fetch Record"
	ErrorFailedToUnmarshallRecord = "Failed To Unmarshall Record"
)

type User struct{
	Email string `json:"email"`
	FirstName string `json:"firstName"`
	LastName string `json:"lastName"`
}

func FetchUser(email, tableName string, dynaClient dynamodbiface.DynamoDBAPI)(*User, error){
	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"email":{
				S: aws.String(email)
			}
		},
		TableName: aws.String(tableName)
	}

	result, err := dynaClient.GetItem(input)

	if err!=nil{
		return nil, errors.New(ErrorFailedToFetchRecord)
	}

	item := new(User)
	err = dynamodbattribute.UnmarshallMap(result.Item, item)
	if err!=nil{
		return nil, errors.New(ErrorFailedToUnmarshallRecord)
	}
	return item, nil

}

func FetchUsers(tableName string, dynaClient dynamodbiface.DynamoDBAPI)(*[]User, error){
	input := &dynamodb.ScanInput{
		TableName: aws.String(tableName)
	}

	result := dynaClient.Scan(input)
	if err!=nill{
		return nil, errors.New(ErrorFailedToFetchRecord)
	}
	item := new([]User)
	err = dynamodbattribute.UnmarshallMap(result.Items, item)
	return item, nil
}

func CreateUser()(){

}

func UpdateUser()(){

}

func DeleteUser()error{

}