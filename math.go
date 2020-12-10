package main

import (
	"math"
)

type Vec2f struct {
	x, y float64
}

func ZeroVec() *Vec2f {
	return &Vec2f{0.0, 0.0}
}

func (vec *Vec2f) Clone() *Vec2f {
	return &Vec2f{x: vec.x, y: vec.y}
}

func (vec *Vec2f) Add(other *Vec2f) *Vec2f {
	vec.x += other.x
	vec.y += other.y
	return vec
}

func (vec *Vec2f) AddScalar(scalar float64) *Vec2f {
	vec.x += scalar
	vec.y += scalar
	return vec
}

func (vec *Vec2f) Sub(other *Vec2f) *Vec2f {
	vec.x -= other.x
	vec.y -= other.y
	return vec
}

func (vec *Vec2f) SubScalar(scalar float64) *Vec2f {
	vec.x -= scalar
	vec.y -= scalar
	return vec
}

func (vec *Vec2f) Floor() *Vec2f {
	vec.x = math.Floor(vec.x)
	vec.y = math.Floor(vec.y)
	return vec
}

func (vec *Vec2f) Ceil() *Vec2f {
	vec.x = math.Ceil(vec.x)
	vec.y = math.Ceil(vec.y)
	return vec
}

func (vec *Vec2f) Normalize() *Vec2f {
	len := vec.Length()
	if len != 0 {
		vec.x /= len
		vec.y /= len
	}
	return vec
}

func (vec *Vec2f) Length() float64 {
	return math.Sqrt(vec.x*vec.x + vec.y*vec.y)
}

func (vec *Vec2f) Scale(s float64) *Vec2f {
	vec.x *= s
	vec.y *= s
	return vec
}

func VecMin(v0, v1 *Vec2f) *Vec2f {
	return &Vec2f{
		x: math.Min(v0.x, v1.x),
		y: math.Min(v0.y, v1.y),
	}
}

func VecMax(v0, v1 *Vec2f) *Vec2f {
	return &Vec2f{
		x: math.Max(v0.x, v1.x),
		y: math.Max(v0.y, v1.y),
	}
}
