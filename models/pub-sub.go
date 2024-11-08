package models

import notification "pub-sub-service/sns"

type CreateTopicInput struct {
	TopicName string `json:"topicName"`
}

type SubscribeEmailToTopicInput struct {
	Email string `json:"email"`
}

type SubscribeQueueToTopicInput struct {
	QueueName string `json:"queueName"`
}

type UnsubscribeFromTopicInput struct {
	SubscriptionID string `json:"subscriptionID"`
}

type PublishMessageInput struct {
	Message string `json:"message"`
}

func ListTopics() (*Response, error) {
	res, err := notification.ListTopics()
	if err != nil {
		return &Response{
			Ok: false,
			Response: nil,
		}, err
	}

	return &Response{
		Ok: true,
		Response: res,
	}, nil
}

func ListSubscriptions(topicARN string) (*Response, error) {
	res, err := notification.ListSubscriptions(&topicARN)
	if err != nil {
		return &Response{
			Ok: false,
			Response: nil,
		}, err
	}

	return &Response{
		Ok: true,
		Response: res,
	}, nil
}

func SubscribeEmailToTopic(topicARN string, subscribeEmailToTopicInput SubscribeEmailToTopicInput) (*Response, error) {
	res, err := notification.SubscribeEmailToTopic(&subscribeEmailToTopicInput.Email, &topicARN)
	if err != nil {
		return &Response{
			Ok: false,
			Response: nil,
		}, err
	}

	return &Response{
		Ok: true,
		Response: res,
	}, nil
}

func SubscribeQueueToTopic(topicARN string, subscribeQueueToTopicInput SubscribeQueueToTopicInput) (*Response, error) {
	res, err := notification.SubscribeQueueToTopic(subscribeQueueToTopicInput.QueueName, &topicARN)
	if err != nil {
		return &Response{
			Ok: false,
			Response: nil,
		}, err
	}

	return &Response{
		Ok: true,
		Response: res,
	}, nil
}

func UnsubscribeFromTopic(topicARN string, unsubscribeFromTopicInput UnsubscribeFromTopicInput) (*Response, error) {
	res, err := notification.UnsubscribeFromTopic(&unsubscribeFromTopicInput.SubscriptionID, &topicARN)
	if err != nil {
		return &Response{
			Ok: false,
			Response: nil,
		}, err
	}

	return &Response{
		Ok: true,
		Response: res,
	}, nil
}

func PublishMessageToAllTopicSubscribers(topicARN string, message PublishMessageInput) (*Response, error) {
	res, err := notification.PublishMessageToAllTopicSubscribers(&message.Message, &topicARN)
	if err != nil {
		return &Response{
			Ok: false,
			Response: nil,
		}, err
	}

	return &Response{
		Ok: true,
		Response: res,
	}, nil
}