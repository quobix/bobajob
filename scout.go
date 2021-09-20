package bobajob

import (
	"fmt"
	"github.com/vmware/transport-go/bridge"
	"github.com/vmware/transport-go/bus"
	"github.com/vmware/transport-go/plank/utils"
)

type JobHandler interface {
	HandleJob(job JobEnvelope) string
}

type ScoutConfig struct {
	Name string
	*bridge.BrokerConnectorConfig
}

type Scout struct {
	config  ScoutConfig
	handler JobHandler
	bus.EventBus
	bridgeConnection bridge.Connection
	subscription     bridge.Subscription
	finishedChan     chan bool
}

func NewScout(config ScoutConfig) *Scout {
	if config.Name == "" {
		config.Name = fmt.Sprintf("bobajob-scout-%s", RandomString(5))
	}
	return &Scout{config: config, finishedChan: make(chan bool)}
}

func (s *Scout) RegisterHandler(jobHandler JobHandler) {
	s.handler = jobHandler
}

func (s *Scout) Connect() error {

	s.EventBus = bus.GetBus()

	if s.config.BrokerConnectorConfig == nil {
		return fmt.Errorf("broker config is missing from leader config. unable to connect")
	}
	var err error
	s.bridgeConnection, err = s.EventBus.ConnectBroker(s.config.BrokerConnectorConfig)
	if err != nil {
		return err
	}
	utils.Log.Infof("scout is now connected to broker, sessionID is: %s", s.bridgeConnection.GetId().String())
	return nil
}

func (s *Scout) subscribe() {
	s.subscription, _ = s.bridgeConnection.Subscribe("/queue/" + s.config.Name)
}

func (s *Scout) ListenForJob() {
	s.subscribe()
	for {
		select {
		case <-s.finishedChan:
			utils.Log.Infof("shutting down now.")
			return
		case m := <-s.subscription.GetMsgChannel():

			je, _ := DecodeBase64ToGob(string(m.Payload.([]byte)))
			jr := s.handler.HandleJob(*je)

			utils.Log.Infof("job completed: %v", je.Id)
			var replyTo string
			for _, h := range m.Headers {
				if h.Label == "reply-to" {
					replyTo = h.Value
				}
			}

			je.Payload = jr
			je.HandledBy = fmt.Sprintf("scout-%s-%s", s.config.Name, RandomString(5))
			je.Version++
			je.Hops++

			s.SendCompletedJobBack(replyTo, je)
		}

	}
}

func (s *Scout) SendCompletedJobBack(destination string, je *JobEnvelope) {
	b64, _ := EncodeGobToBase64(*je)
	s.bridgeConnection.SendMessage(destination, "text/plain", []byte(b64))
}

func (s *Scout) Unsubscribe() {
	s.subscription.Unsubscribe()
}
