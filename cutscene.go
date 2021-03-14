package main

import (
	"image"
	"math"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
	"github.com/thetophatdemon/Feta-Feles-Remastered/audio"
	"github.com/thetophatdemon/Feta-Feles-Remastered/vmath"
)

type Cutscene struct {
	bodyType BodyType
	faces    []FaceType
	dialog   []string
	music    string
	voice    string
}

var cutscenes []Cutscene

func init() {
	cutscenes = []Cutscene{
		{
			bodyType: BODY_NONE,
			faces:    []FaceType{FACE_NONE},
			dialog: []string{ //26 chars per line
				"HELLO THERE...",
				"DON'T BE ALARMED! I'M HERE TO HELP.",
				"BUT I NEED YOU TO DO ME A SIMPLE FAVOR.",
				"I WANT YOU TO MAKE ME...  REAL!",
				"LET ME TELL YOU HOW.",
			},
			music: "mystery",
			voice: "voice",
		},
		{
			bodyType: BODY_CAT,
			faces:    []FaceType{FACE_TALK, FACE_SMILE, FACE_TALK, FACE_SMILE, FACE_WINK, FACE_TALK},
			dialog: []string{
				"HAHA! YES!!",
				"IT'S SO WONDERFUL BEING   HERE.",
				"THANK YOU!",
				"BUT...IT'S SUCH A DARK    WORLD THAT YOU LIVE IN...",
				"SURELY WE CAN MAKE THINGS BETTER!",
				"BUT FIRST WE HAVE TO GROW STRONGER.",
			},
			music: "mystery",
			voice: "voice",
		},
		{
			bodyType: BODY_HUMAN,
			faces:    []FaceType{FACE_TALK, FACE_WINK, FACE_SMILE, FACE_SMILE, FACE_TALK},
			dialog: []string{
				"ARE YOU CONFUSED?",
				"WHAT CAN I SAY...I SEEM TO HAVE GROWN ACCUSTOMED TO LIFE.",
				"I FEEL LIKE THIS WAS ALL  MEANT TO HAPPEN.",
				"CAN YOU FEEL THE POWER IN THE AIR? WE CAN DO ALL    SORTS OF GREAT THINGS NOW!",
				"WHY NOT CLIMB EVEN HIGHER?",
			},
			music: "hope",
			voice: "voice",
		},
		{
			bodyType: BODY_ANGEL2,
			faces:    []FaceType{FACE_SCAR_TALK, FACE_SCAR, FACE_SCAR, FACE_SCAR_TALK, FACE_SCAR},
			dialog: []string{
				"WELL DONE.",
				"I CAN'T BELIEVE WE'VE MADE IT THIS FAR!",
				"MY BODY HAS REACHED BEYOND WHAT MOST BEINGS ARE     CAPABLE OF...",
				"IT'S KIND OF UNCOMFORTABLE",
				"DON'T WORRY. THIS IS ALL A NATURAL PART OF THE      PROCESS.",
			},
			music: "hope",
			voice: "voice",
		},
		{
			bodyType: BODY_CORRUPTED,
			faces:    []FaceType{FACE_EMPTY_TALK, FACE_EMPTY, FACE_EMPTY_SAD, FACE_EMPTY_SAD, FACE_EMPTY_SAD, FACE_EMPTY, FACE_EMPTY, FACE_EMPTY_TALK},
			dialog: []string{
				"EXCELLENT WORK!",
				"I'M SO SORRY YOU HAVE TO  SEE ME LIKE THIS.",
				"...",
				"THE VOICES SAY WHAT WE ARE DOING IS GOING TO BRING  GREAT MISERY TO THE WORLD.",
				"THEY DON'T UNDERSTAND.    THEY ARE JEALOUS. AFRAID.",
				"MY TRANSFORMATION IS      ALMOST COMPLETE!",
				"ONCE WE GET THERE, WE WILL FIX EVERYTHING.",
				"KEEP GOING!",
			},
			music: "malform",
			voice: "voice",
		},
		{
			bodyType: BODY_MELTED,
			faces:    []FaceType{FACE_MELTED},
			dialog: []string{
				"..........",
				"I CAN'T GET RID OF THESE  VOICES IN MY HEAD...",
				"THE SOULS THAT YOU'VE     GIVEN TO ME ARE           OVERWHELMING MY BODY.",
				"NGH...I AM FINALLY        BEGINNING TO UNDERSTAND   THIS WORLD",
				"THERE IS SOMETHING LURKING INSIDE OF EVERYONE AND   EVERYTHING.",
				"AH...HHHH...COULD IT BE?  COULD IT BE ME?",
				"WE HAVE TO TAKE CONTROL...",
				"...YOU'RE DOING F..FINE...",
			},
			music: "malform",
			voice: "voice",
		},
		{
			bodyType: BODY_HORROR,
			faces:    []FaceType{FACE_NONE},
			dialog: []string{
				".........",
				"...",
				"....H.....E....L.....P.....",
				"...HE'S...C...COMING...!!",
			},
			music: "",
			voice: "voice",
		},
		{
			bodyType: BODY_NONE,
			faces:    []FaceType{FACE_NONE},
			dialog: []string{
				"HAHAHAHAHAHAHA! YOU IDIOT.",
				"THE CAT WAS MY VESSEL.    SHE'S A GONER.",
				"THANKS TO YOU, I HAVE     BECOME REAL. AND NOW...",
				"I  A M  G O D .",
				"YOU CANNOT REPLACE ME.    I WILL ALWAYS COME BACK.",
				"I WILL RESET THE UNIVERSE OVER AND OVER UNTIL I WIN.",
				"NOW IT IS TIME FOR ME TO  BRING AN END TO EVERYTHING",
				"IT IS YOUR FATE TO BE     DESTROYED. EMBRACE IT!    BASK IN THE GLORY.",
			},
			music: "him",
			voice: "evil_voice",
		},
		{
			bodyType: BODY_NONE,
			faces:    []FaceType{FACE_NONE, FACE_EMPTY_SAD, FACE_EMPTY_TALK, FACE_EMPTY_SAD, FACE_EMPTY_SAD, FACE_EMPTY_SAD, FACE_EMPTY, FACE_EMPTY_TALK, FACE_EMPTY_TALK, FACE_EMPTY_SAD, FACE_MELTED, FACE_NONE, FACE_NONE},
			dialog: []string{
				"HAHAHAHAHAHA! YOU-",
				"WAIT. WHAT?",
				"I'M ALIVE! BUT WAIT...",
				"THE DEMON. WHERE IS IT!?  WHAT A MISTAKE I'VE MADE!!",
				"...",
				"I CAN'T REMAIN MUCH LONGER  THE ENERGY HAS WITHERED AWAY MY BODY.",
				"YOU'VE BEEN SUCH GREAT    HELP, BUT I NEED ONE LAST THING FROM YOU.",
				"TAKE WHAT REMAINS OF ME   AND PUT IT SOMEWHERE SAFE.",
				"WHOEVER LIVES AFTER US    WILL NEED WHAT REMAINS OF MY POWER.",
				"THE DEMON, I'M AFRAID,    WILL HAVE TO BE...D-DEALT WITH LATER.",
				"I H-HOPE THEY WILL FORGIVE ME...",
				"GOODBYE, FRIEND.",
			},
			music: "rescue",
			voice: "voice",
		},
	}
}

type CutsceneState struct {
	feles        *Object
	felesBody    BodyType
	cutscene     *Cutscene
	dialog       []*UIText
	dialogIndex  int
	nextMission  int
	skipTimer    float64
	skipText     *UIText
	transition   FadeMode
	transTimer   float64
	renderTarget *ebiten.Image
	uiRoot       *UINode
	elapsedTime float64
}

func NewCutsceneState(sceneNum int) *CutsceneState {
	state := new(CutsceneState)

	//Change ending if player beat all par times
	if sceneNum >= len(missions) {
		for _, m := range missions {
			if !m.goodEndFlag {
				goto skip
			}
		}
		sceneNum++
	skip:
	}

	state.cutscene = &cutscenes[sceneNum]
	state.feles = MakeFeles(state.cutscene.faces[0], state.cutscene.bodyType, vmath.NewVec(SCR_WIDTH_H, SCR_HEIGHT_H-24.0))
	state.felesBody = state.cutscene.bodyType

	state.uiRoot = EmptyUINode()

	dialogBox := CreateUIBox(image.Rect(112, 40, 136, 48), image.Rect(48, 176, 272, 216), true)
	state.uiRoot.AddChild(&dialogBox.UINode)

	state.dialog = make([]*UIText, len(state.cutscene.dialog))
	for i, s := range state.cutscene.dialog {
		state.dialog[i] = GenerateText(s, image.Rect(8, 8, dialogBox.Width()-8, dialogBox.Height()-8))
		state.dialog[i].fillPos = 0
		state.dialog[i].fillSound = state.cutscene.voice
		if i > 0 {
			state.dialog[i].visible = false
		}
		dialogBox.AddChild(&state.dialog[i].UINode)
	}
	state.dialogIndex = 0
	instructText := GenerateText("SPACE/CLICK: NEXT ... HOLD ENTER: SKIP", image.Rect(4, 228, SCR_WIDTH-4, SCR_HEIGHT))
	state.uiRoot.AddChild(&instructText.UINode)

	state.skipText = GenerateText("SKIPPING...", image.Rect(SCR_WIDTH_H-44, 24, SCR_WIDTH_H+44, 32))
	state.skipText.visible = false
	state.uiRoot.AddChild(&state.skipText.UINode)
	state.skipTimer = 0.0

	state.nextMission = sceneNum
	state.transition = FM_FADE_IN
	state.transTimer = 0.0
	state.renderTarget = ebiten.NewImage(SCR_WIDTH, SCR_HEIGHT)

	audio.PlayMusic(state.cutscene.music)

	return state
}

const CUTSCENE_FADE_SPEED = 1.0

func (ct *CutsceneState) Update(deltaTime float64) {
	ct.elapsedTime += deltaTime
	if ct.transition == FM_NO_FADE {
		//Dialog advancements
		dlg := ct.dialog[ct.dialogIndex]
		dlg.Update(deltaTime)
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			if dlg.fillPos >= len(dlg.text) {
				ct.dialogIndex++
				if ct.dialogIndex >= len(ct.dialog) {
					ct.dialogIndex--
					ct.transition = FM_FADE_OUT
					audio.PlayMusic("")
					ct.dialog[ct.dialogIndex].visible = false
					ct.dialog[ct.dialogIndex].parent.visible = false
				} else {
					faceIdx := int(math.Min(float64(ct.dialogIndex), float64(len(ct.cutscene.faces)-1)))
					ct.feles = MakeFeles(ct.cutscene.faces[faceIdx], ct.felesBody, ct.feles.pos)
					ct.dialog[ct.dialogIndex-1].visible = false
					ct.dialog[ct.dialogIndex].visible = true
				}
			} else {
				dlg.fillPos = len(dlg.text)
			}
		}
		//Cutscene skipping
		if ebiten.IsKeyPressed(ebiten.KeyEnter) {
			ct.skipText.visible = true
			ct.skipTimer += deltaTime
			if ct.skipTimer > 0.5 {
				ct.skipTimer = 0.0
				ct.transition = FM_FADE_OUT
				audio.PlayMusic("")
			}
		} else {
			ct.skipText.visible = false
			ct.skipTimer = 0.0
		}
	} else {
		ct.transTimer += deltaTime
		if ct.transTimer > CUTSCENE_FADE_SPEED {
			ct.transTimer = 0.0
			if ct.transition == FM_FADE_OUT {
				if ct.nextMission >= 0 && ct.nextMission < len(missions) {
					ChangeAppState(NewGame(ct.nextMission))
				} else {
					ts := new(TitleScreen)
					ts.badEnd = (ct.nextMission == len(missions))
					ts.goodEnd = (ct.nextMission == len(missions)+1)
					ChangeAppState(ts)
				}
			} else {
				ct.transition = FM_NO_FADE
			}
		}
	}
}

func (ct *CutsceneState) Draw(screen *ebiten.Image) {
	if ct.nextMission == len(missions) {
		ct.DrawEvilBackground(screen, math.Sin(ct.elapsedTime) * 80.0, math.Cos(ct.elapsedTime) * 80.0, (math.Cos(ct.elapsedTime + math.Pi / 6.0) * 0.4) + 1.25)
		ct.DrawEvilBackground(screen, math.Sin(ct.elapsedTime) * 64.0, math.Cos(ct.elapsedTime) * 64.0, (math.Sin(ct.elapsedTime + math.Pi / 3.0) * 0.25) + 1.5)
	}
	ct.feles.DrawAllSprites(screen, nil)
	
	ct.uiRoot.Draw(screen, nil)


	ct.renderTarget.Clear()
	ct.renderTarget.DrawImage(screen, nil)
	if ct.transition != FM_NO_FADE {
		op := &ebiten.DrawRectShaderOptions{}
		if ct.transition == FM_FADE_IN {
			op.Uniforms = map[string]interface{}{
				"Coverage": float32(1.0 - (ct.transTimer / CUTSCENE_FADE_SPEED)),
			}
		} else if ct.transition == FM_FADE_OUT {
			op.Uniforms = map[string]interface{}{
				"Coverage": float32(ct.transTimer / CUTSCENE_FADE_SPEED),
			}
		}

		op.Images[0] = ct.renderTarget
		op.Images[1] = noiseImg
		screen.DrawRectShader(SCR_WIDTH, SCR_HEIGHT, whiteFadeShader, op)
	}
}

func (ct *CutsceneState) DrawEvilBackground(screen *ebiten.Image, x, y, scale float64) {
	op := &ebiten.DrawImageOptions{}
	imgSize := float64(GetGraphics().Bounds().Dx())
	op.ColorM.ChangeHSV(math.Pi / 4.0, 0.1, 0.05)
	op.GeoM.Translate(-SCR_WIDTH_H, -SCR_HEIGHT_H)
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(x + SCR_WIDTH_H, y + SCR_HEIGHT_H)
	screen.DrawImage(GetGraphics(), op)
	op.GeoM.Translate(imgSize, 0.0)
	screen.DrawImage(GetGraphics(), op)
	op.GeoM.Translate(0.0, imgSize)
	screen.DrawImage(GetGraphics(), op)
	op.GeoM.Translate(-imgSize, 0.0)
	screen.DrawImage(GetGraphics(), op)
	op.GeoM.Translate(-imgSize, 0.0)
	screen.DrawImage(GetGraphics(), op)
	op.GeoM.Translate(0.0, -imgSize)
	screen.DrawImage(GetGraphics(), op)
	op.GeoM.Translate(0.0, -imgSize)
	screen.DrawImage(GetGraphics(), op)
	op.GeoM.Translate(imgSize, 0.0)
	screen.DrawImage(GetGraphics(), op)
	op.GeoM.Translate(imgSize, 0.0)
	screen.DrawImage(GetGraphics(), op)
}

func (ct *CutsceneState) Enter() {}

func (ct *CutsceneState) Leave() {
	ct.renderTarget.Dispose()
	ct.uiRoot.Unlink()
}
