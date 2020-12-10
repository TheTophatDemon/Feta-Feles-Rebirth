package main

import (
	"image"

	"github.com/hajimehoshi/ebiten"
)

const (
	MAX_LOVE      = 100
	PL_SHOOT_FREQ = 0.33
)

type PlayerState int32

const (
	PS_NORMAL = iota
	PS_HURT
	PS_ASCENDED
)

type Player struct {
	*Actor
	love       int
	state      PlayerState
	shootTimer float64
}

var plSpriteNormal *Sprite
var plSpriteShoot *Sprite

func AddPlayer(game *Game, x, y float64) *Player {
	player := &Player{
		Actor: NewActor(120.0, 500_000.0, 50_000.0),
		love:  0,
		state: PS_NORMAL,
	}

	plSpriteNormal = NewSprite(image.Rect(0, 0, 16, 16), &Vec2f{-8.0, -8.0}, false, false, 0)
	plSpriteShoot = NewSprite(image.Rect(16, 0, 32, 16), &Vec2f{-8.0, -8.0}, false, false, 0)

	game.objects.PushBack(&Object{
		pos: &Vec2f{x, y}, radius: 8.0, colType: CT_PLAYER,
		sprites: []*Sprite{
			plSpriteNormal,
		},
		components: []Component{player},
	})

	return player
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
			player.shootTimer = PL_SHOOT_FREQ
		}
	}

	if player.shootTimer > 0.0 {
		obj.sprites[0] = plSpriteShoot
		player.shootTimer -= game.deltaTime
	} else {
		obj.sprites[0] = plSpriteNormal
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

	game.camPos = obj.pos
}
