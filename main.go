package main

/*
TODO:
-Sprite sorting
-Level transition
-Function to spawn things offscreen
-Teleporting / Laser?
-Wrap around level?
-Blargh
-Gopnik
-Barrels
-Worm
-Loading screen?
-Fix audio loading issue
-Feles
-Music
*/

import (
	"bytes"
	"container/list"
	"fmt"
	"image"
	_ "image/color"
	_ "image/png"
	"log"
	"math/rand"
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

type Mission struct {
	loveQuota  int
	knightFreq float64
	blarghFreq float64
	gopnikFreq float64
	wormFreq   float64
	barrelFreq float64
	crossFreq  float64
}

var missions [6]Mission

func init() {
	missions[0] = Mission{
		loveQuota:  100,
		knightFreq: 0.5,
		blarghFreq: 0.5,
		gopnikFreq: 0.5,
		wormFreq:   0.5,
		barrelFreq: 0.5,
		crossFreq:  0.0,
	}
}

var cheatText string = ""

//To mark points visually for inspection of collision detection
type DebugSpot struct {
	pos *Vec2f
	spr *Sprite
}

var __debugSpots []*DebugSpot

func (g *Game) Update() error {
	now := time.Now()
	g.deltaTime = now.Sub(g.lastTime).Seconds()
	g.lastTime = now

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

	//Set camera to player position
	g.camPos = g.playerObj.pos

	//Animate UI
	if g.winTimer > 0.0 {
		g.hud.winText.Update(g.deltaTime)
		g.winTimer -= g.deltaTime
	}

	//Update UI with love amount
	lbbRect := g.hud.loveBarBorder.rect
	barRect := image.Rect(lbbRect.Min.X+3, lbbRect.Min.Y+3, lbbRect.Max.X-3, lbbRect.Max.Y-3)
	barRect.Max.X = barRect.Min.X + int(float64(barRect.Dx())*float64(g.love)/float64(g.mission.loveQuota))
	g.hud.loveBar = SpriteFromScaledImg(g.hud.loveBar.subImg, barRect, 0)

	return nil
}

//Draw ...
func (g *Game) Draw(screen *ebiten.Image) {
	camMat := &ebiten.GeoM{}
	camMat.Translate(-g.camPos.x+SCR_WIDTH_H, -g.camPos.y+SCR_HEIGHT_H)

	g.level.Draw(g, screen, camMat)
	for objE := g.objects.Front(); objE != nil; objE = objE.Next() {
		obj := objE.Value.(*Object)
		if !obj.hidden {
			objM := &ebiten.DrawImageOptions{}
			objM.GeoM.Concat(*camMat)
			objM.GeoM.Translate(obj.pos.x, obj.pos.y)
			for _, spr := range obj.sprites {
				spr.Draw(screen, &objM.GeoM)
			}
		}
	}

	g.hud.loveBarBorder.Draw(screen)
	g.hud.loveBar.Draw(screen, nil)
	if g.winTimer > 0.0 {
		g.hud.winBox.Draw(screen)
		g.hud.winText.Draw(screen)
	}

	for _, spot := range __debugSpots {
		o := &ebiten.DrawImageOptions{}
		o.GeoM.Concat(*camMat)
		o.GeoM.Translate(spot.pos.x, spot.pos.y)
		spot.spr.Draw(screen, &o.GeoM)
	}

	GenerateText(fmt.Sprintf("FPS: %.2f", ebiten.CurrentTPS()), image.Rect(SCR_WIDTH-80, 0, SCR_WIDTH, 64)).Draw(screen)
	//ebitenutil.DebugPrint(screen, fmt.Sprint(ebiten.CurrentFPS()))
}

//Adds or removes from the love counter. Returns true if the operations causes the quota to be met.
func (g *Game) AddLoveCounter(amt int) bool {
	if g.love == g.mission.loveQuota {
		return true
	}
	g.love += amt
	if g.love < 0 {
		g.love = 0
	} else if g.love >= g.mission.loveQuota {
		g.love = g.mission.loveQuota
		g.Win()
		return true
	}
	return false
}

func (g *Game) Win() {
	if g.winTimer <= 0.0 {
		g.winTimer = 8.0
		g.hud.winText.fillPos = 0
		AddCat(g)
	}
}

//Layout ...
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return SCR_WIDTH, SCR_HEIGHT
}

var __graphics *ebiten.Image

//Returns the graphics page and loads it if it isn't there
func GetGraphics() *ebiten.Image {
	if __graphics == nil {
		img, _, err := image.Decode(bytes.NewReader(assets.PNG_GRAPHICS))
		if err != nil {
			log.Fatal(err)
		}
		__graphics = ebiten.NewImageFromImage(img)
	}
	return __graphics
}

var game *Game

type Game struct {
	objects   *list.List
	level     *Level
	deltaTime float64
	lastTime  time.Time
	camPos    *Vec2f
	winTimer  float64
	hud       GameHUD
	mission   *Mission
	playerObj *Object
	love      int
}

type GameHUD struct {
	loveBarBorder UIBox
	loveBar       *Sprite
	winText       *Text
	winBox        UIBox
}

const MAX_LOVE = 100

func AddDebugSpot(x, y float64, color int) {
	var spr *Sprite
	switch color {
	case 0:
		spr = NewSprite(image.Rect(112, 40, 116, 44), &Vec2f{-3.0, -3.0}, false, false, 0)
	case 1:
		spr = NewSprite(image.Rect(136, 40, 140, 44), &Vec2f{-2.0, -2.0}, false, false, 0)
	case 2:
		spr = NewSprite(image.Rect(104, 40, 108, 44), &Vec2f{-2.0, -2.0}, false, false, 0)
	}
	__debugSpots = append(__debugSpots, &DebugSpot{&Vec2f{x, y}, spr})
}

func ClearDebugSpots() {
	__debugSpots = make([]*DebugSpot, 0, 10)
}

func main() {
	rand.Seed(time.Now().UnixNano())
	//Init game
	game = &Game{
		objects:  list.New(),
		level:    GenerateLevel(64, 64),
		lastTime: time.Now(),
		camPos:   ZeroVec(),
		hud: GameHUD{
			loveBarBorder: CreateUIBox(image.Rect(64, 40, 88, 48), image.Rect(4, 4, 4+160, 4+16)),
			loveBar:       SpriteFromScaledImg(GetGraphics().SubImage(image.Rect(104, 40, 112, 48)).(*ebiten.Image), image.Rect(4+8, 4+8, 4+160-8, 4+16-8), 0),
			winText:       GenerateText("  EXCELLENT. NOW...     GO GET THE CAT!", image.Rect(SCR_WIDTH_H-84, SCR_HEIGHT_H-56, SCR_WIDTH_H+84, SCR_HEIGHT_H-36)),
			winBox:        CreateUIBox(image.Rect(112, 40, 136, 48), image.Rect(SCR_WIDTH_H-88, SCR_HEIGHT_H-64, SCR_WIDTH_H+88, SCR_HEIGHT_H-32)),
		},
		mission: &missions[0],
	}

	__debugSpots = make([]*DebugSpot, 0, 10)

	center := func(x int) float64 {
		return float64(x)*TILE_SIZE + 8.0
	}
	for _, sp := range game.level.spawns {
		switch sp.spawnType {
		case SP_PLAYER:
			game.playerObj = AddPlayer(game, center(sp.ix), center(sp.iy))
		case SP_ENEMY:
			AddKnight(game, center(sp.ix), center(sp.iy))
		}
	}

	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Feta Feles Remake")
	//ebiten.SetRunnableOnUnfocused(true)
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
