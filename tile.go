/*
Copyright (C) 2021 Alexander Lunsford

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package main

import (
	"image"
	"math"
	"math/rand"

	"github.com/thetophatdemon/Feta-Feles-Remastered/vmath"
)

const (
	TILE_SIZE   = 16.0
	TILE_SIZE_H = TILE_SIZE / 2.0
)

type Outline int

const (
	OUTLINE_NONE   Outline = 0
	OUTLINE_TOP    Outline = 1 << 0
	OUTLINE_BOTTOM Outline = 1 << 1
	OUTLINE_LEFT   Outline = 1 << 2
	OUTLINE_RIGHT  Outline = 1 << 3
)

type Tile struct {
	tt                       TileType
	spr                      *Sprite
	outline                  Outline //Bitmask for drawing outlines on edges. Set by level generator.
	gridX, gridY             int
	left, right, top, bottom float64 //Coordinates of tile boundaries in world space / pixels
	centerX, centerY         float64 //In world space/pixels
	modified                 bool    //Is true when the tile has changed and needs its sprite regenerated
	space                    *Space  //Body of empty space the tile has been assigned to, if any
}

func (t *Tile) IsSolid() bool {
	return t.tt&TT_SOLIDS > 0
}

func (t *Tile) IsSlope() bool {
	return t.tt&TT_SLOPES > 0
}

func (t *Tile) IsTerrain() bool {
	return t.tt&TT_TERRAIN > 0
}

func (t *Tile) GetSlopeNormal() *vmath.Vec2f {
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
		rect := tileTypeRects[t.tt]

		switch t.tt {
		case TT_SLOPE_45, TT_TENTACLE_RIGHT:
			orient = 1
		case TT_SLOPE_315, TT_TENTACLE_DOWN:
			orient = 2
		case TT_SLOPE_225, TT_TENTACLE_LEFT:
			orient = 3
		case TT_RUNE: //Randomize rune sprite
			orient = rand.Intn(4)
			x := rand.Intn(4) * 16
			rect = image.Rect(rect.Min.X+x, rect.Min.Y, rect.Max.X+x, rect.Max.Y)
		case TT_BLOCK: //Alternate between different sprites
			orient = rand.Intn(4)
			x := rand.Intn(2) * 16
			rect = image.Rect(rect.Min.X+x, rect.Min.Y, rect.Max.X+x, rect.Max.Y)
		}

		t.spr = NewSprite(rect, vmath.NewVec(t.left, t.top), false, false, orient)
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

	TT_SOLIDS    TileType = TT_BLOCK | TT_SLOPES | TT_TENTACLES | TT_PYLON | TT_RUNE
	TT_SLOPES    TileType = TT_SLOPE_45 | TT_SLOPE_135 | TT_SLOPE_225 | TT_SLOPE_315
	TT_TERRAIN   TileType = TT_SLOPES | TT_BLOCK | TT_TENTACLES | TT_RUNE
	TT_TENTACLES TileType = TT_TENTACLE_UP | TT_TENTACLE_DOWN | TT_TENTACLE_LEFT | TT_TENTACLE_RIGHT
)

var tileTypeRects map[TileType]image.Rectangle
var slopeNormals map[TileType]*vmath.Vec2f

func init() {
	tileTypeRects = map[TileType]image.Rectangle{
		TT_BLOCK:          image.Rect(16, 96, 32, 112),
		TT_SLOPE_45:       image.Rect(0, 96, 16, 112),
		TT_SLOPE_135:      image.Rect(0, 96, 16, 112),
		TT_SLOPE_225:      image.Rect(0, 96, 16, 112),
		TT_SLOPE_315:      image.Rect(0, 96, 16, 112),
		TT_TENTACLE_UP:    image.Rect(64, 96, 80, 112),
		TT_TENTACLE_DOWN:  image.Rect(64, 96, 80, 112),
		TT_TENTACLE_LEFT:  image.Rect(64, 96, 80, 112),
		TT_TENTACLE_RIGHT: image.Rect(64, 96, 80, 112),
		TT_RUNE:           image.Rect(0, 112, 16, 128),
		TT_PYLON:          image.Rect(48, 96, 64, 112),
	}
	//45 & 225 are backwards
	slopeNormals = map[TileType]*vmath.Vec2f{
		TT_SLOPE_45:  vmath.NewVec(-math.Cos(3.0*math.Pi/4.0), -math.Sin(3.0*math.Pi/4.0)),
		TT_SLOPE_135: vmath.NewVec(math.Cos(5.0*math.Pi/4.0), math.Sin(5.0*math.Pi/4.0)),
		TT_SLOPE_225: vmath.NewVec(-math.Cos(7.0*math.Pi/4.0), -math.Sin(7.0*math.Pi/4.0)),
		TT_SLOPE_315: vmath.NewVec(math.Cos(9.0*math.Pi/4.0), math.Sin(9.0*math.Pi/4.0)),
	}
}
