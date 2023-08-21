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
	"math/rand"

	"github.com/thetophatdemon/feta-feles-rebirth/audio"
	"github.com/thetophatdemon/feta-feles-rebirth/vmath"
)

type Blargh struct {
	Mob
	shootTimer float64
}

var sprBlarghNormal *Sprite
var sprBlarghShoot *Sprite
var sprBlarghHurt *Sprite
var sprBlarghDie []*Sprite

func init() {
	sprBlarghNormal = NewSprite(image.Rect(0, 48, 16, 64), vmath.NewVec(-8.0, -8.0), false, false, 0)
	sprBlarghShoot = NewSprite(image.Rect(16, 48, 32, 64), vmath.NewVec(-8.0, -8.0), false, false, 0)
	sprBlarghHurt = NewSprite(image.Rect(32, 48, 48, 64), vmath.NewVec(-8.0, -8.0), false, false, 0)
	sprBlarghDie = NewSprites(vmath.NewVec(-8.0, -8.0), image.Rect(32, 48, 48, 64), image.Rect(48, 48, 64, 64))
}

var blarghCtr *ObjCtr

func init() {
	blarghCtr = NewObjCtr()
}

func AddBlargh(game *Game, x, y float64) *Object {
	blargh := &Blargh{
		Mob: Mob{
			Actor:             NewActor(50.0, 100_000.0, 50_000.0),
			health:            5,
			currAnim:          nil,
			lastSeenPlayerPos: vmath.ZeroVec(),
			vecToPlayer:       vmath.ZeroVec(),
		},
		shootTimer: rand.Float64()/2.0 + 0.5,
	}
	blarghCtr.Inc()
	return game.AddObject(&Object{
		pos: vmath.NewVec(x, y), radius: 7.0, colType: CT_ENEMY,
		sprites:    []*Sprite{sprBlarghNormal},
		components: []Component{blargh},
	})
}

const (
	BLARGH_SHOOT_INTERVAL  float64 = 2.0
	BLARGH_SHOOT_THRESHOLD float64 = 0.5 //The actual shot is made _ seconds before the timer reaches 0, so the animation can look better
)

func (bl *Blargh) Update(game *Game, obj *Object) {
	if bl.hurtTimer > 0.0 {
		obj.sprites[0] = sprBlarghHurt
	} else if bl.shootTimer < BLARGH_SHOOT_THRESHOLD {
		obj.sprites[0] = sprBlarghShoot
	} else {
		obj.sprites[0] = sprBlarghNormal
	}

	bl.Mob.Update(game, obj)

	if bl.hunting {
		if bl.shootTimer > BLARGH_SHOOT_THRESHOLD && bl.shootTimer-game.deltaTime < BLARGH_SHOOT_THRESHOLD {
			AddBouncyShot(game, obj.pos.Clone(), bl.vecToPlayer.Clone(), 80.0, true, 2)
		}
		//Move for 0.5 seconds after timer starts
		if bl.shootTimer < BLARGH_SHOOT_INTERVAL-0.5 {
			bl.Move(bl.vecToPlayer.X, bl.vecToPlayer.Y)
		} else {
			bl.Move(0.0, 0.0)
		}
		if bl.shootTimer < BLARGH_SHOOT_INTERVAL {
			bl.shootTimer -= game.deltaTime
			if bl.shootTimer < 0.0 {
				bl.shootTimer = BLARGH_SHOOT_INTERVAL
			}
		} else {
			if bl.seesPlayer { //Restart ticking once player is spotted again
				bl.shootTimer -= game.deltaTime
			}
		}
	}

	bl.Actor.Update(game, obj)
}

func (bl *Blargh) OnCollision(game *Game, obj, other *Object) {
	bl.Mob.OnCollision(game, obj, other)

	//Death
	if bl.health <= 0 && !bl.dead {
		bl.dead = true
		audio.PlaySound("enemy_die")
		bl.currAnim = &Anim{
			frames: sprBlarghDie,
			speed:  0.15,
			callback: func(anm *Anim) {
				if anm.finished {
					obj.removeMe = true
					blarghCtr.Dec()
					AddLove(game, 4, obj.pos.X, obj.pos.Y)
				}
			},
		}
	}
}
