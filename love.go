package main

import (
	"image"
	"math"
	"math/rand"
)

const (
	LOVE_FRICTION = 20_000.0
	LOVE_SPEED    = 120.0
)

type Love struct {
	vel       *Vec2f
	blinkAnim *Anim
	life      float64
}

var sprLoveBlink []*Sprite

func init() {
	sprLoveBlink = NewSprites(&Vec2f{-4.0, -4.0}, image.Rect(64, 96, 64+8, 96+8), image.Rect(64+8, 96, 64+16, 96+8))
}

func AddLove(game *Game, count int, x, y float64) {
	//Since they all spawn at the same time, they can share an anim struct
	anim := &Anim{
		frames: sprLoveBlink,
		speed:  100.0,
		loop:   true,
	}
	angle := rand.Float64() * math.Pi * 2.0
	for i := 0; i < count; i++ {
		lv := &Love{
			vel:       (&Vec2f{math.Cos(angle), math.Sin(angle)}).Scale(LOVE_SPEED),
			blinkAnim: anim,
			life:      6.0,
		}
		angle += rand.Float64() * math.Pi * 0.666
		game.AddObject(&Object{
			pos: &Vec2f{x, y}, radius: 4.0, colType: CT_ITEM,
			drawPriority: -1,
			sprites: []*Sprite{
				sprLoveBlink[0],
			},
			components: []Component{lv},
		})
	}
}

func (lv *Love) Update(game *Game, obj *Object) {
	//Friction
	spd := lv.vel.Length()
	if spd != 0.0 {
		lv.vel.Sub(lv.vel.Clone().Scale(1.0 / spd).Scale(math.Min(spd, LOVE_FRICTION*game.deltaTime*game.deltaTime)))
	}
	//Movement
	obj.pos.Add(lv.vel.Clone().Scale(game.deltaTime))

	lv.blinkAnim.Update(game.deltaTime)
	obj.sprites[0] = lv.blinkAnim.GetSprite()

	lv.life -= game.deltaTime
	if lv.life < 3.0 {
		lv.blinkAnim.speed = 0.5
	}
	if lv.life <= 0.0 {
		obj.removeMe = true
	}
}

func (lv *Love) OnCollision(game *Game, obj, other *Object) {
	if other.colType == CT_PLAYER {
		PlaySound("love_get")
		obj.removeMe = true
	}
}
