package routes

import (
	"net/http"
	"pub-sub-service/models"

	"github.com/gin-gonic/gin"
)

func listTopics(context *gin.Context) {
	res, err := models.ListTopics()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not list topics"})
		return
	}

	context.JSON(http.StatusOK, res)
}

func createTopic(context *gin.Context) {
	var createTopicInput models.CreateTopicInput

	err := context.ShouldBindJSON(&createTopicInput)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "could not parse request body"})
		return
	}

	res, err := models.CreateTopic(createTopicInput)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not create topic"})
		return
	}

	context.JSON(http.StatusOK, res)
}

func listSubscriptions(context *gin.Context) {
	topicARN := context.Param("topicARN")

	res, err := models.ListSubscriptions(topicARN)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not list subscriptions to topic"})
		return
	}

	context.JSON(http.StatusOK, res)
}

func subscribeEmailToTopic(context *gin.Context) {
	topicARN := context.Param("topicARN")

	var subscribeEmailToTopicInput models.SubscribeEmailToTopicInput

	err := context.ShouldBindJSON(&subscribeEmailToTopicInput)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "could not parse request body"})
		return
	}

	res, err := models.SubscribeEmailToTopic(topicARN, subscribeEmailToTopicInput)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not subscribe email to topic"})
		return
	}

	context.JSON(http.StatusOK, res)
}

func subscribeQueueToTopic(context *gin.Context) {
	topicARN := context.Param("topicARN")

	var subscribeQueueToTopicInput models.SubscribeQueueToTopicInput

	err := context.ShouldBindJSON(&subscribeQueueToTopicInput)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "could not parse request body"})
		return
	}

	res, err := models.SubscribeQueueToTopic(topicARN, subscribeQueueToTopicInput)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not subscribe queue to topic"})
		return
	}

	context.JSON(http.StatusOK, res)
}

func unsubscribeFromTopic(context *gin.Context) {
	topicARN := context.Param("topicARN")

	var unsubscribeFromTopicInput models.UnsubscribeFromTopicInput

	err := context.ShouldBindJSON(&unsubscribeFromTopicInput)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "could not parse request body"})
		return
	}

	res, err := models.UnsubscribeFromTopic(topicARN, unsubscribeFromTopicInput)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not unsubscribe subscription ID from topic"})
		return
	}

	context.JSON(http.StatusOK, res)
}

func publishMessageToAllTopicSubscribers(context *gin.Context) {
	topicARN := context.Param("topicARN")

	var publishMessageInput models.PublishMessageInput

	err := context.ShouldBindJSON(&publishMessageInput)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "could not parse request body"})
		return
	}

	res, err := models.PublishMessageToAllTopicSubscribers(topicARN, publishMessageInput)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not publish message"})
		return
	}

	context.JSON(http.StatusOK, res)
}