package main

import (
	"image"
	"math"
	"math/rand"
)

type Cat struct {
	Mob
	meowTimer float64
}

var sprCatRunLeft []*Sprite
var sprCatRunRight []*Sprite
var sprCatDie []*Sprite

func init() {
	sprCatRunLeft = NewSprites(&Vec2f{-8.0, -8.0}, image.Rect(0, 16, 16, 32), image.Rect(16, 16, 32, 32))
	sprCatRunRight = CloneSprites(sprCatRunLeft)
	for _, spr := range sprCatRunRight {
		spr.Flip(true, false)
	}
	sprCatDie = NewSprites(&Vec2f{-8.0, -8.0}, image.Rect(32, 16, 48, 32), image.Rect(48, 16, 64, 32))
}

func AddCat(game *Game, x, y float64) (*Cat, *Object) {
	cat := &Cat{
		Mob: Mob{
			Actor:  NewActor(120.0, 100_000.0, 75_000.0),
			health: 3,
			currAnim: &Anim{
				frames: sprCatRunLeft,
				speed:  0.1,
				loop:   true,
			},
		},
		meowTimer: rand.Float64() * 5.0,
	}
	obj := &Object{
		pos: &Vec2f{x, y}, radius: 6.0, colType: CT_CAT,
		sprites:    []*Sprite{sprCatRunLeft[0]},
		components: []Component{cat},
	}
	game.AddObject(obj)
	//Move in random direction
	d := RandomDirection()
	cat.Move(d.x, d.y)
	//Spawn poofs
	ang := rand.Float64() * math.Pi * 2.0
	for i := 0.0; i < math.Pi*2.0; i += math.Pi / 4.0 {
		ox := math.Cos(ang+i) * 12.0
		oy := math.Sin(ang+i) * 12.0
		AddPoof(game, x+ox, y+oy)
	}
	return cat, obj
}

func (cat *Cat) Update(game *Game, obj *Object) {
	if !cat.dead {
		hit, normal, _ := game.level.SphereIntersects(obj.pos.Clone().Add(cat.velocity.Clone().Scale(game.deltaTime*4.0)), obj.radius)
		if hit {
			reflect := normal.Clone()
			reflect.Add((&Vec2f{normal.y, -normal.x}).Scale((rand.Float64() * 2.0) - 1.0))
			reflect.Normalize()
			cat.Move(reflect.x, reflect.y)
		}

		cat.meowTimer += game.deltaTime
		if cat.meowTimer > 5.0 {
			game.PlaySoundAttenuated("cat_meow", obj.pos.x, obj.pos.y, 256.0)
			cat.meowTimer = 0.0
		}

		//Flip the sprites in the animation to match movement direction
		if cat.currAnim != nil {
			if cat.movement.x > 0 {
				cat.currAnim.frames = sprCatRunRight
			} else {
				cat.currAnim.frames = sprCatRunLeft
			}
		}
	} else {
		cat.Move(0.0, 0.0)
	}

	cat.Mob.Update(game, obj)
	cat.Actor.Update(game, obj)
}

var __dudShots int

func (cat *Cat) OnCollision(game *Game, obj, other *Object) {
	//Make the cat immune to non-bouncy shots by skipping the mob's default behavior
	if other.HasColType(CT_BOUNCYSHOT) || !other.HasColType(CT_PLAYERSHOT) {
		cat.Mob.OnCollision(game, obj, other)
		if other.HasColType(CT_CAT) {
			reflect := obj.pos.Clone().Sub(other.pos)
			reflect.Add((&Vec2f{reflect.y, -reflect.x}).Scale((rand.Float64() * 2.0) - 1.0))
			reflect.Normalize()
			cat.Move(reflect.x, reflect.y)
		}
	} else if other.HasColType(CT_PLAYERSHOT) {
		__dudShots++
		if __dudShots%16 == 0 {
			Emit_Signal(SIGNAL_CAT_RULE, obj, nil)
		}
	}

	//Death
	if cat.health <= 0 && !cat.dead {
		cat.Move(0.0, 0.0)
		cat.dead = true
		PlaySound("cat_die")
		cat.currAnim = &Anim{
			frames: sprCatDie,
			speed:  0.5,
			callback: func(anm *Anim) {
				if anm.finished {
					game.BeginEndTransition()
				}
			},
		}
	}
}
