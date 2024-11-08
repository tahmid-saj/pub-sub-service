package main

import (
	"os"
	"pub-sub-service/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// CreateQueue
	// res, err := queue.CreateQueue("test-queue")
	// if err != nil {
	// 	return
	// }
	// fmt.Println(res)

	// ListQueues
	// res, err := queue.ListQueues()
	// if err != nil {
	// 	return
	// }
	// fmt.Println(res)

	// GetQueueURL
	// res, err := queue.GetQueueURL("test-queue")
	// if err != nil {
	// 	return
	// }
	// fmt.Println(res)

	// SendMessage
	// message := queue.Message{
	// 	Subject: "test-subject",
	// 	Body: "test-body",
	// 	Timestamp: time.Now(),
	// }
	// res, err := queue.SendMessage("test-queue", message)
	// if err != nil {
	// 	return
	// }
	// fmt.Println(res)

	// ReceiveMessage
	// res, err := queue.ReceiveMessage("test-queue", 10)
	// if err != nil {
	// 	return
	// }
	// fmt.Println(res)

	// DeleteMessage
	// res, err := queue.DeleteMessage("test-queue", "")
	// if err != nil {
	// 	return
	// }
	// fmt.Println(res)

	// ListTopics
	// topics, err := notification.ListTopics()
	// if err != nil {
	// 	return
	// }
	// fmt.Println(topics)

	// ListSubscriptions
	// topicARN := ""
	// subscriptions, err := notification.ListSubscriptions(&topicARN)
	// if err != nil {
	// 	return
	// }
	// fmt.Println(subscriptions)

	// SubscribeEmailToTopic
	// email := ""
	// topicARN := ""
	// subscription, err := notification.SubscribeEmailToTopic(&email, &topicARN)
	// if err != nil {
	// 	return
	// }
	// fmt.Println(subscription)

	// SubscribeQueueToTopic
	// topicARN := ""
	// res, err := notification.SubscribeQueueToTopic("test-queue", &topicARN)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// fmt.Println(res)

	// PublishMessageToAllTopicSubscribers
	// message := "test message"
	// topicARN := ""
	// publishOutput, err := notification.PublishMessageToAllTopicSubscribers(&message, &topicARN)
	// if err != nil {
	// 	return
	// }
	// fmt.Println(publishOutput)

	godotenv.Load()

	server := gin.Default()

	routes.RegisterRoutes(server)

	server.Run(os.Getenv("PORT"))
}