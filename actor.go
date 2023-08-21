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
	"math"

	"github.com/thetophatdemon/feta-feles-rebirth/vmath"
)

//Component that allows somewhat physically based movement
type Actor struct {
	velocity     *vmath.Vec2f
	movement     *vmath.Vec2f //Unit vector representing desired movement direction
	facing       *vmath.Vec2f //Represents the last direction the actor faced when moving
	maxSpeed     float64
	acceleration float64 //Rate of acceleration in units per seconds squared
	friction     float64 //Rate of deceleration in units per seconds squared
	ignoreBounds bool    //True if actor is to not collide with level boundaries (to allow warping)
}

func NewActor(maxSpeed, acceleration, friction float64) *Actor {
	return &Actor{
		vmath.ZeroVec(),
		vmath.ZeroVec(),
		vmath.NewVec(0.0, 1.0),
		maxSpeed,
		acceleration,
		friction,
		false,
	}
}

func (actor *Actor) Update(game *Game, obj *Object) {
	//Accelerate in direction of desired movement
	actor.velocity.Add(actor.movement.Clone().Scale(game.deltaTime * game.deltaTime * actor.acceleration))
	//Should actually only be multiplying by game.deltaTime once. It's too late to go back and tweak all the accelerations, though.

	//Cap out at maximum speed
	speed := actor.velocity.Length()
	if speed > actor.maxSpeed {
		actor.velocity.Normalize().Scale(actor.maxSpeed)
	}

	//Apply friction
	actor.velocity.Sub(
		actor.velocity.Clone().Normalize().Scale(
			math.Min(game.deltaTime*game.deltaTime*actor.friction, speed)))

	actor.ApplyMovement(game, obj, actor.velocity.Clone().Scale(game.deltaTime))
}

func (actor *Actor) ApplyMovement(game *Game, obj *Object, vel *vmath.Vec2f) {
	vel = vel.Clone()
	newPos := obj.pos.Clone().Add(vel)

	//Iterate over portion of the level grid that roughly covers the area between the object and its destination
	gridMin, gridMax := game.level.GetGridAreaOverCapsule(obj.pos, newPos, obj.radius, true)

	for j := int(gridMin.Y); j < int(gridMax.Y); j++ {
		for i := int(gridMin.X); i < int(gridMax.X); i++ {
			t := game.level.GetTile(i, j, true)
			if t != nil && t.IsSolid() {
				dest := obj.pos.Clone().Add(vel)
				proj := game.level.ProjectPosOntoTile(dest, t)
				diff := dest.Clone().Sub(proj)
				push := obj.radius - diff.Length()
				if push > 0 {
					diff.Normalize().Scale(push)
					vel.Add(diff)
				}
			}
		}
	}

	//Collide against level boundaries
	if !actor.ignoreBounds {
		if obj.pos.X+vel.X-obj.radius < 0 && vel.X < 0.0 {
			vel.X = 0.0
		}
		if obj.pos.X+vel.X+obj.radius > game.level.pixelWidth && vel.X > 0.0 {
			vel.X = 0.0
		}
		if obj.pos.Y+vel.Y-obj.radius < 0 && vel.Y < 0 {
			vel.Y = 0.0
		}
		if obj.pos.Y+vel.Y+obj.radius > game.level.pixelHeight && vel.Y > 0.0 {
			vel.Y = 0.0
		}
	}

	obj.pos.Add(vel)
}

func (actor *Actor) Move(dx, dy float64) {
	actor.movement = vmath.NewVec(dx, dy)
	len := actor.movement.Length()
	if len > 0.001 {
		actor.movement.Scale(1.0 / len) //Normalization
		actor.facing = actor.movement.Clone()
	}
}

//Rotate movement direction by 'da' in radians
func (actor *Actor) Turn(da, deltaTime float64) {
	/*
		[cos -sin] [x] = [xcos-ysin]
		[sin  cos] [y]   [xsin+ycos]
	*/
	angle := da * deltaTime
	ca := math.Cos(angle)
	sa := math.Sin(angle)
	mx := actor.movement.X
	my := actor.movement.Y

	actor.Move((mx*ca)-(my*sa), (mx*sa)+(my*ca))
}
