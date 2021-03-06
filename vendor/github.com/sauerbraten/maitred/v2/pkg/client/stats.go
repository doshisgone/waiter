package client

import (
	"fmt"
	"log"
	"strings"

	"github.com/sauerbraten/maitred/v2/pkg/protocol"
)

type StatsClient struct {
	*VanillaClient
	onSuccess func(reqID uint32)
	onFailure func(reqID uint32, reason string)
}

func NewStats(c *VanillaClient, onSuccess func(uint32), onFailure func(uint32, string)) *StatsClient {
	return &StatsClient{
		VanillaClient: c,
		onSuccess:     onSuccess,
		onFailure:     onFailure,
	}
}

func (c *StatsClient) Handle(msg string) {
	cmd := strings.Split(msg, " ")[0]
	args := strings.TrimSpace(msg[len(cmd):])

	switch cmd {
	case protocol.FailStats:
		c.handleFailStats(args)

	case protocol.SuccStats:
		c.handleSuccStats(args)

	default:
		c.VanillaClient.Handle(msg)
	}
}

func (c *StatsClient) handleFailStats(args string) {
	r := strings.NewReader(strings.TrimSpace(args))

	var reqID uint32
	_, err := fmt.Fscanf(r, "%d", &reqID)
	if err != nil {
		log.Printf("malformed %s message from stats server: '%s': %v", protocol.FailStats, args, err)
		return
	}

	reason := args[len(args)-r.Len():] // unread portion of args
	reason = strings.TrimSpace(reason)

	c.onFailure(reqID, reason)
}

func (c *StatsClient) handleSuccStats(args string) {
	var reqID uint32
	_, err := fmt.Sscanf(args, "%d", &reqID)
	if err != nil {
		log.Printf("malformed %s message from stats server: '%s': %v", protocol.SuccStats, args, err)
		return
	}

	c.onSuccess(reqID)
}
