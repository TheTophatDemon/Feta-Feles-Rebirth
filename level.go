package main

import (
	"container/list"
	"fmt"
	"math/rand"

	"github.com/hajimehoshi/ebiten"
)

const (
	TILE_SIZE = 16.0
)

type Level struct {
	tiles          [][]*Tile
	rows, cols     int
	pixelW, pixelH int
}

type TileType int

const (
	TT_EMPTY TileType = iota
	TT_BLOCK
	TT_SLOPE_45 //Number refers to angle in degrees of normal vector relative to positive x axis
	TT_SLOPE_135
	TT_SLOPE_225
	TT_SLOPE_315
	TT_TENTACLE_UP
	TT_TENTACLE_DOWN
	TT_TENTACLE_LEFT
	TT_TENTACLE_RIGHT
	TT_RUNE
	TT_DECOR
	TT_MAX
)

//Defines, for a single tile type, which tile types are allowed to be placed next to it
type TConstraints struct {
	left   []TileType
	right  []TileType
	top    []TileType
	bottom []TileType
}

var tileConstraints map[TileType]TConstraints
var tileTypeRects map[TileType]Rect
var allTileTypes []TileType

func init() {
	allTileTypes = []TileType{
		TT_EMPTY,
		TT_BLOCK,
		TT_SLOPE_45,
		TT_SLOPE_135,
		TT_SLOPE_225,
		TT_SLOPE_315,
		TT_TENTACLE_UP,
		TT_TENTACLE_DOWN,
		TT_TENTACLE_LEFT,
		TT_TENTACLE_RIGHT,
		TT_RUNE,
		TT_DECOR,
	}
	tileTypeRects = map[TileType]Rect{
		TT_BLOCK:          {x: 16, y: 96, w: 16, h: 16},
		TT_SLOPE_45:       {x: 0, y: 96, w: 16, h: 16},
		TT_SLOPE_135:      {x: 0, y: 96, w: 16, h: 16},
		TT_SLOPE_225:      {x: 0, y: 96, w: 16, h: 16},
		TT_SLOPE_315:      {x: 0, y: 96, w: 16, h: 16},
		TT_TENTACLE_UP:    {x: 32, y: 96, w: 16, h: 16},
		TT_TENTACLE_DOWN:  {x: 32, y: 96, w: 16, h: 16},
		TT_TENTACLE_LEFT:  {x: 32, y: 96, w: 16, h: 16},
		TT_TENTACLE_RIGHT: {x: 32, y: 96, w: 16, h: 16},
		TT_RUNE:           {x: 0, y: 112, w: 16, h: 16},
		TT_DECOR:          {x: 48, y: 96, w: 16, h: 16},
	}
	tileConstraints = map[TileType]TConstraints{
		TT_EMPTY: {
			left:   []TileType{TT_EMPTY, TT_BLOCK, TT_TENTACLE_RIGHT, TT_SLOPE_45, TT_SLOPE_315, TT_SLOPE_225, TT_SLOPE_135, TT_RUNE, TT_DECOR},
			right:  []TileType{TT_EMPTY, TT_BLOCK, TT_TENTACLE_LEFT, TT_SLOPE_45, TT_SLOPE_315, TT_SLOPE_225, TT_SLOPE_135, TT_RUNE, TT_DECOR},
			top:    []TileType{TT_EMPTY, TT_BLOCK, TT_TENTACLE_DOWN, TT_SLOPE_45, TT_SLOPE_315, TT_SLOPE_225, TT_SLOPE_135, TT_RUNE, TT_DECOR},
			bottom: []TileType{TT_EMPTY, TT_BLOCK, TT_TENTACLE_UP, TT_SLOPE_45, TT_SLOPE_315, TT_SLOPE_225, TT_SLOPE_135, TT_RUNE, TT_DECOR},
		},
		TT_BLOCK: {
			left:   []TileType{TT_EMPTY, TT_BLOCK, TT_TENTACLE_LEFT, TT_DECOR, TT_RUNE, TT_SLOPE_135, TT_SLOPE_225},
			right:  []TileType{TT_EMPTY, TT_BLOCK, TT_TENTACLE_RIGHT, TT_DECOR, TT_RUNE, TT_SLOPE_45, TT_SLOPE_315},
			top:    []TileType{TT_EMPTY, TT_BLOCK, TT_TENTACLE_UP, TT_DECOR, TT_RUNE, TT_SLOPE_45, TT_SLOPE_135},
			bottom: []TileType{TT_EMPTY, TT_BLOCK, TT_TENTACLE_DOWN, TT_DECOR, TT_RUNE, TT_SLOPE_225, TT_SLOPE_315},
		},
		TT_SLOPE_45: {
			left:   []TileType{TT_BLOCK, TT_SLOPE_135},
			right:  []TileType{TT_EMPTY, TT_DECOR},
			top:    []TileType{TT_EMPTY, TT_DECOR},
			bottom: []TileType{TT_BLOCK, TT_SLOPE_315},
		},
		TT_SLOPE_135: {
			left:   []TileType{TT_EMPTY, TT_DECOR},
			right:  []TileType{TT_BLOCK, TT_SLOPE_45},
			top:    []TileType{TT_EMPTY, TT_DECOR},
			bottom: []TileType{TT_BLOCK, TT_SLOPE_225},
		},
		TT_SLOPE_225: {
			left:   []TileType{TT_EMPTY, TT_DECOR},
			right:  []TileType{TT_BLOCK, TT_SLOPE_315},
			top:    []TileType{TT_BLOCK, TT_SLOPE_135},
			bottom: []TileType{TT_EMPTY, TT_DECOR},
		},
		TT_SLOPE_315: {
			left:   []TileType{TT_BLOCK, TT_SLOPE_225},
			right:  []TileType{TT_EMPTY, TT_DECOR},
			top:    []TileType{TT_BLOCK, TT_SLOPE_45},
			bottom: []TileType{TT_EMPTY, TT_DECOR},
		},
		TT_TENTACLE_UP: {
			left:   []TileType{TT_EMPTY, TT_TENTACLE_UP},
			right:  []TileType{TT_EMPTY, TT_TENTACLE_UP},
			top:    []TileType{TT_EMPTY},
			bottom: []TileType{TT_BLOCK, TT_RUNE},
		},
		TT_TENTACLE_DOWN: {
			left:   []TileType{TT_EMPTY, TT_TENTACLE_DOWN},
			right:  []TileType{TT_EMPTY, TT_TENTACLE_DOWN},
			top:    []TileType{TT_BLOCK, TT_RUNE},
			bottom: []TileType{TT_EMPTY},
		},
		TT_TENTACLE_LEFT: {
			left:   []TileType{TT_EMPTY},
			right:  []TileType{TT_BLOCK},
			top:    []TileType{TT_EMPTY, TT_TENTACLE_LEFT},
			bottom: []TileType{TT_EMPTY, TT_TENTACLE_LEFT},
		},
		TT_TENTACLE_RIGHT: {
			left:   []TileType{TT_BLOCK},
			right:  []TileType{TT_EMPTY},
			top:    []TileType{TT_EMPTY, TT_TENTACLE_RIGHT},
			bottom: []TileType{TT_EMPTY, TT_TENTACLE_RIGHT},
		},
		TT_RUNE: {
			left:   []TileType{TT_EMPTY, TT_RUNE, TT_BLOCK},
			right:  []TileType{TT_EMPTY, TT_RUNE, TT_BLOCK},
			top:    []TileType{TT_EMPTY, TT_BLOCK, TT_TENTACLE_UP},
			bottom: []TileType{TT_EMPTY, TT_BLOCK, TT_TENTACLE_DOWN},
		},
		TT_DECOR: {
			left:   []TileType{TT_EMPTY, TT_BLOCK, TT_DECOR, TT_RUNE},
			right:  []TileType{TT_EMPTY, TT_BLOCK, TT_DECOR, TT_RUNE},
			top:    []TileType{TT_EMPTY, TT_BLOCK, TT_DECOR, TT_RUNE},
			bottom: []TileType{TT_EMPTY, TT_BLOCK, TT_DECOR, TT_RUNE},
		},
	}
}

type Tile struct {
	*Sprite
	ty TileType
}

func NewTile(ty TileType) *Tile {
	orient := 0

	switch ty {
	case TT_SLOPE_45, TT_TENTACLE_RIGHT:
		orient = 1
	case TT_SLOPE_315, TT_TENTACLE_DOWN:
		orient = 2
	case TT_SLOPE_225, TT_TENTACLE_LEFT:
		orient = 3
	}

	var spr *Sprite
	if ty != TT_EMPTY {
		spr = &Sprite{
			src:    tileTypeRects[ty],
			ofs:    ZeroVec(),
			flipH:  false,
			flipV:  false,
			orient: orient,
		}
	}
	return &Tile{
		spr,
		ty,
	}
}

func PropagateTile(t [][]TileType, x, y int) {
	cst := tileConstraints[t[y][x]]
	h := len(t)
	w := len(t[0])
	if x > 0 {
		if t[y][x-1] == -1 {
			t[y][x-1] = cst.left[rand.Intn(len(cst.left))]
			defer PropagateTile(t, x-1, y)
		}
	}
	if x < w-1 {
		if t[y][x+1] == -1 {
			t[y][x+1] = cst.right[rand.Intn(len(cst.right))]
			defer PropagateTile(t, x+1, y)
		}
	}
	if y > 0 {
		if t[y-1][x] == -1 {
			t[y-1][x] = cst.top[rand.Intn(len(cst.top))]
			defer PropagateTile(t, x, y-1)
		}
	}
	if y < h-1 {
		if t[y+1][x] == -1 {
			t[y+1][x] = cst.bottom[rand.Intn(len(cst.bottom))]
			defer PropagateTile(t, x, y+1)
		}
	}
}

func GenerateLevel(w, h int) *Level {
	t := make([][]TileType, h)
	for j := 0; j < h; j++ {
		t[j] = make([]TileType, w)
		for i := 0; i < w; i++ {
			t[j][i] = -1
		}
	}

	x, y := rand.Intn(w), rand.Intn(h)
	t[y][x] = TT_BLOCK
	PropagateTile(t, x, y)

	//Correct tiles
	for k := 0; k < 10; k++ {
		for j := 0; j < h; j++ {
			for i := 0; i < w; i++ {
				potentialSet := list.New()
				for _, tt := range allTileTypes {
					potentialSet.PushBack(tt)
				}
				intersect := func(tc []TileType) {
					if potentialSet.Len() != 0 {
						for tt := potentialSet.Front(); tt != nil; tt = tt.Next() {
							t := tt.Value.(TileType)
							//tn := tt.Next()
							for _, c := range tc {
								if t == c {
									goto valid
								}
							}
							//When the type in the potential set is not valid in the tile's constraints
							potentialSet.Remove(tt)
							//tt = tn
						valid:
						}
					}
				}
				if i > 0 {
					intersect(tileConstraints[t[j][i-1]].right)
				}
				if i < w-1 {
					intersect(tileConstraints[t[j][i+1]].left)
				}
				if j > 0 {
					intersect(tileConstraints[t[j-1][i]].bottom)
				}
				if j < h-1 {
					intersect(tileConstraints[t[j+1][i]].top)
				}
				if potentialSet.Len() > 0 {
					//Create array of possible set
					ps := make([]TileType, 0, potentialSet.Len())
					for tte := potentialSet.Front(); tte != nil; tte = tte.Next() {
						tt := tte.Value.(TileType)
						ps = append(ps, tt)
					}
					t[j][i] = ps[rand.Intn(len(ps))]
				} else {
					fmt.Printf("Impossible tile at %d, %d", i, j)
					t[j][i] = TT_EMPTY
				}
			}
		}
	}

	tiles := make([][]*Tile, h)
	for j := 0; j < h; j++ {
		tiles[j] = make([]*Tile, w)
		for i := 0; i < w; i++ {
			tiles[j][i] = NewTile(t[j][i])
		}
	}
	return &Level{tiles: tiles, rows: h, cols: w, pixelW: w * TILE_SIZE, pixelH: h * TILE_SIZE}
}

func (lev *Level) Draw(screen *ebiten.Image, pt *ebiten.GeoM) {
	op := &ebiten.DrawImageOptions{}
	for j := 0; j < lev.rows; j++ {
		op.GeoM.Concat(*pt)
		op.GeoM.Translate(0.0, float64(j)*TILE_SIZE)
		for i := 0; i < lev.cols; i++ {
			if lev.tiles[j][i].Sprite != nil {
				lev.tiles[j][i].Draw(screen, &op.GeoM)
			}
			op.GeoM.Translate(TILE_SIZE, 0.0)
		}
		op.GeoM.Reset()
	}
}
