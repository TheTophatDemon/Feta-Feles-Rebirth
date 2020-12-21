package main

import (
	"math"
	"math/rand"
)

func RectsIntersect(min0, max0, min1, max1 *Vec2f) bool {
	return max0.x > min1.x && min0.x < max1.x && max0.y > min1.y && min0.y < max1.y
}

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

func VecDot(v0, v1 *Vec2f) float64 {
	return (v0.x * v1.x) + (v0.y * v1.y)
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

func RandomDirection() *Vec2f {
	return (&Vec2f{
		rand.Float64() - 0.5,
		rand.Float64() - 0.5,
	}).Normalize()
}

func (vec *Vec2f) Equals(other *Vec2f) bool {
	return vec.x == other.x && vec.y == other.y
}

func (vec *Vec2f) Lerp(other *Vec2f, t float64) *Vec2f {
	vec.x += (other.x - vec.x) * t
	vec.y += (other.y - vec.y) * t
	return vec
}
