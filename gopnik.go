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

	"github.com/thetophatdemon/Feta-Feles-Remastered/audio"
	"github.com/thetophatdemon/Feta-Feles-Remastered/vmath"
)

type Gopnik struct {
	Mob
	shootTimer, shootAngle float64
}

var sprGopnikNormal []*Sprite
var sprGopnikHurt *Sprite
var sprGopnikDie []*Sprite

func init() {
	sprGopnikNormal = NewSprites(vmath.NewVec(-8.0, -8.0), image.Rect(0, 64, 16, 80), image.Rect(16, 64, 32, 80))
	sprGopnikHurt = NewSprite(image.Rect(32, 64, 48, 80), vmath.NewVec(-8.0, -8.0), false, false, 0)
	sprGopnikDie = NewSprites(vmath.NewVec(-8.0, -8.0), image.Rect(32, 64, 48, 80), image.Rect(48, 64, 64, 80))
}

var gopnikCtr *ObjCtr

func init() {
	gopnikCtr = NewObjCtr()
}

func AddGopnik(game *Game, x, y float64) *Object {
	gopnik := &Gopnik{
		Mob: Mob{
			Actor:  NewActor(50.0, 10_000.0, 100_000.0),
			health: 6,
			currAnim: &Anim{
				frames: sprGopnikNormal,
				speed:  0.5,
				loop:   true,
			},
			lastSeenPlayerPos: vmath.ZeroVec(),
			vecToPlayer:       vmath.ZeroVec(),
		},
		shootTimer: rand.Float64()/2.0 + 0.5,
		shootAngle: rand.Float64() * math.Pi * 2.0,
	}
	gopnikCtr.Inc()
	return game.AddObject(&Object{
		pos: vmath.NewVec(x, y), radius: 7.0, colType: CT_ENEMY,
		sprites:    []*Sprite{sprGopnikNormal[0]},
		components: []Component{gopnik},
	})
}

func (gp *Gopnik) Update(game *Game, obj *Object) {
	gp.Mob.Update(game, obj)
	gp.Actor.Update(game, obj)

	if gp.hurtTimer > 0.0 {
		obj.sprites[0] = sprGopnikHurt
	}

	if gp.hunting {
		gp.shootTimer += game.deltaTime
		if gp.shootTimer > 1.0 {
			gp.shootTimer = 0.0
			for a := 0.0; a < math.Pi*2.0; a += math.Pi / 2.0 {
				AddShot(game, obj.pos.Clone(), vmath.VecFromAngle(gp.shootAngle+a, 1.0), 40.0, true)
			}
			gp.shootAngle += math.Pi / 8.0
		}
	}
}

func (gp *Gopnik) OnCollision(game *Game, obj, other *Object) {
	gp.Mob.OnCollision(game, obj, other)

	//Death
	if gp.health <= 0 && !gp.dead {
		gp.dead = true
		audio.PlaySound("enemy_die")
		gp.currAnim = &Anim{
			frames: sprGopnikDie,
			speed:  0.15,
			callback: func(anm *Anim) {
				if anm.finished {
					obj.removeMe = true
					gopnikCtr.Dec()
					AddLove(game, 5, obj.pos.X, obj.pos.Y)
				}
			},
		}
	}
}
