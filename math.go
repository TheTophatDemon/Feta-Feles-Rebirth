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

func (vec *Vec2f) Sub(other *Vec2f) *Vec2f {
	vec.x -= other.x
	vec.y -= other.y
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
