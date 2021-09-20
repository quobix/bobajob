package main

import (
	"github.com/quobix/bobajob"
	"github.com/vmware/transport-go/bridge"
	"github.com/vmware/transport-go/plank/utils"
)

type myJob struct {
	data string
}

func (job myJob) GetInput() string {
	return job.data
}

func (job myJob) HandleOutput(result string) {
	utils.Log.Infof("the size of my response is %v bytes", len(result))
}

func main() {

	config := &bridge.BrokerConnectorConfig{
		Username:   "guest",
		Password:   "guest",
		UseWS:      false,
		ServerAddr: "localhost:61613"}

	leader := bobajob.NewLeader(bobajob.LeaderConfig{
		BrokerConnectorConfig: config,
	})

	err := leader.Connect()
	if err != nil {
		utils.Log.Errorf("unable to connect. %s", err.Error())
		return
	}

	troop := bobajob.NewTroop("reverse-string")

	for i := 0; i < 100; i++ {
		troop.AddJob(bobajob.CreateJob(myJob{data: bobajob.RandomString(11000000)}))
	}

	leader.AddTroop(troop)

	leader.Run()
	leader.Disconnect()
}
