package main

import (
	"math/rand"

	"github.com/hajimehoshi/ebiten"
)

const (
	TILE_SIZE = 16.0
)

type Level struct {
	tiles          [][]*Sprite
	rows, cols     int
	pixelW, pixelH int
}

func GenerateLevel(w, h int) *Level {

	tiles := make([][]*Sprite, h)
	for j := 0; j < h; j++ {
		tiles[j] = make([]*Sprite, w)
		for i := 0; i < w; i++ {
			tiles[j][i] = &Sprite{
				src: Rect{
					float64(rand.Intn(4) * TILE_SIZE), float64(rand.Intn(2)*TILE_SIZE) + 96.0, TILE_SIZE, TILE_SIZE,
				},
				ofs:    Vec2f{0.0, 0.0},
				flipH:  rand.Float32() > 0.5,
				flipV:  rand.Float32() > 0.5,
				orient: rand.Intn(4),
			}
		}
	}
	return &Level{tiles: tiles, rows: h, cols: w, pixelW: w * 16, pixelH: h * 16}
}

func (lev *Level) Draw(screen *ebiten.Image, pt *ebiten.GeoM) {
	op := &ebiten.DrawImageOptions{}
	for j := 0; j < lev.rows; j++ {
		op.GeoM.Concat(*pt)
		op.GeoM.Translate(0.0, float64(j)*TILE_SIZE)
		for i := 0; i < lev.cols; i++ {
			if lev.tiles[j][i] != nil {
				lev.tiles[j][i].Draw(screen, &op.GeoM)
			}
			op.GeoM.Translate(TILE_SIZE, 0.0)
		}
		op.GeoM.Reset()
	}
}
