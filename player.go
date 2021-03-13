package main

import (
	"image"

	"github.com/hajimehoshi/ebiten"
	"github.com/thetophatdemon/Feta-Feles-Remastered/audio"
	"github.com/thetophatdemon/Feta-Feles-Remastered/vmath"
)

const (
	PL_SHOOT_FREQ = 0.2
)

type Player struct {
	*Actor
	hurt, ascended   bool
	shootTimer       float64
	hurtTimer        float64
	lastShootDir     *vmath.Vec2f
	moveAmt, shotAmt int
	warpX, warpY     float64 //Indicates if the player is warping across the map, and from what direction
}

var plSpriteNormal *Sprite
var plSpriteShoot *Sprite
var plSpriteHurt *Sprite
var plSpriteAscended *Sprite

func init() {
	plSpriteNormal = NewSprite(image.Rect(0, 0, 16, 16), vmath.NewVec(-8.0, -8.0), false, false, 0)
	plSpriteShoot = NewSprite(image.Rect(16, 0, 32, 16), vmath.NewVec(-8.0, -8.0), false, false, 0)
	plSpriteHurt = NewSprite(image.Rect(32, 0, 32+16, 16), vmath.NewVec(-8.0, -8.0), false, false, 0)
	plSpriteAscended = NewSprite(image.Rect(48, 0, 64, 16), vmath.NewVec(-8.0, -8.0), false, false, 0)
}

func AddPlayer(game *Game, x, y float64) *Object {
	player := &Player{
		Actor:        NewActor(120.0, 500_000.0, 50_000.0),
		hurt:         false,
		ascended:     false,
		lastShootDir: vmath.ZeroVec(),
	}

	obj := &Object{
		pos: vmath.NewVec(x, y), radius: 6.0, colType: CT_PLAYER,
		sprites: []*Sprite{
			plSpriteNormal,
		},
		components:   []Component{player},
		drawPriority: 10,
	}
	game.AddObject(obj)

	return obj
}

func (player *Player) Update(game *Game, obj *Object) {
	if player.warpX == 0.0 && player.warpY == 0.0 {
		//Attack
		if player.shootTimer <= 0.0 {
			//Set direction
			var dir *vmath.Vec2f
			if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) { //Shoot in direction of mouse click
				cx, cy := ebiten.CursorPosition()
				rPos := obj.pos.Clone().Sub(game.camPos).Add(vmath.NewVec(SCR_WIDTH_H, SCR_HEIGHT_H))
				dir = (vmath.NewVec(float64(cx), float64(cy))).Sub(rPos)
				player.lastShootDir = dir
			} else if ebiten.IsKeyPressed(ebiten.KeySpace) { //Or shoot in direction of last movement
				if player.lastShootDir == nil {
					player.lastShootDir = player.facing.Clone()
				}
				dir = player.lastShootDir
			} else {
				player.lastShootDir = nil
			}
			//Add shot
			if dir != nil {
				if player.ascended {
					AddBouncyShot(game, obj.pos, dir, 240.0, false, 2)
					player.shootTimer = PL_SHOOT_FREQ / 2.0
				} else {
					AddShot(game, obj.pos, dir, 240.0, false)
					player.shootTimer = PL_SHOOT_FREQ
				}
				audio.PlaySound("player_shot")
				player.shotAmt++
				if player.shotAmt == 8 {
					Emit_Signal(SIGNAL_PLAYER_SHOT, obj, nil)
				}
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
				player.hurt = false
				obj.hidden = false
			}
		}

		//Set sprite
		if player.hurt {
			obj.sprites[0] = plSpriteHurt
		} else {
			if player.ascended {
				obj.sprites[0] = plSpriteAscended
			} else {
				if player.shootTimer > 0.0 {
					obj.sprites[0] = plSpriteShoot
				} else {
					obj.sprites[0] = plSpriteNormal
				}
			}
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

		if dx != 0.0 || dy != 0.0 {
			player.moveAmt++
			if player.moveAmt == 100 {
				Emit_Signal(SIGNAL_PLAYER_MOVED, obj, nil)
			}
		}

		player.Actor.Move(dx, dy)
		player.Actor.Update(game, obj)
		//Activate screen wrap
		if player.hurtTimer <= 0.0 { //This indicates that the player hasn't just warped
			newPos := vmath.NewVec(0.0, 0.0)
			switch {
			case obj.pos.X-obj.radius <= 1:
				player.warpX = -1.0
				newPos.X = game.level.pixelWidth - obj.radius
				newPos.Y = obj.pos.Y
			case obj.pos.X+obj.radius >= game.level.pixelWidth-1:
				player.warpX = 1.0
				newPos.X = obj.radius
				newPos.Y = obj.pos.Y
			case obj.pos.Y-obj.radius <= 1:
				player.warpY = -1.0
				newPos.X = obj.pos.X
				newPos.Y = game.level.pixelHeight - obj.radius
			case obj.pos.Y+obj.radius >= game.level.pixelHeight-1:
				player.warpY = 1.0
				newPos.X = obj.pos.X
				newPos.Y = obj.radius
			}
			if player.warpX != 0.0 || player.warpY != 0.0 {
				//Cancel if there's a solid on the other side of the map
				if hit, _, _ := game.level.SphereIntersects(newPos, obj.radius); hit {
					player.warpX, player.warpY = 0.0, 0.0
				} else {
					Emit_Signal(SIGNAL_PLAYER_WARP, obj, map[string]interface{}{})
				}
			}
		}
	} else {
		//Implement screen wrapping effect
		//Move in the original movement direction until the "border" has been cleared, then set position to the other side.
		player.hurtTimer = 1.0 //Activate invincibility in case there's an enemy on the other side of the warp
		if player.warpX < 0.0 && obj.pos.X < -obj.radius {
			obj.pos.X += game.level.pixelWidth
			player.warpX, player.warpY = 0.0, 0.0
		}
		if player.warpX > 0.0 && obj.pos.X > game.level.pixelWidth+obj.radius {
			obj.pos.X -= game.level.pixelWidth
			player.warpX, player.warpY = 0.0, 0.0
		}
		if player.warpY < 0.0 && obj.pos.Y < -obj.radius {
			obj.pos.Y += game.level.pixelHeight
			player.warpX, player.warpY = 0.0, 0.0
		}
		if player.warpY > 0.0 && obj.pos.Y > game.level.pixelHeight+obj.radius {
			obj.pos.Y -= game.level.pixelHeight
			player.warpX, player.warpY = 0.0, 0.0
		}
		obj.pos.X += player.maxSpeed / 4.0 * game.deltaTime * player.warpX
		obj.pos.Y += player.maxSpeed / 4.0 * game.deltaTime * player.warpY
	}
}

func (player *Player) OnCollision(game *Game, obj, other *Object) {
	switch {
	case other.HasColType(CT_ITEM):
		ascend := game.IncLoveCounter(1)
		if ascend {
			if !player.ascended {
				Emit_Signal(SIGNAL_PLAYER_ASCEND, obj, nil)
			}
			player.ascended = true
		}
	case other.HasColType(CT_ENEMY | CT_ENEMYSHOT | CT_EXPLOSION):
		if !player.hurt && player.hurtTimer <= 0.0 {
			player.hurt = true
			player.hurtTimer = 1.0
			lost := false
			if other.colType == CT_EXPLOSION {
				lost = game.DecLoveCounter(20)
			} else {
				if player.ascended {
					lost = game.DecLoveCounter(20)
				} else {
					lost = game.DecLoveCounter(10)
				}
			}
			if lost {
				player.ascended = false
				audio.PlaySound("descend")
			} else {
				audio.PlaySound("player_hurt")
			}
		}
	}
}
