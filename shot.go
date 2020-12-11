package main

import "image"

type Shot struct {
	vel   *Vec2f  //Velocity
	life  float64 //Time in seconds until it disappears
	enemy bool    //Will this shot hurt the player?
}

const (
	SHOT_COLMASK ColType = CT_ENEMY
)

func AddShot(game *Game, pos, dir *Vec2f, speed float64, enemy bool) *Shot {
	shot := &Shot{
		vel:   dir.Clone().Normalize().Scale(speed),
		life:  5.0,
		enemy: enemy,
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
	obj.pos.Add(shot.vel.Clone().Scale(game.deltaTime))

	shot.life -= game.deltaTime
	if shot.life < 0.0 || game.level.ObjIntersects(obj) {
		obj.removeMe = true
	}
}

func (shot *Shot) OnCollision(game *Game, obj, other *Object) {
	if other.colType&SHOT_COLMASK > 0 {
		obj.removeMe = true
	}
}
