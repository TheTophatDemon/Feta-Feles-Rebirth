package main

import (
	"container/list"
	"image/color"
	"log"
	"math"
	"math/rand"
	"runtime"
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
	"github.com/thetophatdemon/Feta-Feles-Remastered/audio"
	"github.com/thetophatdemon/Feta-Feles-Remastered/vmath"
)

type Game struct {
	objects                *list.List
	level                  *Level
	deltaTime              float64
	lastTime               time.Time
	camPos, camMin, camMax *vmath.Vec2f
	winTimer               float64
	hud                    GameHUD
	mission                *Mission
	missionNumber          int
	playerObj              *Object
	love                   int
	respawnTimer           float64
	fade                   FadeMode
	fadeTimer              float64
	fadeStage              int
	renderTarget           *ebiten.Image
	strobeTimer            float64
	strobeForward          bool
	strobeSpeed            float64
	bgColor                color.RGBA
	elapsedTime            float64
	pause                  bool
}

type FadeMode int

const (
	FM_FADE_OUT FadeMode = -1
	FM_NO_FADE  FadeMode = 0
	FM_FADE_IN  FadeMode = 1
	FADE_STAGES int      = 8
)

var __totalGameTime float64

func NewGame(mission int) *Game {
	if mission < 0 || mission >= len(missions) {
		log.Println("Invalid mission number!")
		mission = int(math.Max(0, math.Min(float64(len(missions)-1), float64(mission))))
	}
	if mission == 0 {
		__totalGameTime = 0.0
	}
	game := &Game{
		objects:       list.New(),
		lastTime:      time.Now(),
		camPos:        vmath.ZeroVec(),
		camMin:        vmath.ZeroVec(),
		camMax:        vmath.ZeroVec(),
		mission:       &missions[mission],
		missionNumber: mission,
		fade:          FM_FADE_IN,
		renderTarget:  ebiten.NewImage(SCR_WIDTH, SCR_HEIGHT),
		strobeSpeed:   6.0,
		strobeTimer:   0.0,
		strobeForward: true,
		bgColor:       missions[mission].bgColor1,
		hud:           *CreateGameHUD(),
	}
	Emit_Signal(SIGNAL_GAME_INIT, game, nil)
	game.level = GenerateLevel(missions[mission].mapWidth, missions[mission].mapHeight, mission <= 1)

	//Spawn entities
	playerSpawn := game.level.FindCenterSpawnPoint(game)
	game.playerObj = AddPlayer(game, playerSpawn.centerX, playerSpawn.centerY)

	game.CenterCameraOn(game.playerObj, true)

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
	for i := 0; i < missions[mission].maxWorms; i++ {
		spawn := game.level.FindOffscreenSpawnPoint(game)
		AddWorm(game, spawn.centerX, spawn.centerY)
	}

	for i := 0; i < missions[mission].maxBarrels; i++ {
		spawn := game.level.FindOffscreenSpawnPoint(game)
		AddBarrel(game, spawn.centerX, spawn.centerY)
	}

	audio.PlaySound("intro_chime")

	if mission == 0 {
		Listen_Signal(SIGNAL_PLAYER_MOVED, game)
		Listen_Signal(SIGNAL_PLAYER_SHOT, game)
	}
	Listen_Signal(SIGNAL_PLAYER_ASCEND, game)
	Listen_Signal(SIGNAL_CAT_RULE, game)
	Listen_Signal(SIGNAL_CAT_DIE, game)
	Listen_Signal(SIGNAL_GAME_START, game)
	Listen_Signal(SIGNAL_PLAYER_WARP, game)

	return game
}

func (g *Game) Enter() {}

func (g *Game) Leave() {
	g.hud.root.Unlink()
}

var cheatText string = ""
var debugDraw bool

func (g *Game) Update(deltaTime float64) {
	g.deltaTime = deltaTime
	if g.fade == FM_NO_FADE {
		if !g.pause {
			g.elapsedTime += deltaTime
			__totalGameTime += deltaTime

			cheatText += strings.ToLower(string(ebiten.InputChars()))
			//Cheat codes
			if strings.Contains(cheatText, "tdnepotis") {
				g.love = g.mission.loveQuota - 1
				cheatText = ""
			}
			if strings.Contains(cheatText, "tdnyaah") {
				cheatText = ""
				AddCat(g, g.playerObj.pos.X, g.playerObj.pos.Y)
			}
			if strings.Contains(cheatText, "tdnyaaaah") {
				cheatText = ""
				for i := 0; i < 32; i++ {
					AddCat(g, g.playerObj.pos.X, g.playerObj.pos.Y)
				}
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
				AddExplosion(g, g.playerObj.pos.X, g.playerObj.pos.Y)
			}
			if strings.Contains(cheatText, "tdascend") {
				cheatText = ""
				g.love = g.mission.loveQuota
				ply := g.playerObj.components[0].(*Player)
				ply.ascended = true
				Emit_Signal(SIGNAL_PLAYER_ASCEND, g.playerObj, nil)
			}
			if strings.Contains(cheatText, "tdgottam") {
				cheatText = ""
				ChangeAppState(NewGame(g.missionNumber + 1))
				return
			}
			if strings.Contains(cheatText, "tdspicy") {
				cheatText = ""
				AddWorm(g, g.playerObj.pos.X, g.playerObj.pos.Y)
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
					S_WORM
				)

				pool := make([]int, 0, 4)
				if knightCtr.count < g.mission.maxKnights {
					pool = append(pool, S_KNIGHT)
				}
				if blarghCtr.count < g.mission.maxBlarghs {
					pool = append(pool, S_BLARGH)
				}
				if gopnikCtr.count < g.mission.maxGopniks {
					pool = append(pool, S_GOPNIK)
				}
				if barrelCtr.count < g.mission.maxBarrels {
					pool = append(pool, S_BARREL)
				}
				if wormCtr.count < g.mission.maxWorms {
					pool = append(pool, S_WORM)
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
					case S_WORM:
						AddWorm(g, spawn.centerX, spawn.centerY)
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
				}
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

			//Strobe background color by incrementing the timer in a "ping pong" motion.
			if g.strobeForward {
				g.strobeTimer += g.deltaTime
				if g.strobeTimer > g.strobeSpeed {
					g.strobeTimer = g.strobeSpeed
					g.strobeForward = false
				}
			} else {
				g.strobeTimer -= g.deltaTime
				if g.strobeTimer < 0.0 {
					g.strobeTimer = 0.0
					g.strobeForward = true
				}
			}

			{
				//Linearly interpolate between the two background colors using the timer variable
				r1, g1, b1 := float64(g.mission.bgColor1.R), float64(g.mission.bgColor1.G), float64(g.mission.bgColor1.B)
				r2, g2, b2 := float64(g.mission.bgColor2.R), float64(g.mission.bgColor2.G), float64(g.mission.bgColor2.B)
				t := (g.strobeTimer / g.strobeSpeed)
				g.bgColor = color.RGBA{
					uint8(r1 + (r2-r1)*t),
					uint8(g1 + (g2-g1)*t),
					uint8(b1 + (b2-b1)*t),
					255,
				}
			}

		}
		g.hud.Update(g)
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			audio.PlaySound("menu")
			g.pause = !g.pause
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
					ChangeAppState(NewCutsceneState(g.missionNumber + 1))
					return
				} else {
					runtime.GC() //Get rid of all that level generation memory
					Emit_Signal(SIGNAL_GAME_START, g, nil)
					audio.PlayMusic(g.mission.music)
				}
				g.fade = FM_NO_FADE
			}
		}
	}

	g.CenterCameraOn(g.playerObj, false)
}

const CAM_TRACK_SPEED = 2.0

func (g *Game) CenterCameraOn(obj *Object, instant bool) {
	topLeft := vmath.NewVec(SCR_WIDTH_H, SCR_HEIGHT_H)
	bottomRight := vmath.NewVec(g.level.pixelWidth-SCR_WIDTH_H, g.level.pixelHeight-SCR_HEIGHT_H)
	targetPos := vmath.VecMax(topLeft, vmath.VecMin(bottomRight, obj.pos))
	hscr := vmath.NewVec(SCR_WIDTH_H, SCR_HEIGHT_H)

	camMove := targetPos.Clone().Sub(g.camPos)
	//Scroll slowly when moving across large distances
	if camMove.Length() < 16.0 || instant {
		g.camPos.X = targetPos.X
		g.camPos.Y = targetPos.Y
	} else {
		g.camPos.Add(camMove.Clone().Normalize().Scale(g.deltaTime * math.Max(g.level.pixelHeight, g.level.pixelWidth) * CAM_TRACK_SPEED))
	}
	g.camMin = g.camPos.Clone().Sub(hscr)
	g.camMax = g.camPos.Clone().Add(hscr)
}

func (g *Game) Draw(screen *ebiten.Image) {
	//Background
	screen.Fill(g.bgColor)

	camMat := &ebiten.GeoM{}
	camMat.Translate(math.Floor(-g.camPos.X+SCR_WIDTH_H), math.Floor(-g.camPos.Y+SCR_HEIGHT_H))

	g.level.Draw(g, screen, camMat)
	for objE := g.objects.Front(); objE != nil; objE = objE.Next() {
		obj := objE.Value.(*Object)
		if !obj.hidden && g.SquareOnScreen(obj.pos.X, obj.pos.Y, obj.radius) {
			objM := &ebiten.DrawImageOptions{}
			objM.GeoM.Concat(*camMat)
			objM.GeoM.Translate(math.Floor(obj.pos.X), math.Floor(obj.pos.Y))
			for _, spr := range obj.sprites {
				spr.Draw(screen, &objM.GeoM)
			}
		}
	}
	if g.fade == FM_NO_FADE {
		g.hud.Draw(screen)
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

func (g *Game) HandleSignal(kind Signal, src interface{}, params map[string]interface{}) {
	if g.missionNumber == 0 {
		switch kind {
		case SIGNAL_PLAYER_MOVED:
			g.hud.DisplayMessage("HOLDING CLICK/SPACE WILL SHOOT. ENTER   WILL PAUSE.", 5.0)
		case SIGNAL_PLAYER_SHOT:
			g.hud.DisplayMessage("THE MONSTERS PRODUCE FUEL FOR ASCENTION. FILL THE BAR!", 5.0)
		case SIGNAL_GAME_START:
			g.hud.DisplayMessage("MOVE WITH WASD KEYS OR ARROWS.", 4.0)
		}
	}
	switch kind {
	case SIGNAL_PLAYER_ASCEND:
		spawn := g.level.FindOffscreenSpawnPoint(g)
		AddCat(g, spawn.centerX, spawn.centerY)
		AddStarBurst(g, g.playerObj.pos.X, g.playerObj.pos.Y)
		audio.PlaySound("ascend")
		if g.missionNumber == 0 {
			g.hud.DisplayMessage("  EXCELLENT. NOW...     GO GET THE CAT!", 4.0)
		}
	case SIGNAL_CAT_RULE:
		g.hud.DisplayMessage("YOU MUST ASCEND TO  SLAY THE CAT", 4.0)
	case SIGNAL_CAT_DIE:
		g.fade = FM_FADE_OUT
		audio.PlaySound("outro_chime")
		g.mission.goodEndFlag = g.elapsedTime < float64(g.mission.parTime)
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

func (g *Game) SquareOnScreen(x, y, radius float64) bool {
	return x+radius > g.camMin.X && x-radius < g.camMax.X && y+radius > g.camMin.Y && y-radius < g.camMax.Y
}
