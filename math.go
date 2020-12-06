package main

import (
	"image"
	"math"
)

type Rect struct {
	x, y, w, h float64
}

func (r *Rect) ToImgRect() image.Rectangle {
	return image.Rect(int(r.x), int(r.y), int(r.x+r.w), int(r.y+r.h))
}

type Vec2f struct {
	x, y float64
}

func (vec Vec2f) Clone() Vec2f {
	return vec
}

func (vec *Vec2f) Add(other Vec2f) {
	vec.x += other.x
	vec.y += other.y
}

func (vec Vec2f) Added(other Vec2f) Vec2f {
	return Vec2f{vec.x + other.x, vec.y + other.y}
}

func (vec *Vec2f) Sub(other Vec2f) {
	vec.x -= other.x
	vec.y -= other.y
}

func (vec Vec2f) Subbed(other Vec2f) Vec2f {
	return Vec2f{vec.x - other.x, vec.y - other.y}
}

func (vec *Vec2f) Normalize() {
	len := vec.Length()
	if len != 0 {
		vec.x /= len
		vec.y /= len
	}
}

func (vec Vec2f) Normalized() Vec2f {
	vec.Normalize()
	return vec
}

func (vec *Vec2f) Length() float64 {
	return math.Sqrt(vec.x*vec.x + vec.y*vec.y)
}

func (vec *Vec2f) Scale(s float64) {
	vec.x *= s
	vec.y *= s
}

func (vec Vec2f) Scaled(s float64) Vec2f {
	vec.Scale(s)
	return vec
}
