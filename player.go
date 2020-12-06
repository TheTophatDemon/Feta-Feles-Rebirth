package main

import (
	"github.com/hajimehoshi/ebiten"
)

const (
	MAX_LOVE         = 100
	MAX_PLAYER_SPEED = 80.0
	PLAYER_FRICTION  = 512.0
)

type PlayerState int32

const (
	PS_NORMAL = iota
	PS_HURT
	PS_ASCENDED
)

type Player struct {
	love  int
	state PlayerState
	vel   *Vec2f
}

func MakePlayer(game *Game, x, y float64) *Player {
	player := &Player{
		love:  0,
		state: PS_NORMAL,
		vel:   &Vec2f{0.0, 0.0},
	}

	game.objects.PushBack(&Object{
		pos: &Vec2f{x, y}, radius: 8.0, colType: CT_PLAYER,
		sprites: []*Sprite{
			{
				src:    Rect{0, 0, 16, 16},
				ofs:    Vec2f{0.0, 0.0},
				flipH:  false,
				flipV:  false,
				orient: 0,
			},
		},
		components: []Component{player},
	})

	return player
}

func (player *Player) Update(game *Game, obj *Object) {

	d := Vec2f{0.0, 0.0}
	//Movement
	if ebiten.IsKeyPressed(ebiten.KeyUp) || ebiten.IsKeyPressed(ebiten.KeyW) {
		d.y = -1.0
	} else if ebiten.IsKeyPressed(ebiten.KeyDown) || ebiten.IsKeyPressed(ebiten.KeyS) {
		d.y = 1.0
	}

	if ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		d.x = 1.0
	} else if ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		d.x = -1.0
	}

	d.Normalize()
	if d.x != 0.0 || d.y != 0.0 {
		player.vel = d.Clone().Scale(MAX_PLAYER_SPEED * game.deltaTime)
	} else {
		player.vel.Add(player.vel.Clone().Scale(-1.0 * game.deltaTime * game.deltaTime * PLAYER_FRICTION))
	}

	obj.pos.Add(player.vel)

	game.camPos = obj.pos
}
