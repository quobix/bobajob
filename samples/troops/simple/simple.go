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
	utils.Log.Infof("my job result for %s is: %s", job.data, result)
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

	troop := bobajob.NewTroop("reverse-string")

	job1 := bobajob.CreateJob(myJob{data: "kittens"})
	job2 := bobajob.CreateJob(myJob{data: "lemons"})
	job3 := bobajob.CreateJob(myJob{data: "racecar"})
	job4 := bobajob.CreateJob(myJob{data: "potatoes"})
	job5 := bobajob.CreateJob(myJob{data: "yummymummy"})
	job6 := bobajob.CreateJob(myJob{data: "waterpark"})

	troop.AddJobs([]*bobajob.Job{job1, job2, job3, job4, job5, job6})

	leader.AddTroop(troop)

	err := leader.Connect()
	if err != nil {
		utils.Log.Errorf("unable to connect. %s", err.Error())
		return
	}

	leader.Run()
	leader.Disconnect()
}
