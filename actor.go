package main

import "math"

//Component that allows somewhat physically based movement
type Actor struct {
	velocity     *Vec2f
	movement     *Vec2f //Unit vector representing desired movement direction
	maxSpeed     float64
	acceleration float64 //Rate of acceleration in units per seconds squared
	friction     float64 //Rate of deceleration in units per seconds squared
}

func NewActor(maxSpeed, acceleration, friction float64) *Actor {
	return &Actor{
		ZeroVec(),
		ZeroVec(),
		maxSpeed,
		acceleration,
		friction,
	}
}

func (actor *Actor) Update(game *Game, obj *Object) {
	//Accelerate in direction of desired movement
	actor.movement.Normalize()
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

	obj.pos.Add(actor.velocity.Clone().Scale(game.deltaTime))
}

func (actor *Actor) Move(dx, dy float64) {
	actor.movement.x = dx
	actor.movement.y = dy
}
