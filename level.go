package main

import (
	"container/list"
	"image"
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/thetophatdemon/Feta-Feles-Remastered/vmath"
)

type Level struct {
	tiles                   [][]Tile
	spaces                  []*Space
	rows, cols              int
	pixelWidth, pixelHeight float64
	recalcEdges             bool //Flag for when edges need to be recalculated
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

	pixelWidth := float64(cols * TILE_SIZE)
	pixelHeight := float64(rows * TILE_SIZE)

	return &Level{tiles, make([]*Space, 0, 10), rows, cols, pixelWidth, pixelHeight, false}
}

func (level *Level) WrapGridCoords(x, y int) (int, int) {
	x = x % level.cols
	y = y % level.rows
	if x < 0 {
		x += level.cols
	}
	if y < 0 {
		y += level.rows
	}
	return x, y
}

func (level *Level) WrapPixelCoords(x, y float64) (float64, float64) {
	if x < 0.0 {
		x += level.pixelWidth
	} else if x >= level.pixelWidth {
		x -= level.pixelWidth
	}
	if y < 0.0 {
		y += level.pixelHeight
	} else if y >= level.pixelHeight {
		y -= level.pixelHeight
	}
	return x, y
}

//Sets the tile at the coordinate to specified type. Returns true if coordinate is valid. If wrap is set, out of bounds coordinates will be offset to the other side of the level.
func (level *Level) SetTile(x, y int, newType TileType, wrap bool) bool {
	if wrap {
		x, y = level.WrapGridCoords(x, y)
	} else if x < 0 || y < 0 || x >= level.cols || y >= level.rows {
		return false
	}
	level.tiles[y][x].SetType(newType)
	return true
}

//Removes a solid tile and reshapes the surrounding terrain to make the deformation smooth
func (level *Level) DestroyTile(t *Tile) {
	if t.IsSolid() {
		level.recalcEdges = true
	}
	t.SetType(TT_EMPTY)
}

//Gets a reference to the tile at the coordinates. Returns nil if out of bounds unless wrap is enabled.
func (level *Level) GetTile(x, y int, wrap bool) *Tile {
	if wrap {
		x, y = level.WrapGridCoords(x, y)
	} else if x < 0 || y < 0 || x >= level.cols || y >= level.rows {
		return nil
	}
	return &level.tiles[y][x]
}

//Randomly chooses an empty tile.
func (level *Level) FindSpawnPoint() *Tile {
	emptyTiles := make([]*Tile, 0, 1024)
	for _, sp := range level.spaces {
		for _, t := range sp.tiles {
			if t.tt == TT_EMPTY {
				//Check is neccesary because pylon placement is done after space mapping
				emptyTiles = append(emptyTiles, t)
			}
		}
	}
	return emptyTiles[rand.Intn(len(emptyTiles))]
}

//Randomly chooses an empty tile that is off screen
func (level *Level) FindOffscreenSpawnPoint(game *Game) *Tile {
	emptyTiles := make([]*Tile, 0, 1024)
	for _, sp := range level.spaces {
		for _, t := range sp.tiles {
			//Find empty and off-screen tiles
			if t.tt == TT_EMPTY && !game.SquareOnScreen(t.centerX, t.centerY, TILE_SIZE_H) {
				//Skip tile if something's already there
				for e := game.objects.Front(); e != nil; e = e.Next() {
					obj := e.Value.(*Object)
					if obj.pos.Clone().Sub(vmath.NewVec(t.centerX, t.centerY)).Length() < obj.radius*2.0 {
						goto skip
					}
				}
				emptyTiles = append(emptyTiles, t)
			skip:
			}
		}
	}
	if len(emptyTiles) == 0 {
		return nil
	}
	return emptyTiles[rand.Intn(len(emptyTiles))]
}

//Randomly chooses an empty tile that is somewhat near the center
func (level *Level) FindCenterSpawnPoint(game *Game) *Tile {
	emptyTiles := make([]*Tile, 0, 1024)
	for _, sp := range level.spaces {
		for _, t := range sp.tiles {
			if t.tt == TT_EMPTY && t.centerX > SCR_WIDTH_H &&
				t.centerY > SCR_WIDTH_H && t.centerX < level.pixelWidth-SCR_WIDTH_H && t.centerY < level.pixelHeight-SCR_HEIGHT_H {
				emptyTiles = append(emptyTiles, t)
			}
		}
	}
	if len(emptyTiles) == 0 {
		return nil
	}
	return emptyTiles[rand.Intn(len(emptyTiles))]
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

//Struct represents a glob of contiguous empty space (Not taking into account screen wrapping)
type Space struct {
	tiles            []*Tile
	frontier         []*Tile //These tiles are on the border of the space\
	centerX, centerY float64 //The average position of all its tiles, used for approximation
}

var spaceImg *ebiten.Image
var spaceColors [255]color.RGBA
var spaceCenterImg *ebiten.Image

func init() {
	spaceImg = ebiten.NewImage(8, 8)
	spaceImg.Fill(color.RGBA{255, 255, 255, 255})
	for i := 0; i < 255; i++ {
		spaceColors[i] = color.RGBA{uint8(rand.Intn(256)), uint8(rand.Intn(256)), uint8(rand.Intn(256)), 255}
	}
	spaceCenterImg = GetGraphics().SubImage(image.Rect(96, 96, 112, 112)).(*ebiten.Image)
}

//Draws a colored square for debugging
func (space *Space) Draw(screen *ebiten.Image, pt *ebiten.GeoM, clr color.RGBA) {
	for _, t := range space.tiles {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(-4.0, -4.0)
		op.GeoM.Translate(t.centerX, t.centerY)
		op.GeoM.Concat(*pt)
		op.ColorM.Scale(0.0, 0.0, 0.0, 1.0)
		fr, fg, fb := float64(clr.R), float64(clr.G), float64(clr.B)
		op.ColorM.Translate(fr/255.0, fg/255.0, fb/255.0, 0.0)
		screen.DrawImage(spaceImg, op)
	}
	//Draw frontier tiles with small inner squares
	for _, t := range space.frontier {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(-4.0, -4.0)
		op.GeoM.Scale(0.5, 0.5)
		op.GeoM.Translate(t.centerX, t.centerY)
		op.GeoM.Concat(*pt)
		op.ColorM.Scale(0.0, 0.0, 0.0, 1.0)
		fr, fg, fb := 1.0-float64(clr.R), 1.0-float64(clr.G), 1.0-float64(clr.B)
		op.ColorM.Translate(fr/255.0, fg/255.0, fb/255.0, 0.0)
		screen.DrawImage(spaceImg, op)
	}
	//Draw symbol at center
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-8.0, -8.0)
	op.GeoM.Translate(space.centerX, space.centerY)
	op.GeoM.Concat(*pt)
	screen.DrawImage(spaceCenterImg, op)
}

func (level *Level) FindSpaces() {
	//Clear existing space data
	level.spaces = make([]*Space, 0, 10)

	tilesLeft := list.New() //List of empty tiles to be assigned to spaces
	for j := 0; j < level.rows; j++ {
		for i := 0; i < level.cols; i++ {
			t := level.GetTile(i, j, false)
			if t.tt == TT_EMPTY {
				t.space = nil
				tilesLeft.PushBack(t)
			}
		}
	}
	for e := tilesLeft.Front(); e != nil; e = e.Next() {
		t := e.Value.(*Tile)
		if t.space == nil {
			space := new(Space)
			level.PropagateSpace(t, space)
			level.spaces = append(level.spaces, space)
		}
	}
	//Calculate center
	for _, sp := range level.spaces {
		sp.centerX, sp.centerY = 0.0, 0.0
		for _, t := range sp.tiles {
			sp.centerX += t.centerX
			sp.centerY += t.centerY
		}
		sp.centerX /= float64(len(sp.tiles))
		sp.centerY /= float64(len(sp.tiles))
	}
}

//Recursively adds to the space's domain by checking neighbors and propagating to empty neighbors (within the level bounds)
func (level *Level) PropagateSpace(tile *Tile, space *Space) {
	tile.space = space
	space.tiles = append(space.tiles, tile)
	neighbors := []*Tile{
		level.GetTile(tile.gridX-1, tile.gridY, false),
		level.GetTile(tile.gridX, tile.gridY-1, false),
		level.GetTile(tile.gridX+1, tile.gridY, false),
		level.GetTile(tile.gridX, tile.gridY+1, false),
		//level.GetTile(tile.gridX+1, tile.gridY+1, false),
		//level.GetTile(tile.gridX-1, tile.gridY+1, false),
		//level.GetTile(tile.gridX+1, tile.gridY-1, false),
		//level.GetTile(tile.gridX-1, tile.gridY-1, false),
	}
	for _, n := range neighbors {
		if n != nil {
			if n.IsSolid() {
				space.frontier = append(space.frontier, tile)
			} else if n.space == nil {
				level.PropagateSpace(n, space)
			}
		}
	}
}

func (level *Level) Draw(game *Game, screen *ebiten.Image, pt *ebiten.GeoM) {
	//Reorient tiles if some have been deleted
	if level.recalcEdges {
		level.SmoothEdges()
		level.recalcEdges = false
	}

	//Determine the area of the grid that is on screen
	gridMin := game.camPos.Clone().Sub(vmath.NewVec(SCR_WIDTH_H, SCR_HEIGHT_H)).Scale(1.0 / TILE_SIZE).Floor()
	gridMax := game.camPos.Clone().Add(vmath.NewVec(SCR_WIDTH_H, SCR_HEIGHT_H)).Scale(1.0 / TILE_SIZE).Ceil()
	//Draw only the tiles in that area
	for j := int(gridMin.Y); j < int(gridMax.Y); j++ {
		for i := int(gridMin.X); i < int(gridMax.X); i++ {
			t := level.GetTile(i, j, true)
			if t.modified {
				t.RegenSprite()
				t.modified = false
			}
			if t.spr != nil {
				mat := *pt
				if i < 0 {
					mat.Translate(-level.pixelWidth, 0.0)
				} else if i >= level.cols {
					mat.Translate(level.pixelWidth, 0.0)
				}
				if j < 0 {
					mat.Translate(0.0, -level.pixelHeight)
				} else if j >= level.rows {
					mat.Translate(0.0, level.pixelHeight)
				}
				t.spr.Draw(screen, &mat)
				//Draw outlines for square tiles
				if t.outline != OUTLINE_NONE {
					ofsx := mat.Element(0, 2)
					ofsy := mat.Element(1, 2)
					col := color.Black
					if t.outline&OUTLINE_TOP > 0 {
						ebitenutil.DrawLine(screen, ofsx+t.left, ofsy+t.top, ofsx+t.right, ofsy+t.top, col)
					}
					if t.outline&OUTLINE_BOTTOM > 0 {
						ebitenutil.DrawLine(screen, ofsx+t.left, ofsy+t.bottom, ofsx+t.right, ofsy+t.bottom, col)
					}
					if t.outline&OUTLINE_LEFT > 0 {
						ebitenutil.DrawLine(screen, ofsx+t.left+1, ofsy+t.top+1, ofsx+t.left+1, ofsy+t.bottom, col)
					}
					if t.outline&OUTLINE_RIGHT > 0 {
						ebitenutil.DrawLine(screen, ofsx+t.right, ofsy+t.top, ofsx+t.right, ofsy+t.bottom, col)
					}
				}
			}
		}
	}

	if debugDraw {
		//Draw spaces
		for i, sp := range level.spaces {
			sp.Draw(screen, pt, spaceColors[i%len(spaceColors)])
		}
	}
}

//If t is nil, then position is projected onto level boundaries
func (level *Level) ProjectPosOntoTile(pos *vmath.Vec2f, t *Tile) *vmath.Vec2f {
	tileMin := vmath.NewVec(t.left, t.top)
	tileMax := vmath.NewVec(t.right, t.bottom)

	//Project onto a box by clamping the destination to the box boundaries
	proj := vmath.VecMax(tileMin, vmath.VecMin(tileMax, pos))
	if t != nil && t.IsSlope() {
		//Project onto a diagonal plane using the dot product if positing is coming from the right direction
		cDiff := pos.Clone().Sub(vmath.NewVec(t.centerX, t.centerY))
		planeDist := vmath.VecDot(t.GetSlopeNormal(), cDiff)
		if planeDist > 0.0 {
			proj = pos.Clone().Sub(t.GetSlopeNormal().Scale(planeDist))
			proj = vmath.VecMax(tileMin, vmath.VecMin(tileMax, proj))
		}
	}

	return proj
}

func (level *Level) GetGridAreaOverCapsule(start, dest *vmath.Vec2f, radius float64, clamp bool) (gridMin, gridMax *vmath.Vec2f) {
	gridMin = vmath.VecMin(start, dest).SubScalar(radius).Scale(1.0 / TILE_SIZE).Floor()
	if clamp {
		gridMin = vmath.VecMax(vmath.ZeroVec(), gridMin)
	}
	gridMax = vmath.VecMax(start, dest).AddScalar(radius).Scale(1.0 / TILE_SIZE).Ceil()
	if clamp {
		gridMax = vmath.VecMin(vmath.NewVec(float64(level.cols), float64(level.rows)), gridMax)
	}
	return gridMin, gridMax
}

func (level *Level) GetTilesWithinRadius(pos *vmath.Vec2f, radius float64) []*Tile {
	gridMin := pos.Clone().SubScalar(radius).Scale(1.0 / TILE_SIZE).Floor()
	gridMax := pos.Clone().AddScalar(radius).Scale(1.0 / TILE_SIZE).Ceil()
	result := make([]*Tile, 0, int(radius*2.0*radius*2.0))
	for i := int(gridMin.X); i < int(gridMax.X); i++ {
		for j := int(gridMin.Y); j < int(gridMax.Y); j++ {
			if t := level.GetTile(i, j, true); t != nil {
				diff := (vmath.NewVec(t.centerX, t.centerY)).Sub(pos)
				if diff.Length() < radius {
					result = append(result, t)
				}
			}
		}
	}
	return result
}

//Determines if sphere intersects a solid tile. If so, the normal and the collided tile is returned
func (level *Level) SphereIntersects(pos *vmath.Vec2f, radius float64) (bool, *vmath.Vec2f, *Tile) {
	//Check against level borders
	if pos.X-radius < 0 {
		return true, vmath.NewVec(1.0, 0.0), nil
	} else if pos.X+radius > level.pixelWidth {
		return true, vmath.NewVec(-1.0, 0.0), nil
	}
	if pos.Y-radius < 0 {
		return true, vmath.NewVec(0.0, 1.0), nil
	} else if pos.Y+radius > level.pixelHeight {
		return true, vmath.NewVec(0.0, -1.0), nil
	}

	gridMin, gridMax := level.GetGridAreaOverCapsule(pos, pos, radius, true)
	for j := int(gridMin.Y); j < int(gridMax.Y); j++ {
		for i := int(gridMin.X); i < int(gridMax.X); i++ {
			t := level.GetTile(i, j, true)
			if t.IsSolid() {
				diff := pos.Clone().Sub(level.ProjectPosOntoTile(pos, t))
				dLen := diff.Length()
				if dLen < radius {
					if dLen != 0.0 {
						diff.Scale(1.0 / dLen)
					}
					return true, diff, t
				}
			}
		}
	}

	return false, nil, nil
}

type RaycastResult struct {
	hit      bool
	pos      *vmath.Vec2f
	distance float64
	tile     *Tile
}

func (level *Level) Raycast(pos *vmath.Vec2f, dir *vmath.Vec2f, maxDist float64) *RaycastResult {
	var rx, ry, rdx, rdy, tan float64
	if dir.X != 0.0 {
		tan = dir.Y / dir.X
	}

	castRay := func(x, y, dx, dy float64, vert bool) (*Tile, float64, float64) {
		ox, oy := x, y
		fauxDist := (vmath.NewVec(pos.X-x, pos.Y-y)).Length()
		fauxStep := (vmath.NewVec(dx, dy)).Length() //The approximate distance the ray travels each step
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
			//ix, iy = level.WrapGridCoords(ix, iy)
			if ix < 0 || iy < 0 || ix >= level.cols || iy >= level.rows {
				return nil, x, y
			}

			t := level.tiles[iy][ix]
			if t.IsSlope() {
				//Test against slopes
				slopeNormal := level.tiles[iy][ix].GetSlopeNormal()
				//Calculate intersection point
				wx, wy := level.WrapPixelCoords(x, y)
				t := (slopeNormal.X*(t.centerX-wx) + slopeNormal.Y*(t.centerY-wy)) /
					((slopeNormal.X * dx) + (slopeNormal.Y * dy))
				px, py := wx+dx*t, wy+dy*t
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
	if dir.X > 0 {
		rx = math.Ceil(pos.X/TILE_SIZE) * TILE_SIZE
		rdx = TILE_SIZE
	} else {
		rx = math.Floor(pos.X/TILE_SIZE) * TILE_SIZE
		rdx = -TILE_SIZE
	}
	ry = pos.Y + (rx-pos.X)*tan
	rdy = rdx * tan
	//Raycast loop, etc.
	var vertX, vertY float64
	var vTile *Tile
	if dir.X != 0.0 {
		vTile, vertX, vertY = castRay(rx, ry, rdx, rdy, true)
	}

	//Horizontal line phase (moving y)
	if dir.Y > 0 {
		ry = math.Ceil(pos.Y/TILE_SIZE) * TILE_SIZE
		rdy = TILE_SIZE
	} else {
		ry = math.Floor(pos.Y/TILE_SIZE) * TILE_SIZE
		rdy = -TILE_SIZE
	}
	if tan == 0.0 {
		rx = pos.X
		rdx = 0.0
	} else {
		rx = pos.X + (ry-pos.Y)/tan
		rdx = rdy / tan
	}
	//Raycast loop, etc.
	var horzX, horzY float64
	var hTile *Tile
	if dir.Y != 0.0 {
		hTile, horzX, horzY = castRay(rx, ry, rdx, rdy, false)
	}
	//hHit, horzX, horzY := false, 0.0, 0.0

	vDist := math.Pow(vertX-pos.X, 2.0) + math.Pow(vertY-pos.Y, 2.0)
	hDist := math.Pow(horzX-pos.X, 2.0) + math.Pow(horzY-pos.Y, 2.0)
	if hDist < vDist {
		return &RaycastResult{
			hit:      hTile != nil || horzX <= 0.0 || horzY <= 0.0 || horzX >= level.pixelWidth || horzY >= level.pixelHeight,
			pos:      vmath.NewVec(horzX, horzY),
			distance: math.Sqrt(hDist),
			tile:     hTile,
		}
	} else {
		return &RaycastResult{
			hit:      vTile != nil || horzX <= 0.0 || horzY <= 0.0 || horzX >= level.pixelWidth || horzY >= level.pixelHeight,
			pos:      vmath.NewVec(vertX, vertY),
			distance: math.Sqrt(vDist),
			tile:     vTile,
		}
	}
}
