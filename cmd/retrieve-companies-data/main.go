package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"github.com/google/uuid"
)

func uploadNewFileHandler(e events.S3Event) error {
	eventJsonData, _ := json.Marshal(e)

	log.Println("CompaniesDataTableName: " + os.Getenv("CompaniesDataTableName"))
	log.Println("handler was called - " + string(eventJsonData))

	sess := session.Must(session.NewSession())
	downloader := s3manager.NewDownloader(sess)
	dynamodbSvc := dynamodb.New(sess)

	r := regexp.MustCompile(`([\W]+)`)

	for _, rec := range e.Records {
		log.Println("File", rec.S3.Object.Key, "processing", "...")

		filename := "/tmp/" + r.ReplaceAllString(rec.S3.Object.Key, "_") + "--" + fmt.Sprintf("%d", time.Now().UTC().UnixNano())
		file, err := os.Create(filename)
		if err != nil {
			return fmt.Errorf("failed to create file %q, %v", filename, err)
		}

		n, err := downloader.Download(file, &s3.GetObjectInput{
			Bucket: aws.String(rec.S3.Bucket.Name),
			Key:    aws.String(rec.S3.Object.Key),
		})
		if err != nil {
			return fmt.Errorf("failed to download file, %v", err)
		}
		fmt.Printf("file %s downloaded, %d bytes\n", filename, n)

		// TODO: parse file
		err = parseCompDataFile(dynamodbSvc, file)
		if err != nil {
			return fmt.Errorf("failed to parse file, %v", err)
		}

		_ = file.Close()
		_ = os.Remove(filename)
	}

	return nil
}

type Item struct {
	Id    string `dynamodbav:"id"`
	Title string `dynamodbav:"title"`
	Year  int    `dynamodbav:"year"`
}

func parseCompDataFile(db *dynamodb.DynamoDB, file *os.File) error {
	reader := bufio.NewReader(file)

	tableName := os.Getenv("CompaniesDataTableName")

	for {
		line, _, err := reader.ReadLine()
		if nil != err {
			break
		}
		log.Println(string(line))

		data := strings.Split(string(line), ";")

		if 2 != len(data) {
			return fmt.Errorf("wrond format of string \"%s\"", string(line))
		}

		companyTitle := data[0]
		companyFoundationYear, err := strconv.Atoi(data[1])
		if err != nil {
			return err
		}

		titleCond := &dynamodb.Condition{}
		titleCond.SetAttributeValueList([]*dynamodb.AttributeValue{
			{S: aws.String(companyTitle)},
		})
		titleCond.SetComparisonOperator(dynamodb.ComparisonOperatorEq)

		params := &dynamodb.ScanInput{
			TableName:  aws.String(tableName),
			ScanFilter: map[string]*dynamodb.Condition{"title": titleCond},
		}
		result, err := db.Scan(params)
		if err != nil {
			return fmt.Errorf("failed to make Query API call, %v", err)
		}
		if *result.Count > 0 {
			log.Println("Company", companyTitle, "already exists")
			continue
		}

		item := Item{
			Id:    uuid.New().String(),
			Title: companyTitle,
			Year:  companyFoundationYear,
		}
		av, err := dynamodbattribute.MarshalMap(item)
		if err != nil {
			return err
		}

		input := &dynamodb.PutItemInput{
			Item:      av,
			TableName: aws.String(tableName),
		}

		_, err = db.PutItem(input)
		if err != nil {
			return err
		}

	}

	return nil
}

func main() {
	lambda.Start(uploadNewFileHandler)
}
