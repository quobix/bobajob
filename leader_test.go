package bobajob

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewLeader(t *testing.T) {
	l := NewLeader(LeaderConfig{})
	assert.Contains(t, l.config.name, "bobajob-leader")
}

func TestLeader_Run_NoTroops(t *testing.T) {
	l := NewLeader(LeaderConfig{})
	ch, err := l.Run()
	assert.Error(t, err)
	assert.Nil(t, ch)
}

func TestLeader_AddTroop(t *testing.T) {
	l := NewLeader(LeaderConfig{})
	l.AddTroop(NewTroop("pizza"))
	assert.Len(t, l.troops, 1)
	l.AddTroop(NewTroop("cake"))
	assert.Len(t, l.troops, 2)
}

func TestLeader_Connect_NoBrokerConfig(t *testing.T) {
	l := NewLeader(LeaderConfig{})
	e := l.Connect()
	assert.Error(t, e)
}