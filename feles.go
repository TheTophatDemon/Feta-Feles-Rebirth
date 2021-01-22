package main

import (
	"image"
)

type FaceType int

const (
	FACE_SMILE FaceType = iota
	FACE_WINK
	FACE_TALK
	FACE_SCAR
	FACE_SCAR_TALK
	FACE_EMPTY
	FACE_EMPTY_TALK
	FACE_MELTED
	FACE_NONE
)

type BodyType int

const (
	BODY_CAT BodyType = iota
	BODY_HUMAN
	BODY_ANGEL
	BODY_ANGEL2
	BODY_CORRUPTED
	BODY_MELTED
	BODY_HORROR
)

//Assembles the sprites for Feles for the appropriate phase of the story
func MakeFeles(ft FaceType, bt BodyType, pos *Vec2f) *Object {
	const (
		FACE_OFS_X = -24.0
		FACE_OFS_Y = -32.0
	)

	sprites := make([]*Sprite, 0, 8)

	normalTailRect := image.Rect(160, 96, 176, 112)
	doubleTailRect := image.Rect(176, 96, 208, 112)
	quadTailRect := image.Rect(128, 80, 160, 112)
	//====================================
	//Tail (Sprite farthest behind)
	//===============================
	switch bt {
	case BODY_CAT:
		sprites = append(sprites, NewSprite(normalTailRect,
			&Vec2f{FACE_OFS_X + 4.0, FACE_OFS_Y + 32.0}, false, false, 0))
	case BODY_HUMAN, BODY_ANGEL:
		sprites = append(sprites, NewSprite(normalTailRect,
			&Vec2f{FACE_OFS_X + 4.0, FACE_OFS_Y + 52.0}, false, false, 0))
	case BODY_ANGEL2:
		sprites = append(sprites, NewSprite(doubleTailRect,
			&Vec2f{FACE_OFS_X - 16.0, FACE_OFS_Y + 52.0}, false, false, 0))
	case BODY_CORRUPTED:
		sprites = append(sprites, NewSprite(quadTailRect,
			&Vec2f{FACE_OFS_X - 16.0, FACE_OFS_Y + 42.0}, false, false, 0))
		sprites = append(sprites, NewSprite(quadTailRect,
			&Vec2f{FACE_OFS_X + 32.0, FACE_OFS_Y + 42.0}, true, false, 0))
	}

	//==================================
	//Wings
	//=================================
	angelWingRect := image.Rect(160, 64, 192, 80)
	corruptWingRect := image.Rect(160, 80, 192, 96)
	switch bt {
	case BODY_ANGEL, BODY_ANGEL2:
		sprites = append(sprites, NewSprite(angelWingRect,
			&Vec2f{FACE_OFS_X - 16.0, FACE_OFS_Y + 30.0}, false, false, 0))
		sprites = append(sprites, NewSprite(angelWingRect,
			&Vec2f{FACE_OFS_X + 32.0, FACE_OFS_Y + 30.0}, true, false, 0))
	case BODY_CORRUPTED:
		sprites = append(sprites, NewSprite(corruptWingRect,
			&Vec2f{FACE_OFS_X - 16.0, FACE_OFS_Y + 30.0}, false, false, 0))
		sprites = append(sprites, NewSprite(corruptWingRect,
			&Vec2f{FACE_OFS_X + 32.0, FACE_OFS_Y + 30.0}, true, false, 0))
	}

	//=========================
	//Body
	//=========================
	switch bt {
	case BODY_CAT:
		sprites = append(sprites, NewSprite(image.Rect(192, 64, 208, 96), &Vec2f{FACE_OFS_X + 8.0, FACE_OFS_Y + 27.0}, false, false, 0))       //Left half
		sprites = append(sprites, NewSprite(image.Rect(192, 64, 208, 96), &Vec2f{FACE_OFS_X + 8.0 + 16.0, FACE_OFS_Y + 27.0}, true, false, 0)) //Right Half
	case BODY_HUMAN, BODY_ANGEL:
		sprites = append(sprites, NewSprite(image.Rect(192, 0, 208, 64), &Vec2f{FACE_OFS_X + 8.0, FACE_OFS_Y + 28.0}, false, false, 0))       //Left half
		sprites = append(sprites, NewSprite(image.Rect(192, 0, 208, 64), &Vec2f{FACE_OFS_X + 8.0 + 16.0, FACE_OFS_Y + 28.0}, true, false, 0)) //Right Half
	case BODY_ANGEL2:
		sprites = append(sprites, NewSprite(image.Rect(176, 0, 192, 64), &Vec2f{FACE_OFS_X + 8.0, FACE_OFS_Y + 28.0}, false, false, 0))       //Left half
		sprites = append(sprites, NewSprite(image.Rect(176, 0, 192, 64), &Vec2f{FACE_OFS_X + 8.0 + 16.0, FACE_OFS_Y + 28.0}, true, false, 0)) //Right Half
	case BODY_CORRUPTED:
		sprites = append(sprites, NewSprite(image.Rect(160, 0, 176, 64), &Vec2f{FACE_OFS_X + 8.0, FACE_OFS_Y + 28.0}, false, false, 0))       //Left half
		sprites = append(sprites, NewSprite(image.Rect(160, 0, 176, 64), &Vec2f{FACE_OFS_X + 8.0 + 16.0, FACE_OFS_Y + 28.0}, true, false, 0)) //Right Half
	case BODY_MELTED:
		sprites = append(sprites, NewSprite(image.Rect(128, 208, 208, 256), &Vec2f{FACE_OFS_X - 17.0, FACE_OFS_Y + 28.0}, false, false, 0))
	case BODY_HORROR:
		sprites = append(sprites, NewSprite(image.Rect(0, 160, 128, 256), &Vec2f{-60.0, -48.0}, false, false, 0))
	}

	//===================================
	//Face
	//====================================
	if ft != FACE_NONE {
		var faceRect image.Rectangle
		switch ft {
		case FACE_SMILE:
			faceRect = image.Rect(208, 0, 256, 32)
		case FACE_WINK:
			faceRect = image.Rect(208, 32, 256, 64)
		case FACE_TALK:
			faceRect = image.Rect(208, 64, 256, 96)
		case FACE_SCAR:
			faceRect = image.Rect(208, 96, 256, 128)
		case FACE_SCAR_TALK:
			faceRect = image.Rect(208, 128, 256, 160)
		case FACE_EMPTY:
			faceRect = image.Rect(208, 160, 256, 192)
		case FACE_EMPTY_TALK:
			faceRect = image.Rect(208, 192, 256, 224)
		case FACE_MELTED:
			faceRect = image.Rect(208, 224, 256, 256)
		}
		sprites = append(sprites, NewSprite(faceRect, &Vec2f{FACE_OFS_X, FACE_OFS_Y}, false, false, 0))
	}

	return &Object{
		pos:          pos.Clone(),
		colType:      CT_NONE,
		radius:       0.0,
		drawPriority: 666,
		components:   []Component{},
		sprites:      sprites,
	}
}
