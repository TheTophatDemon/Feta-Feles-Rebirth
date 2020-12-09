package main

import (
	"image"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten"
)

const (
	TILE_SIZE = 16.0
)

type Level struct {
	tiles      [][]TileType
	sprites    [][]*Sprite
	spawns     []*Spawn
	rows, cols int
	positions  [][]ebiten.GeoM
}

type SpawnType int

const (
	SP_PLAYER SpawnType = iota
	SP_BARREL
)

type Spawn struct {
	spawnType SpawnType
	ix, iy    int
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
)

var tileTypeRects map[TileType]image.Rectangle

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
}

func SpriteFromTile(tt TileType) *Sprite {
	if tt != TT_EMPTY {
		orient := 0

		switch tt {
		case TT_SLOPE_45, TT_TENTACLE_RIGHT:
			orient = 1
		case TT_SLOPE_315, TT_TENTACLE_DOWN:
			orient = 2
		case TT_SLOPE_225, TT_TENTACLE_LEFT:
			orient = 3
		}

		if tt == TT_RUNE {
			x := int(math.Floor(rand.Float64()*4.0)) * 16
			tileTypeRects[tt] = image.Rect(x, 112, x+16, 128)
		}

		return NewSprite(tileTypeRects[tt], ZeroVec(), false, false, orient)
	} else {
		return nil
	}
}

func (level *Level) IsOccupied(x, y int) bool {
	if x <= 0 || y <= 0 || x >= level.cols || y >= level.rows {
		return true
	}
	return level.tiles[y][x]&TT_SOLIDS > 0
}

func (level *Level) FindEmptySpace(r int) (x, y int) {
	for {
		x, y = rand.Intn(level.cols), rand.Intn(level.rows)
		for j := y - r; j <= y+r; j++ {
			for i := x - r; i <= x+r; i++ {
				if i > 0 && j > 0 && i < level.cols-1 && j < level.rows-1 && level.IsOccupied(i, j) {
					goto reject
				}
			}
		}
		break
	reject:
	}
	return x, y
}

func (level *Level) FindFullSpace() (x, y int) {
	for {
		x, y = rand.Intn(level.cols), rand.Intn(level.rows)
		if level.IsOccupied(x, y) {
			break
		}
	}
	return x, y
}

func GenerateLevel(w, h int) *Level {
	t := make([][]TileType, h)
	for j := 0; j < h; j++ {
		t[j] = make([]TileType, w)
		for i := 0; i < w; i++ {
			t[j][i] = TT_EMPTY
		}
	}

	//Generate solid structures
	var p func(x, y int, f float32, fd float32)
	p = func(x, y int, f, fd float32) {
		t[y][x] = TT_BLOCK
		//Create a blob by recursively adding tiles adjacent to the one just placed with decreasing probability
		if f > 0.0 {
			if x > 0 && rand.Float32() < f {
				p(x-1, y, f-fd, fd)
			}
			if x < w-1 && rand.Float32() < f {
				p(x+1, y, f-fd, fd)
			}
			if y > 0 && rand.Float32() < f {
				p(x, y-1, f-fd, fd)
			}
			if y < h-1 && rand.Float32() < f {
				p(x, y+1, f-fd, fd)
			}
		}
	}
	for k := 0; k < w*h/64; k++ {
		x, y := rand.Intn(w-16)+8, rand.Intn(h-16)+8
		p(x, y, 1.0, 0.25)
	}
	//Smooth edges and create gaps
	for j := 0; j < h; j++ {
		for i := 0; i < w; i++ {
			nn := 0
			ln := false
			if i > 0 {
				ln = t[j][i-1]&TT_SOLIDS > 0
				if ln {
					nn++
				}
			}
			rn := false
			if i < w-1 {
				rn = t[j][i+1]&TT_SOLIDS > 0
				if rn {
					nn++
				}
			}
			tn := false
			if j > 0 {
				tn = t[j-1][i]&TT_SOLIDS > 0
				if tn {
					nn++
				}
			}
			bn := false
			if j < h-1 {
				bn = t[j+1][i]&TT_SOLIDS > 0
				if bn {
					nn++
				}
			}
			//Remove random holes
			if nn == 4 {
				t[j][i] = TT_BLOCK
			}
			//Turn poking structures into tentacles
			if nn == 1 && t[j][i] == TT_BLOCK {
				switch {
				case bn:
					if t[j+1][i] == TT_BLOCK {
						t[j][i] = TT_TENTACLE_UP
					}
				case tn:
					if t[j-1][i] == TT_BLOCK {
						t[j][i] = TT_TENTACLE_DOWN
					}
				case ln:
					if t[j][i-1] == TT_BLOCK {
						t[j][i] = TT_TENTACLE_RIGHT
					}
				case rn:
					if t[j][i+1] == TT_BLOCK {
						t[j][i] = TT_TENTACLE_LEFT
					}
				}
			}
			if t[j][i] == TT_BLOCK {
				//Turn into slope?
				if ln && bn && !tn && !rn {
					t[j][i] = TT_SLOPE_45
				} else if rn && bn && !ln && !tn {
					t[j][i] = TT_SLOPE_135
				} else if rn && tn && !ln && !bn {
					t[j][i] = TT_SLOPE_225
				} else if ln && tn && !rn && !bn {
					t[j][i] = TT_SLOPE_315
				}
			}
		}
	}
	level := &Level{tiles: t, rows: h, cols: w}

	//Add rune bars
	var rune func(x, y, d, l int)
	rune = func(x, y, d, l int) {
		level.tiles[y][x] = TT_RUNE
		if l > 0 {
			if d == 2 && level.IsOccupied(x-1, y) && x > 0 {
				rune(x-1, y, d, l-1)
			} else if d == 0 && level.IsOccupied(x+1, y) && x < w-1 {
				rune(x+1, y, d, l-1)
			} else if d == 1 && level.IsOccupied(x, y-1) && y > 0 {
				rune(x, y-1, d, l-1)
			} else if d == 3 && level.IsOccupied(x, y+1) && y < h-1 {
				rune(x, y+1, d, l-1)
			}
			if rand.Float32() < 0.2 {
				var nd int
				if d == 2 || d == 0 {
					if rand.Float32() > 0.5 {
						nd = 1
					} else {
						nd = 3
					}
				} else {
					if rand.Float32() > 0.5 {
						nd = 2
					} else {
						nd = 0
					}
				}
				rune(x, y, nd, l)
			}
		}
	}
	for i := 0; i < w*h/1024; i++ {
		rx, ry := level.FindFullSpace()
		for j := 0; j < 4; j++ {
			rune(rx, ry, j, 4)
		}
	}

	//Add pylons
	for i := 0; i < w*h/48; i++ {
		pix, piy := level.FindEmptySpace(1)
		level.tiles[piy][pix] = TT_PYLON
	}

	sprites := make([][]*Sprite, h)
	for j := 0; j < h; j++ {
		sprites[j] = make([]*Sprite, w)
		for i := 0; i < w; i++ {
			sprites[j][i] = SpriteFromTile(t[j][i])
		}
	}
	spawns := make([]*Spawn, 0, 10)
	px, py := level.FindEmptySpace(0)
	spawns = append(spawns, &Spawn{spawnType: SP_PLAYER, ix: px, iy: py})

	level.sprites = sprites
	level.spawns = spawns

	level.positions = make([][]ebiten.GeoM, h)
	mat := new(ebiten.GeoM)
	for j := 0; j < h; j++ {
		level.positions[j] = make([]ebiten.GeoM, w)
		mat.Translate(0.0, float64(j)*TILE_SIZE)
		for i := 0; i < w; i++ {
			level.positions[j][i] = *mat
			mat.Translate(TILE_SIZE, 0.0)
		}
		mat.Reset()
	}

	return level
}

func (lev *Level) Draw(screen *ebiten.Image, pt *ebiten.GeoM) {
	for j := 0; j < lev.rows; j++ {
		for i := 0; i < lev.cols; i++ {
			if lev.sprites[j][i] != nil {
				mat := lev.positions[j][i]
				mat.Concat(*pt)
				lev.sprites[j][i].Draw(screen, &mat)
			}
		}
	}
}
