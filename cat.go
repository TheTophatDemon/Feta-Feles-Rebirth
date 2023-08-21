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

	"github.com/thetophatdemon/feta-feles-rebirth/audio"
	"github.com/thetophatdemon/feta-feles-rebirth/vmath"
)

type Cat struct {
	Mob
	meowTimer    float64
	stuckTimer   float64
	walkDistance float64
}

var sprCatRunLeft []*Sprite
var sprCatRunRight []*Sprite
var sprCatDie []*Sprite

func init() {
	sprCatRunLeft = NewSprites(vmath.NewVec(-8.0, -8.0), image.Rect(0, 16, 16, 32), image.Rect(16, 16, 32, 32))
	sprCatRunRight = CloneSprites(sprCatRunLeft)
	for _, spr := range sprCatRunRight {
		spr.Flip(true, false)
	}
	sprCatDie = NewSprites(vmath.NewVec(-8.0, -8.0), image.Rect(32, 16, 48, 32), image.Rect(48, 16, 64, 32))
}

func AddCat(game *Game, x, y float64) (*Cat, *Object) {
	cat := &Cat{
		Mob: Mob{
			Actor:  NewActor(120.0, 100_000.0, 75_000.0),
			health: game.mission.catHealth,
			currAnim: &Anim{
				frames: sprCatRunLeft,
				speed:  0.1,
				loop:   true,
			},
		},
		meowTimer: rand.Float64() * 5.0,
	}
	obj := &Object{
		pos: vmath.NewVec(x, y), radius: 6.0, colType: CT_CAT,
		sprites:    []*Sprite{sprCatRunLeft[0]},
		components: []Component{cat},
	}
	game.AddObject(obj)
	//Move in random direction
	d := vmath.RandomDirection()
	cat.Move(d.X, d.Y)
	//Spawn poofs
	ang := rand.Float64() * math.Pi * 2.0
	for i := 0.0; i < math.Pi*2.0; i += math.Pi / 4.0 {
		ox := math.Cos(ang+i) * 12.0
		oy := math.Sin(ang+i) * 12.0
		AddPoof(game, x+ox, y+oy)
	}
	return cat, obj
}

func (cat *Cat) Update(game *Game, obj *Object) {
	//Another fail-safe. Apparently if there are too many cats on screen at once they will occasionally be stuck in NaNspace
	if math.IsNaN(cat.walkDistance) {
		obj.removeMe = true
		spawn := game.level.FindOffscreenSpawnPoint(game)
		AddCat(game, spawn.centerX, spawn.centerY)
	}

	if !cat.dead {
		cat.Wander(game, obj, 64.0, math.Pi)

		cat.meowTimer += game.deltaTime
		if cat.meowTimer > 5.0 {
			audio.PlaySoundAttenuated("cat_meow", 256.0, obj.pos, game.camMin, game.camMax)
			cat.meowTimer = 0.0
		}

		//Flip the sprites in the animation to match movement direction
		if cat.currAnim != nil {
			if cat.movement.X > 0 {
				cat.currAnim.frames = sprCatRunRight
			} else {
				cat.currAnim.frames = sprCatRunLeft
			}
		}
	} else {
		cat.Move(0.0, 0.0)
	}

	//Death
	if cat.health <= 0 && !cat.dead {
		cat.Move(0.0, 0.0)
		cat.dead = true
		audio.PlaySound("cat_die")
		cat.currAnim = &Anim{
			frames: sprCatDie,
			speed:  0.5,
			callback: func(anm *Anim) {
				if anm.finished {
					Emit_Signal(SIGNAL_CAT_DIE, obj, nil)
				}
			},
		}
	}

	walkDiff := obj.pos.Clone()

	cat.Mob.Update(game, obj)
	cat.Actor.Update(game, obj)

	walkDiff.Sub(obj.pos.Clone())
	walkDelta := walkDiff.Length()
	cat.walkDistance += walkDelta //Keep track of how much the cat moves

	//Respawns the cat when it spends 10 seconds without moving much
	//A failsafe in case I haven't actually fixed that elusive bug
	if walkDelta < 8.0/60.0 {
		cat.stuckTimer += game.deltaTime
		if cat.stuckTimer > 5.0 {
			//Spawn poofs
			ang := rand.Float64() * math.Pi * 2.0
			for i := 0.0; i < math.Pi*2.0; i += math.Pi / 4.0 {
				ox := math.Cos(ang+i) * 12.0
				oy := math.Sin(ang+i) * 12.0
				AddPoof(game, obj.pos.X+ox, obj.pos.Y+oy)
			}
			obj.removeMe = true
			spawn := game.level.FindOffscreenSpawnPoint(game)
			AddCat(game, spawn.centerX, spawn.centerY)
		}
	} else {
		cat.stuckTimer = 0.0
	}
}

var __dudShots int

func (cat *Cat) OnCollision(game *Game, obj, other *Object) {
	//Make the cat immune to non-bouncy shots by skipping the mob's default behavior
	if other.HasColType(CT_BOUNCYSHOT) || !other.HasColType(CT_PLAYERSHOT) {
		cat.Mob.OnCollision(game, obj, other)
		if other.HasColType(CT_CAT) {
			reflect := obj.pos.Clone().Sub(other.pos)
			reflect.Add((vmath.NewVec(reflect.Y, -reflect.X)).Scale((rand.Float64() * 2.0) - 1.0))
			reflect.Normalize()
			cat.Move(reflect.X, reflect.Y)
		}
	} else if other.HasColType(CT_PLAYERSHOT) {
		__dudShots++
		if __dudShots%16 == 0 {
			Emit_Signal(SIGNAL_CAT_RULE, obj, nil)
		}
	}
}
