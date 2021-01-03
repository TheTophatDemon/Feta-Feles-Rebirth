package main

import (
	"math"
	"math/rand"
)

const (
	SPREAD_DELTA = 0.25
)

//Generates a blob of tiles by recursively and randomly setting tiles adjacent to the last tile placed at decreasing frequency
func PropagateBlob(level *Level, x, y int, spreadChance float64) {
	level.SetTile(x, y, TT_BLOCK, true)
	if spreadChance > 0.0 {
		if rand.Float64() < spreadChance {
			PropagateBlob(level, x-1, y, spreadChance-SPREAD_DELTA)
		}
		if rand.Float64() < spreadChance {
			PropagateBlob(level, x+1, y, spreadChance-SPREAD_DELTA)
		}
		if rand.Float64() < spreadChance {
			PropagateBlob(level, x, y-1, spreadChance-SPREAD_DELTA)
		}
		if rand.Float64() < spreadChance {
			PropagateBlob(level, x, y+1, spreadChance-SPREAD_DELTA)
		}
	}
}

func PropagateRune(level *Level, x, y int, dir int, life int) {
	level.SetTile(x, y, TT_RUNE, true)
	if life > 0 {
		if dir == 2 && level.GetTile(x-1, y, true).tt == TT_BLOCK {
			PropagateRune(level, x-1, y, dir, life-1)
		} else if dir == 0 && level.GetTile(x+1, y, true).tt == TT_BLOCK {
			PropagateRune(level, x+1, y, dir, life-1)
		} else if dir == 1 && level.GetTile(x, y-1, true).tt == TT_BLOCK {
			PropagateRune(level, x, y-1, dir, life-1)
		} else if dir == 3 && level.GetTile(x, y+1, true).tt == TT_BLOCK {
			PropagateRune(level, x, y+1, dir, life-1)
		}
		if rand.Float32() < 0.2 {
			var nd int
			if dir == 2 || dir == 0 {
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
			PropagateRune(level, x, y, nd, life-1)
		}
	}
}

//Replaces tiles with slopes wherever it makes sense. Removes oddities after terrain deformation.
func (level *Level) SmoothEdges() {
	for j := 0; j < level.rows; j++ {
		for i := 0; i < level.cols; i++ {
			t := level.GetTile(i, j, false)
			ln := level.GetTile(i-1, j, true) //Left neighbor
			lns := ln.IsTerrain()
			rn := level.GetTile(i+1, j, true) //Right neighbor
			rns := rn.IsTerrain()
			tn := level.GetTile(i, j-1, true) //Top neighbor
			tns := tn.IsTerrain()
			bn := level.GetTile(i, j+1, true) //Bottom neighbor
			bns := bn.IsTerrain()
			if t.tt == TT_BLOCK {
				//Turn poking structures into tentacles
				if bns && !tns && !rns && !lns && bn.tt == TT_BLOCK {
					level.SetTile(i, j, TT_TENTACLE_UP, false)
				} else if tns && !bns && !rns && !lns && tn.tt == TT_BLOCK {
					level.SetTile(i, j, TT_TENTACLE_DOWN, false)
				} else if lns && !rns && !tns && !bns && ln.tt == TT_BLOCK {
					level.SetTile(i, j, TT_TENTACLE_RIGHT, false)
				} else if rns && !lns && !tns && !bns && rn.tt == TT_BLOCK {
					level.SetTile(i, j, TT_TENTACLE_LEFT, false)
				}
				//Turn into slope?
				if lns && bns && !tns && !rns {
					level.SetTile(i, j, TT_SLOPE_45, false)
				} else if rns && bns && !lns && !tns {
					level.SetTile(i, j, TT_SLOPE_135, false)
				} else if rns && tns && !lns && !bns {
					level.SetTile(i, j, TT_SLOPE_225, false)
				} else if lns && tns && !rns && !bns {
					level.SetTile(i, j, TT_SLOPE_315, false)
				}
			} else if t.tt&TT_TENTACLES > 0 {
				//Remove lonely tentacles
				if !lns && !rns && !tns && !bns {
					level.SetTile(i, j, TT_EMPTY, false)
				}
			}
		}
	}
}

//Ensures that each space in the level is accessible to the player by digging tunnels
func (level *Level) ConnectCaves() {
	//Find the largest space, which is the main "area" of the gameplay
	maxSize := -1
	var mainSpace *Space
	for _, sp := range level.spaces {
		if len(sp.tiles) > maxSize {
			maxSize = len(sp.tiles)
			mainSpace = sp
		}
	}

	//Determine start and endpoints of tunnels from smaller spaces to main space
	for _, space := range level.spaces {
		if space == mainSpace {
			continue
		}

		maxDist := -1.0
		var tileA *Tile //The tile on the space's frontier that is farthest away from the center
		for _, tile := range space.frontier {
			dist := math.Pow(tile.centerX-space.centerX, 2.0) + math.Pow(tile.centerY-space.centerY, 2.0)
			if dist > maxDist {
				maxDist = dist
				tileA = tile
			}
		}

		minDist := math.MaxFloat64
		var tileB *Tile //The tile on the main space's frontier that is closest to our current space
		for _, tile := range mainSpace.frontier {
			dist := math.Pow(tile.centerX-space.centerX, 2.0) + math.Pow(tile.centerY-space.centerY, 2.0)
			if dist < minDist {
				minDist = dist
				tileB = tile
			}
		}

		level.DigTunnel(tileA, tileB)
	}
}

//Empties tiles in a linear-ish path from the start tile to the end tile
func (level *Level) DigTunnel(start *Tile, end *Tile) {
	dx := end.gridX - start.gridX
	dy := end.gridY - start.gridY
	dxCount := int(math.Abs(float64(dx)))
	dyCount := int(math.Abs(float64(dy)))
	var sdx, sdy int
	if dx > 0 {
		sdx = 1
	} else if dx < 0 {
		sdx = -1
	}
	if dy > 0 {
		sdy = 1
	} else if dy < 0 {
		sdy = -1
	}
	currTile := start
	for currTile != end {
		currTile.SetType(TT_EMPTY)

		var direction bool //False for moving in x, true for moving in Y
		if dxCount > 0 && dyCount > 0 {
			direction = rand.Float64() < 0.5
		} else if dyCount > 0 {
			direction = true
		} else if dxCount > 0 {
			direction = false
		}

		if direction {
			currTile = level.GetTile(currTile.gridX, currTile.gridY+sdy, true)
			dyCount--
			//Add some unevenness
			if rand.Float64() < 0.25 {
				level.SetTile(currTile.gridX+rand.Intn(3)-1, currTile.gridY, TT_EMPTY, true)
			}
		} else {
			currTile = level.GetTile(currTile.gridX+sdx, currTile.gridY, true)
			dxCount--
			//Add some unevenness
			if rand.Float64() < 0.25 {
				level.SetTile(currTile.gridX, currTile.gridY+rand.Intn(3)-1, TT_EMPTY, true)
			}
		}
	}
}

func GenerateLevel(w, h int) *Level {
	level := NewLevel(w, h)

	//Generate blobs
	for k := 0; k < w*h/32; k++ {
		x, y := rand.Intn(w), rand.Intn(h)
		PropagateBlob(level, x, y, 1.0)
	}

	//Remove random little holes
	for j := 0; j < h; j++ {
		for i := 0; i < w; i++ {
			lns := level.GetTile(i-1, j, true).IsSolid() //Left neighbor
			rns := level.GetTile(i+1, j, true).IsSolid() //Right neighbor
			tns := level.GetTile(i, j-1, true).IsSolid() //Top neighbor
			bns := level.GetTile(i, j+1, true).IsSolid() //Bottom neighbor
			if lns && rns && tns && bns {
				level.SetTile(i, j, TT_BLOCK, false)
			}
		}
	}

	//Add rune bars
	for i := 0; i < w*h/720; i++ {
		t := level.FindFullSpace(0)
		for j := 0; j < 4; j++ {
			PropagateRune(level, t.gridX, t.gridY, j, 4)
		}
	}

	level.FindSpaces()

	//Add pylons
	for i := 0; i < w*h/48; i++ {
		pylonSpawn := level.FindSpawnPoint()
		x, y := pylonSpawn.gridX, pylonSpawn.gridY
		if level.GetTile(x-1, y, true).tt != TT_PYLON &&
			level.GetTile(x+1, y, true).tt != TT_PYLON &&
			level.GetTile(x, y-1, true).tt != TT_PYLON &&
			level.GetTile(x, y+1, true).tt != TT_PYLON {
			pylonSpawn.SetType(TT_PYLON)
		}
	}

	level.FindSpaces()
	level.ConnectCaves()
	level.SmoothEdges()

	return level
}
