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

	"github.com/thetophatdemon/feta-feles-rebirth/vmath"
)

type Barrel struct {
	health int
}

var sprBarrel *Sprite
var sprBarrelDamaged *Sprite

func init() {
	sprBarrel = NewSprite(image.Rect(16, 128, 32, 144), vmath.NewVec(-8.0, -8.0), false, false, 0)
	sprBarrelDamaged = NewSprite(image.Rect(0, 144, 16, 160), vmath.NewVec(-8.0, -8.0), false, false, 0)
}

var barrelCtr *ObjCtr

func init() {
	barrelCtr = NewObjCtr()
}

func AddBarrel(game *Game, x, y float64) *Object {
	barrelCtr.Inc()
	return game.AddObject(&Object{
		pos: vmath.NewVec(x, y), radius: 6.0, colType: CT_BARREL,
		sprites:      []*Sprite{sprBarrel},
		drawPriority: -1,
		components: []Component{&Barrel{
			health: 40,
		}},
	})
}

func (brl *Barrel) Update(game *Game, obj *Object) {
	if brl.health < 20 {
		obj.sprites[0] = sprBarrelDamaged
	} else {
		obj.sprites[0] = sprBarrel
	}
}

func (brl *Barrel) OnCollision(game *Game, obj, other *Object) {
	if brl.health > 0 {
		if other.HasColType(CT_SHOT) {
			brl.health -= 10
		}
		if other.HasColType(CT_EXPLOSION) {
			brl.health -= 20
		}
	}
	if brl.health <= 0 && !obj.removeMe {
		obj.removeMe = true
		barrelCtr.Dec()
		AddExplosion(game, obj.pos.X, obj.pos.Y)
	}
}
