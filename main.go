package main

import (
	"container/list"
	"image"
	_ "image/color"
	_ "image/png"
	"log"
	"os"

	"time"

	"github.com/hajimehoshi/ebiten"
)

const (
	SCR_WIDTH    = 320
	SCR_WIDTH_H  = SCR_WIDTH / 2
	SCR_HEIGHT   = 240
	SCR_HEIGHT_H = SCR_HEIGHT / 2
)

var (
	graphics  *ebiten.Image
	startTime = time.Now()
)

//Game ...
type Game struct {
	objects     *list.List
	level       *Level
	deltaTime   float64
	elapsedTime float64
	camPos      *Vec2f
}

//Update ...
func (g *Game) Update() error {
	gt := time.Since(startTime).Seconds()
	g.deltaTime = gt - g.elapsedTime
	g.elapsedTime = gt

	for objE := g.objects.Front(); objE != nil; objE = objE.Next() {
		obj := objE.Value.(*Object)
		for _, c := range obj.components {
			c.Update(g, obj)
		}
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
		objM := &ebiten.DrawImageOptions{}
		objM.GeoM.Concat(*camMat)
		objM.GeoM.Translate(obj.pos.x, obj.pos.y)
		for _, spr := range obj.sprites {
			spr.Draw(screen, &objM.GeoM)
		}
	}
}

//Layout ...
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return SCR_WIDTH, SCR_HEIGHT
}

func main() {

	//Init game
	game := new(Game)

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
	graphics = ebiten.NewImageFromImage(img)

	//Initialize world
	game.objects = list.New()

	MakePlayer(game, 32.0, 32.0)

	game.level = GenerateLevel(32, 32)

	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Feta Feles Remake")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
