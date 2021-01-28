package main

import (
	"math"
)

//Component that allows somewhat physically based movement
type Actor struct {
	velocity     *Vec2f
	movement     *Vec2f //Unit vector representing desired movement direction
	facing       *Vec2f //Represents the last direction the actor faced when moving
	maxSpeed     float64
	acceleration float64 //Rate of acceleration in units per seconds squared
	friction     float64 //Rate of deceleration in units per seconds squared
}

func NewActor(maxSpeed, acceleration, friction float64) *Actor {
	return &Actor{
		ZeroVec(),
		ZeroVec(),
		&Vec2f{0.0, 1.0},
		maxSpeed,
		acceleration,
		friction,
	}
}

func (actor *Actor) Update(game *Game, obj *Object) {
	//Accelerate in direction of desired movement
	actor.velocity.Add(actor.movement.Clone().Scale(game.deltaTime * game.deltaTime * actor.acceleration))

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

func (actor *Actor) ApplyMovement(game *Game, obj *Object, vel *Vec2f) {
	vel = vel.Clone()
	newPos := obj.pos.Clone().Add(vel)

	//Iterate over portion of the level grid that roughly covers the area between the object and its destination
	gridMin, gridMax := game.level.GetGridAreaOverCapsule(obj.pos, newPos, obj.radius, true)

	for j := int(gridMin.y); j < int(gridMax.y); j++ {
		for i := int(gridMin.x); i < int(gridMax.x); i++ {
			t := game.level.GetTile(i, j, false)
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

	obj.pos.Add(vel)
}

func (actor *Actor) Move(dx, dy float64) {
	actor.movement = &Vec2f{dx, dy}
	len := actor.movement.Length()
	if len > 0.001 {
		actor.movement.Scale(1.0 / len) //Normalization
		actor.facing = actor.movement.Clone()
	}
}
