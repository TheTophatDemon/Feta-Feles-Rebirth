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

type UINodeRenderer interface {
	Draw(target *ebiten.Image, globalTransform *ebiten.GeoM)
	Regen() //Called when node has changed and sprites should be updated
}

type UINode struct {
	children *list.List
	parent   *UINode
	dest     image.Rectangle //Area on screen relative to parent
	render   UINodeRenderer
	visible  bool
}

func EmptyUINode() *UINode {
	return &UINode{
		children: list.New(),
		parent:   nil,
		dest:     image.Rect(0, 0, SCR_WIDTH, SCR_HEIGHT),
		render:   nil,
		visible:  true,
	}
}

func NewUINode(dest image.Rectangle, render UINodeRenderer) UINode {
	return UINode{
		children: list.New(),
		parent:   nil,
		dest:     dest,
		render:   render,
		visible:  true,
	}
}

func (node *UINode) Width() int {
	return node.dest.Size().X
}

func (node *UINode) Height() int {
	return node.dest.Size().Y
}

//Returns true if user clicked inside of the node's area
func (node *UINode) Clicked() bool {
	mx, my := ebiten.CursorPosition()
	//Get global position by subtracting parent positions
	for n := node.parent; n != nil; n = n.parent {
		mx -= n.dest.Min.X
		my -= n.dest.Min.Y
	}
	return inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) &&
		mx > node.dest.Min.X && mx < node.dest.Max.X && my > node.dest.Min.Y && my < node.dest.Max.Y
}

//Recursively severs relationships between all nodes and its children
func (node *UINode) Unlink() {
	node.parent = nil
	for e := node.children.Front(); e != nil; e = e.Next() {
		child, ok := e.Value.(*UINode)
		if ok {
			child.Unlink()
		}
		node.children.Remove(e)
	}
}

func (node *UINode) AddChild(child *UINode) {
	child.parent = node
	node.children.PushBack(child)
}

//Recursively draws the node and all of its children
func (node *UINode) Draw(target *ebiten.Image, parentTransform *ebiten.GeoM) {
	globalTransform := &ebiten.GeoM{}
	globalTransform.Translate(float64(node.dest.Min.X), float64(node.dest.Min.Y))
	if parentTransform != nil {
		globalTransform.Concat(*parentTransform)
	}
	if node.visible {
		if node.render != nil {
			node.render.Draw(target, globalTransform)
		}
		for e := node.children.Front(); e != nil; e = e.Next() {
			child := e.Value.(*UINode)
			child.Draw(target, globalTransform)
		}
	}
}

type UIBox struct {
	UINode
	sprites []*Sprite
	src     image.Rectangle //Area on graphics page
	border  bool            //True if the image is supposed to form a border. Otherwise, a single image is stretched over the box's area.
}

//Creates 9 streched sprites that form a border around the box's area
func (box *UIBox) Regen() {
	src := box.src
	sx, sy := box.Width(), box.Height()
	var sprites []*Sprite
	if box.border {
		sprites = make([]*Sprite, 9)
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
	} else {
		sprites = []*Sprite{
			SpriteFromScaledImg(GetGraphics().SubImage(src).(*ebiten.Image), image.Rect(0, 0, sx, sy), 0),
		}
	}
	box.sprites = sprites
}

func CreateUIBox(src, dest image.Rectangle, border bool) *UIBox {
	box := &UIBox{
		src:    src,
		border: border,
	}
	box.UINode = NewUINode(dest, box)
	box.Regen()
	return box
}

//Arranges the box's children centered into vertically or horizontally arranged sections for each child
func (box *UIBox) ArrangeChildren(padding image.Rectangle, vertical bool) {
	if box.children.Len() <= 0 {
		return
	}
	domainWidth := box.Width() - padding.Size().X
	domainHeight := box.Height() - padding.Size().Y
	var sectionSize int
	if vertical {
		sectionSize = domainHeight / box.children.Len()
	} else {
		sectionSize = domainWidth / box.children.Len()
	}
	sectionOfs := 0
	for elem := box.children.Front(); elem != nil; elem = elem.Next() {
		child := elem.Value.(*UINode)
		childW := child.dest.Size().X
		childH := child.dest.Size().Y

		var childX, childY int
		if vertical {
			childX = (domainWidth - childW) / 2
			childY = padding.Min.Y + sectionOfs + (sectionSize / 2) - (childH / 2)
		} else {
			childX = padding.Min.X + sectionOfs + (sectionSize / 2) - (childW / 2)
			childY = (domainHeight - childH) / 2
		}
		child.dest = image.Rect(childX, childY, childX+childW, childY+childH)

		if child.render != nil {
			child.render.Regen()
		}

		sectionOfs += sectionSize
	}
}

func (ui *UIBox) Draw(target *ebiten.Image, globalTransform *ebiten.GeoM) {
	op := &ebiten.DrawImageOptions{}
	for _, sp := range ui.sprites {
		if sp != nil {
			op.GeoM = *sp.matrix
			if globalTransform != nil {
				op.GeoM.Concat(*globalTransform)
			}
			target.DrawImage(sp.subImg, op)
		}
	}
}

type UIText struct {
	UINode
	sprites   []*Sprite
	text      string
	fillPos   int //The character before which to stop rendering. Used for 'typing in' effect
	fillTimer float64
	fillSpeed float64
	fillSound string
}

func (text *UIText) Regen() {
	text.sprites = make([]*Sprite, len(text.text))
	lineLen := text.Width() / 8
	for i := 0; i < len(text.text); i++ {
		r, _ := utf8.DecodeRuneInString(text.text[i:])
		if r > ' ' && r <= 'Z' {
			charDestX := (i % lineLen) * 8
			charDestY := (i / lineLen) * 8
			charSrcX := (int(r-' ')%12)*8 + 64
			charSrcY := (int(r-' ') / 12) * 8
			text.sprites[i] = NewSprite(image.Rect(charSrcX, charSrcY, charSrcX+8, charSrcY+8), vmath.NewVec(float64(charDestX), float64(charDestY)), false, false, 0)
		}
	}
}

func GenerateText(text string, dest image.Rectangle) *UIText {
	text = strings.ToUpper(text)
	uiText := &UIText{
		text:      text,
		fillPos:   len(text),
		fillTimer: 0.0,
		fillSpeed: 0.04,
		fillSound: "voice",
	}
	uiText.UINode = NewUINode(dest, uiText)
	uiText.Regen()
	return uiText
}

func (text *UIText) Update(deltaTime float64) {
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

func (text *UIText) Draw(target *ebiten.Image, globalTransform *ebiten.GeoM) {
	op := &ebiten.DrawImageOptions{}
	for i := 0; i < text.fillPos; i++ {
		sp := text.sprites[i]
		if sp != nil {
			op.GeoM = *sp.matrix
			if globalTransform != nil {
				op.GeoM.Concat(*globalTransform)
			}
			target.DrawImage(sp.subImg, op)
		}
	}
}
