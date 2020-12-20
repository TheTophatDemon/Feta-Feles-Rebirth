package main

import (
	"image"
	"math"
	"math/rand"
)

const (
	TILE_SIZE   = 16.0
	TILE_SIZE_H = TILE_SIZE / 2.0
)

type Tile struct {
	tt                       TileType
	spr                      *Sprite
	gridX, gridY             int
	left, right, top, bottom float64 //Coordinates of tile boundaries in world space / pixels
	centerX, centerY         float64 //In world space/ pixels
	modified                 bool    //Is true when the tile has changed and needs its sprite regenerated
}

func (t *Tile) IsSolid() bool {
	return t.tt&TT_SOLIDS > 0
}

func (t *Tile) IsSlope() bool {
	return t.tt&TT_SLOPES > 0
}

func (t *Tile) GetSlopeNormal() *Vec2f {
	if t.IsSlope() {
		return slopeNormals[t.tt].Clone()
	}
	return nil
}

func (t *Tile) SetType(newType TileType) {
	t.tt = newType
	t.modified = true
}

func (t *Tile) RegenSprite() {
	if t.tt != TT_EMPTY {
		orient := 0

		switch t.tt {
		case TT_SLOPE_45, TT_TENTACLE_RIGHT:
			orient = 1
		case TT_SLOPE_315, TT_TENTACLE_DOWN:
			orient = 2
		case TT_SLOPE_225, TT_TENTACLE_LEFT:
			orient = 3
		}

		if t.tt == TT_RUNE {
			x := int(math.Floor(rand.Float64()*4.0)) * 16
			tileTypeRects[t.tt] = image.Rect(x, 112, x+16, 128)
		}

		t.spr = NewSprite(tileTypeRects[t.tt], &Vec2f{t.left, t.top}, false, false, orient)
	} else {
		t.spr = nil
	}
}

type TileType int

const (
	TT_EMPTY          TileType = 0
	TT_BLOCK          TileType = 1 << 0
	TT_SLOPE_45       TileType = 1 << 1 //Slope number refers to angle in degrees of normal vector relative to positive x axis
	TT_SLOPE_135      TileType = 1 << 2
	TT_SLOPE_225      TileType = 1 << 3
	TT_SLOPE_315      TileType = 1 << 4
	TT_TENTACLE_UP    TileType = 1 << 5
	TT_TENTACLE_DOWN  TileType = 1 << 6
	TT_TENTACLE_LEFT  TileType = 1 << 7
	TT_TENTACLE_RIGHT TileType = 1 << 8
	TT_RUNE           TileType = 1 << 9
	TT_PYLON          TileType = 1 << 10

	TT_SOLIDS TileType = TT_BLOCK | TT_SLOPE_45 | TT_SLOPE_135 | TT_SLOPE_225 | TT_SLOPE_315 | TT_TENTACLE_DOWN | TT_TENTACLE_LEFT | TT_TENTACLE_UP | TT_TENTACLE_RIGHT | TT_PYLON | TT_RUNE
	TT_SLOPES TileType = TT_SLOPE_45 | TT_SLOPE_135 | TT_SLOPE_225 | TT_SLOPE_315
)

var tileTypeRects map[TileType]image.Rectangle
var slopeNormals map[TileType]*Vec2f

func init() {
	tileTypeRects = map[TileType]image.Rectangle{
		TT_BLOCK:          image.Rect(16, 96, 32, 112),
		TT_SLOPE_45:       image.Rect(0, 96, 16, 112),
		TT_SLOPE_135:      image.Rect(0, 96, 16, 112),
		TT_SLOPE_225:      image.Rect(0, 96, 16, 112),
		TT_SLOPE_315:      image.Rect(0, 96, 16, 112),
		TT_TENTACLE_UP:    image.Rect(32, 96, 48, 112),
		TT_TENTACLE_DOWN:  image.Rect(32, 96, 48, 112),
		TT_TENTACLE_LEFT:  image.Rect(32, 96, 48, 112),
		TT_TENTACLE_RIGHT: image.Rect(32, 96, 48, 112),
		TT_RUNE:           image.Rect(0, 112, 16, 128),
		TT_PYLON:          image.Rect(48, 96, 64, 112),
	}

	slopeNormals = map[TileType]*Vec2f{
		TT_SLOPE_45:  &Vec2f{math.Cos(3.0 * math.Pi / 4.0), math.Sin(3.0 * math.Pi / 4.0)},
		TT_SLOPE_135: &Vec2f{math.Cos(5.0 * math.Pi / 4.0), math.Sin(5.0 * math.Pi / 4.0)},
		TT_SLOPE_225: &Vec2f{math.Cos(7.0 * math.Pi / 4.0), math.Sin(7.0 * math.Pi / 4.0)},
		TT_SLOPE_315: &Vec2f{math.Cos(9.0 * math.Pi / 4.0), math.Sin(9.0 * math.Pi / 4.0)},
	}
}
