package app

import (
	"sync"

	"github.com/streadway/amqp"

	"github.com/trustedanalytics/tap-go-common/queue"
	"github.com/trustedanalytics/tap-go-common/util"
	"github.com/trustedanalytics/tap-image-factory/models"
)

func StartConsumer(waitGroup *sync.WaitGroup) {
	waitGroup.Add(1)

	ch, conn := queue.GetConnectionChannel()
	queue.CreateExchangeWithQueueByRoutingKeys(ch, models.IMAGE_FACTORY_QUEUE_NAME, []string{models.IMAGE_FACTORY_IMAGE_ROUTING_KEY})
	queue.ConsumeMessages(ch, handleMessage, models.IMAGE_FACTORY_QUEUE_NAME)

	defer conn.Close()
	defer ch.Close()

	logger.Info("Consuming stopped. Queue:", models.IMAGE_FACTORY_QUEUE_NAME)
	waitGroup.Done()
}

func handleMessage(msg amqp.Delivery) {
	buildImageRequest := BuildImagePostRequest{}
	err := util.ReadJsonFromByte(msg.Body, &buildImageRequest)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	if err := BuildAndPushImage(buildImageRequest); err != nil {
		logger.Error("Building image error:", err)
	}
}
