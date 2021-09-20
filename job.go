package bobajob

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/vmware/transport-go/bridge"
	"sync"
)

type JobEnvelope struct {
	Id        string
	Payload   string
	SentBy    string
	HandledBy string
	Version   int
	Hops      int
}

type Bob interface {
	GetInput() string
	HandleOutput(result string)
}

type Job struct {
	Id         string
	Envelope   JobEnvelope
	sequenceId int
	completed  bool
	troop      *Troop
	bob        Bob
	sub        bridge.Subscription
}

func CreateJob(bob Bob) *Job {
	id := uuid.New()
	return &Job{
		bob: bob,
		Id:  id.String(),
	}
}

func (job *Job) Run(completedChannel chan JobEnvelope) {
	var wg sync.WaitGroup
	wg.Add(1)

	jobReplyDest := fmt.Sprintf("/temp-queue/%s-%s", job.troop.Name, RandomString(6))
	job.sub, _ = job.troop.bridgeConnection.SubscribeReplyDestination(jobReplyDest)

	job.Envelope = JobEnvelope{
		Id:      job.Id,
		Payload: job.bob.GetInput(),
		SentBy:  fmt.Sprintf("job-%s", job.Id),
		Version: 1,
		Hops:    0,
	}
	b64, _ := EncodeGobToBase64(job.Envelope)

	go func() {
		for {
			m := <-job.sub.GetMsgChannel()
			val, _ := DecodeBase64ToGob(string(m.Payload.([]byte)))
			job.Envelope = *val
			job.bob.HandleOutput(job.Envelope.Payload)
			wg.Done()
			return
		}
	}()

	job.troop.bridgeConnection.SendMessageWithReplyDestination("/queue/"+job.troop.Name,
		jobReplyDest, "text/plain", []byte(b64), nil)

	wg.Wait()
	job.completed = true
	completedChannel <- job.Envelope
}
