package queue

import (
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type Message struct {
	Subject string
	Body string
	Timestamp time.Time
}

// Message operations
func SendMessage(queueName string, message Message) (bool, error) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := sqs.New(sess)

	result, err := svc.GetQueueUrl(&sqs.GetQueueUrlInput{
    QueueName: &queueName,
	})
	if err != nil {
		log.Println(err)
		return false, err
	}

	queueUrl := result.QueueUrl

	_, err = svc.SendMessage(&sqs.SendMessageInput{
		DelaySeconds: aws.Int64(10),
		MessageAttributes: map[string]*sqs.MessageAttributeValue{
			"Subject": &sqs.MessageAttributeValue{
				DataType: aws.String("String"),
				StringValue: aws.String(message.Subject),
			},
			"Timestamp": &sqs.MessageAttributeValue{
				DataType: aws.String("String"),
				StringValue: aws.String(message.Timestamp.String()),
			},
		},
		MessageBody: aws.String(message.Body),
		QueueUrl: queueUrl,
	})
	if err != nil {
		log.Println(err)
		return false, err
	}

	return true, nil
}

func ReceiveMessage(queueName string, visibilityTimeout int) (*sqs.Message, error) {
	if visibilityTimeout < 0 { visibilityTimeout = 0 }
	if visibilityTimeout > 12 * 60 * 60 { visibilityTimeout = 12 * 60 * 60 }

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := sqs.New(sess)

	result, err := svc.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: &queueName,
	})
	if err != nil {
		log.Println(err)
		return nil, err
	}

	queueUrl := result.QueueUrl

	messageResult, err := svc.ReceiveMessage(&sqs.ReceiveMessageInput{
		AttributeNames: []*string{
			aws.String(sqs.MessageSystemAttributeNameSentTimestamp),
		},
		MessageAttributeNames: []*string{
			aws.String(sqs.QueueAttributeNameAll),
		},
		QueueUrl: queueUrl,
		MaxNumberOfMessages: aws.Int64(1),
		VisibilityTimeout: aws.Int64(int64(visibilityTimeout)),
	})
	if err != nil {
		log.Println(err)
		return nil, err
	}

	message := *messageResult.Messages[0]
	fmt.Println("Message: " + *message.ReceiptHandle)
	
	return &message, nil
}

func DeleteMessage(queueName, receiptHandle string) (bool, error) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := sqs.New(sess)

	result, err := svc.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: &queueName,
	})
	if err != nil {
		log.Println(err)
		return false, err
	}

	queueUrl := result.QueueUrl

	_, err = svc.DeleteMessage(&sqs.DeleteMessageInput{
		QueueUrl: queueUrl,
		ReceiptHandle: &receiptHandle,
	})
	if err != nil {
		log.Println(err)
		return false, err
	}

	return true, nil
}