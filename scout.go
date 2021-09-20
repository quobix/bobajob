package bobajob

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"fmt"
	"github.com/vmware/transport-go/bridge"
	"github.com/vmware/transport-go/bus"
	"github.com/vmware/transport-go/plank/utils"
)

type JobHandler interface {
	HandleJob(job JobEnvelope) JobResult
}

type ScoutConfig struct {
	Name string
	*bridge.BrokerConnectorConfig
}


type Scout struct {
	config ScoutConfig
	handler JobHandler
	bus.EventBus
	bridgeConnection bridge.Connection
	subscription bridge.Subscription
	finishedChan chan bool
}


func NewScout(config ScoutConfig) *Scout {
	if config.Name == "" {
		config.Name = fmt.Sprintf("bobajob-scout-%s", RandomString(5))
	}
	return &Scout{ config: config, finishedChan: make(chan bool)}
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
	s.Subscribe()
	return nil
}

func (s *Scout) Subscribe() {
	s.subscription, _ = s.bridgeConnection.Subscribe("/queue/" + s.config.Name)
}

func (s *Scout) ListenForJob() {
	for {
		select {
			case <-s.finishedChan:
				utils.Log.Infof("shutting down now.")
				return
			case m := <- s.subscription.GetMsgChannel():

				var b64 = m.Payload.([]byte)
				gobBits, err := base64.StdEncoding.DecodeString(string(b64))
				if err != nil {
					utils.Log.Fatal(err.Error())
					return
				}

				var buf bytes.Buffer
				buf.Write(gobBits)

				enc := gob.NewDecoder(&buf)

				var je JobEnvelope
				err = enc.Decode(&je)
				if err != nil {
					utils.Log.Fatal("failed")
				}
				jr := s.handler.HandleJob(je)

				utils.Log.Infof("render completed.%s.%v", jr.Payload, je.Id)
				var replyTo string
				for _, h := range m.Headers {
					if h.Label == "reply-to" {
						replyTo = h.Value
					}
				}
				s.SendCompletedJobBack(replyTo)
		}

	}
}

func (s *Scout) SendCompletedJobBack(destination string) {
	s.bridgeConnection.SendMessage(destination, "text/plain", []byte("rice"))
}

func (s *Scout) Unsubscribe() {
	s.subscription.Unsubscribe()
}

