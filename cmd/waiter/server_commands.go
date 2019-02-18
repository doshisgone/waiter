package main

import (
	"strconv"
	"strings"

	"github.com/sauerbraten/waiter/internal/definitions/nmc"
	"github.com/sauerbraten/waiter/internal/definitions/role"
	"github.com/sauerbraten/waiter/pkg/protocol/cubecode"
)

func (s *Server) HandleCommand(c *Client, msg string) {
	msg = strings.TrimSpace(msg)
	parts := strings.Split(msg, " ")
	cmd := parts[0]

	switch cmd {
	case "help":
		switch c.Role {
		case role.None:
			c.Send(nmc.ServerMessage, "available commands: keepteams (=persist), queuemap")
		}

	case "persist", "persistteams", "keepteams":
		if c.Role == role.None {
			return
		}
		changed := false
		if len(parts) > 1 {
			val, err := strconv.Atoi(parts[1])
			if err != nil || (val != 0 && val != 1) {
				return
			}
			changed = s.KeepTeams != (val == 1)
			s.KeepTeams = val == 1
		}
		if changed {
			if s.KeepTeams {
				s.Clients.Broadcast(nil, nmc.ServerMessage, "teams will be kept")
			} else {
				s.Clients.Broadcast(nil, nmc.ServerMessage, "teams will be shuffled")
			}
		} else {
			if s.KeepTeams {
				c.Send(nmc.ServerMessage, "teams will be kept")
			} else {
				c.Send(nmc.ServerMessage, "teams will be shuffled")
			}
		}

	case "queuemap", "queuedmap", "queuemaps", "queuedmaps", "mapqueue", "mapsqueue":
		if c.Role == role.None {
			return
		}
		for _, mapp := range parts[1:] {
			err := s.MapRotation.queueMap(mapp)
			if err != "" {
				c.Send(nmc.ServerMessage, cubecode.Fail(err))
			}
		}
		switch len(s.MapRotation.queue) {
		case 0:
			c.Send(nmc.ServerMessage, "no maps queued")
		case 1:
			c.Send(nmc.ServerMessage, "queued map: "+s.MapRotation.queue[0])
		default:
			c.Send(nmc.ServerMessage, "queued maps: "+strings.Join(s.MapRotation.queue, ", "))
		}

	default:
		c.Send(nmc.ServerMessage, cubecode.Fail("unknown command"))
	}
}
