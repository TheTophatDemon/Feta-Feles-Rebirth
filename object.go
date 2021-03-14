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
	"github.com/hajimehoshi/ebiten"
	"github.com/thetophatdemon/Feta-Feles-Remastered/vmath"
)

type ColType int32

const (
	CT_NONE       ColType = 0
	CT_PLAYER     ColType = 1
	CT_ENEMY      ColType = 1 << 1
	CT_SHOT       ColType = 1 << 2
	CT_PLAYERSHOT ColType = 1 << 3
	CT_ENEMYSHOT  ColType = 1 << 4
	CT_BOUNCYSHOT ColType = 1 << 5
	CT_ITEM       ColType = 1 << 6
	CT_CAT        ColType = 1 << 7
	CT_EXPLOSION  ColType = 1 << 8
	CT_BARREL     ColType = 1 << 9
)

type Component interface {
	Update(game *Game, obj *Object)
}

type Collidable interface {
	OnCollision(game *Game, obj *Object, other *Object)
}

//Object ...
type Object struct {
	pos          *vmath.Vec2f
	radius       float64
	colType      ColType
	sprites      []*Sprite
	components   []Component
	drawPriority int
	removeMe     bool
	hidden       bool
}

func (obj *Object) Intersects(other *Object) bool {
	return obj.pos.Clone().Sub(other.pos).Length() < obj.radius+other.radius
}

func (obj *Object) HasColType(target ColType) bool {
	return (obj.colType & target) > 0
}

func (obj *Object) DrawAllSprites(screen *ebiten.Image, pt *ebiten.GeoM) {
	var objT ebiten.GeoM
	if pt != nil {
		objT = *pt
	}
	objT.Translate(obj.pos.X, obj.pos.Y)
	for _, sp := range obj.sprites {
		sp.Draw(screen, &objT)
	}
}

//Helper for allowing objects to keep track of how many of them are on the playing field
type ObjCtr struct {
	count int
}

func (ctr *ObjCtr) HandleSignal(kind Signal, src interface{}, params map[string]interface{}) {
	if kind == SIGNAL_GAME_INIT {
		ctr.count = 0
	}
}

func NewObjCtr() *ObjCtr {
	ctr := &ObjCtr{}
	Listen_Signal(SIGNAL_GAME_INIT, ctr)
	return ctr
}

func (ctr *ObjCtr) Inc() {
	ctr.count++
}

func (ctr *ObjCtr) Dec() {
	ctr.count--
}
