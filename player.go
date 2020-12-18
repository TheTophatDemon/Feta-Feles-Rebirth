package main

import (
	"image"

	"github.com/hajimehoshi/ebiten"
)

const (
	PL_SHOOT_FREQ = 0.2
)

type PlayerState int32

const (
	PS_NORMAL = iota
	PS_HURT
	PS_ASCENDED
)

type Player struct {
	*Actor
	state      PlayerState
	shootTimer float64
	hurtTimer  float64
}

var plSpriteNormal *Sprite
var plSpriteShoot *Sprite
var plSpriteHurt *Sprite
var plSpriteAscended *Sprite

func init() {
	plSpriteNormal = NewSprite(image.Rect(0, 0, 16, 16), &Vec2f{-8.0, -8.0}, false, false, 0)
	plSpriteShoot = NewSprite(image.Rect(16, 0, 32, 16), &Vec2f{-8.0, -8.0}, false, false, 0)
	plSpriteHurt = NewSprite(image.Rect(32, 0, 32+16, 16), &Vec2f{-8.0, -8.0}, false, false, 0)
	plSpriteAscended = NewSprite(image.Rect(48, 0, 64, 16), &Vec2f{-8.0, -8.0}, false, false, 0)
}

func AddPlayer(game *Game, x, y float64) *Object {
	player := &Player{
		Actor: NewActor(120.0, 500_000.0, 50_000.0),
		state: PS_NORMAL,
	}

	obj := &Object{
		pos: &Vec2f{x, y}, radius: 6.0, colType: CT_PLAYER,
		sprites: []*Sprite{
			plSpriteNormal,
		},
		components: []Component{player},
	}
	game.objects.PushBack(obj)

	return obj
}

func (player *Player) Update(game *Game, obj *Object) {
	//Attack
	if player.shootTimer <= 0.0 {
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) || ebiten.IsKeyPressed(ebiten.KeySpace) {
			var dir *Vec2f
			if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
				cx, cy := ebiten.CursorPosition()
				rPos := obj.pos.Clone().Sub(game.camPos).Add(&Vec2f{SCR_WIDTH_H, SCR_HEIGHT_H})
				dir = (&Vec2f{float64(cx), float64(cy)}).Sub(rPos)
			} else if ebiten.IsKeyPressed(ebiten.KeySpace) {
				dir = player.facing.Clone()
			}
			AddShot(game, obj.pos, dir, 240.0, false)
			PlaySound("player_shot")
			player.shootTimer = PL_SHOOT_FREQ
		}
	} else {
		player.shootTimer -= game.deltaTime
	}

	if player.hurtTimer > 0.0 {
		player.hurtTimer -= game.deltaTime
		if int(player.hurtTimer/0.125)%2 == 0 {
			obj.hidden = false
		} else {
			obj.hidden = true
		}
		if player.hurtTimer <= 0.0 {
			player.state = PS_NORMAL
			obj.hidden = false
		}
	}

	//Set sprite
	switch player.state {
	case PS_NORMAL:
		if player.shootTimer > 0.0 {
			obj.sprites[0] = plSpriteShoot
		} else {
			obj.sprites[0] = plSpriteNormal
		}
	case PS_HURT:
		obj.sprites[0] = plSpriteHurt
	case PS_ASCENDED:
		obj.sprites[0] = plSpriteAscended
	}

	//Movement
	var dx, dy float64
	if ebiten.IsKeyPressed(ebiten.KeyUp) || ebiten.IsKeyPressed(ebiten.KeyW) {
		dy = -1.0
	} else if ebiten.IsKeyPressed(ebiten.KeyDown) || ebiten.IsKeyPressed(ebiten.KeyS) {
		dy = 1.0
	}

	if ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		dx = 1.0
	} else if ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		dx = -1.0
	}

	player.Actor.Move(dx, dy)
	player.Actor.Update(game, obj)
}

func (player *Player) OnCollision(game *Game, obj, other *Object) {
	switch other.colType {
	case CT_ITEM:
		if game.AddLoveCounter(1) && player.state != PS_ASCENDED {
			player.state = PS_ASCENDED
		}
	case CT_ENEMY, CT_ENEMYSHOT:
		if player.state == PS_NORMAL {
			player.state = PS_HURT
			player.hurtTimer = 1.0
			game.AddLoveCounter(-10)
			PlaySound("player_hurt")
		}
	}
}
