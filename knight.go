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

	"github.com/thetophatdemon/Feta-Feles-Remastered/audio"
	"github.com/thetophatdemon/Feta-Feles-Remastered/vmath"
)

type Knight struct {
	Mob
	chargeTimer float64
}

var sprKnightNormal *Sprite
var sprKnightCharge *Sprite
var sprKnightHurt *Sprite
var sprKnightDie []*Sprite

func init() {
	sprKnightNormal = NewSprite(image.Rect(16, 32, 32, 48), vmath.NewVec(-8.0, -8.0), false, false, 0)
	sprKnightCharge = NewSprite(image.Rect(0, 32, 16, 48), vmath.NewVec(-8.0, -8.0), false, false, 0)
	sprKnightHurt = NewSprite(image.Rect(32, 32, 48, 48), vmath.NewVec(-8.0, -8.0), false, false, 0)
	sprKnightDie = NewSprites(vmath.NewVec(-8.0, -8.0), image.Rect(32, 32, 48, 48), image.Rect(48, 32, 64, 48))
}

var knightCtr *ObjCtr

func init() {
	knightCtr = NewObjCtr()
}

func AddKnight(game *Game, x, y float64) *Knight {
	knight := &Knight{
		Mob: Mob{
			Actor:             NewActor(game.mission.knightSpeed, 200_000.0, 25_000.0),
			health:            3,
			currAnim:          nil,
			lastSeenPlayerPos: vmath.ZeroVec(),
			vecToPlayer:       vmath.ZeroVec(),
		},
		chargeTimer: rand.Float64(),
	}
	game.AddObject(&Object{
		pos: vmath.NewVec(x, y), radius: 6.0, colType: CT_ENEMY,
		sprites:    []*Sprite{sprKnightNormal},
		components: []Component{knight},
	})
	knightCtr.Inc()
	return knight
}

func (kn *Knight) Update(game *Game, obj *Object) {
	if kn.hurtTimer > 0.0 {
		obj.sprites[0] = sprKnightHurt
	} else if kn.chargeTimer < 1.0 {
		obj.sprites[0] = sprKnightCharge
	} else {
		obj.sprites[0] = sprKnightNormal
	}

	kn.Mob.Update(game, obj)

	if kn.hunting {
		kn.chargeTimer += game.deltaTime
		if kn.chargeTimer > 2.0 {
			kn.chargeTimer = 0.0
			diff := kn.lastSeenPlayerPos.Clone().Sub(obj.pos)
			kn.Move(diff.X, diff.Y)
		} else if kn.chargeTimer > 0.25 {
			kn.Move(0.0, 0.0)
		}
	} else {
		kn.Move(0.0, 0.0)
	}

	kn.Actor.Update(game, obj)
}

func (kn *Knight) OnCollision(game *Game, obj, other *Object) {
	kn.Mob.OnCollision(game, obj, other)

	//Death
	if kn.health <= 0 && !kn.dead {
		kn.dead = true
		audio.PlaySound("enemy_die")
		kn.currAnim = &Anim{
			frames: sprKnightDie,
			speed:  0.15,
			callback: func(anm *Anim) {
				if anm.finished {
					obj.removeMe = true
					knightCtr.Dec()
					AddLove(game, 3, obj.pos.X, obj.pos.Y)
				}
			},
		}
	}
}
