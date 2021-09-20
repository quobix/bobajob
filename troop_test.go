package bobajob

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewTroop(t *testing.T) {
	tr := NewTroop("ember")
	assert.NotNil(t, tr.CompletedChan)
}

func TestTroop_Run_NoJobs(t *testing.T) {

	tr := NewTroop("pizza")

	go func() {
		tr.Run()
	}()

	for {
		select {
		case _, ok := <-tr.CompletedChan:
			assert.False(t, ok)
			return
		}
	}

}


func TestTroop_Run_AFewJobs(t *testing.T) {

	tr := NewTroop("cake")

	tr.AddJob(&Job{})
	tr.AddJob(&Job{})
	tr.AddJob(&Job{})

	go tr.Run()

	for {
		select {
		case _, ok := <-tr.CompletedChan:
			assert.True(t, ok)
			completed, total := tr.JobsCompleted()
			assert.Equal(t, 3, completed)
			assert.Equal(t, 3, total)
			break
		}
		return
	}
}