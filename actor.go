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
	gridMax = VecMin(&Vec2f{x: float64(game.level.cols), y: float64(game.level.rows)}, gridMax)

	for j := int(gridMin.y); j < int(gridMax.y); j++ {
		for i := int(gridMin.x); i < int(gridMax.x); i++ {
			if game.level.tiles[j][i]&TT_SOLIDS > 0 {
				tileMin := &Vec2f{x: float64(i) * TILE_SIZE, y: float64(j) * TILE_SIZE}
				tileMax := &Vec2f{x: float64(i+1) * TILE_SIZE, y: float64(j+1) * TILE_SIZE}
				dest := obj.pos.Clone().Add(vel)
				//Project object's destination onto the tile boundary
				var proj *Vec2f
				if game.level.tiles[j][i]&TT_SLOPES == 0 {
					//Project onto a box by clamping the destination to the box boundaries
					proj = VecMax(tileMin, VecMin(tileMax, dest))
				} else {
					//Project onto a diagonal plane using the dot product
					tileCenter := tileMin.Clone().AddScalar(TILE_SIZE / 2.0)
					cDiff := dest.Clone().Sub(tileCenter)
					var angle float64
					switch game.level.tiles[j][i] {
					case TT_SLOPE_45:
						angle = math.Pi / 4.0
					case TT_SLOPE_135:
						angle = 3.0 * math.Pi / 4.0
					case TT_SLOPE_225:
						angle = 5.0 * math.Pi / 4.0
					case TT_SLOPE_315:
						angle = 7.0 * math.Pi / 4.0
					}
					angle += math.Pi / 2.0
					normal := &Vec2f{math.Cos(angle), math.Sin(angle)}
					planeDist := VecDot(normal, cDiff)
					proj = dest.Clone().Sub(normal.Clone().Scale(planeDist))
					proj = VecMax(tileMin, VecMin(tileMax, proj))
				}
				//debugSpot = proj.Clone()
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
