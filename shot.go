package main

import (
	"image"
)

type Shot struct {
	vel     *Vec2f  //Velocity
	life    float64 //Time in seconds until it disappears
	enemy   bool    //Will this shot hurt the player?
	bounces int     //Number of times shot can hit the wall before dying
	bouncy  bool
	anim    *Anim
}

func AddShot(game *Game, pos, dir *Vec2f, speed float64, enemy bool) *Shot {
	return AddBouncyShot(game, pos, dir, speed, enemy, 0)
}

var sprShotPlayer *Sprite
var sprShotEnemy *Sprite
var sprShotEnemyBouncy []*Sprite
var sprShotPlayerBouncy []*Sprite

func init() {
	sprShotPlayer = NewSprite(image.Rect(72, 104, 80, 112), &Vec2f{-4.0, -4.0}, false, false, 0)
	sprShotEnemy = NewSprite(image.Rect(64, 104, 72, 112), &Vec2f{-4.0, -4.0}, false, false, 0)
	sprShotEnemyBouncy = NewSprites(&Vec2f{-4.0, -4.0}, image.Rect(0, 128, 8, 136), image.Rect(8, 128, 16, 136))
	sprShotPlayerBouncy = NewSprites(&Vec2f{-4.0, -4.0}, image.Rect(0, 136, 8, 144), image.Rect(8, 136, 16, 144))
}

func AddBouncyShot(game *Game, pos, dir *Vec2f, speed float64, enemy bool, bounces int) *Shot {
	shot := &Shot{
		vel:     dir.Clone().Normalize().Scale(speed),
		life:    5.0,
		enemy:   enemy,
		bounces: bounces,
		anim:    nil,
	}
	var spr *Sprite
	ct := CT_SHOT
	//Set animation & Collision
	if bounces > 0 {
		shot.anim = &Anim{
			loop:  true,
			speed: 0.5,
		}
		ct |= CT_BOUNCYSHOT
	}
	if enemy {
		ct |= CT_ENEMYSHOT
		if shot.anim != nil {
			shot.anim.frames = sprShotEnemyBouncy
			spr = shot.anim.frames[0]
		} else {
			spr = sprShotEnemy
		}
	} else {
		ct |= CT_PLAYERSHOT
		if shot.anim != nil {
			shot.anim.frames = sprShotPlayerBouncy
			spr = shot.anim.frames[0]
		} else {
			spr = sprShotPlayer
		}
	}

	game.objects.PushBack(&Object{
		pos: pos.Clone(), radius: 4.0, colType: ct,
		sprites:    []*Sprite{spr},
		components: []Component{shot},
	})
	return shot
}

func (shot *Shot) Update(game *Game, obj *Object) {
	if shot.anim != nil {
		shot.anim.Update(game.deltaTime)
		obj.sprites[0] = shot.anim.GetSprite()
	}
	//Wall bounce
	hit, normal, hitTile := game.level.SphereIntersects(obj.pos.Clone().Add(shot.vel.Clone().Scale(game.deltaTime)), obj.radius)
	if hit {
		if hitTile.tt == TT_RUNE {
			hitTile.SetType(TT_EMPTY)
			AddExplosion(game, hitTile.centerX, hitTile.centerY)
		}
		if shot.bounces > 0 {
			if normal.x != 0.0 || normal.y != 0.0 {
				shot.vel = normal.Scale(shot.vel.Length())
			}
			shot.bounces--
		} else {
			obj.removeMe = true
		}
	}

	obj.pos.Add(shot.vel.Clone().Scale(game.deltaTime))

	shot.life -= game.deltaTime
	if shot.life < 0.0 {
		obj.removeMe = true
	}
}

func (shot *Shot) OnCollision(game *Game, obj, other *Object) {
	if (other.HasColType(CT_ENEMY) && !shot.enemy) || (other.HasColType(CT_PLAYER) && shot.enemy) || other.HasColType(CT_CAT) {
		obj.removeMe = true
	}
}
