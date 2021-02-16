package main

import (
	"image/color"
)

type Mission struct {
	loveQuota           int
	maxKnights          int
	maxBlarghs          int
	maxGopniks          int
	maxWorms            int
	maxBarrels          int
	catHealth           int
	mapWidth, mapHeight int
	bgColor1, bgColor2  color.RGBA
	parTime             int  //Time in seconds under which the mission must be completed in order to get the good ending
	goodEndFlag         bool //Set after mission completion
}

var missions []Mission

func init() {
	missions = []Mission{
		{ //Tutorial
			loveQuota:  25,
			maxKnights: 3,
			catHealth:  3,
			mapWidth:   32, mapHeight: 32,
			bgColor1: color.RGBA{91, 110, 225, 255},
			bgColor2: color.RGBA{21, 52, 225, 255},
			parTime:  60,
		},
		{ //1 (Cat)
			loveQuota:  70,
			maxKnights: 3, maxBlarghs: 3, maxBarrels: 3,
			catHealth: 3,
			mapWidth:  32, mapHeight: 32,
			bgColor1: color.RGBA{91, 110, 225, 255},
			bgColor2: color.RGBA{48, 96, 130, 255},
			parTime:  60 + 45,
		},
		{ //2 (Human)
			loveQuota:  100,
			maxKnights: 15, maxBlarghs: 10, maxGopniks: 2, maxBarrels: 5,
			catHealth: 3,
			mapWidth:  64, mapHeight: 64,
			bgColor1: color.RGBA{48, 96, 130, 255},
			bgColor2: color.RGBA{48, 96, 130, 255},
			parTime:  (2 * 60) + 30,
		},
		{ //3 (Angel)
			loveQuota:  100,
			maxKnights: 15, maxBlarghs: 15, maxGopniks: 7, maxBarrels: 10,
			catHealth: 5,
			mapWidth:  48, mapHeight: 48,
			bgColor1: color.RGBA{50, 60, 57, 255},
			bgColor2: color.RGBA{89, 86, 82, 255},
			parTime:  (3 * 60) + 45,
		},
		{ //4 (Corrupt)
			loveQuota:  150,
			maxKnights: 20, maxBlarghs: 20, maxGopniks: 16, maxBarrels: 15, maxWorms: 1,
			catHealth: 5,
			mapWidth:  64, mapHeight: 64,
			bgColor1: color.RGBA{34, 32, 32, 255},
			bgColor2: color.RGBA{0, 0, 0, 255},
			parTime:  (4 * 60) + 15,
		},
		{ //5 (Melting)
			loveQuota:  150,
			maxKnights: 25, maxBlarghs: 25, maxGopniks: 20, maxBarrels: 20, maxWorms: 5,
			catHealth: 7,
			mapWidth:  72, mapHeight: 72,
			bgColor1: color.RGBA{0, 0, 0, 255},
			bgColor2: color.RGBA{0, 0, 0, 255},
			parTime:  (5 * 60),
		},
		{ //6 (Monster)
			loveQuota:  200,
			maxKnights: 30, maxBlarghs: 30, maxGopniks: 25, maxBarrels: 30, maxWorms: 10,
			catHealth: 7,
			mapWidth:  48, mapHeight: 72,
			bgColor1: color.RGBA{0, 0, 0, 255},
			bgColor2: color.RGBA{186, 32, 32, 255},
			parTime:  (5 * 60) + 30,
		},
	}
}
