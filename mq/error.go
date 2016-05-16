package mq

import "fmt"

func NewJsonMarshalError(err error, v interface{}) error {
	return fmt.Errorf("Unable to marshal %+v to JSON: %s", v, err)
}

func NewQueueSendError(err error, json []byte, queueName string) error {
	return fmt.Errorf("Unable to send `%s` to queue `%s`: %s", string(json), queueName, err)
}

func NewQueueReceiveError(err error, queueName string) error {
	return fmt.Errorf("Unable to receive on queue `%s`: %s", queueName, err)
}
