package mq

import "github.com/streadway/amqp"

type DeployQueueMessage struct {
	BuildHandle               string `json:"build_handle"`
	ProjectHandle             string `json:"project_handle"`
	SourceProvider            string `json:"source_provider"`
	SourceProviderAccessToken string `json:"source_provider_access_token"`
}

func (rc *McflyChannel) SendDeployQueueMessage(deployMsg *DeployQueueMessage) error {
	return rc.Send(rc.DeployQueue, deployMsg)
}

func (rc *McflyChannel) ReceiveDeployQueueMessage() (<-chan amqp.Delivery, error) {
	return rc.Receive(rc.DeployQueue)
}
