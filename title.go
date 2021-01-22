package main

import (
	"image"
	"math/rand"

	"github.com/hajimehoshi/ebiten"
)

type TitleScreen struct {
	title       *Object
	logo        *Object
	feles       *Object
	link        *Text
	enterText   *Text
	flinchTimer float64
	blinkTimer  float64
}

func (ts *TitleScreen) Enter() {
	ts.title = ts.GenerateTitle()
	ts.logo = &Object{
		sprites: []*Sprite{
			NewSprite(image.Rect(64, 48, 80, 80), &Vec2f{SCR_WIDTH_H - 16, SCR_HEIGHT - 48}, false, false, 0),
			NewSprite(image.Rect(64, 48, 80, 80), &Vec2f{SCR_WIDTH_H, SCR_HEIGHT - 48}, true, false, 0),
		},
		pos: ZeroVec(),
	}
	ts.link = GenerateText("tophatdemon.com", image.Rect(SCR_WIDTH_H-60, SCR_HEIGHT-16, SCR_WIDTH_H+64, SCR_HEIGHT))
	ts.feles = MakeFeles(FACE_WINK, BODY_CAT, &Vec2f{SCR_WIDTH_H, SCR_HEIGHT_H})
	ts.enterText = GenerateText("PRESS ENTER TO BEGIN", image.Rect(SCR_WIDTH_H-10*8, SCR_HEIGHT_H+40.0, SCR_WIDTH_H+10*8, SCR_HEIGHT_H+56.0))
}

func (ts *TitleScreen) Leave() {

}

func (ts *TitleScreen) Update(deltaTime float64) {
	ts.flinchTimer += deltaTime
	if ts.flinchTimer > 0.25 {
		ts.title = ts.GenerateTitle()
		ts.flinchTimer = 0.0
	}
	ts.blinkTimer += deltaTime
	if ts.enterText.fillPos > 0 {
		if ts.blinkTimer > 0.75 {
			ts.blinkTimer = 0.0
			ts.enterText.fillPos = 0
		}
	} else {
		if ts.blinkTimer > 0.25 {
			ts.blinkTimer = 0.0
			ts.enterText.fillPos = len(ts.enterText.text)
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyEnter) {
		ChangeAppState(NewGame(0))
	}
}

func (ts *TitleScreen) Draw(screen *ebiten.Image) {
	ts.title.DrawAllSprites(screen, nil)
	ts.logo.DrawAllSprites(screen, nil)
	ts.link.Draw(screen)
	ts.feles.DrawAllSprites(screen, nil)
	ts.enterText.Draw(screen)
}

func (ts *TitleScreen) GenerateTitle() *Object {
	titleLetters := []image.Rectangle{
		image.Rect(96, 48, 112, 64),    //F
		image.Rect(112, 48, 128, 64),   //E
		image.Rect(128, 48, 144, 64),   //T
		image.Rect(144, 48, 160, 64),   //A
		image.Rect(176, 128, 192, 144), //
		image.Rect(96, 48, 112, 64),    //F
		image.Rect(112, 48, 128, 64),   //E
		image.Rect(112, 64, 128, 80),   //L
		image.Rect(112, 48, 128, 64),   //E
		image.Rect(96, 64, 112, 80),    //S
	}
	sprites := make([]*Sprite, len(titleLetters))
	ofsX := float64(SCR_WIDTH_H - (len(titleLetters) * 16 / 2))
	ofsY := 16.0
	for i, l := range titleLetters {
		r := 0
		if rand.Float64() < 0.025 {
			r = rand.Intn(4)
		}
		sprites[i] = NewSprite(l, &Vec2f{ofsX + float64(i)*16.0, ofsY}, false, false, r)
	}
	subtitleLetters := []image.Rectangle{
		image.Rect(128, 64, 136, 72), //R
		image.Rect(136, 64, 144, 72), //E
		image.Rect(144, 64, 152, 72), //B
		image.Rect(152, 64, 160, 72), //I
		image.Rect(128, 64, 136, 72), //R
		image.Rect(128, 72, 136, 80), //T
		image.Rect(136, 72, 144, 80), //H
	}
	ofsX = float64(SCR_WIDTH_H - (len(subtitleLetters) * 8 / 2))
	ofsY = 40.0
	for i, l := range subtitleLetters {
		r := 0
		if rand.Float64() < 0.025 {
			r = rand.Intn(4)
		}
		sprites = append(sprites, NewSprite(l, &Vec2f{ofsX + float64(i)*8.0, ofsY}, false, false, r))
	}
	return &Object{
		sprites: sprites,
		pos:     ZeroVec(),
	}
}
