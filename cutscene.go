package main

import (
	"image"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
)

type Cutscene struct {
	bodyType BodyType
	faces    []FaceType
	dialog   []string
	music    string
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
		},
		{
			bodyType: BODY_HUMAN,
			faces:    []FaceType{FACE_TALK, FACE_WINK, FACE_SMILE, FACE_SMILE, FACE_TALK},
			dialog: []string{
				"BOY DO YOU LOOK CONFUSED.",
				"WHAT CAN I SAY? I SEEM TO HAVE GROWN ACCUSTOMED TO  LIFE.",
				"I FEEL LIKE THIS WAS ALL  MEANT TO HAPPEN.",
				"CAN YOU FEEL THE POWER IN THE AIR? WE CAN DO ALL    SORTS OF GREAT THINGS NOW!",
				"LET'S CLIMB EVEN HIGHER!",
			},
			music: "maiden",
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
			music: "maiden",
		},
		{
			bodyType: BODY_CORRUPTED,
			faces:    []FaceType{FACE_EMPTY_TALK, FACE_EMPTY, FACE_EMPTY_SAD, FACE_EMPTY_SAD, FACE_EMPTY_SAD, FACE_EMPTY, FACE_EMPTY, FACE_EMPTY_TALK},
			dialog: []string{
				"EXCELLENT WORK!",
				"I'M SO SORRY YOU HAVE TO  SEE ME LIKE THIS.",
				"...",
				"I AM TOLD THAT WHAT WE ARE DOING IS GOING TO BRING  GREAT MISERY TO THE WORLD.",
				"THEY DON'T UNDERSTAND.    THEY ARE JEALOUS. AFRAID.",
				"MY TRANSFORMATION IS      ALMOST COMPLETE!",
				"ONCE WE GET THERE, WE WILL FIX EVERYTHING.",
				"KEEP GOING!",
			},
			music: "malform",
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
			music: "monster",
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
			music: "monster",
		},
	}
}

type CutsceneState struct {
	feles        *Object
	felesBody    BodyType
	cutscene     *Cutscene
	dialogBox    UIBox
	dialog       []*Text
	dialogIndex  int
	nextMission  int
	instructText *Text
	skipTimer    float64
	skipText     *Text
	transition   FadeMode
	transTimer   float64
	renderTarget *ebiten.Image
}

func NewCutsceneState(sceneNum int) *CutsceneState {
	ctscn := new(CutsceneState)
	ctscn.cutscene = &cutscenes[sceneNum]
	bodies := [8]BodyType{
		BODY_NONE, BODY_CAT, BODY_HUMAN, BODY_ANGEL2, BODY_CORRUPTED, BODY_MELTED, BODY_HORROR,
	}
	ctscn.feles = MakeFeles(ctscn.cutscene.faces[0], bodies[sceneNum], &Vec2f{SCR_WIDTH_H, SCR_HEIGHT_H - 24.0})
	ctscn.felesBody = bodies[sceneNum]
	ctscn.dialog = make([]*Text, len(ctscn.cutscene.dialog))
	for i, s := range ctscn.cutscene.dialog {
		ctscn.dialog[i] = GenerateText(s, image.Rect(56, 184, 264, 208))
		ctscn.dialog[i].fillPos = 0
	}
	ctscn.dialogBox = CreateUIBox(image.Rect(112, 40, 136, 48), image.Rect(48, 176, 272, 216))
	ctscn.dialogIndex = 0
	ctscn.nextMission = sceneNum
	ctscn.instructText = GenerateText("SPACE/CLICK: NEXT ... HOLD ENTER: SKIP", image.Rect(4, 228, SCR_WIDTH-4, SCR_HEIGHT))
	ctscn.skipTimer = 0.0
	ctscn.skipText = GenerateText("SKIPPING...", image.Rect(SCR_WIDTH_H-44, 24, SCR_WIDTH_H+44, 32))
	ctscn.transition = FM_FADE_IN
	ctscn.transTimer = 0.0
	ctscn.renderTarget = ebiten.NewImage(SCR_WIDTH, SCR_HEIGHT)
	return ctscn
}

const CUTSCENE_FADE_SPEED = 1.0

func (ct *CutsceneState) Update(deltaTime float64) {
	if ct.transition == FM_NO_FADE {
		//Dialog advancement
		dlg := ct.dialog[ct.dialogIndex]
		dlg.Update(deltaTime)
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			if dlg.fillPos >= len(dlg.text) {
				ct.dialogIndex++
				if ct.dialogIndex >= len(ct.dialog) {
					ct.dialogIndex--
					ct.transition = FM_FADE_OUT
				} else {
					faceIdx := int(math.Min(float64(ct.dialogIndex), float64(len(ct.cutscene.faces)-1)))
					ct.feles = MakeFeles(ct.cutscene.faces[faceIdx], ct.felesBody, ct.feles.pos)
				}
			} else {
				dlg.fillPos = len(dlg.text)
			}
		}
		//Cutscene skipping
		if ebiten.IsKeyPressed(ebiten.KeyEnter) {
			ct.skipText.fillPos = len(ct.skipText.text)
			ct.skipTimer += deltaTime
			if ct.skipTimer > 0.5 {
				ct.skipTimer = 0.0
				ct.transition = FM_FADE_OUT
			}
		} else {
			ct.skipText.fillPos = 0
		}
	} else {
		ct.transTimer += deltaTime
		if ct.transTimer > CUTSCENE_FADE_SPEED {
			ct.transTimer = 0.0
			if ct.transition == FM_FADE_OUT {
				ChangeAppState(NewGame(ct.nextMission))
			} else {
				ct.transition = FM_NO_FADE
			}
		}
	}
}

func (ct *CutsceneState) Draw(screen *ebiten.Image) {
	ct.dialogBox.Draw(screen, nil)
	ct.instructText.Draw(screen)
	ct.feles.DrawAllSprites(screen, nil)
	if ct.skipTimer > 0.0 {
		ct.skipText.Draw(screen)
	}
	ct.dialog[ct.dialogIndex].Draw(screen)
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

func (ct *CutsceneState) Enter() {}

func (ct *CutsceneState) Leave() {
	ct.renderTarget.Dispose()
}

var whiteFadeShader *ebiten.Shader
var noiseImg *ebiten.Image

func init() {
	noiseImg = ebiten.NewImage(SCR_WIDTH, SCR_HEIGHT)
	noisePixels := make([]byte, SCR_WIDTH*SCR_HEIGHT*4)
	for i := 0; i < SCR_HEIGHT*SCR_WIDTH; i++ {
		noisePixels[i*4+0] = byte(rand.Intn(255))
		noisePixels[i*4+1] = byte(rand.Intn(255))
		noisePixels[i*4+2] = byte(rand.Intn(255))
		noisePixels[i*4+3] = 255
	}
	noiseImg.ReplacePixels(noisePixels)

	var err error
	whiteFadeShader, err = ebiten.NewShader([]byte(`
		package main

		var Coverage float

		func Fragment(position vec4, texCoord vec2, color vec4) vec4 {
			diffuse := imageSrc0UnsafeAt(texCoord)
			noise := imageSrc1UnsafeAt(texCoord)
			mask := step(1.0 - Coverage, noise.r)
			return min(diffuse + mask, vec4(0.9, 0.9, 0.9, 1.0))
		}
	`))
	if err != nil {
		println(err)
		panic(err)
	}
}
