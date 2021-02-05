package main

import (
	"image/color"
)

type Mission struct {
	loveQuota               int
	maxKnights, knightCount int
	maxBlarghs, blarghCount int
	maxGopniks, gopnikCount int
	maxBarrels, barrelCount int
	mapWidth, mapHeight     int
	bgColor1, bgColor2      color.RGBA
}

var missions []Mission

func init() {
	missions = []Mission{
		{ //Tutorial
			loveQuota:  10,
			maxKnights: 1,
			mapWidth:   32, mapHeight: 32,
			bgColor1: color.RGBA{91, 110, 225, 255},
			bgColor2: color.RGBA{21, 52, 225, 255},
		},
		{ //1 (Cat)
			loveQuota:  70,
			maxKnights: 3, maxBlarghs: 3, maxBarrels: 3,
			mapWidth: 32, mapHeight: 32,
			bgColor1: color.RGBA{91, 110, 225, 255},
			bgColor2: color.RGBA{48, 96, 130, 255},
		},
		{ //2 (Human)
			loveQuota:  100,
			maxKnights: 15, maxBlarghs: 10, maxGopniks: 3, maxBarrels: 5,
			mapWidth: 64, mapHeight: 64,
			bgColor1: color.RGBA{48, 96, 130, 255},
			bgColor2: color.RGBA{48, 96, 130, 255},
		},
		{ //3 (Angel)
			loveQuota:  100,
			maxKnights: 15, maxBlarghs: 15, maxGopniks: 7, maxBarrels: 10,
			mapWidth: 48, mapHeight: 48,
			bgColor1: color.RGBA{50, 60, 57, 255},
			bgColor2: color.RGBA{89, 86, 82, 255},
		},
		{ //4 (Corrupt)
			loveQuota:  150,
			maxKnights: 20, maxBlarghs: 20, maxGopniks: 20, maxBarrels: 15,
			mapWidth: 64, mapHeight: 64,
			bgColor1: color.RGBA{34, 32, 32, 255},
			bgColor2: color.RGBA{0, 0, 0, 255},
		},
		{ //5 (Melting)
			loveQuota:  150,
			maxKnights: 30, maxBlarghs: 30, maxGopniks: 25, maxBarrels: 25,
			mapWidth: 72, mapHeight: 72,
			bgColor1: color.RGBA{0, 0, 0, 255},
			bgColor2: color.RGBA{0, 0, 0, 255},
		},
		{ //6 (Monster)
			loveQuota:  200,
			maxKnights: 35, maxBlarghs: 35, maxGopniks: 25, maxBarrels: 30,
			mapWidth: 48, mapHeight: 72,
			bgColor1: color.RGBA{0, 0, 0, 255},
			bgColor2: color.RGBA{186, 32, 32, 255},
		},
	}
}
