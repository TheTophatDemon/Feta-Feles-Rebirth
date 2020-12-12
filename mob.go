package main

import (
	"image"
)

type Mob struct {
	*Actor
	health    int
	love      int
	hurtTimer float64
	currAnim  *Anim
	dead      bool
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
	if mb.currAnim != nil {
		mb.currAnim.Update(game.deltaTime)
		obj.sprites[0] = mb.currAnim.GetSprite()
	}
}

func (mb *Mob) OnCollision(game *Game, obj *Object, other *Object) {
	if mb.hurtTimer <= 0.0 && other.colType == CT_PLAYERSHOT {
		mb.health--
		if mb.health > 0 {
			mb.hurtTimer = 0.5
			PlaySound("enemy_hurt")
		} else if !mb.dead {
			PlaySound("enemy_die")
		}
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
	sprKnightDie = NewSprites(&Vec2f{-8.0, -8.0}, image.Rect(32, 32, 48, 48), image.Rect(48, 32, 64, 48))
}

func AddKnight(game *Game, x, y float64) *Knight {
	knight := &Knight{
		Actor:    NewActor(200.0, 1_000_000.0, 10_000.0),
		health:   5,
		love:     3,
		currAnim: nil,
	}
	game.objects.PushBack(&Object{
		pos: &Vec2f{x, y}, radius: 8.0, colType: CT_ENEMY,
		sprites:    []*Sprite{sprKnightNormal},
		components: []Component{knight},
	})
	return knight
}

func (kn *Knight) Update(game *Game, obj *Object) {
	if kn.hurtTimer > 0.0 {
		obj.sprites[0] = sprKnightHurt
	} else {
		obj.sprites[0] = sprKnightNormal
	}

	kn.Actor.Update(game, obj)
	(*Mob)(kn).Update(game, obj)
}

func (kn *Knight) OnCollision(game *Game, obj, other *Object) {
	(*Mob)(kn).OnCollision(game, obj, other)

	//Death
	if kn.health <= 0 && !kn.dead {
		kn.dead = true
		kn.currAnim = &Anim{
			frames: sprKnightDie,
			rate:   0.15,
			callback: func(anm *Anim) {
				obj.removeMe = true
			},
		}
	}
}
