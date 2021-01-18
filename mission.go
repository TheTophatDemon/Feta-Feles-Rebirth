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
		Mission{
			loveQuota:  5,
			maxKnights: 2,
			maxBlarghs: 0,
			maxGopniks: 0,
			maxBarrels: 0,
			mapWidth:   32,
			mapHeight:  32,
		},
		Mission{
			loveQuota:  50,
			maxKnights: 5,
			maxBlarghs: 5,
			maxGopniks: 0,
			maxBarrels: 0,
			mapWidth:   32,
			mapHeight:  32,
		},
		Mission{
			loveQuota:  100,
			maxKnights: 15,
			maxBlarghs: 10,
			maxGopniks: 5,
			maxBarrels: 5,
			mapWidth:   64,
			mapHeight:  64,
		},
		Mission{
			loveQuota:  100,
			maxKnights: 15,
			maxBlarghs: 15,
			maxGopniks: 15,
			maxBarrels: 10,
			mapWidth:   64,
			mapHeight:  64,
		},
	}
}
