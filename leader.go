package bobajob

import (
	"fmt"
	"github.com/vmware/transport-go/bridge"
	"github.com/vmware/transport-go/bus"
	"github.com/vmware/transport-go/plank/utils"
)

type LeaderConfig struct {
	name string
	*bridge.BrokerConnectorConfig
}

type LeaderReport struct {
	Success bool
	Errors []string
}

type Leader struct {
	troops []*Troop
	config LeaderConfig
	bus.EventBus
	bridgeConnection bridge.Connection
}

func NewLeader(config LeaderConfig) *Leader {
	if config.name == "" {
		config.name = fmt.Sprintf("bobajob-leader-%s", RandomString(5))
	}
	return &Leader{ config: config }
}

func (l *Leader) Connect() error {

	l.EventBus = bus.GetBus()

	if l.config.BrokerConnectorConfig == nil {
		return fmt.Errorf("broker config is missing from leader config. unable to connect")
	}
	var err error
	l.bridgeConnection, err = l.EventBus.ConnectBroker(l.config.BrokerConnectorConfig)
	if err != nil {
		return err
	}
	utils.Log.Infof("leader is connected to broker, sessionID is: %s", l.bridgeConnection.GetId().String())
	l.subscribeTroops()
	return nil
}

func (l *Leader) subscribeTroops() {
	for _, t := range l.troops {
		t.Subscribe(l.bridgeConnection)
	}
}


func (l *Leader) AddTroop(t *Troop) {
	l.troops = append(l.troops, t)
}

func (l *Leader) AddTroops(t []*Troop) {
	l.troops = append(l.troops, t...)
}

func (l *Leader) Run() (chan *LeaderReport, error) {
	if len(l.troops) <= 0 {
		return nil, fmt.Errorf("cannot run, no troops defined for leader %s", l.config.name)
	}
	for _, t := range l.troops {
		t.Run()
	}




	lrChan := make(chan *LeaderReport, 1)
	lrChan <- &LeaderReport{ Success: false}
	return lrChan, nil
}