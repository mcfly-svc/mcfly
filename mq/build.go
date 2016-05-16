package mq

import "github.com/streadway/amqp"

type BuildMessage struct {
	BuildHandle    string `json:"build_handle"`
	SourceProvider string `json:"source_provider"`
}

func (rc *MsplChannel) StartDeploy(buildHandle, sourceProvider string) error {
	buildMsg := &BuildMessage{
		BuildHandle:    buildHandle,
		SourceProvider: sourceProvider,
	}
	return rc.Send(rc.BuildQueue, buildMsg)

	// NEXT: add mq channel to router/handlers
	// call StartBuild after webhook for project update is received, for each build.
	// create another go program (griswold) to process build jobs

}

func (rc *MsplChannel) ReceiveBuildMessage() (<-chan amqp.Delivery, error) {
	return rc.Receive(rc.BuildQueue)
}
