package routes

import "github.com/gin-gonic/gin"

func RegisterRoutes(server *gin.Engine) {
	// ListTopics
	server.GET("/topics")

	// CreateTopic
	server.POST("/topics")

	// ListSubscriptions
	server.POST("/topics/:topicARN/subscriptions")

	// SubscribeEmailToTopic
	server.PUT("/topics/:topicARN/subscribe/email")

	// SubscribeQueueToTopic
	server.PUT("/topics/:topicARN/subscribe/queue")

	// UnsubscribeFromTopic
	server.PUT("/topics/:topicARN/unsubscribe")

	// PublishMessageToAllTopicSubscribers
	server.POST("/topics/:topicARN")
}