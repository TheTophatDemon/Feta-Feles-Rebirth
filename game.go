package main

import (
	"container/list"
	"fmt"
	"image"
	"log"
	"math"
	"math/rand"
	"runtime"
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten"
)

type Game struct {
	objects                *list.List
	level                  *Level
	deltaTime              float64
	lastTime               time.Time
	camPos, camMin, camMax *Vec2f
	winTimer               float64
	hud                    GameHUD
	mission                Mission
	missionNumber          int
	playerObj              *Object
	love                   int
	respawnTimer           float64
	fade                   FadeMode
	fadeTimer              float64
	fadeStage              int
	renderTarget           *ebiten.Image
}

type GameHUD struct {
	loveBarBorder UIBox
	loveBar       *Sprite
	winText       *Text
	winBox        UIBox
	winTextTimer  float64
	mapBorder     UIBox
}

type FadeMode int

const (
	FM_FADE_OUT FadeMode = -1
	FM_NO_FADE  FadeMode = 0
	FM_FADE_IN  FadeMode = 1
	FADE_STAGES int      = 8
)

func NewGame(mission int) *Game {
	if mission < 0 || mission >= len(missions) {
		log.Println("Invalid mission number!")
		mission = int(math.Max(0, math.Min(float64(len(missions)-1), float64(mission))))
	}
	game := &Game{
		objects:  list.New(),
		level:    GenerateLevel(missions[mission].mapWidth, missions[mission].mapHeight),
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
		mission:       missions[mission],
		missionNumber: mission,
		fade:          FM_FADE_IN,
		renderTarget:  ebiten.NewImage(SCR_WIDTH, SCR_HEIGHT),
	}

	__debugSpots = make([]*DebugSpot, 0, 10)

	//Spawn entities
	playerSpawn := game.level.FindSpawnPoint()
	game.playerObj = AddPlayer(game, playerSpawn.centerX, playerSpawn.centerY)
	game.CenterCameraOn(game.playerObj) //Neccessary for FindOffscreenSpawnPoint

	for i := 0; i < missions[mission].maxKnights; i++ {
		spawn := game.level.FindOffscreenSpawnPoint(game)
		AddKnight(game, spawn.centerX, spawn.centerY)
	}
	for i := 0; i < missions[mission].maxBlarghs; i++ {
		spawn := game.level.FindOffscreenSpawnPoint(game)
		AddBlargh(game, spawn.centerX, spawn.centerY)
	}
	for i := 0; i < missions[mission].maxGopniks; i++ {
		spawn := game.level.FindOffscreenSpawnPoint(game)
		AddGopnik(game, spawn.centerX, spawn.centerY)
	}

	for i := 0; i < missions[mission].maxBarrels; i++ {
		spawn := game.level.FindOffscreenSpawnPoint(game)
		AddBarrel(game, spawn.centerX, spawn.centerY)
	}

	PlaySound("intro_chime")

	return game
}

func (g *Game) Enter() {}
func (g *Game) Leave() {}

var cheatText string = ""
var debugDraw bool

func (g *Game) Update(deltaTime float64) {
	g.deltaTime = deltaTime
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
			ChangeAppState(NewGame(g.missionNumber))
			return
		}
		if strings.Contains(cheatText, "tdcruoris") {
			cheatText = ""
			debugDraw = !debugDraw
		}
		if strings.Contains(cheatText, "tdasplode") {
			cheatText = ""
			AddExplosion(g, g.playerObj.pos.x, g.playerObj.pos.y)
		}

		//Prevent the game from going AWOL when the window is moved
		if g.deltaTime > 0.25 {
			return
		}

		//Respawn monsters/barrels offscreen to maintain gameplay intensity
		g.respawnTimer += g.deltaTime
		if g.respawnTimer > 4.0 {
			g.respawnTimer = 0.0

			const (
				S_KNIGHT = iota
				S_BLARGH
				S_GOPNIK
				S_BARREL
			)

			pool := make([]int, 0, 4)
			if g.mission.knightCount < g.mission.maxKnights {
				pool = append(pool, S_KNIGHT)
			}
			if g.mission.blarghCount < g.mission.maxBlarghs {
				pool = append(pool, S_BLARGH)
			}
			if g.mission.gopnikCount < g.mission.maxGopniks {
				pool = append(pool, S_GOPNIK)
			}
			if g.mission.barrelCount < g.mission.maxBarrels {
				pool = append(pool, S_BARREL)
			}

			if len(pool) > 0 {
				spawn := g.level.FindOffscreenSpawnPoint(g)
				c := pool[rand.Intn(len(pool))]
				switch c {
				case S_KNIGHT:
					AddKnight(g, spawn.centerX, spawn.centerY)
				case S_BLARGH:
					AddBlargh(g, spawn.centerX, spawn.centerY)
				case S_GOPNIK:
					AddGopnik(g, spawn.centerX, spawn.centerY)
				case S_BARREL:
					AddBarrel(g, spawn.centerX, spawn.centerY)
				}
			}
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
				//Update relevant mission counters
				switch obj.components[0].(type) {
				case *Knight:
					g.mission.knightCount--
				case *Blargh:
					g.mission.blarghCount--
				case *Gopnik:
					g.mission.gopnikCount--
				case *Barrel:
					g.mission.barrelCount--
				}
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
					ChangeAppState(NewGame(g.missionNumber + 1))
					return
				} else {
					runtime.GC() //Get rid of all that level generation memory
				}
				g.fade = FM_NO_FADE
			}
		}
	}

	//Center camera on player
	g.CenterCameraOn(g.playerObj)

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
}

func (g *Game) CenterCameraOn(obj *Object) {
	g.camPos = VecMax(&Vec2f{SCR_WIDTH_H, SCR_HEIGHT_H}, VecMin(&Vec2f{g.level.pixelWidth - SCR_WIDTH_H, g.level.pixelHeight - SCR_HEIGHT_H}, obj.pos))
	hscr := &Vec2f{SCR_WIDTH_H, SCR_HEIGHT_H}
	g.camMin = g.camPos.Clone().Sub(hscr)
	g.camMax = g.camPos.Clone().Add(hscr)
}

func (g *Game) Draw(screen *ebiten.Image) {
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
	//Update relevant mission counters
	switch newObj.components[0].(type) {
	case *Knight:
		g.mission.knightCount++
	case *Blargh:
		g.mission.blarghCount++
	case *Gopnik:
		g.mission.gopnikCount++
	case *Barrel:
		g.mission.barrelCount++
	}
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

//Adds to the love counter. Returns true if the operations causes the quota to be met.
func (g *Game) IncLoveCounter(amt int) bool {
	if g.love == g.mission.loveQuota {
		return true
	}
	if amt < 0 {
		log.Println("Invalid amt. Use DecLoveCounter instead?")
		return false
	}
	g.love += amt
	if g.love >= g.mission.loveQuota {
		g.love = g.mission.loveQuota
		if g.hud.winTextTimer <= 0.0 {
			if (g.playerObj.components[0].(*Player)).ascended == false {
				//Quota has been met. Trigger the endgame sequence
				g.hud.winTextTimer = 8.0
				g.hud.winText.fillPos = 0
				AddCat(g)
			}
		}
		return true
	}
	return false
}

//Subtracts from the love counter. Returns true if the operation causes to counter to hit zero.
func (g *Game) DecLoveCounter(amt int) bool {
	if g.love == 0 {
		return true
	}
	if amt < 0 {
		log.Println("Invalid amt. Use IncLoveCounter instead?")
		return false
	}
	g.love -= amt
	if g.love <= 0 {
		g.love = 0
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
