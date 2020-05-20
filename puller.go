package puller

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/landonp1203/goUtils/loggly"
	"github.com/robfig/cron"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
)

type ItemInfo struct {
	Zero1Symbol       string `json:"01. symbol"`
	Zero3High         string `json:"03. high"`
	Zero4Low          string `json:"04. low"`
	Zero5Price        string `json:"05. price"`
	Zero6Volume       string `json:"06. volume"`
	Zero9Change       string `json:"09. change"`
	One0ChangePercent string `json:"10. change percent"`
}

type Res struct {
	GlobalQuote struct {
		Zero1Symbol           string `json:"01. symbol"`
		Zero2Open             string `json:"02. open"`
		Zero3High             string `json:"03. high"`
		Zero4Low              string `json:"04. low"`
		Zero5Price            string `json:"05. price"`
		Zero6Volume           string `json:"06. volume"`
		Zero7LatestTradingDay string `json:"07. latest trading day"`
		Zero8PreviousClose    string `json:"08. previous close"`
		Zero9Change           string `json:"09. change"`
		One0ChangePercent     string `json:"10. change percent"`
	} `json:"Global Quote"`
}

func main() {
	//test

	response, err := http.Get("https://www.alphavantage.co/query?function=GLOBAL_QUOTE&symbol=amzn&apikey="+"API_KEY")
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(responseData))

	var responseObject Res
	json.Unmarshal(responseData, &responseObject)

	fmt.Println(responseObject.GlobalQuote)
	loggly.Trace(responseObject.GlobalQuote)

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

	length := len(responseData)

	fmt.Println(length)

	for i := 0; i < length; i++ {
		res := responseObject

		e := ItemInfo{
			Zero1Symbol:       res.GlobalQuote.Zero1Symbol,
			Zero3High:         res.GlobalQuote.Zero3High,
			Zero4Low:          res.GlobalQuote.Zero4Low,
			Zero5Price:        res.GlobalQuote.Zero5Price,
			Zero6Volume:       res.GlobalQuote.Zero6Volume,
			Zero9Change:       res.GlobalQuote.Zero9Change,
			One0ChangePercent: res.GlobalQuote.One0ChangePercent,
		}

		av, err := dynamodbattribute.MarshalMap(e)

		input := &dynamodb.PutItemInput{
			Item:      av,
			TableName: aws.String("AmazonStock"),
		}

		_, err = svc.PutItem(input)

		if err != nil {
			fmt.Println("Got error calling PutItem:")
			fmt.Println(err.Error())
			os.Exit(1)
		}

		fmt.Println("success")

	}

	c := cron.New()
	c.AddFunc("@every 5m", func() { loggly.Trace(responseObject.GlobalQuote) })
	c.Start()

	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt, os.Kill)
	<-sig

}
