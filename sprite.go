package main

import (
	"image"
	"math"
	"strings"
	"unicode/utf8"

	"github.com/hajimehoshi/ebiten"
)

type Sprite struct {
	subImg *ebiten.Image
	matrix *ebiten.GeoM
}

func NewSprite(src image.Rectangle, ofs *Vec2f, flipH, flipV bool, orient int) *Sprite {
	subImg := GetGraphics().SubImage(src).(*ebiten.Image)
	return NewSpriteFromSubImg(subImg, ofs, flipH, flipV, orient)
}

func NewSpriteFromSubImg(subImg *ebiten.Image, ofs *Vec2f, flipH, flipV bool, orient int) *Sprite {
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

func NewSprites(ofs *Vec2f, rects ...image.Rectangle) []*Sprite {
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

type UIBox struct {
	sprites []*Sprite
	rect    image.Rectangle
}

func CreateUIBox(src, dest image.Rectangle) UIBox {
	sprites := make([]*Sprite, 9)
	bodyImg := GetGraphics().SubImage(image.Rect(src.Min.X+16, src.Min.Y, src.Min.X+24, src.Max.Y)).(*ebiten.Image)
	sprites[0] = SpriteFromScaledImg(bodyImg, image.Rect(dest.Min.X+8, dest.Min.Y+8, dest.Max.X-8, dest.Max.Y-8), 0) //Middle
	edgeImg := GetGraphics().SubImage(image.Rect(src.Min.X+8, src.Min.Y, src.Min.X+16, src.Max.Y)).(*ebiten.Image)
	sprites[1] = SpriteFromScaledImg(edgeImg, image.Rect(dest.Min.X+8, dest.Min.Y, dest.Max.X-8, dest.Min.Y+8), 0) //Top
	sprites[2] = SpriteFromScaledImg(edgeImg, image.Rect(dest.Max.X-8, dest.Min.Y+8, dest.Max.X, dest.Max.Y-8), 1) //Right
	sprites[3] = SpriteFromScaledImg(edgeImg, image.Rect(dest.Min.X+8, dest.Max.Y-8, dest.Max.X-8, dest.Max.Y), 2) //Bottom
	sprites[4] = SpriteFromScaledImg(edgeImg, image.Rect(dest.Min.X, dest.Min.Y+8, dest.Min.X+8, dest.Max.Y-8), 3) //Left
	cornerImg := GetGraphics().SubImage(image.Rect(src.Min.X, src.Min.Y, src.Min.X+8, src.Max.Y)).(*ebiten.Image)
	sprites[5] = SpriteFromScaledImg(cornerImg, image.Rect(dest.Min.X, dest.Min.Y, dest.Min.X+8, dest.Min.Y+8), 0) //Top left
	sprites[6] = SpriteFromScaledImg(cornerImg, image.Rect(dest.Max.X-8, dest.Min.Y, dest.Max.X, dest.Min.Y+8), 1) //Top right
	sprites[7] = SpriteFromScaledImg(cornerImg, image.Rect(dest.Max.X-8, dest.Max.Y-8, dest.Max.X, dest.Max.Y), 2) //Bottom right
	sprites[8] = SpriteFromScaledImg(cornerImg, image.Rect(dest.Min.X, dest.Max.Y-8, dest.Min.X+8, dest.Max.Y), 3) //Bottom left
	return UIBox{sprites, dest}
}

func (ui UIBox) Draw(target *ebiten.Image, pt *ebiten.GeoM) {
	op := &ebiten.DrawImageOptions{}
	for _, sp := range ui.sprites {
		if sp != nil {
			op.GeoM = *sp.matrix
			if pt != nil {
				op.GeoM.Concat(*pt)
			}
			target.DrawImage(sp.subImg, op)
		}
	}
}

type Text struct {
	UIBox
	text      string
	fillPos   int //The character before which to stop rendering. Used for 'typing in' effect
	fillTimer float64
	fillSpeed float64
	fillSound string
}

func GenerateText(text string, dest image.Rectangle) *Text {
	text = strings.ToUpper(text)
	sprites := make([]*Sprite, len(text))
	lineLen := dest.Dx() / 8
	for i := 0; i < len(text); i++ {
		r, _ := utf8.DecodeRuneInString(text[i:])
		if r > ' ' && r <= 'Z' {
			charDestX := (i%lineLen)*8 + dest.Min.X
			charDestY := (i/lineLen)*8 + dest.Min.Y
			charSrcX := (int(r-' ')%12)*8 + 64
			charSrcY := (int(r-' ') / 12) * 8
			sprites[i] = NewSprite(image.Rect(charSrcX, charSrcY, charSrcX+8, charSrcY+8), &Vec2f{float64(charDestX), float64(charDestY)}, false, false, 0)
		}
	}
	return &Text{
		UIBox:     UIBox{sprites, dest},
		text:      text,
		fillPos:   len(text),
		fillTimer: 0.0,
		fillSpeed: 0.04,
		fillSound: "voice",
	}
}

func (text *Text) Update(deltaTime float64) {
	text.fillTimer += deltaTime
	if text.fillTimer > text.fillSpeed {
		text.fillTimer = 0.0
		if text.fillPos < len(text.text) {
			if text.text[text.fillPos] > 'A' && text.text[text.fillPos] < 'z' {
				PlaySound(text.fillSound)
			}
			text.fillPos++
		}
	}
}

func (text *Text) Draw(target *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	for i := 0; i < text.fillPos; i++ {
		sp := text.sprites[i]
		if sp != nil {
			op.GeoM = *sp.matrix
			target.DrawImage(sp.subImg, op)
		}
	}
}
