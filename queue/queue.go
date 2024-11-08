package queue

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

// Queue operations
func ListQueues() ([]string, error) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := sqs.New(sess)

	result, err := svc.ListQueues(nil)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var queueUrls []string
	for i, url := range result.QueueUrls {
		queueUrls = append(queueUrls, *url)
		fmt.Printf("%d: %s\n", i, *url)
	}

	return queueUrls, nil
}

func CreateQueue(queueName string) (bool, error) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := sqs.New(sess)

	result, err := svc.CreateQueue(&sqs.CreateQueueInput{
		QueueName: &queueName,
		Attributes: map[string]*string{
			"DelaySeconds": aws.String("60"),
			"MessageRetentionPeriod": aws.String("86400"),
		},
	})
	if err != nil {
		log.Println(err)
		return false, err
	}

	fmt.Println("URL: " + *result.QueueUrl)
	return true, nil
}

func GetQueueURL(queueName string) (string, error) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := sqs.New(sess)

	result, err := svc.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: &queueName,
	})
	if err != nil {
		log.Println(err)
		return "", err
	}

	fmt.Println("URL: " + *result.QueueUrl)
	return *result.QueueUrl, nil
}

func DeleteQueue(queueName string) (bool, error) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := sqs.New(sess)

	queueUrl, err := svc.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: &queueName,
	})
	if err != nil {
		log.Println(err)
		return false, err
	}

	_, err = svc.DeleteQueue(&sqs.DeleteQueueInput{
		QueueUrl: queueUrl.QueueUrl,
	})
	if err != nil {
		log.Println(err)
		return false, err
	}

	return true, nil
}

func ConfigureVisibilityTimeout(queueName, receiptHandle string, visibilityDuration int) (bool, error) {
	if visibilityDuration < 0 { visibilityDuration = 0 }
	if visibilityDuration > 12 * 60 * 60 { visibilityDuration = 12 * 60 * 60 }

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

	_, err = svc.ChangeMessageVisibility(&sqs.ChangeMessageVisibilityInput{
		ReceiptHandle:     &receiptHandle,
		QueueUrl:          queueUrl,
		VisibilityTimeout: aws.Int64(int64(visibilityDuration)),
	})
	if err != nil {
		log.Println(err)
		return false, err
	}

	return true, nil
}
