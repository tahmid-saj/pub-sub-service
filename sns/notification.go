package notification

import (
	"errors"
	"log"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sqs"
)

func ListTopics() ([]*sns.Topic, error) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := sns.New(sess)

	result, err := svc.ListTopics(nil)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	var topics []*sns.Topic
	topics = append(topics, result.Topics...)

	return topics, nil
}

func CreateTopic(topicName string) (*sns.CreateTopicOutput, error) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := sns.New(sess)

	result, err := svc.CreateTopic(&sns.CreateTopicInput{
		Name: aws.String(topicName),
	})

	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return result, nil
}

func ListSubscriptions(topicPtr *string) ([]*sns.Subscription, error) {
	if *topicPtr == "" {
		fmt.Println("You must supply a topic ARN")
		return nil, errors.New("must supply email and topic")
	}

	// Initialize a session that the SDK will use to load
	// credentials from the shared credentials file. (~/.aws/credentials).
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := sns.New(sess)
	var previousToken *string

	result, err := svc.ListSubscriptionsByTopic(&sns.ListSubscriptionsByTopicInput{
		NextToken: previousToken,
		TopicArn: topicPtr,
	})
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return result.Subscriptions, nil
}

func SubscribeEmailToTopic(emailPtr *string, topicPtr *string) (*sns.SubscribeOutput, error) {
	if *emailPtr == "" || *topicPtr == "" {
		fmt.Println("You must supply an email address and topic ARN")
		return nil, errors.New("must supply email and topic")
	}

	// Initialize a session that the SDK will use to load
	// credentials from the shared credentials file. (~/.aws/credentials).
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := sns.New(sess)

	result, err := svc.Subscribe(&sns.SubscribeInput{
		Endpoint:              emailPtr,
		Protocol:              aws.String("email"),
		ReturnSubscriptionArn: aws.Bool(true), // Return the ARN, even if user has yet to confirm
		TopicArn:              topicPtr,
	})

	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return result, nil
}

func SubscribeQueueToTopic(queueName string, topicPtr *string) (bool, error) {
  if queueName == "" || topicPtr == nil || *topicPtr == "" {
    return false, errors.New("must supply both queue name and topic ARN")
  }

  // Initialize a session to load AWS credentials and configuration from the shared config.
  sess := session.Must(session.NewSessionWithOptions(session.Options{
    SharedConfigState: session.SharedConfigEnable,
  }))

  // Create SNS and SQS clients
  snsSvc := sns.New(sess)
  sqsSvc := sqs.New(sess)

  // Use the provided topic ARN
  topicArn := *topicPtr

  // Get the SQS queue URL and ARN
  queueUrlOutput, err := sqsSvc.GetQueueUrl(&sqs.GetQueueUrlInput{
    QueueName: aws.String(queueName),
  })
  if err != nil {
    return false, fmt.Errorf("unable to get SQS queue URL: %v", err)
  }

  queueAttrsOutput, err := sqsSvc.GetQueueAttributes(&sqs.GetQueueAttributesInput{
    QueueUrl:       queueUrlOutput.QueueUrl,
    AttributeNames: []*string{aws.String("QueueArn")},
  })
  if err != nil {
    return false, fmt.Errorf("unable to get SQS queue attributes: %v", err)
  }

  queueArn := queueAttrsOutput.Attributes["QueueArn"]

  // Subscribe the SQS queue to the SNS topic
  _, err = snsSvc.Subscribe(&sns.SubscribeInput{
    Protocol: aws.String("sqs"),
    TopicArn: aws.String(topicArn),
    Endpoint: aws.String(*queueArn),
  })
  if err != nil {
    return false, fmt.Errorf("unable to subscribe SQS queue to SNS topic: %v", err)
  }

  // Set policy to allow SNS to send messages to the SQS queue
  policy := fmt.Sprintf(`{
    "Version": "2012-10-17",
    "Statement": [
      {
        "Effect": "Allow",
        "Principal": "*",
        "Action": "SQS:SendMessage",
        "Resource": "%s",
        "Condition": {
          "ArnEquals": {
            "aws:SourceArn": "%s"
          }
        }
      }
    ]
  }`, *queueArn, topicArn)

  _, err = sqsSvc.SetQueueAttributes(&sqs.SetQueueAttributesInput{
    QueueUrl: queueUrlOutput.QueueUrl,
    Attributes: map[string]*string{
      "Policy": aws.String(policy),
    },
  })
  if err != nil {
    return false, fmt.Errorf("unable to set SQS queue policy: %v", err)
  }

  log.Printf("Successfully subscribed SQS queue %s to SNS topic %s", queueName, *topicPtr)
  return true, nil
}


func PublishMessageToAllTopicSubscribers(messagePtr *string, topicPtr *string) (*sns.PublishOutput, error) {
  if messagePtr == nil || topicPtr == nil || *messagePtr == "" || *topicPtr == "" {
    return nil, errors.New("must supply both a message and topic ARN")
  }

  // Initialize a session to load AWS credentials and configuration from the shared config.
  sess := session.Must(session.NewSessionWithOptions(session.Options{
    SharedConfigState: session.SharedConfigEnable,
  }))

  // Create SNS client
  svc := sns.New(sess)

  // Publish the message to the SNS topic
  result, err := svc.Publish(&sns.PublishInput{
    Message:  messagePtr,
    TopicArn: topicPtr,
  })

  if err != nil {
    return nil, fmt.Errorf("failed to publish message to topic %s: %v", *topicPtr, err)
  }

  log.Printf("Successfully published message to topic %s", *topicPtr)
  return result, nil
}
