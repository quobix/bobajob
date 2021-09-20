package bobajob

import (
	"fmt"
	"github.com/vmware/transport-go/bridge"
	"github.com/vmware/transport-go/plank/utils"
)

type Troop struct {
	Jobs []*Job
	Name string
	CompletedChan chan bool
	jobChannel chan JobEnvelope
	bridgeConnection bridge.Connection
	replySubscription bridge.Subscription
	troopSubscription bridge.Subscription
}

func NewTroop(name string) *Troop {
	return &Troop{
		Name: name,
		CompletedChan: make(chan bool, 1),
		jobChannel: make(chan JobEnvelope),
	}
}

func (t *Troop) AddJob(job *Job) {
	job.troop = t
	t.Jobs = append(t.Jobs, job)
}

func (t *Troop) AddJobs(jobs []*Job) {
	for _, j := range jobs {
		j.troop = t
	}
	t.Jobs = append(t.Jobs, jobs...)
}

func (t *Troop) Subscribe(bc bridge.Connection) {
	t.bridgeConnection = bc
	utils.Log.Infof("troop subscribing to reply queue %s", t.Name)
	t.replySubscription, _ = t.bridgeConnection.SubscribeReplyDestination("/temp-queue/" + t.Name)
	//t.troopSubscription, _ = t.bridgeConnection.Subscribe("/queue/" + t.Name)
}

func (t *Troop) Run() {
	if len(t.Jobs) <= 0 {
		close(t.CompletedChan)
	}
	for _, j := range t.Jobs {
		go j.Run(t.jobChannel)
	}
	var count = 0
	for {
		select {
		case <-t.jobChannel:
			count++
			if count == len(t.Jobs) {
				fmt.Printf("we have %v jobs completed\n", count)
				t.CompletedChan <- true
				break
			}
		}
	}
}

func (t *Troop) JobsCompleted() (int,int) {
	completed := 0
	for _, i := range t.Jobs {
		if i.completed {
			completed++
		}
	}
	return completed, len(t.Jobs)
}

