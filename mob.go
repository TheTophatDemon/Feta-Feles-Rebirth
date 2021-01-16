package main

type Mob struct {
	*Actor
	health            int
	hurtTimer         float64
	currAnim          *Anim
	dead              bool
	lastSeenPlayerPos *Vec2f
	vecToPlayer       *Vec2f
	distToPlayer      float64
	seesPlayer        bool
	hunting           bool
}

func (mb *Mob) Update(game *Game, obj *Object) {
	mb.vecToPlayer = game.playerObj.pos.Clone().Sub(obj.pos)
	mb.distToPlayer = mb.vecToPlayer.Length()
	if raycast := game.level.Raycast(obj.pos.Clone(), mb.vecToPlayer, SCR_HEIGHT); raycast != nil {
		if raycast.distance >= mb.vecToPlayer.Length() {
			mb.lastSeenPlayerPos = game.playerObj.pos.Clone()
			mb.seesPlayer = true
			mb.hunting = true
		} else {
			mb.seesPlayer = false
		}
	}

	if mb.hurtTimer > 0.0 {
		mb.hurtTimer -= game.deltaTime
		if int(mb.hurtTimer/0.125)%2 == 0 {
			obj.hidden = false
		} else {
			obj.hidden = true
		}
		if mb.hurtTimer < 0.0 {
			obj.hidden = false
			mb.hurtTimer = 0.0
		}
	}
	if mb.currAnim != nil {
		mb.currAnim.Update(game.deltaTime)
		obj.sprites[0] = mb.currAnim.GetSprite()
	}
}

func (mb *Mob) OnCollision(game *Game, obj *Object, other *Object) {
	if mb.hurtTimer <= 0.0 && other.HasColType(CT_PLAYERSHOT|CT_EXPLOSION) {
		mb.health--
		if other.HasColType(CT_EXPLOSION) {
			mb.health -= 2
		}
		if mb.health > 0 {
			mb.hurtTimer = 0.5
			PlaySound("enemy_hurt")
		} else if !mb.dead {
			PlaySound("enemy_die")
		}
	}
	if other.colType == obj.colType {
		diff := obj.pos.Clone().Sub(other.pos)
		diffL := diff.Length()
		if diffL != 0.0 {
			mb.velocity.Add(diff.Scale((obj.radius + other.radius - diffL) / diffL / game.deltaTime))
		}
	}
}
