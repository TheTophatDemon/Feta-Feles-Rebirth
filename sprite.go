package main

import (
	"image"
	"math"

	"github.com/hajimehoshi/ebiten"
	"github.com/thetophatdemon/Feta-Feles-Remastered/vmath"
)

type Sprite struct {
	subImg *ebiten.Image
	matrix *ebiten.GeoM
}

func NewSprite(src image.Rectangle, ofs *vmath.Vec2f, flipH, flipV bool, orient int) *Sprite {
	subImg := GetGraphics().SubImage(src).(*ebiten.Image)
	return NewSpriteFromSubImg(subImg, ofs, flipH, flipV, orient)
}

func NewSpriteFromSubImg(subImg *ebiten.Image, ofs *vmath.Vec2f, flipH, flipV bool, orient int) *Sprite {
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
	matrix.Translate(ofs.X, ofs.Y)

	return &Sprite{subImg, matrix}
}

func CloneSprite(org *Sprite) *Sprite {
	return &Sprite{
		org.subImg,
		&(*org.matrix),
	}
}

func CloneSprites(orgs []*Sprite) []*Sprite {
	out := make([]*Sprite, len(orgs))
	for i, spr := range orgs {
		out[i] = CloneSprite(spr)
	}
	return out
}

func (spr *Sprite) Draw(target *ebiten.Image, pt *ebiten.GeoM) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM = *spr.matrix
	if pt != nil {
		op.GeoM.Concat(*pt)
	}
	target.DrawImage(spr.subImg, op)
}

func NewSprites(ofs *vmath.Vec2f, rects ...image.Rectangle) []*Sprite {
	sprites := make([]*Sprite, len(rects))
	for i, rect := range rects {
		sprites[i] = NewSprite(rect, ofs, false, false, 0)
	}
	return sprites
}

func (spr *Sprite) Flip(horz, vert bool) *Sprite {
	matrix := new(ebiten.GeoM)

	//Perform rotation and scaling with respect to the center
	hw := float64(spr.subImg.Bounds().Dx()) / 2.0
	hh := float64(spr.subImg.Bounds().Dy()) / 2.0
	matrix.Translate(-hw, -hh)
	scx, scy := 1.0, 1.0
	if horz {
		scx = -1.0
	}
	if vert {
		scy = -1.0
	}
	matrix.Scale(scx, scy)
	matrix.Translate(hw, hh)

	matrix.Concat(*spr.matrix)
	spr.matrix = matrix
	return spr
}

//Creates a sprite that fills a given rectangle on the screen
func SpriteFromScaledImg(subImg *ebiten.Image, dest image.Rectangle, orient int) *Sprite {
	matrix := new(ebiten.GeoM)

	//Perform rotation
	hw := float64(subImg.Bounds().Dx()) / 2.0
	hh := float64(subImg.Bounds().Dy()) / 2.0
	matrix.Translate(-hw, -hh)
	matrix.Rotate(float64(orient) * math.Pi / 2.0)
	matrix.Translate(hw, hh)

	//Perform scaling
	matrix.Scale(
		float64(dest.Dx())/float64(subImg.Bounds().Dx()),
		float64(dest.Dy())/float64(subImg.Bounds().Dy()))
	//Offset
	matrix.Translate(float64(dest.Min.X), float64(dest.Min.Y))

	return &Sprite{subImg, matrix}
}
