package maprotation

import (
	"github.com/sauerbraten/waiter/internal/definitions/gamemode"
	"github.com/sauerbraten/waiter/internal/utils"
)

// temporary set of maps used in development phase
var (
	instaMaps = []string{
		"hashi",
		"turbine",
		"ot",
		"memento",
		"kffa",
	}
	efficMaps    = instaMaps
	tacticsMaps  = instaMaps
	efficCTFMaps = []string{
		"reissen",
		"forge",
		"haste",
		"dust2",
		"redemption",
	}
	captureMaps = []string{
		"nmp8",
		"nmp9",
		"nmp4",
		"nevil_c",
		"serenity",
	}
	mr = map[gamemode.ID][]string{
		gamemode.Insta:    instaMaps,
		gamemode.Effic:    efficMaps,
		gamemode.Tactics:  tacticsMaps,
		gamemode.EfficCTF: efficCTFMaps,
		gamemode.Capture:  captureMaps,
	}
)

func NextMap(mode gamemode.ID, currentMap string) string {
	for i, m := range mr[mode] {
		if m == currentMap {
			return mr[mode][(i+1)%len(mr[mode])]
		}
	}

	// current map wasn't found in map rotation, return random map in rotation
	return mr[mode][utils.RNG.Intn(len(mr[mode]))]
}
