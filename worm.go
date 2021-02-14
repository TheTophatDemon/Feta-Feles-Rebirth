package main

import (
	"image"
	"math"
	"math/rand"
)

var sprWormHead []*Sprite //First is left facing, second is right facing
var sprWormHeadHurt []*Sprite
var sprWormHeadDie []*Sprite
var sprWormBody []*Sprite
var sprWormBodyDie []*Sprite

func init() {
	sprWormHead = make([]*Sprite, 2)
	sprWormHead[0] = NewSprite(image.Rect(0, 80, 16, 96), &Vec2f{-8.0, -8.0}, false, false, 0)
	sprWormHead[1] = CloneSprite(sprWormHead[0]).Flip(true, false)
	sprWormHeadHurt = make([]*Sprite, 2)
	sprWormHeadHurt[0] = NewSprite(image.Rect(16, 80, 32, 96), &Vec2f{-8.0, -8.0}, false, false, 0)
	sprWormHeadHurt[1] = CloneSprite(sprWormHeadHurt[0]).Flip(true, false)
	sprWormHeadDie = []*Sprite{
		sprWormHead[0], sprWormHeadHurt[0], NewSprite(image.Rect(32, 80, 48, 96), &Vec2f{-8.0, -8.0}, false, false, 0),
	}
	sprWormBody = []*Sprite{
		NewSprite(image.Rect(48, 80, 64, 96), &Vec2f{-8.0, -8.0}, false, false, 0),
		NewSprite(image.Rect(48, 80, 64, 96), &Vec2f{-8.0, -8.0}, false, false, 1),
		NewSprite(image.Rect(48, 80, 64, 96), &Vec2f{-8.0, -8.0}, false, false, 2),
	}
	sprWormBodyDie = []*Sprite{sprWormBody[0], NewSprite(image.Rect(48, 48, 64, 64), &Vec2f{-8.0, -8.0}, false, false, 0)}
}

var wormCtr ObjCtr

func init() {
	wormCtr = *NewObjCtr()
}

const WORM_NSEGS = 6    //Number of body segments
const WORM_QDIST = 12.0 //Distance between queue updates

type Worm struct {
	Mob
	segs                 [WORM_NSEGS]*Object //Body segments, including tail
	segTargets           [WORM_NSEGS]*Vec2f  //Queue of previous head positions that the segments move towards
	enqDistCtr           float64             //Measures distance traveled since last enqueue, up to WORM_QDIST
	segDeathTimer        float64             //Timer for destroying segments in the death animation
	turnSpeed, turnTimer float64
}

func AddWorm(game *Game, x, y float64) (obj *Object, worm *Worm) {
	worm = &Worm{
		Mob: Mob{
			Actor:             NewActor(100.0, 100_000.0, 50_000.0),
			health:            10,
			currAnim:          nil,
			lastSeenPlayerPos: ZeroVec(),
			vecToPlayer:       ZeroVec(),
		},
		turnSpeed: math.Pi,
		turnTimer: rand.Float64() * 5.0,
	}
	dir := RandomDirection()
	worm.Move(dir.x, dir.y)
	for i := 0; i < WORM_NSEGS; i++ {
		spr := sprWormBody[rand.Intn(len(sprWormBody))] //Select random body segment sprite
		effect := &Effect{                              //Effect component is added so segments can be animated
			anim: Anim{frames: []*Sprite{spr}},
		}
		worm.segs[i] = &Object{
			pos: &Vec2f{x, y}, radius: 6.0, colType: CT_ENEMY,
			sprites:    []*Sprite{spr},
			components: []Component{effect},
		}
		game.AddObject(worm.segs[i])
	}
	//Worm code is attached to the head object
	obj = &Object{
		pos: &Vec2f{x, y}, radius: 7.0, colType: CT_ENEMY,
		sprites:    []*Sprite{sprWormHead[0]},
		components: []Component{worm},
	}
	game.AddObject(obj)
	wormCtr.Inc()
	return
}

func (worm *Worm) Update(game *Game, obj *Object) {
	//Update sprites
	if worm.movement.x > 0.0 {
		if worm.hurtTimer > 0.0 || worm.dead {
			obj.sprites[0] = sprWormHeadHurt[1]
		} else {
			obj.sprites[0] = sprWormHead[1]
		}
	} else {
		if worm.hurtTimer > 0.0 || worm.dead {
			obj.sprites[0] = sprWormHeadHurt[0]
		} else {
			obj.sprites[0] = sprWormHead[0]
		}
	}

	if !worm.dead {
		//Occasionally reverse the direction of turning to ensure it doesn't get stuck in circles
		worm.turnTimer -= game.deltaTime
		if worm.turnTimer < 0.0 {
			worm.turnTimer = rand.Float64() * 5.0
			worm.turnSpeed = -worm.turnSpeed
		}
		worm.Wander(game, obj, 64.0, worm.turnSpeed)
		if worm.enqDistCtr > WORM_QDIST {
			worm.enqDistCtr = 0.0
			//Add to front of position queue and shift the rest backward
			for i := WORM_NSEGS - 1; i > 0; i-- {
				worm.segTargets[i] = worm.segTargets[i-1]
			}
			worm.segTargets[0] = obj.pos.Clone()
		}
		//Move body segments towards desired positions
		for i, seg := range worm.segs {
			if seg != nil && worm.segTargets[i] != nil {
				diff := worm.segTargets[i].Clone().Sub(seg.pos)
				mvSpd := worm.Actor.velocity.Length() * game.deltaTime
				if diff.Length() < mvSpd {
					seg.pos.x = worm.segTargets[i].x
					seg.pos.y = worm.segTargets[i].y
				} else {
					seg.pos.Add(diff.Normalize().Scale(mvSpd))
				}
			}
		}
	} else {
		worm.Move(0.0, 0.0)
		//Destroy segments one at a time
		worm.segDeathTimer += game.deltaTime
		if worm.segDeathTimer > 0.25 {
			worm.segDeathTimer = 0.0
			var i int
			for i = WORM_NSEGS - 1; i >= 0; i-- {
				//Find furthest segment not yet being destroyed
				if worm.segs[i] != nil && worm.segs[i].removeMe == false {
					break
				}
			}
			PlaySound("enemy_die")
			//Destroy head when segments are gone
			if i < 0 {
				worm.currAnim = &Anim{
					frames: sprWormHeadDie,
					speed:  0.1,
					callback: func(a *Anim) {
						if a.finished {
							obj.removeMe = true
							wormCtr.Dec()
							AddLove(game, 2, obj.pos.x, obj.pos.y)
						}
					},
				}
				worm.segDeathTimer = -1000.0
			} else {
				segObj := worm.segs[i]
				fx := segObj.components[0].(*Effect)
				fx.anim = Anim{
					frames: sprWormBodyDie,
					speed:  0.1,
					callback: func(a *Anim) {
						if a.finished {
							segObj.removeMe = true
							AddLove(game, 2, segObj.pos.x, segObj.pos.y)
						}
					},
				}
				worm.segs[i] = nil
			}
		}
	}

	displace := obj.pos.Clone()

	worm.Mob.Update(game, obj)
	worm.Actor.Update(game, obj)

	displace.Sub(obj.pos)
	worm.enqDistCtr += displace.Length()
}

func (worm *Worm) OnCollision(game *Game, obj, other *Object) {
	//Do not collide with body segments
	for _, seg := range worm.segs {
		if seg == other {
			goto skip
		}
	}
	worm.Mob.OnCollision(game, obj, other)
skip:
	//Death
	if worm.health <= 0 && !worm.dead {
		worm.dead = true
	}
}
