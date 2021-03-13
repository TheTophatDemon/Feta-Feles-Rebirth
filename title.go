package main

import (
	"image"
	"math/rand"
	"strings"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
	"github.com/thetophatdemon/Feta-Feles-Remastered/vmath"
)

type TitleScreen struct {
	title           *Object
	logo            *Object
	feles           *Object
	uiRoot          *UINode
	link            *UIText
	enterText       *UIText
	flinchTimer     float64
	blinkTimer      float64
	missionSelect   bool
	goodEnd, badEnd bool //Flags for when you return to the title screen after beating the game
}

func (ts *TitleScreen) Enter() {
	ts.title = ts.GenerateTitle()
	ts.logo = &Object{
		sprites: []*Sprite{
			NewSprite(image.Rect(64, 48, 80, 80), vmath.NewVec(SCR_WIDTH_H-16, SCR_HEIGHT-48), false, false, 0),
			NewSprite(image.Rect(64, 48, 80, 80), vmath.NewVec(SCR_WIDTH_H, SCR_HEIGHT-48), true, false, 0),
		},
		pos: vmath.ZeroVec(),
	}
	ts.uiRoot = EmptyUINode()
	ts.link = GenerateText("tophatdemon.com", image.Rect(SCR_WIDTH_H-60, SCR_HEIGHT-16, SCR_WIDTH_H+64, SCR_HEIGHT))
	ts.uiRoot.AddChild(&ts.link.UINode)
	ts.enterText = GenerateText("CLICK OR SPACE TO BEGIN", image.Rect(SCR_WIDTH_H-10*8-12, SCR_HEIGHT_H+40.0, SCR_WIDTH_H+10*8+12, SCR_HEIGHT_H+56.0))
	ts.uiRoot.AddChild(&ts.enterText.UINode)
	if ts.goodEnd {
		ts.feles = MakeFeles(FACE_SMILE, BODY_ANGEL, vmath.NewVec(SCR_WIDTH_H, SCR_HEIGHT_H-32.0))
	} else if ts.badEnd {
		ts.feles = MakeFeles(FACE_EMPTY, BODY_CAT, vmath.NewVec(SCR_WIDTH_H, SCR_HEIGHT_H-32.0))
	} else {
		ts.feles = MakeFeles(FACE_WINK, BODY_CAT, vmath.NewVec(SCR_WIDTH_H, SCR_HEIGHT_H-32.0))
	}
}

func (ts *TitleScreen) Leave() {
	ts.uiRoot.Unlink()
}

func (ts *TitleScreen) Update(deltaTime float64) {
	if !ts.missionSelect {
		//Mission select cheat
		cheatText += strings.ToLower(string(ebiten.InputChars()))
		if strings.Contains(cheatText, "tdyeehaw") {
			cheatText = ""
			ts.missionSelect = true
			ts.enterText.text = "PRESS MISSION NUMBER"
			ts.enterText.fillPos = len(ts.enterText.text)
			ts.enterText.Regen()
			ts.enterText.visible = true
		}
		//Good ending cheat
		if strings.Contains(cheatText, "tdbutter") {
			cheatText = ""
			for i := range missions {
				missions[i].goodEndFlag = true
			}
			ChangeAppState(NewCutsceneState(7))
		}
		//Bad ending cheat
		if strings.Contains(cheatText, "tdblyat") {
			cheatText = ""
			for i := range missions {
				missions[i].goodEndFlag = false
			}
			ChangeAppState(NewCutsceneState(7))
		}

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

		if inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			ChangeAppState(NewCutsceneState(0))
		}
	} else {
		switch {
		case inpututil.IsKeyJustPressed(ebiten.Key0):
			ChangeAppState(NewGame(0))
		case inpututil.IsKeyJustPressed(ebiten.Key1):
			ChangeAppState(NewGame(1))
		case inpututil.IsKeyJustPressed(ebiten.Key2):
			ChangeAppState(NewGame(2))
		case inpututil.IsKeyJustPressed(ebiten.Key3):
			ChangeAppState(NewGame(3))
		case inpututil.IsKeyJustPressed(ebiten.Key4):
			ChangeAppState(NewGame(4))
		case inpututil.IsKeyJustPressed(ebiten.Key5):
			ChangeAppState(NewGame(5))
		case inpututil.IsKeyJustPressed(ebiten.Key6):
			ChangeAppState(NewGame(6))
		}
	}
}

func (ts *TitleScreen) Draw(screen *ebiten.Image) {
	ts.title.DrawAllSprites(screen, nil)
	ts.logo.DrawAllSprites(screen, nil)
	ts.feles.DrawAllSprites(screen, nil)
	ts.uiRoot.Draw(screen, nil)
}

func (ts *TitleScreen) GenerateTitle() *Object {
	var titleLetters []image.Rectangle
	if ts.goodEnd {
		titleLetters = []image.Rectangle{
			image.Rect(80, 96, 96, 112),    //T (SGA)
			image.Rect(96, 96, 112, 112),   //H (SGA)
			image.Rect(112, 96, 128, 112),  //E (SGA)
			image.Rect(176, 128, 192, 144), //Space
			image.Rect(144, 96, 160, 112),  //I (SGA)
			image.Rect(160, 80, 176, 96),   //N (SGA)
			image.Rect(112, 80, 128, 96),   //V (SGA)
			image.Rect(128, 96, 144, 112),  //A (SGA)
			image.Rect(144, 80, 160, 96),   //S (SGA)
			image.Rect(144, 96, 160, 112),  //I (SGA)
			image.Rect(160, 96, 176, 112),  //O (SGA)
			image.Rect(160, 80, 176, 96),   //N (SGA)
			image.Rect(176, 128, 192, 144), //Space
			image.Rect(96, 80, 112, 96),    //B (SGA)
			image.Rect(112, 96, 128, 112),  //E (SGA)
			image.Rect(128, 80, 144, 96),   //G (SGA)
			image.Rect(144, 96, 160, 112),  //I (SGA)
			image.Rect(160, 80, 176, 96),   //N (SGA)
			image.Rect(144, 80, 160, 96),   //S (SGA)
		}
	} else {
		titleLetters = []image.Rectangle{
			image.Rect(96, 48, 112, 64),    //F
			image.Rect(112, 48, 128, 64),   //E
			image.Rect(128, 48, 144, 64),   //T
			image.Rect(144, 48, 160, 64),   //A
			image.Rect(176, 128, 192, 144), //Space
			image.Rect(96, 48, 112, 64),    //F
			image.Rect(112, 48, 128, 64),   //E
			image.Rect(112, 64, 128, 80),   //L
			image.Rect(112, 48, 128, 64),   //E
			image.Rect(96, 64, 112, 80),    //S
		}
	}
	sprites := make([]*Sprite, len(titleLetters))
	ofsX := float64(SCR_WIDTH_H - (len(titleLetters) * 16 / 2))
	ofsY := 16.0
	for i, l := range titleLetters {
		r := 0
		if rand.Float64() < 0.025 {
			r = rand.Intn(4)
		}
		sprites[i] = NewSprite(l, vmath.NewVec(ofsX+float64(i)*16.0, ofsY), false, false, r)
	}

	//Subtitle letters
	if !ts.goodEnd {
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
			sprites = append(sprites, NewSprite(l, vmath.NewVec(ofsX+float64(i)*8.0, ofsY), false, false, r))
		}
	}
	return &Object{
		sprites: sprites,
		pos:     vmath.ZeroVec(),
	}
}
