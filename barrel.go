package main

import (
	"image"
)

type Barrel struct {
	health int
}

var sprBarrel *Sprite
var sprBarrelDamaged *Sprite

func init() {
	sprBarrel = NewSprite(image.Rect(16, 128, 32, 144), &Vec2f{-8.0, -8.0}, false, false, 0)
	sprBarrelDamaged = NewSprite(image.Rect(0, 144, 16, 160), &Vec2f{-8.0, -8.0}, false, false, 0)
}

func AddBarrel(game *Game, x, y float64) *Object {
	return game.AddObject(&Object{
		pos: &Vec2f{x, y}, radius: 6.0, colType: CT_BARREL,
		sprites:      []*Sprite{sprBarrel},
		drawPriority: -1,
		components: []Component{&Barrel{
			health: 40,
		}},
	})
}

func (brl *Barrel) Update(game *Game, obj *Object) {
	if brl.health < 20 {
		obj.sprites[0] = sprBarrelDamaged
	} else {
		obj.sprites[0] = sprBarrel
	}
}

func (brl *Barrel) OnCollision(game *Game, obj, other *Object) {
	if brl.health > 0 {
		if other.HasColType(CT_SHOT) {
			brl.health -= 10
		}
		if other.HasColType(CT_EXPLOSION) {
			brl.health -= 20
		}
	}
	if brl.health <= 0 && !obj.removeMe {
		obj.removeMe = true
		AddExplosion(game, obj.pos.x, obj.pos.y)
	}
}
