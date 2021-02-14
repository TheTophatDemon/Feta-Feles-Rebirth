package main

import (
	"image"
	"math"
	"math/rand"
)

type Gopnik struct {
	Mob
	shootTimer, shootAngle float64
}

var sprGopnikNormal []*Sprite
var sprGopnikHurt *Sprite
var sprGopnikDie []*Sprite

func init() {
	sprGopnikNormal = NewSprites(&Vec2f{-8.0, -8.0}, image.Rect(0, 64, 16, 80), image.Rect(16, 64, 32, 80))
	sprGopnikHurt = NewSprite(image.Rect(32, 64, 48, 80), &Vec2f{-8.0, -8.0}, false, false, 0)
	sprGopnikDie = NewSprites(&Vec2f{-8.0, -8.0}, image.Rect(32, 64, 48, 80), image.Rect(48, 64, 64, 80))
}

var gopnikCtr ObjCtr

func init() {
	gopnikCtr = *NewObjCtr()
}

func AddGopnik(game *Game, x, y float64) *Object {
	gopnik := &Gopnik{
		Mob: Mob{
			Actor:  NewActor(50.0, 10_000.0, 100_000.0),
			health: 6,
			currAnim: &Anim{
				frames: sprGopnikNormal,
				speed:  0.5,
				loop:   true,
			},
			lastSeenPlayerPos: ZeroVec(),
			vecToPlayer:       ZeroVec(),
		},
		shootTimer: rand.Float64()/2.0 + 0.5,
		shootAngle: rand.Float64() * math.Pi * 2.0,
	}
	gopnikCtr.Inc()
	return game.AddObject(&Object{
		pos: &Vec2f{x, y}, radius: 7.0, colType: CT_ENEMY,
		sprites:    []*Sprite{sprGopnikNormal[0]},
		components: []Component{gopnik},
	})
}

func (gp *Gopnik) Update(game *Game, obj *Object) {
	gp.Mob.Update(game, obj)
	gp.Actor.Update(game, obj)

	if gp.hurtTimer > 0.0 {
		obj.sprites[0] = sprGopnikHurt
	}

	if gp.hunting {
		gp.shootTimer += game.deltaTime
		if gp.shootTimer > 1.0 {
			gp.shootTimer = 0.0
			for a := 0.0; a < math.Pi*2.0; a += math.Pi / 2.0 {
				AddShot(game, obj.pos.Clone(), VecFromAngle(gp.shootAngle+a, 1.0), 40.0, true)
			}
			gp.shootAngle += math.Pi / 8.0
		}
	}
}

func (gp *Gopnik) OnCollision(game *Game, obj, other *Object) {
	gp.Mob.OnCollision(game, obj, other)

	//Death
	if gp.health <= 0 && !gp.dead {
		gp.dead = true
		PlaySound("enemy_die")
		gp.currAnim = &Anim{
			frames: sprGopnikDie,
			speed:  0.15,
			callback: func(anm *Anim) {
				if anm.finished {
					obj.removeMe = true
					gopnikCtr.Dec()
					AddLove(game, 5, obj.pos.x, obj.pos.y)
				}
			},
		}
	}
}
