package main

import (
	"image"
	"math"

	"github.com/hajimehoshi/ebiten"
)

type Sprite struct {
	subImg *ebiten.Image
	matrix *ebiten.GeoM
}

func NewSprite(src image.Rectangle, ofs *Vec2f, flipH, flipV bool, orient int) *Sprite {
	subImg := GetGraphics().SubImage(src).(*ebiten.Image)
	matrix := new(ebiten.GeoM)

	//Perform rotation and scaling with respect to the center
	hw := float64(subImg.Bounds().Dx()) / 2.0
	hh := float64(subImg.Bounds().Dy()) / 2.0
	matrix.Translate(-hw, -hh)
	scx := 1.0
	if flipH {
		scx = -1.0
	}
	scy := 1.0
	if flipV {
		scy = -1.0
	}
	matrix.Scale(scx, scy)
	matrix.Rotate(float64(orient) * math.Pi / 2.0)
	matrix.Translate(hw, hh)
	matrix.Translate(ofs.x, ofs.y)

	return &Sprite{subImg, matrix}
}

func (spr *Sprite) Draw(target *ebiten.Image, pt *ebiten.GeoM) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM = *spr.matrix
	op.GeoM.Concat(*pt)
	target.DrawImage(spr.subImg, op)
}
