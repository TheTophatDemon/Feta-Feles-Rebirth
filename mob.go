package main

import (
	"image"
)

type Mob struct {
	*Actor
	health    int
	love      int
	hurtTimer float64
}

func (mb *Mob) Update(game *Game, obj *Object) {
	if mb.hurtTimer > 0.0 {
		mb.hurtTimer -= game.deltaTime
		if int(mb.hurtTimer/0.125)%2 == 0 {
			obj.hidden = false
		} else {
			obj.hidden = true
		}
		if mb.hurtTimer < 0.0 {
			obj.hidden = false
			mb.hurtTimer = 0.0
		}
	}
}

func (mb *Mob) OnCollision(game *Game, obj *Object, other *Object) {
	if mb.hurtTimer <= 0.0 && other.colType == CT_PLAYERSHOT {
		mb.hurtTimer = 0.5
		mb.health--
	}
}

type Knight Mob

var sprKnightNormal *Sprite
var sprKnightCharge *Sprite
var sprKnightHurt *Sprite
var sprKnightDie []*Sprite

func init() {
	sprKnightNormal = NewSprite(image.Rect(16, 32, 32, 48), &Vec2f{-8.0, -8.0}, false, false, 0)
	sprKnightCharge = NewSprite(image.Rect(0, 32, 16, 48), &Vec2f{-8.0, -8.0}, false, false, 0)
	sprKnightHurt = NewSprite(image.Rect(32, 32, 48, 48), &Vec2f{-8.0, -8.0}, false, false, 0)
	sprKnightDie = NewSprites(&Vec2f{-8.0, -8.0}, image.Rect(32, 32, 48, 48), image.Rect(48, 32, 64, 32))
}

func AddKnight(game *Game, x, y float64) *Knight {
	knight := &Knight{
		Actor:  NewActor(200.0, 1_000_000.0, 10_000.0),
		health: 5,
		love:   3,
	}
	game.objects.PushBack(&Object{
		pos: &Vec2f{x, y}, radius: 8.0, colType: CT_ENEMY,
		sprites:    []*Sprite{sprKnightNormal},
		components: []Component{knight},
	})
	return knight
}

func (kn *Knight) Update(game *Game, obj *Object) {
	kn.Actor.Update(game, obj)
	(*Mob)(kn).Update(game, obj)
}

func (kn *Knight) OnCollision(game *Game, obj, other *Object) {
	(*Mob)(kn).OnCollision(game, obj, other)
}
