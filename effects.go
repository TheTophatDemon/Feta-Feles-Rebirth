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
	"image"
	"math"
	"math/rand"

	"github.com/thetophatdemon/feta-feles-rebirth/audio"
	"github.com/thetophatdemon/feta-feles-rebirth/vmath"
)

type Effect struct {
	anim         Anim
	velocity     *vmath.Vec2f
	acceleration *vmath.Vec2f
}

func (fx *Effect) Update(game *Game, obj *Object) {
	if fx.acceleration != nil {
		fx.velocity.Add(fx.acceleration.Clone().Scale(game.deltaTime))
	}
	if fx.velocity != nil {
		obj.pos.Add(fx.velocity.Clone().Scale(game.deltaTime))
	}

	fx.anim.Update(game.deltaTime)
	obj.sprites[0] = fx.anim.GetSprite()
}

var sprExplosion []*Sprite

func init() {
	sprExplosion = make([]*Sprite, 5)
	sprExplosion[0] = NewSprite(image.Rect(16, 144, 32, 160), vmath.NewVec(-8.0, -8.0), false, false, 0)
	sprExplosion[1] = NewSprite(image.Rect(32, 128, 64, 160), vmath.NewVec(-16.0, -16.0), false, false, 0)
	sprExplosion[2] = NewSprite(image.Rect(64, 112, 112, 160), vmath.NewVec(-24.0, -24.0), false, false, 0)
	sprExplosion[3] = NewSprite(image.Rect(112, 112, 160, 160), vmath.NewVec(-24.0, -24.0), false, false, 0)
	sprExplosion[4] = NewSprite(image.Rect(160, 112, 208, 160), vmath.NewVec(-24.0, -24.0), false, false, 0)
}

func AddExplosion(game *Game, x, y float64) *Object {
	obj := &Object{
		pos:          vmath.NewVec(x, y),
		radius:       8.0,
		colType:      CT_EXPLOSION,
		sprites:      []*Sprite{sprExplosion[0]},
		drawPriority: 20,
	}
	effect := new(Effect)
	effect.anim = Anim{
		frames: sprExplosion,
		speed:  0.1,
		callback: func(anm *Anim) {
			if anm.finished {
				obj.removeMe = true
			} else {
				//Expand collision shape along with sprite
				if anm.frame == 1 {
					obj.radius = 16
				} else if anm.frame > 1 {
					obj.radius = 24
				}
			}
			//Destroy adjacent tiles
			tiles := game.level.GetTilesWithinRadius(obj.pos, obj.radius)
			for _, t := range tiles {
				//Make runes spawn more explosions
				if t.tt == TT_RUNE {
					AddExplosion(game, t.centerX, t.centerY)
				}
				game.level.DestroyTile(t)
			}
		},
	}
	obj.components = []Component{effect}
	game.AddObject(obj)
	audio.PlaySoundAttenuated("explode", 256.0, obj.pos, game.camMin, game.camMax)
	return obj
}

var sprPoof []*Sprite

func init() {
	sprPoof = NewSprites(vmath.NewVec(-4.0, -4.0), image.Rect(80, 48, 88, 56), image.Rect(88, 48, 96, 56))
}

func AddPoof(game *Game, x, y float64) *Object {
	obj := &Object{
		pos:          vmath.NewVec(x, y),
		radius:       0.0,
		colType:      CT_NONE,
		sprites:      []*Sprite{sprPoof[0]},
		drawPriority: 5,
	}
	effect := new(Effect)
	effect.anim = Anim{
		frames: sprPoof,
		speed:  0.1,
		callback: func(anm *Anim) {
			if anm.finished {
				obj.removeMe = true
			}
		},
	}
	obj.components = []Component{effect}
	game.AddObject(obj)
	return obj
}

var sprStars []*Sprite

func init() {
	sprStars = NewSprites(vmath.NewVec(-4.0, -4.0), image.Rect(80, 48, 88, 56), image.Rect(80, 56, 88, 64), image.Rect(88, 56, 96, 64), image.Rect(88, 56, 96, 64), image.Rect(88, 56, 96, 64), image.Rect(88, 48, 96, 56))
}

func AddStarBurst(game *Game, x, y float64) {
	angle := rand.Float64() * math.Pi * 2.0
	for a := 0.0; a < math.Pi*2.0; a += (rand.Float64() * math.Pi / 4.0) + math.Pi/8.0 {
		obj := &Object{
			pos:          vmath.NewVec(x, y),
			radius:       0.0,
			colType:      CT_NONE,
			sprites:      []*Sprite{sprStars[0]},
			drawPriority: 25,
		}
		effect := new(Effect)
		effect.anim = Anim{
			frames: sprStars,
			speed:  0.1,
			callback: func(anm *Anim) {
				if anm.finished {
					obj.removeMe = true
				}
			},
		}
		const SPEED = 50.0
		effect.velocity = vmath.NewVec(math.Cos(angle+a)*SPEED, math.Sin(angle+a)*SPEED)
		effect.acceleration = effect.velocity.Clone().Scale(-0.5)
		obj.components = []Component{effect}
		game.AddObject(obj)
	}
}
