package main

import (
	"flag"
	"fmt"
	"os"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
)

var wg sync.WaitGroup
var messages = 1

func main() {
	fileAsString := readFileAsString("./resources/256kb.json")
	message := flag.String("m", fileAsString, "The message payload")
	topic := flag.String("t", "arn:aws:sns:us-east-1:"+os.Getenv("AWS_ACCOUNT_ID")+":document-saved", "The ARN of the topic to which the user subscribes")
	flag.Parse()
	wg.Add(messages)
	for i := 0; i < messages; i++ {
		go sendMessageToTopic(message, topic)
	}
	wg.Wait()
}

func readFileAsString(fileName string) string {
	file, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("Successfully opened file")
	fileAsString := string(file)
	fmt.Println("Successfully converted file to string")
	return fileAsString
}

func sendMessageToTopic(message *string, topic *string) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Profile: "default",
		Config: aws.Config{
			Region: aws.String("us-east-1"),
		},
	}))
	svc := sns.New(sess)
	result, err := svc.Publish(&sns.PublishInput{
		Message:  message,
		TopicArn: topic,
	})
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	fmt.Println(*result.MessageId)
	wg.Done()
}
