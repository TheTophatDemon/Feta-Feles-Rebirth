package main

import "image"

type Shot struct {
	vel     *Vec2f  //Velocity
	life    float64 //Time in seconds until it disappears
	enemy   bool    //Will this shot hurt the player?
	bounces int     //Number of times shot can hit the wall before dying
}

func AddShot(game *Game, pos, dir *Vec2f, speed float64, enemy bool) *Shot {
	return AddBouncyShot(game, pos, dir, speed, enemy, 0)
}

func AddBouncyShot(game *Game, pos, dir *Vec2f, speed float64, enemy bool, bounces int) *Shot {
	shot := &Shot{
		vel:     dir.Clone().Normalize().Scale(speed),
		life:    3.0,
		enemy:   enemy,
		bounces: bounces,
	}
	var rect image.Rectangle
	var ct ColType
	if enemy {
		rect = image.Rect(64, 104, 72, 112)
		ct = CT_ENEMYSHOT
	} else {
		rect = image.Rect(72, 104, 80, 112)
		ct = CT_PLAYERSHOT
	}
	game.objects.PushBack(&Object{
		pos: pos.Clone(), radius: 4.0, colType: ct,
		sprites: []*Sprite{
			NewSprite(rect, &Vec2f{-4.0, -4.0}, false, false, 0),
		},
		components: []Component{shot},
	})
	return shot
}

func (shot *Shot) Update(game *Game, obj *Object) {

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
	if (other.colType == CT_ENEMY && !shot.enemy) || (other.colType == CT_PLAYER && shot.enemy) || other.colType == CT_CAT {
		obj.removeMe = true
	}
}
