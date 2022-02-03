/*
Copyright (C) 2021 Alexander Lunsford

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

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
	knightSpeed			float64
	mapWidth, mapHeight int
	bgColor1, bgColor2  color.RGBA
	music               string
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
			knightSpeed: 150.0,
			mapWidth:   32, mapHeight: 32,
			bgColor1: color.RGBA{91, 110, 225, 255},
			bgColor2: color.RGBA{21, 52, 225, 255},
			parTime:  90,
			music:    "mystery_ingame",
		},
		{ //1 (Cat)
			loveQuota:  50,
			maxKnights: 3, maxBlarghs: 3, maxBarrels: 6,
			catHealth: 3,
			knightSpeed: 150.0,
			mapWidth:  32, mapHeight: 32,
			bgColor1: color.RGBA{91, 110, 225, 255},
			bgColor2: color.RGBA{48, 96, 130, 255},
			parTime:  120,
			music:    "mystery_ingame",
		},
		{ //2 (Human)
			loveQuota:  75,
			maxKnights: 15, maxBlarghs: 10, maxGopniks: 2, maxBarrels: 7,
			catHealth: 6,
			knightSpeed: 175.0,
			mapWidth:  64, mapHeight: 64,
			bgColor1: color.RGBA{48, 96, 130, 255},
			bgColor2: color.RGBA{48, 96, 130, 255},
			parTime:  (3 * 60),
			music:    "hope_ingame",
		},
		{ //3 (Angel)
			loveQuota:  75,
			maxKnights: 15, maxBlarghs: 15, maxGopniks: 7, maxBarrels: 10,
			catHealth: 8,
			knightSpeed: 175.0,
			mapWidth:  48, mapHeight: 48,
			bgColor1: color.RGBA{160, 0, 160, 255},
			bgColor2: color.RGBA{160, 15, 160, 255},
			parTime:  (4 * 60),
			music:    "hope_ingame",
		},
		{ //4 (Corrupt)
			loveQuota:  85,
			maxKnights: 20, maxBlarghs: 20, maxGopniks: 16, maxBarrels: 15, maxWorms: 1,
			catHealth: 8,
			knightSpeed: 175.0,
			mapWidth:  64, mapHeight: 64,
			bgColor1: color.RGBA{34, 32, 32, 255},
			bgColor2: color.RGBA{0, 0, 0, 255},
			parTime:  (4 * 60) + 30,
			music:    "malform_ingame",
		},
		{ //5 (Melting)
			loveQuota:  100,
			maxKnights: 25, maxBlarghs: 25, maxGopniks: 20, maxBarrels: 20, maxWorms: 5,
			catHealth: 10,
			knightSpeed: 175.0,
			mapWidth:  72, mapHeight: 72,
			bgColor1: color.RGBA{0, 0, 0, 255},
			bgColor2: color.RGBA{0, 0, 0, 255},
			parTime:  (5 * 60),
			music:    "malform_ingame",
		},
		{ //6 (Monster)
			loveQuota:  100,
			maxKnights: 30, maxBlarghs: 30, maxGopniks: 25, maxBarrels: 30, maxWorms: 10,
			catHealth: 10,
			knightSpeed: 175.0,
			mapWidth:  48, mapHeight: 72,
			bgColor1: color.RGBA{0, 0, 0, 255},
			bgColor2: color.RGBA{186, 32, 32, 255},
			parTime:  (5 * 60) + 30,
		},
	}
}
