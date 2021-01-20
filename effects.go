package main

import (
	"image"
)

type Effect struct {
	anim Anim
}

func (fx *Effect) Update(game *Game, obj *Object) {
	fx.anim.Update(game.deltaTime)
	obj.sprites[0] = fx.anim.GetSprite()
}

var sprExplosion []*Sprite

func init() {
	sprExplosion = make([]*Sprite, 5)
	sprExplosion[0] = NewSprite(image.Rect(16, 144, 32, 160), &Vec2f{-8.0, -8.0}, false, false, 0)
	sprExplosion[1] = NewSprite(image.Rect(32, 128, 64, 160), &Vec2f{-16.0, -16.0}, false, false, 0)
	sprExplosion[2] = NewSprite(image.Rect(64, 112, 112, 160), &Vec2f{-24.0, -24.0}, false, false, 0)
	sprExplosion[3] = NewSprite(image.Rect(112, 112, 160, 160), &Vec2f{-24.0, -24.0}, false, false, 0)
	sprExplosion[4] = NewSprite(image.Rect(160, 112, 208, 160), &Vec2f{-24.0, -24.0}, false, false, 0)
}

func AddExplosion(game *Game, x, y float64) *Object {
	obj := &Object{
		pos:          &Vec2f{x, y},
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
	PlaySound("explode")
	return obj
}

var sprPoof []*Sprite

func init() {
	sprPoof = NewSprites(&Vec2f{-4.0, -4.0}, image.Rect(152, 0, 160, 8), image.Rect(152, 32, 152+8, 32+8))
}

func AddPoof(game *Game, x, y float64) *Object {
	obj := &Object{
		pos:          &Vec2f{x, y},
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
