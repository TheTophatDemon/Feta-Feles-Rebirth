package main

import (
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

func GenerateLevel(w, h int) *Level {
	level := NewLevel(w, h)

	//Generate blobs
	for k := 0; k < w*h/64; k++ {
		x, y := rand.Intn(w), rand.Intn(h)
		PropagateBlob(level, x, y, 1.0)
	}

	//Smooth edges and create gaps
	for j := 0; j < h; j++ {
		for i := 0; i < w; i++ {
			t := level.GetTile(i, j, false)
			ln := level.GetTile(i-1, j, true) //Left neighbor
			lns := ln.IsSolid()
			rn := level.GetTile(i+1, j, true) //Right neighbor
			rns := rn.IsSolid()
			tn := level.GetTile(i, j-1, true) //Top neighbor
			tns := tn.IsSolid()
			bn := level.GetTile(i, j+1, true) //Bottom neighbor
			bns := bn.IsSolid()
			//Remove random holes
			if lns && rns && tns && bns {
				level.SetTile(i, j, TT_BLOCK, false)
			}
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
			}
		}
	}

	//Add rune bars
	for i := 0; i < w*h/1024; i++ {
		t := level.FindFullSpace(0)
		for j := 0; j < 4; j++ {
			PropagateRune(level, t.gridX, t.gridY, j, 4)
		}
	}

	//Add pylons
	for i := 0; i < w*h/48; i++ {
		pylonSpawn := level.FindEmptySpace(1)
		pylonSpawn.SetType(TT_PYLON)
	}

	return level
}
