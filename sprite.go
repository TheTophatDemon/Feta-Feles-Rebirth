package main

import (
	"math"

	"github.com/hajimehoshi/ebiten"
)

//Todo: Optimize by keeping track of previous state and recalculating transform matrix
//Todo: Remove middleman and store sub-image instead of src

type Sprite struct {
	subImg *ebiten.Image
	ofs    *Vec2f
	flipH  bool
	flipV  bool
	orient int
}

func (spr *Sprite) Draw(target *ebiten.Image, pt *ebiten.GeoM) {
	op := &ebiten.DrawImageOptions{}
	//Perform rotation and scaling with respect to the center
	hw := float64(spr.subImg.Bounds().Dx()) / 2.0
	hh := float64(spr.subImg.Bounds().Dy()) / 2.0
	op.GeoM.Translate(-hw, -hh)
	scx := 1.0
	if spr.flipH {
		scx = -1.0
	}
	scy := 1.0
	if spr.flipV {
		scy = -1.0
	}
	op.GeoM.Scale(scx, scy)
	op.GeoM.Rotate(float64(spr.orient) * math.Pi / 2.0)
	op.GeoM.Translate(hw, hh)
	op.GeoM.Translate(spr.ofs.x, spr.ofs.y)
	op.GeoM.Concat(*pt)
	target.DrawImage(spr.subImg, op)
}
