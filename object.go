package main

type ColType int32

const (
	CT_NONE       ColType = 0
	CT_PLAYER     ColType = 1
	CT_ENEMY      ColType = 1 << 1
	CT_PLAYERSHOT ColType = 1 << 2
	CT_ENEMYSHOT  ColType = 1 << 3
	CT_ITEM       ColType = 1 << 4
)

type Component interface {
	Update(game *Game, obj *Object)
}

type Collidable interface {
	OnCollision(game *Game, obj *Object, other *Object)
}

//Object ...
type Object struct {
	pos        *Vec2f
	radius     float64
	colType    ColType
	sprites    []*Sprite
	components []Component
	removeMe   bool
}

func (obj *Object) Intersects(other *Object) bool {
	return obj.pos.Clone().Sub(other.pos).Length() < obj.radius+other.radius
}
