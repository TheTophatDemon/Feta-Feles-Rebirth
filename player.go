package main

import (
	"image"
	"strings"

	"github.com/hajimehoshi/ebiten"
)

const (
	MAX_LOVE      = 100
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
	love       int
	state      PlayerState
	shootTimer float64
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

func AddPlayer(game *Game, x, y float64) *Player {
	player := &Player{
		Actor: NewActor(120.0, 500_000.0, 50_000.0),
		love:  0,
		state: PS_NORMAL,
	}

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
			PlaySound("player_shot")
			player.shootTimer = PL_SHOOT_FREQ
		}
	} else {
		player.shootTimer -= game.deltaTime
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

	game.camPos = obj.pos

	//Love cheat code
	if strings.Contains(strings.ToLower(cheatText), "tdnepotis") {
		player.love = game.mission.loveQuota - 1
		cheatText = ""
	}
	//Update UI with love amount
	barRect := image.Rect(game.loveBarBorder.rect.Min.X+3, game.loveBarBorder.rect.Min.Y+3, game.loveBarBorder.rect.Max.X-3, game.loveBarBorder.rect.Max.Y-3)
	barRect.Max.X = barRect.Min.X + int(float64(barRect.Dx())*float64(player.love)/float64(game.mission.loveQuota))
	game.loveBar = SpriteFromScaledImg(game.loveBar.subImg, barRect, 0)

}

func (player *Player) OnCollision(game *Game, obj, other *Object) {
	if other.colType == CT_ITEM {
		player.love++
		if player.love >= game.mission.loveQuota {
			player.love = game.mission.loveQuota
			if player.state != PS_ASCENDED {
				player.state = PS_ASCENDED
				game.OnWin()
			}
		}
	}
}
