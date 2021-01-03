package main

/*
TODO:
-Enemies getting spawned in walls?
-Blargh
-Gopnik
-Barrels
-Mob Director
-Powerups in little caverns?
-Worm
-Loading screen?
-Mission 0: Show basics of game
-Feles
-Music
*/

/*s
Bugs / Oddities
-Disappearing enemies glitch...?
-Jittering at beginning
*/

import (
	"bytes"
	"container/list"
	"fmt"
	"image"
	_ "image/color"
	_ "image/png"
	"log"
	"math"
	"math/rand"
	"runtime"
	"strings"

	"time"

	"github.com/hajimehoshi/ebiten"
	"github.com/thetophatdemon/Feta-Feles-Remastered/assets"
)

const (
	SCR_WIDTH    = 320
	SCR_WIDTH_H  = SCR_WIDTH / 2
	SCR_HEIGHT   = 240
	SCR_HEIGHT_H = SCR_HEIGHT / 2
)

type App int //Dummy type implementing interface for Ebiten to plug into

func (a *App) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return SCR_WIDTH, SCR_HEIGHT
}

var cheatText string = ""
var debugDraw bool

func (a *App) Update() error {
	g := game
	now := time.Now()
	g.deltaTime = now.Sub(g.lastTime).Seconds()
	g.lastTime = now

	if g.fade == FM_NO_FADE {

		cheatText += strings.ToLower(string(ebiten.InputChars()))

		//Cheat codes
		if strings.Contains(cheatText, "tdnepotis") {
			g.love = g.mission.loveQuota - 1
			cheatText = ""
		}
		if strings.Contains(cheatText, "tdnyaah") {
			cheatText = ""
			_, catObj := AddCat(g)
			catObj.pos.x = g.playerObj.pos.x
			catObj.pos.y = g.playerObj.pos.y
		}
		if strings.Contains(cheatText, "tdsanic") {
			cheatText = ""
			ply := g.playerObj.components[0].(*Player)
			if ply.maxSpeed <= 120.0 {
				ply.maxSpeed = 400.0
			} else {
				ply.maxSpeed = 120.0
			}
		}
		if strings.Contains(cheatText, "tdnovymir") {
			cheatText = ""
			NewGame(0)
			return nil
		}
		if strings.Contains(cheatText, "tdcruoris") {
			cheatText = ""
			debugDraw = !debugDraw
		}
		if strings.Contains(cheatText, "tdasplode") {
			cheatText = ""
			AddExplosion(game, game.playerObj.pos.x, game.playerObj.pos.y)
		}

		//Prevent the game from going AWOL when the window is moved
		if g.deltaTime > 0.25 {
			return nil
		}

		//Update objects
		toRemove := make([]*list.Element, 0, 4)
		for objE := g.objects.Front(); objE != nil; objE = objE.Next() {
			obj := objE.Value.(*Object)
			//Update components
			for _, c := range obj.components {
				if c != nil {
					c.Update(g, obj)
				}
			}
			//Objects are removed later so that they doesn't interfere with collision events
			if obj.removeMe {
				toRemove = append(toRemove, objE)
			}
			//Wrap objects around the map if they exit its boundaries
			obj.pos.x, obj.pos.y = g.level.WrapPixelCoords(obj.pos.x, obj.pos.y)
		}
		//Resolve inter-object collisions
		for objE := g.objects.Front(); objE != nil; objE = objE.Next() {
			obj := objE.Value.(*Object)
			if obj.colType != CT_NONE {
				//O(n^2)...Bleh!
				for obj2E := g.objects.Front(); obj2E != nil; obj2E = obj2E.Next() {
					obj2 := obj2E.Value.(*Object)
					if obj2.colType != CT_NONE && obj2 != obj {
						if obj.Intersects(obj2) {
							for _, c := range obj.components {
								col, ok := c.(Collidable)
								if ok {
									//An equivalent event will be sent for the other object when it is evaluated in the outer loop
									col.OnCollision(g, obj, obj2)
								}
							}
						}
					}
				}
			}
		}
		//Remove objects flagged for removal
		for _, objE := range toRemove {
			g.objects.Remove(objE)
		}

	} else { //Handle level transition FX
		g.fadeTimer += g.deltaTime
		if g.fadeTimer > 0.25 {
			g.fadeTimer = 0.0
			g.fadeStage++
			g.renderTarget.Clear()
			if g.fadeStage >= FADE_STAGES {
				g.fadeStage = 0
				//If the level is ending, start a new game
				if g.fade == FM_FADE_OUT {
					NewGame(0)
					return nil
				}
				g.fade = FM_NO_FADE
			}
		}
	}

	//Set camera to player position
	g.camPos = VecMax(&Vec2f{SCR_WIDTH_H, SCR_HEIGHT_H}, VecMin(&Vec2f{g.level.pixelWidth - SCR_WIDTH_H, g.level.pixelHeight - SCR_HEIGHT_H}, g.playerObj.pos))
	hscr := &Vec2f{SCR_WIDTH_H, SCR_HEIGHT_H}
	g.camMin = g.camPos.Clone().Sub(hscr)
	g.camMax = g.camPos.Clone().Add(hscr)

	//Animate UI
	if g.hud.winTextTimer > 0.0 {
		g.hud.winText.Update(g.deltaTime)
		g.hud.winTextTimer -= g.deltaTime
	}

	//Update UI with love amount
	lbbRect := g.hud.loveBarBorder.rect
	barRect := image.Rect(lbbRect.Min.X+3, lbbRect.Min.Y+3, lbbRect.Max.X-3, lbbRect.Max.Y-3)
	barRect.Max.X = barRect.Min.X + int(float64(barRect.Dx())*float64(g.love)/float64(g.mission.loveQuota))
	g.hud.loveBar = SpriteFromScaledImg(g.hud.loveBar.subImg, barRect, 0)

	return nil
}

//Draw ...
func (a *App) Draw(screen *ebiten.Image) {
	g := game
	camMat := &ebiten.GeoM{}
	camMat.Translate(-g.camPos.x+SCR_WIDTH_H, -g.camPos.y+SCR_HEIGHT_H)

	g.level.Draw(g, screen, camMat)
	for objE := g.objects.Front(); objE != nil; objE = objE.Next() {
		obj := objE.Value.(*Object)
		if !obj.hidden && g.SquareOnScreen(obj.pos.x, obj.pos.y, obj.radius) {
			objM := &ebiten.DrawImageOptions{}
			objM.GeoM.Concat(*camMat)
			objM.GeoM.Translate(obj.pos.x, obj.pos.y)
			for _, spr := range obj.sprites {
				spr.Draw(screen, &objM.GeoM)
			}
		}
	}

	g.hud.loveBarBorder.Draw(screen, nil)
	g.hud.loveBar.Draw(screen, nil)
	if g.hud.winTextTimer > 0.0 {
		g.hud.winBox.Draw(screen, nil)
		g.hud.winText.Draw(screen)
	}

	for _, spot := range __debugSpots {
		o := &ebiten.DrawImageOptions{}
		o.GeoM.Concat(*camMat)
		o.GeoM.Translate(spot.pos.x, spot.pos.y)
		spot.spr.Draw(screen, &o.GeoM)
	}

	if debugDraw {
		GenerateText(fmt.Sprintf("FPS: %.2f", ebiten.CurrentTPS()), image.Rect(SCR_WIDTH-80, 0, SCR_WIDTH, 64)).Draw(screen)
	}

	if g.fade != FM_NO_FADE {
		op := &ebiten.DrawImageOptions{}
		var stage float64
		if g.fade == FM_FADE_IN {
			stage = float64(FADE_STAGES - g.fadeStage)
		} else if g.fade == FM_FADE_OUT {
			stage = float64(1 + g.fadeStage)
		}
		op.GeoM.Scale(1.0/stage, 1.0/stage)
		g.renderTarget.DrawImage(screen, op)
		op.GeoM.Reset()
		op.GeoM.Scale(stage, stage)
		screen.Clear()
		screen.DrawImage(g.renderTarget, op)
	}

}

//Adds the object to the game, sorted by its draw priority, and returns the object
func (g *Game) AddObject(newObj *Object) *Object {
	for e := g.objects.Front(); e != nil; e = e.Next() {
		obj := e.Value.(*Object)
		if obj.drawPriority > newObj.drawPriority {
			g.objects.InsertBefore(newObj, e)
			return newObj
		}
	}
	g.objects.PushBack(newObj)
	return newObj
}

//Adds or removes from the love counter. Returns true if the operations causes the quota to be met.
func (g *Game) IncLoveCounter(amt int) bool {
	if g.love == g.mission.loveQuota {
		return true
	}
	g.love += amt
	if g.love < 0 {
		g.love = 0
	} else if g.love >= g.mission.loveQuota {
		g.love = g.mission.loveQuota
		if g.hud.winTextTimer <= 0.0 {
			//Quota has been met. Trigger the endgame sequence
			g.hud.winTextTimer = 8.0
			g.hud.winText.fillPos = 0
			AddCat(g)
		}
		return true
	}
	return false
}

func (g *Game) BeginEndTransition() {
	g.fade = FM_FADE_OUT
	PlaySound("outro_chime")
}

func (g *Game) SquareOnScreen(x, y, radius float64) bool {
	return x+radius > g.camMin.x && x-radius < g.camMax.x && y+radius > g.camMin.y && y-radius < g.camMax.y
}

var __graphics *ebiten.Image

//Returns the graphics page and loads it if it isn't there
func GetGraphics() *ebiten.Image {
	if __graphics == nil {
		img, _, err := image.Decode(bytes.NewReader(assets.Parse(assets.PNG_GRAPHICS)))
		if err != nil {
			log.Fatal(err)
		}
		__graphics = ebiten.NewImageFromImage(img)
	}
	return __graphics
}

type FadeMode int

const (
	FM_FADE_OUT FadeMode = -1
	FM_NO_FADE  FadeMode = 0
	FM_FADE_IN  FadeMode = 1
	FADE_STAGES int      = 8
)

var game *Game

type Game struct {
	objects      *list.List
	level        *Level
	deltaTime    float64
	lastTime     time.Time
	camPos       *Vec2f
	camMin       *Vec2f
	camMax       *Vec2f
	winTimer     float64
	hud          GameHUD
	mission      *Mission
	playerObj    *Object
	love         int
	fade         FadeMode
	fadeTimer    float64
	fadeStage    int
	renderTarget *ebiten.Image
}

type GameHUD struct {
	loveBarBorder UIBox
	loveBar       *Sprite
	winText       *Text
	winBox        UIBox
	winTextTimer  float64
	mapBorder     UIBox
}

func NewGame(mission int) {
	if mission < 0 || mission >= len(missions) {
		log.Println("Invalid mission number!")
		mission = int(math.Max(0, math.Min(float64(len(missions)-1), float64(mission))))
	}
	game = &Game{
		objects:  list.New(),
		level:    GenerateLevel(64, 64),
		lastTime: time.Now(),
		camPos:   ZeroVec(),
		camMin:   ZeroVec(),
		camMax:   ZeroVec(),
		hud: GameHUD{
			loveBarBorder: CreateUIBox(image.Rect(64, 40, 88, 48), image.Rect(4, 4, 4+160, 4+16)),
			loveBar:       SpriteFromScaledImg(GetGraphics().SubImage(image.Rect(104, 40, 112, 48)).(*ebiten.Image), image.Rect(4+8, 4+8, 4+160-8, 4+16-8), 0),
			winText:       GenerateText("  EXCELLENT. NOW...     GO GET THE CAT!", image.Rect(SCR_WIDTH_H-84, SCR_HEIGHT_H-56, SCR_WIDTH_H+84, SCR_HEIGHT_H-36)),
			winBox:        CreateUIBox(image.Rect(112, 40, 136, 48), image.Rect(SCR_WIDTH_H-88, SCR_HEIGHT_H-64, SCR_WIDTH_H+88, SCR_HEIGHT_H-32)),
			mapBorder:     CreateUIBox(image.Rect(64, 40, 80, 48), image.Rect(-1, -1, 64*int(TILE_SIZE)+1, 64*int(TILE_SIZE)+1)),
		},
		mission:      &missions[mission],
		fade:         FM_FADE_IN,
		renderTarget: ebiten.NewImage(SCR_WIDTH, SCR_HEIGHT),
	}

	__debugSpots = make([]*DebugSpot, 0, 10)

	//Spawn entities
	playerSpawn := game.level.FindSpawnPoint()
	game.playerObj = AddPlayer(game, playerSpawn.centerX, playerSpawn.centerY)

	//TODO: Replace this with a dynamic enemy spawn director
	const (
		ENM_KNIGHT = iota
		ENM_BLARGH
		ENM_MAX
	)
	for i := 0; i < 30; i++ {
		spawn := game.level.FindSpawnPoint()
		switch rand.Intn(ENM_MAX) {
		case ENM_KNIGHT:
			AddKnight(game, spawn.centerX, spawn.centerY)
		case ENM_BLARGH:
			AddBlargh(game, spawn.centerX, spawn.centerY)
		}
	}

	PlaySound("intro_chime")

	runtime.GC()
}

func main() {
	rand.Seed(time.Now().UnixNano())

	NewGame(0)

	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Feta Feles Remake")
	//ebiten.SetRunnableOnUnfocused(true)
	if err := ebiten.RunGame(new(App)); err != nil {
		log.Fatal(err)
	}
}
