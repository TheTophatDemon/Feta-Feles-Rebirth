package main

import (
	"image"

	"github.com/thetophatdemon/Feta-Feles-Remastered/vmath"
)

//To mark points visually for inspection of collision detection
type DebugSpot struct {
	pos *vmath.Vec2f
	spr *Sprite
}

var __debugSpots []*DebugSpot

func AddDebugSpot(x, y float64, color int) {
	var spr *Sprite
	switch color {
	case 0:
		spr = NewSprite(image.Rect(112, 40, 116, 44), vmath.NewVec(-3.0, -3.0), false, false, 0)
	case 1:
		spr = NewSprite(image.Rect(136, 40, 140, 44), vmath.NewVec(-2.0, -2.0), false, false, 0)
	case 2:
		spr = NewSprite(image.Rect(104, 40, 108, 44), vmath.NewVec(-2.0, -2.0), false, false, 0)
	}
	__debugSpots = append(__debugSpots, &DebugSpot{vmath.NewVec(x, y), spr})
}

func ClearDebugSpots() {
	__debugSpots = make([]*DebugSpot, 0, 10)
}
