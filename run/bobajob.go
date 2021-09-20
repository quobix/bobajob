package main

import (
	"github.com/quobix/bobajob"
	"github.com/vmware/transport-go/bridge"
	"github.com/vmware/transport-go/plank/utils"
)

type jobHandler struct {

}

func (jh jobHandler) HandleJob(job bobajob.JobEnvelope) bobajob.JobResult {
	return bobajob.JobResult{ Payload: "nice..."}
}


func main() {

	config := &bridge.BrokerConnectorConfig{
		Username:   "guest",
		Password:   "guest",
		UseWS:      false,
		ServerAddr: "localhost:61613"}


	scout := bobajob.NewScout(bobajob.ScoutConfig{
		Name:                  "star-wars",
		BrokerConnectorConfig: config,
	})





	scout.Connect()
	scout.RegisterHandler(jobHandler{})
	go scout.ListenForJob()





	leader := bobajob.NewLeader(bobajob.LeaderConfig{
		BrokerConnectorConfig: config,
	})









	troop := &bobajob.Troop{Name: "star-wars"}
	troop.AddJob(&bobajob.Job{})

	leader.AddTroop(troop)

	err := leader.Connect()
	if err != nil {
		utils.Log.Errorf("unable to connect. %s", err.Error())
		return
	}

	utils.Log.Infof("ready to run.")

	leader.Run()

	forever := make(chan bool, 1)
	<- forever

	//rId := rand.Intn(10000)

	//forever := make(chan bool, 1)
	//
	//
	//s, _ := c.SubscribeReplyDestination("/temp-queue/reply-all")
	//
	//go func() {
	//	// listen for incoming messages from subscription.
	//	for {
	//		m := <-s.GetMsgChannel()
	//		utils.Log.Errorf("render completed: %v", string(m.Payload.([]byte)))
	//	}
	//}()



	//go func() {
	//
	//
	//
	//	utils.Log.Info("Broadcasting rendering requests....")
	//
	//
	//	for i := 1; i < 10000; i++ {
	//		msg := fmt.Sprintf("hey, myID is %v and this message is %v", rId, i)
	//		c.SendMessageWithReplyDestination("/queue/render-request", "/temp-queue/reply-all","text/plain", []byte(msg),nil)
	//
	//	}
	//
	//
	//}()
	//
	//<- forever
	//c.Disconnect()
	
}
