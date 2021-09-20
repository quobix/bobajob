package main

import (
	"github.com/quobix/bobajob"
	"github.com/vmware/transport-go/bridge"
)

type jobHandler struct{}

func (jh jobHandler) HandleJob(job bobajob.JobEnvelope) string {

	// this will simply reverse the string
	l := len(job.Payload)
	var rev = make([]rune, l)
	n := l - 1
	for _, r := range job.Payload {
		rev[n] = r
		n--
	}
	return string(rev)
}

func main() {

	config := &bridge.BrokerConnectorConfig{
		Username:   "guest",
		Password:   "guest",
		UseWS:      false,
		ServerAddr: "localhost:61613"}

	scout := bobajob.NewScout(bobajob.ScoutConfig{
		Name:                  "reverse-string",
		BrokerConnectorConfig: config,
	})

	scout.Connect()
	scout.RegisterHandler(jobHandler{})
	scout.ListenForJob()

}
