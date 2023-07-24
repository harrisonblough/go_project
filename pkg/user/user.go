package user

import(
	"github.com/harrisonblough/go_project/pkg/validators"
	"encoding/json"
	"errors"
	"github.com/aws/aws-lambda-go/events"
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/service/dynamodb"
    "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
    "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

var(
	ErrorFailedToFetchRecord = "Failed To Fetch Record"
	ErrorFailedToUnmarshalRecord = "Failed To Unmarshall Record"
	ErrorInvalidUserData = "Invalid User Data"
	ErrorInvalidEmail = "Invalid Email"
	ErrorCouldNotMarshalItem = "Could Not Marshal Item"
	ErrorCouldNotDeleteItem = "Could Not Delete Item"
	ErrorCouldNotDynamoPutItem = "Could Not Dynamo Put Item"
	ErrorUserAlreadyExists = "user.User already exists"
	ErrorUserDoesNotExist = "user.User does not exist"
)


type User struct{
	Email 		string	`json:"email"`
	FirstName	string 	`json:"firstName"`
	LastName	string 	`json:"lastName"`
}

func FetchUser(email, tableName string, dynaClient dynamodbiface.DynamoDBAPI)(*User, error){
	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"email":{
				S: aws.String(email),
			},
		},
		TableName: aws.String(tableName),
	}

	result, err := dynaClient.GetItem(input)

	if err!=nil{
		return nil, errors.New(ErrorFailedToFetchRecord)
	}

	item := new(User)
	err = dynamodbattribute.UnmarshalMap(result.Item, item)
	if err!=nil{
		return nil, errors.New(ErrorFailedToUnmarshalRecord)
	}
	return item, nil

}

func FetchUsers(tableName string, dynaClient dynamodbiface.DynamoDBAPI)(*[]User, error){
	input := &dynamodb.ScanInput{
		TableName: aws.String(tableName),
	}

	result,err := dynaClient.Scan(input)
	if err!=nil{
		return nil, errors.New(ErrorFailedToFetchRecord)
	}
	item := new([]User)
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, item)
	return item, nil
}

func CreateUser(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI)(*User, error){
	var u User
	err := json.Unmarshal([]byte(req.Body), &u)
	if err!=nil {
		return nil, errors.New(ErrorInvalidUserData)
	}
	if !validators.IsEmailValid(u.Email){
		return nil, errors.New(ErrorInvalidEmail)
	}

	currentUser, _ := FetchUser(u.Email, tableName, dynaClient)
	if currentUser != nil && len(currentUser.Email) != 0 {
		return nil, errors.New(ErrorUserAlreadyExists)
	}

	av, err := dynamodbattribute.MarshalMap(u)

	if err != nil {
		return nil, errors.New(ErrorCouldNotMarshalItem)
	}

	input := &dynamodb.PutItemInput{
		Item: av,
		TableName: aws.String(tableName),
	}

	_, err = dynaClient.PutItem(input)
	if err != nil {
		return nil, errors.New(ErrorCouldNotDynamoPutItem)
	}
	return &u, nil

}

func UpdateUser(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI)(
	*User, error){

	var u User
	if err:=json.Unmarshal([]byte(req.Body), &u); err!=nil{
		return nil, errors.New(ErrorInvalidEmail)
	}
	currentUser,_ := FetchUser(u.Email, tableName, dynaClient)
	if currentUser!=nil && len(currentUser.Email) == 0 {
		return nil, errors.New(ErrorUserDoesNotExist)
	}

	av,err := dynamodbattribute.MarshalMap(u)

	if err!=nil{
		return nil, errors.New(ErrorCouldNotMarshalItem)
	}

	input := &dynamodb.PutItemInput{
		Item: av,
		TableName: aws.String(tableName),

	}

	_, err = dynaClient.PutItem(input)

	if err!=nil{
		return nil, errors.New(ErrorCouldNotDynamoPutItem)
	}
	return &u, nil
}

func DeleteUser(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI)error{

	email := req.QueryStringParameters["email"]
	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"email":{
				S: aws.String(email),
			},
		},
		TableName: aws.String(tableName),
	}
	_, err := dynaClient.DeleteItem(input)

	if err!=nil{
		return errors.New(ErrorCouldNotDeleteItem)
	}

	return nil
}