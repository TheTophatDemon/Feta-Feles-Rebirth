package main

import "image"

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
			speed:  0.15,
			callback: func(anm *Anim) {
				obj.removeMe = true
			},
		}
		AddLove(game, 3, obj.pos.x, obj.pos.y)
	}
}
