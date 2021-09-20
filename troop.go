package bobajob

import (
	"github.com/vmware/transport-go/bridge"
	"sync"
)

type Troop struct {
	Jobs             []*Job
	Name             string
	CompletedChan    chan bool
	bridgeConnection bridge.Connection
	jobResults       map[string]JobEnvelope
}

func NewTroop(name string) *Troop {
	return &Troop{
		Name:       name,
		jobResults: make(map[string]JobEnvelope),
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

func (t *Troop) Run(wg *sync.WaitGroup) {

	jc := len(t.Jobs)
	if jc <= 0 {
		close(t.CompletedChan)
	}

	jChan := make(chan JobEnvelope)

	for _, j := range t.Jobs {
		go j.Run(jChan)
	}

	var jwg sync.WaitGroup
	jwg.Add(len(t.Jobs))

	go func() {
		for {
			je := <-jChan
			t.jobResults[je.Id] = je
			jwg.Done()
		}
	}()

	jwg.Wait()
	wg.Done()
}

func (t *Troop) GetJobResults() map[string]JobEnvelope {
	return t.jobResults
}

func (t *Troop) JobsCompleted() (int, int) {
	completed := 0
	for _, i := range t.Jobs {
		if i.completed {
			completed++
		}
	}
	return completed, len(t.Jobs)
}
