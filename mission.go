package main

type Mission struct {
	loveQuota               int
	maxKnights, knightCount int
	maxBlarghs, blarghCount int
	maxGopniks, gopnikCount int
	maxBarrels, barrelCount int
	mapWidth, mapHeight     int
}

var missions []Mission

func init() {
	missions = []Mission{
		Mission{ //Tutorial
			loveQuota:  5,
			maxKnights: 2,
			mapWidth:   32, mapHeight: 32,
		},
		Mission{ //1
			loveQuota:  50,
			maxKnights: 3, maxBlarghs: 3, maxBarrels: 3,
			mapWidth: 32, mapHeight: 32,
		},
		Mission{ //2
			loveQuota:  100,
			maxKnights: 15, maxBlarghs: 10, maxGopniks: 3, maxBarrels: 5,
			mapWidth: 64, mapHeight: 64,
		},
		Mission{ //3
			loveQuota:  100,
			maxKnights: 15, maxBlarghs: 15, maxGopniks: 7, maxBarrels: 10,
			mapWidth: 48, mapHeight: 48,
		},
		Mission{ //4
			loveQuota:  150,
			maxKnights: 20, maxBlarghs: 20, maxGopniks: 20, maxBarrels: 15,
			mapWidth: 64, mapHeight: 64,
		},
		Mission{ //5
			loveQuota:  150,
			maxKnights: 30, maxBlarghs: 30, maxGopniks: 25, maxBarrels: 25,
			mapWidth: 72, mapHeight: 72,
		},
		Mission{ //6
			loveQuota:  200,
			maxKnights: 35, maxBlarghs: 35, maxGopniks: 25, maxBarrels: 30,
			mapWidth: 48, mapHeight: 72,
		},
	}
}
