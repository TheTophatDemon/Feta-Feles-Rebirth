package main

import (
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten"
)

type Level struct {
	tiles                   [][]Tile
	rows, cols              int
	pixelWidth, pixelHeight float64
}

func NewLevel(cols, rows int) *Level {
	tiles := make([][]Tile, rows)
	for y := 0; y < rows; y++ {
		tiles[y] = make([]Tile, cols)
		for x := 0; x < cols; x++ {
			tiles[y][x] = Tile{
				tt:      TT_EMPTY,
				spr:     nil,
				gridX:   x,
				gridY:   y,
				left:    float64(x) * TILE_SIZE,
				right:   float64(x)*TILE_SIZE + TILE_SIZE,
				top:     float64(y) * TILE_SIZE,
				bottom:  float64(y)*TILE_SIZE + TILE_SIZE,
				centerX: float64(x)*TILE_SIZE + TILE_SIZE_H,
				centerY: float64(y)*TILE_SIZE + TILE_SIZE_H,
			}
		}
	}

	pixelWidth := float64(cols*TILE_SIZE + TILE_SIZE)
	pixelHeight := float64(rows*TILE_SIZE + TILE_SIZE)

	return &Level{tiles, rows, cols, pixelWidth, pixelHeight}
}

func (level *Level) WrapCoords(x, y int) (int, int) {
	if x < 0 {
		x += level.cols
	}
	if x >= level.cols {
		x -= level.cols
	}
	if y < 0 {
		y += level.rows
	}
	if y >= level.rows {
		y -= level.rows
	}
	return x, y
}

//Sets the tile at the coordinate to specified type. Returns true if coordinate is valid. If wrap is set, out of bounds coordinates will be offset to the other side of the level.
func (level *Level) SetTile(x, y int, newType TileType, wrap bool) bool {
	if wrap {
		x, y = level.WrapCoords(x, y)
	} else if x < 0 || y < 0 || x >= level.cols || y >= level.rows {
		return false
	}
	level.tiles[y][x].tt = newType
	level.tiles[y][x].modified = true
	return true
}

//Gets a reference to the tile at the coordinates. Returns nil if out of bounds unless wrap is enabled.
func (level *Level) GetTile(x, y int, wrap bool) *Tile {
	if wrap {
		x, y = level.WrapCoords(x, y)
	} else if x < 0 || y < 0 || x >= level.cols || y >= level.rows {
		return nil
	}
	return &level.tiles[y][x]
}

//Returns the tile at the center of an empty space of tile radius r
func (level *Level) FindEmptySpace(r int) *Tile {
	for {
		x, y := rand.Intn(level.cols), rand.Intn(level.rows)
		for j := y - r; j <= y+r; j++ {
			for i := x - r; i <= x+r; i++ {
				if level.GetTile(i, j, true).IsSolid() {
					goto reject
				}
			}
		}
		return level.GetTile(x, y, true)
	reject:
	}
}

//Like FindEmptySpace except for finding places inside of the walls
func (level *Level) FindFullSpace(r int) *Tile {
	for {
		x, y := rand.Intn(level.cols), rand.Intn(level.rows)
		for j := y - r; j <= y+r; j++ {
			for i := x - r; i <= x+r; i++ {
				if !level.GetTile(i, j, true).IsSolid() {
					goto reject
				}
			}
		}
		return level.GetTile(x, y, true)
	reject:
	}
}

func (level *Level) Draw(game *Game, screen *ebiten.Image, pt *ebiten.GeoM) {
	//Determine the area of the grid that is on screen
	gridMin := game.camPos.Clone().Sub(&Vec2f{SCR_WIDTH_H, SCR_HEIGHT_H}).Scale(1.0 / TILE_SIZE).Floor()
	gridMax := game.camPos.Clone().Add(&Vec2f{SCR_WIDTH_H, SCR_HEIGHT_H}).Scale(1.0 / TILE_SIZE).Ceil()
	iminx := int(math.Max(0.0, gridMin.x))
	iminy := int(math.Max(0.0, gridMin.y))
	imaxx := int(math.Min(float64(level.cols), gridMax.x))
	imaxy := int(math.Min(float64(level.rows), gridMax.y))
	//Draw only the tiles in that area
	for j := iminy; j < imaxy; j++ {
		for i := iminx; i < imaxx; i++ {
			if level.tiles[j][i].modified {
				level.tiles[j][i].RegenSprite()
				level.tiles[j][i].modified = false
			}
			if level.tiles[j][i].spr != nil {
				level.tiles[j][i].spr.Draw(screen, pt)
			}
		}
	}
}

//Pos is the position. I and J are the tile indices.
func (level *Level) ProjectPosOntoTile(pos *Vec2f, t *Tile) *Vec2f {
	tileMin := &Vec2f{x: t.left, y: t.top}
	tileMax := &Vec2f{x: t.right, y: t.bottom}

	var proj *Vec2f
	if !t.IsSlope() {
		//Project onto a box by clamping the destination to the box boundaries
		proj = VecMax(tileMin, VecMin(tileMax, pos))
	} else {
		//Project onto a diagonal plane using the dot product
		cDiff := pos.Clone().Sub(&Vec2f{x: t.centerX, y: t.centerY})
		planeDist := VecDot(t.GetSlopeNormal(), cDiff)
		proj = pos.Clone().Sub(t.GetSlopeNormal().Scale(planeDist))
		proj = VecMax(tileMin, VecMin(tileMax, proj))
	}

	return proj
}

func (level *Level) GetGridAreaOverCapsule(start, dest *Vec2f, radius float64, clamp bool) (gridMin, gridMax *Vec2f) {
	gridMin = VecMin(start, dest).SubScalar(radius).Scale(1.0 / TILE_SIZE).Floor()
	if clamp {
		//gridMin = VecMax(ZeroVec(), gridMin)
	}
	gridMax = VecMax(start, dest).AddScalar(radius).Scale(1.0 / TILE_SIZE).Ceil()
	if clamp {
		//gridMax = VecMin(&Vec2f{x: float64(level.cols), y: float64(level.rows)}, gridMax)
	}
	return
}

//Determines if sphere intersects a solid tile. If so, the normal is returned.
func (level *Level) SphereIntersects(pos *Vec2f, radius float64) (bool, *Vec2f) {
	gridMin, gridMax := level.GetGridAreaOverCapsule(pos, pos, radius, false)

	for j := int(gridMin.y); j < int(gridMax.y); j++ {
		for i := int(gridMin.x); i < int(gridMax.x); i++ {
			t := game.level.GetTile(i, j, true)
			if t.IsSolid() {
				diff := pos.Clone().Sub(level.ProjectPosOntoTile(pos, t))
				dLen := diff.Length()
				if dLen < radius {
					if dLen != 0.0 {
						diff.Scale(1.0 / dLen)
					}
					return true, diff
				}
			}
		}
	}

	return false, nil
}

type RaycastResult struct {
	hit      bool
	pos      *Vec2f
	distance float64
	tile     *Tile
}

func (level *Level) Raycast(pos *Vec2f, dir *Vec2f, maxDist float64) *RaycastResult {
	var rx, ry, rdx, rdy, tan float64
	if dir.x != 0.0 {
		tan = dir.y / dir.x
	}

	castRay := func(x, y, dx, dy float64, vert bool) (*Tile, float64, float64) {
		ox, oy := x, y
		fauxDist := (&Vec2f{pos.x - x, pos.y - y}).Length()
		fauxStep := (&Vec2f{dx, dy}).Length() //The approximate distance the ray travels each step
		for ; fauxDist+fauxStep < maxDist; fauxDist += fauxStep {
			ix := int(x / TILE_SIZE)
			iy := int(y / TILE_SIZE)

			if vert {
				if dx < 0 {
					ix--
				}
			} else {
				if dy < 0 {
					iy--
				}
			}

			if ix < 0 || iy < 0 || ix >= level.cols || iy >= level.rows {
				return nil, x, y
			}

			t := level.tiles[iy][ix]
			if t.IsSlope() {
				//Test against slopes
				slopeNormal := level.tiles[iy][ix].GetSlopeNormal()
				//Calculate intersection point
				t := (slopeNormal.x*(t.centerX-x) + slopeNormal.y*(t.centerY-y)) /
					((slopeNormal.x * dx) + (slopeNormal.y * dy))
				px, py := x+dx*t, y+dy*t
				//Test if it is within the tile's boundaries
				if px >= float64(ix)*TILE_SIZE && px < float64(ix+1)*TILE_SIZE &&
					py >= float64(iy)*TILE_SIZE && py < float64(iy+1)*TILE_SIZE {
					return &level.tiles[iy][ix], px, py
				}
			} else if t.IsSolid() {
				return &t, x, y
			}
			x += dx
			y += dy
		}
		return nil, ox + (maxDist * dx / fauxStep), oy + (maxDist * dy / fauxStep)
	}

	//Vertical line phase (moving x)
	if dir.x > 0 {
		rx = math.Ceil(pos.x/TILE_SIZE) * TILE_SIZE
		rdx = TILE_SIZE
	} else {
		rx = math.Floor(pos.x/TILE_SIZE) * TILE_SIZE
		rdx = -TILE_SIZE
	}
	ry = pos.y + (rx-pos.x)*tan
	rdy = rdx * tan
	//Raycast loop, etc.
	var vertX, vertY float64
	var vTile *Tile
	if dir.x != 0.0 {
		vTile, vertX, vertY = castRay(rx, ry, rdx, rdy, true)
	}

	//Horizontal line phase (moving y)
	if dir.y > 0 {
		ry = math.Ceil(pos.y/TILE_SIZE) * TILE_SIZE
		rdy = TILE_SIZE
	} else {
		ry = math.Floor(pos.y/TILE_SIZE) * TILE_SIZE
		rdy = -TILE_SIZE
	}
	if tan == 0.0 {
		rx = pos.x
		rdx = 0.0
	} else {
		rx = pos.x + (ry-pos.y)/tan
		rdx = rdy / tan
	}
	//Raycast loop, etc.
	var horzX, horzY float64
	var hTile *Tile
	if dir.y != 0.0 {
		hTile, horzX, horzY = castRay(rx, ry, rdx, rdy, false)
	}
	//hHit, horzX, horzY := false, 0.0, 0.0

	vDist := math.Pow(vertX-pos.x, 2.0) + math.Pow(vertY-pos.y, 2.0)
	hDist := math.Pow(horzX-pos.x, 2.0) + math.Pow(horzY-pos.y, 2.0)
	if hDist < vDist {
		return &RaycastResult{
			hit:      hTile != nil,
			pos:      &Vec2f{horzX, horzY},
			distance: math.Sqrt(hDist),
			tile:     hTile,
		}
	} else {
		return &RaycastResult{
			hit:      vTile != nil,
			pos:      &Vec2f{vertX, vertY},
			distance: math.Sqrt(vDist),
			tile:     vTile,
		}
	}
}
