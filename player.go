package main

import (
	"github.com/hajimehoshi/ebiten"
)

const (
	MAX_LOVE = 100
)

type PlayerState int32

const (
	PS_NORMAL = iota
	PS_HURT
	PS_ASCENDED
)

type Player struct {
	*Actor
	love  int
	state PlayerState
}

func MakePlayer(game *Game, x, y float64) *Player {
	player := &Player{
		Actor: NewActor(120.0, 500_000.0, 50_000.0),
		love:  0,
		state: PS_NORMAL,
	}

	game.objects.PushBack(&Object{
		pos: &Vec2f{x, y}, radius: 8.0, colType: CT_PLAYER,
		sprites: []*Sprite{
			{
				src:    Rect{0, 0, 16, 16},
				ofs:    ZeroVec(),
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

	var dx, dy float64
	//Movement
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
