package main

import (
	"image"
	"math"
)

type Cat Mob

var sprCatRun []*Sprite
var sprCatDie []*Sprite

func init() {
	sprCatRun = NewSprites(&Vec2f{-8.0, -8.0}, image.Rect(0, 16, 16, 32), image.Rect(16, 16, 32, 32))
	sprCatDie = NewSprites(&Vec2f{-8.0, -8.0}, image.Rect(32, 16, 48, 32), image.Rect(48, 16, 64, 32))
}

func AddCat(game *Game) (*Cat, *Object) {
	cat := &Cat{
		Actor:  NewActor(120.0, 100_000.0, 75_000.0),
		health: 3,
		currAnim: &Anim{
			frames: sprCatRun,
			speed:  0.1,
			loop:   true,
		},
	}
	ix, iy := game.level.FindEmptySpace(2)
	obj := &Object{
		pos: &Vec2f{float64(ix)*TILE_SIZE + 8.0, float64(iy)*TILE_SIZE + 8.0}, radius: 8.0, colType: CT_CAT,
		sprites:    []*Sprite{sprCatRun[0]},
		components: []Component{cat},
	}
	game.AddObject(obj)
	d := RandomDirection()
	if d.x > 0 {
		for _, spr := range sprCatRun {
			spr.Flip(true, false)
		}
	}
	cat.Move(d.x, d.y)
	return cat, obj
}

func (cat *Cat) Update(game *Game, obj *Object) {
	pMov := cat.movement.Clone()

	hit, normal := game.level.SphereIntersects(obj.pos.Clone().Add(cat.velocity.Clone().Scale(game.deltaTime*4.0)), obj.radius)
	if hit {
		normal.Lerp(RandomDirection(), 0.25).Normalize()
		cat.Move(normal.x, normal.y)
	}

	//Flip the sprites in the animation to match movement direction
	if math.Signbit(cat.movement.x) != math.Signbit(pMov.x) {
		for _, spr := range cat.currAnim.frames {
			spr.Flip(true, false)
		}
	}

	cat.Actor.Update(game, obj)
	(*Mob)(cat).Update(game, obj)
}

func (cat *Cat) OnCollision(game *Game, obj, other *Object) {
	(*Mob)(cat).OnCollision(game, obj, other)

	//Death
	if cat.health <= 0 && !cat.dead {
		cat.Move(0.0, 0.0)
		cat.dead = true
		cat.currAnim = &Anim{
			frames: sprCatDie,
			speed:  0.5,
			callback: func(anm *Anim) {
				obj.removeMe = true
			},
		}
	}
}
