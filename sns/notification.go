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

func SubscribeQueueToTopic(queueName, topicName string) error {
  if queueName == "" || topicName == "" {
    return errors.New("must supply both queue name and topic name")
  }

  // Initialize a session to load AWS credentials and configuration from the shared config.
  sess := session.Must(session.NewSessionWithOptions(session.Options{
    SharedConfigState: session.SharedConfigEnable,
  }))

  // Create SNS and SQS clients
  snsSvc := sns.New(sess)
  sqsSvc := sqs.New(sess)

  // Get the SNS topic ARN by listing topics
  topicsOutput, err := snsSvc.ListTopics(&sns.ListTopicsInput{})
  if err != nil {
    return fmt.Errorf("unable to list SNS topics: %v", err)
  }

  var topicArn string
  for _, topic := range topicsOutput.Topics {
    // Get topic attributes to check the DisplayName
    attrsOutput, err := snsSvc.GetTopicAttributes(&sns.GetTopicAttributesInput{
      TopicArn: topic.TopicArn,
    })
    if err != nil {
      return fmt.Errorf("unable to get SNS topic attributes: %v", err)
    }

    if *attrsOutput.Attributes["DisplayName"] == topicName {
      topicArn = *topic.TopicArn
      break
    }
  }

  if topicArn == "" {
    return fmt.Errorf("SNS topic %s not found", topicName)
  }

  // Get the SQS queue URL and ARN
  queueUrlOutput, err := sqsSvc.GetQueueUrl(&sqs.GetQueueUrlInput{
    QueueName: aws.String(queueName),
  })
  if err != nil {
    return fmt.Errorf("unable to get SQS queue URL: %v", err)
  }

  queueAttrsOutput, err := sqsSvc.GetQueueAttributes(&sqs.GetQueueAttributesInput{
    QueueUrl:       queueUrlOutput.QueueUrl,
    AttributeNames: []*string{aws.String("QueueArn")},
  })
  if err != nil {
    return fmt.Errorf("unable to get SQS queue attributes: %v", err)
  }

  queueArn := queueAttrsOutput.Attributes["QueueArn"]

  // Subscribe the SQS queue to the SNS topic
  _, err = snsSvc.Subscribe(&sns.SubscribeInput{
    Protocol: aws.String("sqs"),
    TopicArn: aws.String(topicArn),
    Endpoint: aws.String(*queueArn),
  })
  if err != nil {
    return fmt.Errorf("unable to subscribe SQS queue to SNS topic: %v", err)
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
  }`, queueArn, topicArn)

  _, err = sqsSvc.SetQueueAttributes(&sqs.SetQueueAttributesInput{
    QueueUrl: queueUrlOutput.QueueUrl,
    Attributes: map[string]*string{
      "Policy": aws.String(policy),
    },
  })
  if err != nil {
    return fmt.Errorf("unable to set SQS queue policy: %v", err)
  }

  log.Printf("Successfully subscribed SQS queue %s to SNS topic %s", queueName, topicName)
  return nil
}

func PublishMessageToAllTopicSubscribers(messagePtr *string, topicPtr *string) (*sns.PublishOutput, error) {

	if *messagePtr == "" || *topicPtr == "" {
		fmt.Println("You must supply a message and topic ARN")
		return nil, errors.New("must supply message and topic")
	}

	// Initialize a session that the SDK will use to load
	// credentials from the shared credentials file. (~/.aws/credentials).
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := sns.New(sess)

	result, err := svc.Publish(&sns.PublishInput{
		Message:  messagePtr,
		TopicArn: topicPtr,
	})

	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return result, nil
}