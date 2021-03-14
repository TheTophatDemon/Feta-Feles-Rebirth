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
	"github.com/thetophatdemon/Feta-Feles-Remastered/audio"
	"github.com/thetophatdemon/Feta-Feles-Remastered/vmath"
)

type Mob struct {
	*Actor
	health            int
	hurtTimer         float64
	currAnim          *Anim
	dead              bool
	lastSeenPlayerPos *vmath.Vec2f
	vecToPlayer       *vmath.Vec2f
	distToPlayer      float64
	seesPlayer        bool
	hunting           bool //Switched on after monster sees player for the first time
}

func (mb *Mob) Update(game *Game, obj *Object) {
	mb.vecToPlayer = game.playerObj.pos.Clone().Sub(obj.pos)
	mb.distToPlayer = mb.vecToPlayer.Length()
	if raycast := game.level.Raycast(obj.pos.Clone(), mb.vecToPlayer, SCR_HEIGHT); raycast != nil {
		if raycast.distance >= mb.vecToPlayer.Length() {
			mb.lastSeenPlayerPos = game.playerObj.pos.Clone()
			mb.seesPlayer = true
			mb.hunting = true
		} else {
			mb.seesPlayer = false
		}
	}

	if mb.hurtTimer > 0.0 {
		mb.hurtTimer -= game.deltaTime
		if int(mb.hurtTimer/0.125)%2 == 0 {
			obj.hidden = false
		} else {
			obj.hidden = true
		}
		if mb.hurtTimer < 0.0 {
			obj.hidden = false
			mb.hurtTimer = 0.0
		}
	}
	if mb.currAnim != nil {
		mb.currAnim.Update(game.deltaTime)
		obj.sprites[0] = mb.currAnim.GetSprite()
	}
}

func (mb *Mob) OnCollision(game *Game, obj *Object, other *Object) {
	if mb.hurtTimer <= 0.0 && other.HasColType(CT_PLAYERSHOT|CT_EXPLOSION) {
		mb.health--
		if other.HasColType(CT_EXPLOSION) {
			mb.health -= 9
		} else if other.HasColType(CT_BOUNCYSHOT) { //Bouncy shots do double damage
			mb.health--
		}
		if mb.health > 0 {
			mb.hurtTimer = 0.5
			audio.PlaySound("enemy_hurt")
		}
	}
	if other.colType == obj.colType {
		diff := obj.pos.Clone().Sub(other.pos)
		diffL := diff.Length()
		if diffL != 0.0 {
			diff.Normalize()
			mb.velocity.Add(diff.Scale(obj.radius + other.radius/game.deltaTime))
		}
	}
}

//Makes the monster travel around aimlessly
func (mb *Mob) Wander(game *Game, obj *Object, rayDist, turnSpeed float64) {
	//Cast a ray in front of the mob's trajectory
	res := game.level.Raycast(obj.pos.Clone(), mb.movement.Clone(), rayDist)
	if res.hit {
		mb.Turn(turnSpeed, game.deltaTime)
	}
}
