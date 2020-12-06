package main

type ColType int32

const (
	CT_NONE = iota
	CT_PLAYER
	CT_ENEMY
	CT_BULLET
	CT_ITEM
)

type Component interface {
	Update(game *Game, obj *Object)
}

//Object ...
type Object struct {
	pos        Vec2f
	radius     float64
	colType    ColType
	sprites    []*Sprite
	components []Component
}
