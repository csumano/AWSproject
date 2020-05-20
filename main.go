package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/gorilla/mux"
	"net/http"
	"os"
)

//type item struct {
//	Zero1Symbol       string
//	Zero3High         string
//	Zero4Low          string
//	Zero5Price        string
//	Zero6Volume       string
//	Zero9Change       string
//	One0ChangePercent string
//
//}

type statusStruct struct {
	RecordCount int64
	Table       string
}

func allHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Print("test")
	result := all()
	fmt.Fprintf(w, result)
}

func statusHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Print("test")
	result := status()
	fmt.Fprintf(w, result)

}

func status() string {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)
	if err != nil {
		fmt.Println("Error creating session:")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// Create DynamoDB client
	svc := dynamodb.New(sess)
	input := &dynamodb.ScanInput{
		TableName: aws.String("AmazonStocks"),
	}

	result, err := svc.Scan(input)

	if err != nil {
		fmt.Println("Got error calling PutItem:")
		fmt.Println(err.Error())
		os.Exit(1)
	}
	count := result.Count
	status := statusStruct{*count, "AmazonStocks"}
	r, err := json.Marshal(status)
	return string(r)

}

func all() string {

	// Initialize a session that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials
	// and region from the shared configuration file ~/.aws/config.
	//sess, := session.Must(session.NewSessionWithOptions(session.Options{
	//	SharedConfigState: session.SharedConfigEnable,
	//}))

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)
	if err != nil {
		fmt.Println("Error creating session:")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// Create DynamoDB client
	svc := dynamodb.New(sess)
	input := &dynamodb.ScanInput{
		TableName: aws.String("AmazonStocks"),
	}

	result, err := svc.Scan(input)

	if err != nil {
		fmt.Println("Got error calling PutItem:")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	return result.String()

}

func main() {

	r := mux.NewRouter()
	r.HandleFunc("/csumano/all", allHTTP)
	r.HandleFunc("/csumano/status", statusHTTP)
	http.ListenAndServe(":8080", r)
}
