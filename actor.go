package main

import (
	"math"
)

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
	actor.ApplyMovement(obj, actor.velocity.Clone().Scale(game.deltaTime))
}

func (actor *Actor) ApplyMovement(obj *Object, vel *Vec2f) {
	vel = vel.Clone()
	newPos := obj.pos.Clone().Add(vel)

	//Iterate over portion of the level grid that roughly covers the area between the object and its destination
	gridMin := VecMin(obj.pos, newPos).
		SubScalar(obj.radius).Scale(1.0 / TILE_SIZE).Floor()
	gridMin = VecMax(ZeroVec(), gridMin)
	gridMax := VecMax(obj.pos, newPos).
		AddScalar(obj.radius).Scale(1.0 / TILE_SIZE).Ceil()
	gridMax = VecMin(&Vec2f{x: float64(game.level.cols) - 1.0, y: float64(game.level.rows) - 1.0}, gridMax)

	for j := int(gridMin.y); j < int(gridMax.y); j++ {
		for i := int(gridMin.x); i < int(gridMax.x); i++ {
			if game.level.tiles[j][i]&TT_SOLIDS > 0 {
				//Project object's destination onto the tile boundary
				dest := obj.pos.Clone().Add(vel)
				box := VecMax(
					&Vec2f{x: float64(i) * TILE_SIZE, y: float64(j) * TILE_SIZE},
					VecMin(&Vec2f{x: float64(i+1) * TILE_SIZE, y: float64(j+1) * TILE_SIZE}, dest))
				//debugSpot = box.Clone()
				diff := dest.Clone().Sub(box)
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
	actor.movement.x = dx
	actor.movement.y = dy
}
