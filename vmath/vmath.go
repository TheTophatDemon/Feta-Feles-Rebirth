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

package vmath

import (
	"math"
	"math/rand"
)

func RectsIntersect(min0, max0, min1, max1 *Vec2f) bool {
	return max0.X > min1.X && min0.X < max1.X && max0.Y > min1.Y && min0.Y < max1.Y
}

type Vec2f struct {
	X, Y float64
}

func NewVec(x, y float64) *Vec2f {
	return &Vec2f{
		X: x, Y: y,
	}
}

func ZeroVec() *Vec2f {
	return &Vec2f{0.0, 0.0}
}

func (vec *Vec2f) Clone() *Vec2f {
	return &Vec2f{X: vec.X, Y: vec.Y}
}

func (vec *Vec2f) Add(other *Vec2f) *Vec2f {
	vec.X += other.X
	vec.Y += other.Y
	return vec
}

func (vec *Vec2f) AddScalar(scalar float64) *Vec2f {
	vec.X += scalar
	vec.Y += scalar
	return vec
}

func (vec *Vec2f) Sub(other *Vec2f) *Vec2f {
	vec.X -= other.X
	vec.Y -= other.Y
	return vec
}

func (vec *Vec2f) SubScalar(scalar float64) *Vec2f {
	vec.X -= scalar
	vec.Y -= scalar
	return vec
}

func (vec *Vec2f) Floor() *Vec2f {
	vec.X = math.Floor(vec.X)
	vec.Y = math.Floor(vec.Y)
	return vec
}

func (vec *Vec2f) Ceil() *Vec2f {
	vec.X = math.Ceil(vec.X)
	vec.Y = math.Ceil(vec.Y)
	return vec
}

func (vec *Vec2f) Normalize() *Vec2f {
	len := vec.Length()
	if len != 0 {
		vec.X /= len
		vec.Y /= len
	}
	return vec
}

func (vec *Vec2f) Length() float64 {
	return math.Sqrt(vec.X*vec.X + vec.Y*vec.Y)
}

func (vec *Vec2f) Scale(s float64) *Vec2f {
	vec.X *= s
	vec.Y *= s
	return vec
}

func VecDot(v0, v1 *Vec2f) float64 {
	return (v0.X * v1.X) + (v0.Y * v1.Y)
}

//Returns magnitude of 3d cross product between two vectors (treated as on the XY plane)
func VecCross(v0, v1 *Vec2f) float64 {
	//(Axi^ + Ayj^) X (Bxi^ + Byj^)
	//Axi^ X Bxi^ + Axi^ X Byj^ + Ayj^ X Bxi^ + Ayj^ X Byj^
	//Axi^ X Byj^ + Ayj^ X Bxi^ = (AxBy)(i^ X j^) + (AyBx)(j^ X i^) = AxByk^ - AyBxk^
	return (v0.X * v1.Y) - (v0.Y * v1.X)
}

func VecMin(v0, v1 *Vec2f) *Vec2f {
	return &Vec2f{
		X: math.Min(v0.X, v1.X),
		Y: math.Min(v0.Y, v1.Y),
	}
}

func VecMax(v0, v1 *Vec2f) *Vec2f {
	return &Vec2f{
		X: math.Max(v0.X, v1.X),
		Y: math.Max(v0.Y, v1.Y),
	}
}

func RandomDirection() *Vec2f {
	return (&Vec2f{
		rand.Float64() - 0.5,
		rand.Float64() - 0.5,
	}).Normalize()
}

func VecFromAngle(angle, magnitude float64) *Vec2f {
	return &Vec2f{
		math.Cos(angle) * magnitude,
		math.Sin(angle) * magnitude,
	}
}

func (vec *Vec2f) Equals(other *Vec2f) bool {
	return vec.X == other.X && vec.Y == other.Y
}

func (vec *Vec2f) Lerp(other *Vec2f, t float64) *Vec2f {
	vec.X += (other.X - vec.X) * t
	vec.Y += (other.Y - vec.Y) * t
	return vec
}
