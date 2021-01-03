package main

import (
	"image"
	"math/rand"
)

type Blargh struct {
	Mob
	shootTimer float64
}

var sprBlarghNormal *Sprite
var sprBlarghShoot *Sprite
var sprBlarghHurt *Sprite
var sprBlarghDie []*Sprite

func init() {
	sprBlarghNormal = NewSprite(image.Rect(0, 48, 16, 64), &Vec2f{-8.0, -8.0}, false, false, 0)
	sprBlarghShoot = NewSprite(image.Rect(16, 48, 32, 64), &Vec2f{-8.0, -8.0}, false, false, 0)
	sprBlarghHurt = NewSprite(image.Rect(32, 48, 48, 64), &Vec2f{-8.0, -8.0}, false, false, 0)
	sprBlarghDie = NewSprites(&Vec2f{-8.0, -8.0}, image.Rect(32, 48, 48, 64), image.Rect(48, 48, 64, 64))
}

func AddBlargh(game *Game, x, y float64) *Object {
	blargh := &Blargh{
		Mob: Mob{
			Actor:             NewActor(90.0, 100_000.0, 50_000.0),
			health:            5,
			currAnim:          nil,
			lastSeenPlayerPos: ZeroVec(),
			vecToPlayer:       ZeroVec(),
		},
		shootTimer: rand.Float64()/2.0 + 0.5,
	}
	return game.AddObject(&Object{
		pos: &Vec2f{x, y}, radius: 7.0, colType: CT_ENEMY,
		sprites:    []*Sprite{sprBlarghNormal},
		components: []Component{blargh},
	})
}

const (
	BLARGH_SHOOT_COOLDOWN  float64 = 2.0
	BLARGH_SHOOT_THRESHOLD float64 = 0.5
)

func (bl *Blargh) Update(game *Game, obj *Object) {
	if bl.hurtTimer > 0.0 {
		obj.sprites[0] = sprBlarghHurt
	} else if bl.shootTimer < BLARGH_SHOOT_THRESHOLD {
		obj.sprites[0] = sprBlarghShoot
	} else {
		obj.sprites[0] = sprBlarghNormal
	}

	bl.Mob.Update(game, obj)

	if bl.hunting {
		if bl.shootTimer > BLARGH_SHOOT_THRESHOLD && bl.shootTimer-game.deltaTime < BLARGH_SHOOT_THRESHOLD {
			AddBouncyShot(game, obj.pos.Clone(), bl.vecToPlayer.Clone(), 80.0, true, 1)
		}
		if bl.seesPlayer || bl.shootTimer < BLARGH_SHOOT_THRESHOLD {
			bl.shootTimer -= game.deltaTime
			if bl.shootTimer < 0.0 {
				bl.shootTimer = BLARGH_SHOOT_COOLDOWN
			}
		}
	}

	bl.Actor.Update(game, obj)
}

func (bl *Blargh) OnCollision(game *Game, obj, other *Object) {
	bl.Mob.OnCollision(game, obj, other)

	//Death
	if bl.health <= 0 && !bl.dead {
		bl.dead = true
		bl.currAnim = &Anim{
			frames: sprBlarghDie,
			speed:  0.15,
			callback: func(anm *Anim) {
				if anm.finished {
					obj.removeMe = true
					AddLove(game, 5, obj.pos.x, obj.pos.y)
				}
			},
		}
	}
}
