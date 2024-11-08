package routes

import "github.com/gin-gonic/gin"

func RegisterRoutes(server *gin.Engine) {
	// ListTopics
	server.GET("/topics", listTopics)

	// CreateTopic
	server.POST("/topics", createTopic)

	// ListSubscriptions
	server.GET("/topics/:topicARN/subscriptions", listSubscriptions)

	// SubscribeEmailToTopic
	server.PUT("/topics/:topicARN/subscribe/email", subscribeEmailToTopic)

	// SubscribeQueueToTopic
	server.PUT("/topics/:topicARN/subscribe/queue", subscribeQueueToTopic)

	// UnsubscribeFromTopic
	server.PUT("/topics/:topicARN/unsubscribe", unsubscribeFromTopic)

	// PublishMessageToAllTopicSubscribers
	server.POST("/topics/:topicARN", publishMessageToAllTopicSubscribers)
}