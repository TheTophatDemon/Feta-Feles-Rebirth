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

const (
	LOVE_FRICTION = 20_000.0
	LOVE_SPEED    = 120.0
)

type Love struct {
	Actor
	blinkAnim *Anim
	life      float64
}

var sprLoveBlink []*Sprite

func init() {
	sprLoveBlink = NewSprites(vmath.NewVec(-4.0, -4.0), image.Rect(80, 64, 88, 72), image.Rect(88, 64, 96, 72))
}

func AddLove(game *Game, count int, x, y float64) {
	//Since they all spawn at the same time, they can share an anim struct
	anim := &Anim{
		frames: sprLoveBlink,
		speed:  100.0,
		loop:   true,
	}
	angle := rand.Float64() * math.Pi * 2.0
	for i := 0; i < count; i++ {
		lv := &Love{
			Actor:     *NewActor(LOVE_SPEED, 0.0, LOVE_FRICTION),
			blinkAnim: anim,
			life:      6.0,
		}
		lv.velocity = (vmath.NewVec(math.Cos(angle), math.Sin(angle))).Scale(LOVE_SPEED)
		angle += rand.Float64() * math.Pi * 0.666
		game.AddObject(&Object{
			pos: vmath.NewVec(x, y), radius: 4.0, colType: CT_ITEM,
			drawPriority: -1,
			sprites: []*Sprite{
				sprLoveBlink[0],
			},
			components: []Component{lv},
		})
	}
}

func (lv *Love) Update(game *Game, obj *Object) {
	lv.Actor.Update(game, obj)

	lv.blinkAnim.Update(game.deltaTime)
	obj.sprites[0] = lv.blinkAnim.GetSprite()

	lv.life -= game.deltaTime
	if lv.life < 3.0 {
		lv.blinkAnim.speed = 0.5
	}
	if lv.life <= 0.0 {
		obj.removeMe = true
	}
}

func (lv *Love) OnCollision(game *Game, obj, other *Object) {
	if other.HasColType(CT_PLAYER) {
		audio.PlaySound("love_get")
		obj.removeMe = true
	}
}
