package main

/*
TODO:
-Love system
-Player hurting
-Knight
-Blargh
-Gopnik
-Worm
-Barrels
-Cat
-Text rendering
-Fix audio loading issue
-Feles
-Music
*/

import (
	"container/list"
	"fmt"
	"image"
	_ "image/color"
	_ "image/png"
	"log"
	"math/rand"
	"os"

	"time"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

const (
	SCR_WIDTH    = 320
	SCR_WIDTH_H  = SCR_WIDTH / 2
	SCR_HEIGHT   = 240
	SCR_HEIGHT_H = SCR_HEIGHT / 2
)

var (
	__graphics *ebiten.Image
)

//Game ...
type Game struct {
	objects   *list.List
	level     *Level
	deltaTime float64
	lastTime  time.Time
	camPos    *Vec2f
}

var game *Game

//To mark points visually for inspection of collision detection
var debugSpot *Vec2f
var debugSprite *Sprite

//Update ...
func (g *Game) Update() error {
	now := time.Now()
	g.deltaTime = now.Sub(g.lastTime).Seconds()
	g.lastTime = now

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

	return nil
}

//Draw ...
func (g *Game) Draw(screen *ebiten.Image) {
	camMat := &ebiten.GeoM{}
	camMat.Translate(-g.camPos.x+SCR_WIDTH_H, -g.camPos.y+SCR_HEIGHT_H)

	g.level.Draw(screen, camMat)
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

	if debugSpot.x != 0.0 || debugSpot.y != 0.0 {
		o := &ebiten.DrawImageOptions{}
		o.GeoM.Concat(*camMat)
		o.GeoM.Translate(debugSpot.x, debugSpot.y)
		debugSprite.Draw(screen, &o.GeoM)
	}

	ebitenutil.DebugPrint(screen, fmt.Sprint(ebiten.CurrentFPS()))
}

//Layout ...
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return SCR_WIDTH, SCR_HEIGHT
}

//Returns the graphics page and loads it if it isn't there
func GetGraphics() *ebiten.Image {
	if __graphics == nil {
		//Load graphics
		reader, err := os.Open("assets/graphics.png")
		if err != nil {
			log.Fatal(err)
		}
		defer reader.Close()

		img, _, err := image.Decode(reader)
		if err != nil {
			log.Fatal(err)
		}
		__graphics = ebiten.NewImageFromImage(img)
	}
	return __graphics
}

func main() {
	rand.Seed(time.Now().UnixNano())
	//Init game
	game = new(Game)
	game.lastTime = time.Now()
	game.camPos = ZeroVec()

	debugSpot = ZeroVec()
	debugSprite = NewSprite(image.Rect(136, 40, 140, 44), &Vec2f{-2.0, -2.0}, false, false, 0)

	//Initialize world
	game.objects = list.New()
	game.level = GenerateLevel(64, 64)

	center := func(x int) float64 {
		return float64(x)*TILE_SIZE + 8.0
	}
	for _, sp := range game.level.spawns {
		switch sp.spawnType {
		case SP_PLAYER:
			AddPlayer(game, center(sp.ix), center(sp.iy))
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
