package main

import (
	"container/list"
	"image"
	"strings"
	"unicode/utf8"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
	"github.com/thetophatdemon/Feta-Feles-Remastered/audio"
	"github.com/thetophatdemon/Feta-Feles-Remastered/vmath"
)

type UINode interface {
	Draw(target *ebiten.Image, parentTransform *ebiten.GeoM)
}

type UIBox struct {
	sprites   []*Sprite
	src       image.Rectangle //Area on graphics page
	size      image.Point
	dest      image.Rectangle //Area on screen
	transform *ebiten.GeoM
	children  *list.List
	padding   image.Rectangle //Padding of children from border of container
	regen     bool            //True if position has changed and transform needs to be recalculated
}

func (ui *UIBox) GenerateSprites() {
	src := ui.src
	sx, sy := ui.size.X, ui.size.Y
	sprites := make([]*Sprite, 9)
	bodyImg := GetGraphics().SubImage(image.Rect(src.Min.X+16, src.Min.Y, src.Min.X+24, src.Max.Y)).(*ebiten.Image)
	sprites[0] = SpriteFromScaledImg(bodyImg, image.Rect(8, 8, sx-8, sy-8), 0) //Middle
	edgeImg := GetGraphics().SubImage(image.Rect(src.Min.X+8, src.Min.Y, src.Min.X+16, src.Max.Y)).(*ebiten.Image)
	sprites[1] = SpriteFromScaledImg(edgeImg, image.Rect(8, 0, sx-8, 8), 0)     //Top
	sprites[2] = SpriteFromScaledImg(edgeImg, image.Rect(sx-8, 8, sx, sy-8), 1) //Right
	sprites[3] = SpriteFromScaledImg(edgeImg, image.Rect(8, sy-8, sx-8, sy), 2) //Bottom
	sprites[4] = SpriteFromScaledImg(edgeImg, image.Rect(0, 8, 8, sy-8), 3)     //Left
	cornerImg := GetGraphics().SubImage(image.Rect(src.Min.X, src.Min.Y, src.Min.X+8, src.Max.Y)).(*ebiten.Image)
	sprites[5] = SpriteFromScaledImg(cornerImg, image.Rect(0, 0, 8, 8), 0)         //Top left
	sprites[6] = SpriteFromScaledImg(cornerImg, image.Rect(sx-8, 0, sx, 8), 1)     //Top right
	sprites[7] = SpriteFromScaledImg(cornerImg, image.Rect(sx-8, sy-8, sx, sy), 2) //Bottom right
	sprites[8] = SpriteFromScaledImg(cornerImg, image.Rect(0, sy-8, 8, sy), 3)     //Bottom left
	ui.sprites = sprites
}

func CreateUIBox(src, dest image.Rectangle) *UIBox {
	trans := &ebiten.GeoM{}
	return &UIBox{
		sprites:   []*Sprite{},
		src:       src,
		dest:      dest,
		size:      dest.Size(),
		children:  list.New(),
		regen:     true,
		padding:   image.Rect(4, 4, 4, 4),
		transform: trans,
	}
}

func CreateUIBoxSize(src image.Rectangle, corner image.Point, size image.Point) *UIBox {
	return CreateUIBox(src, image.Rectangle{Min: corner, Max: corner.Add(size)})
}

//Arranges the box's children centered into vertically arranged sections for each child
func (ui *UIBox) ArrangeChildren() {
	if ui.children.Len() <= 0 {
		return
	}
	domainWidth := ui.size.X - ui.padding.Size().X
	domainHeight := ui.size.Y - ui.padding.Size().Y
	sectionHeight := domainHeight / ui.children.Len()
	sectionY := 0
	for elem := ui.children.Front(); elem != nil; elem = elem.Next() {
		child, ok := elem.Value.(*UIBox)
		if ok {
			xofs := (domainWidth - child.size.X) / 2
			childX := ui.padding.Min.X + xofs
			childY := ui.padding.Min.Y + sectionY
			child.dest = image.Rect(childX, childY, childX+child.size.X, childY+child.size.Y)
			child.regen = true
			sectionY += sectionHeight
		}
	}
}

func (ui *UIBox) IsClicked() bool {
	mx, my := ebiten.CursorPosition()
	return inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) &&
		mx > ui.dest.Min.X && mx < ui.dest.Max.X && my > ui.dest.Min.Y && my < ui.dest.Max.Y
}

func (ui *UIBox) Draw(target *ebiten.Image, parentTransform *ebiten.GeoM) {
	if ui.regen {
		ui.GenerateSprites()
		ui.transform = &ebiten.GeoM{}
		ui.transform.Translate(float64(ui.dest.Min.X), float64(ui.dest.Min.Y))
		ui.regen = false
	}
	globalTrans := *ui.transform
	if parentTransform != nil {
		globalTrans.Concat(*parentTransform)
	}
	op := &ebiten.DrawImageOptions{}
	for _, sp := range ui.sprites {
		if sp != nil {
			op.GeoM = *sp.matrix
			op.GeoM.Concat(globalTrans)
			target.DrawImage(sp.subImg, op)
		}
	}
	for childElem := ui.children.Front(); childElem != nil; childElem = childElem.Next() {
		child, ok := childElem.Value.(UINode)
		if ok {
			child.Draw(target, &globalTrans)
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
			charDestX := (i % lineLen) * 8
			charDestY := (i / lineLen) * 8
			charSrcX := (int(r-' ')%12)*8 + 64
			charSrcY := (int(r-' ') / 12) * 8
			sprites[i] = NewSprite(image.Rect(charSrcX, charSrcY, charSrcX+8, charSrcY+8), vmath.NewVec(float64(charDestX), float64(charDestY)), false, false, 0)
		}
	}
	return &Text{
		UIBox: UIBox{
			sprites:   sprites,
			src:       image.Rect(0, 0, 0, 0),
			size:      dest.Size(),
			dest:      dest,
			children:  list.New(),
			regen:     true,
			transform: &ebiten.GeoM{},
		},
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
				audio.PlaySound(text.fillSound)
			}
			text.fillPos++
		}
	}
}

func (text *Text) Draw(target *ebiten.Image, parentTransform *ebiten.GeoM) {
	if text.regen {
		text.transform = &ebiten.GeoM{}
		text.transform.Translate(float64(text.dest.Min.X), float64(text.dest.Min.Y))
		text.regen = false
	}
	globalTrans := *text.transform
	if parentTransform != nil {
		globalTrans.Concat(*parentTransform)
	}
	op := &ebiten.DrawImageOptions{}
	for i := 0; i < text.fillPos; i++ {
		sp := text.sprites[i]
		if sp != nil {
			op.GeoM = *sp.matrix
			op.GeoM.Concat(globalTrans)
			target.DrawImage(sp.subImg, op)
		}
	}
}
