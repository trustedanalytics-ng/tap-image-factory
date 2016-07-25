package app

import (
	"fmt"
	"log"

	"bytes"
	"github.com/streadway/amqp"
	"github.com/trustedanalytics/tapng-go-common/util"
)

func failReceiverOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

func StartConsumer(ctx Context) {
	logger.Error(GetQueueConnectionString())
	conn, err := amqp.Dial(GetQueueConnectionString())
	failReceiverOnError(err, "Failed to connect to Queue on address: "+GetQueueConnectionString())
	defer conn.Close()

	ch, err := conn.Channel()
	failReceiverOnError(err, "Failed to open a channel")
	defer ch.Close()

	msgs, err := ch.Consume(
		GetQueueName(), // queue
		"",             // consumer - empty means generate unique id
		true,           // auto-ack
		true,           // exclusive
		false,          // no-local
		false,          // no-wait
		nil,            // args
	)
	failReceiverOnError(err, "Failed to register a consumer")

	go func() {
		for m := range msgs {
			handleMessage(ctx, m)
		}
	}()

	forever := make(chan bool)
	<-forever
}

func handleMessage(c Context, msg amqp.Delivery) {
	msg_json := BuildImagePostRequest{}
	err := util.ReadJsonFromByte(msg.Body, &msg_json)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	imgDetails, err := c.TapCatalogApiConnector.GetImage(msg_json.ImageId)
	if err != nil {
		c.updateImageWithState(msg_json.ImageId, "ERROR")
		logger.Error(err.Error())
		return
	}

	blobBytes, err := c.BlobStoreConnector.GetImageBlob(imgDetails.Id)
	if err != nil {
		c.updateImageWithState(msg_json.ImageId, "ERROR")
		logger.Error(err.Error())
		return
	}
	c.updateImageWithState(msg_json.ImageId, "BUILDING")
	if err != nil {
		c.updateImageWithState(msg_json.ImageId, "ERROR")
		logger.Error(err.Error())
		return
	}

	tag := GetHubAddressWithoutProtocol() + "/" + imgDetails.Id

	err = c.DockerConnector.CreateImage(bytes.NewReader(blobBytes), imgDetails.Type, tag)
	if err != nil {
		c.updateImageWithState(msg_json.ImageId, "ERROR")
		logger.Error(err.Error())
		return
	}
	err = c.DockerConnector.PushImage(tag)
	if err != nil {
		c.updateImageWithState(msg_json.ImageId, "ERROR")
		logger.Error(err.Error())
		return
	}
	c.updateImageWithState(msg_json.ImageId, "READY")
	err = c.BlobStoreConnector.DeleteImageBlob(imgDetails.Id)
	if err != nil {
		logger.Error(err.Error())
		return
	}
}
