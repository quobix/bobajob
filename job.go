package bobajob

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"github.com/google/uuid"
	"github.com/vmware/transport-go/plank/utils"
)


type RunnableJob interface {
	Run()
	//HasCompleted()
}
type JobEnvelope struct {
	Id string
	Payload string
}

type JobResult struct {
	Payload string
}

type Job struct {
	uuid.UUID
	SequenceId int
	RequestPayload []byte
	ResponsePayload []byte
	CompletedBy string
	completed bool
	troop *Troop
}

func (job *Job) Run(completedChannel chan JobEnvelope) {
	done := make(chan bool, 1)


	go func() {
		// listen for incoming messages from subscription.
		for {
			utils.Log.Infof("listening for responses")
			m := <-job.troop.replySubscription.GetMsgChannel()
			utils.Log.Infof("job completed: %v", string(m.Payload.([]byte)))
			done <- true
			return
		}
	}()

	/*
	var network bytes.Buffer        // Stand-in for a network connection
	    enc := gob.NewEncoder(&network) // Will write to network.
	    dec := gob.NewDecoder(&network) // Will read from network.
	    // Encode (send) the value.
	    err := enc.Encode(P{3, 4, 5, "Pythagoras"})
	    if err != nil {
	        log.Fatal("encode error:", err)
	    }
	    // Decode (receive) the value.
	 */

	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	env := JobEnvelope{
		Id:  "123456",
		Payload :"Magic pizza",
	}
	err := enc.Encode(env)
	if err != nil {
		utils.Log.Panicf("failed")
	}

	var b64 = base64.StdEncoding.EncodeToString(buf.Bytes())



	job.troop.bridgeConnection.SendMessageWithReplyDestination("/queue/" + job.troop.Name,
		"/temp-queue/" + job.troop.Name,"text/plain", []byte(b64),nil)

	<- done

	utils.Log.Infof("job done")
	job.completed = true
	completedChannel <- env
}

//func (job *Job) HasCompleted() {
//	return job.completed
//}