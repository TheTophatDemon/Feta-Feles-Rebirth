/*
Copyright (C) 2021 Alexander Lunsford

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package main

type Anim struct {
	frames   []*Sprite
	speed    float64 //Frames per second
	loop     bool
	timer    float64
	frame    int
	finished bool        //True when it is not looped and has reached the end
	callback func(*Anim) //Called when the frame is changed
}

func (anim *Anim) GetSprite() *Sprite {
	return anim.frames[anim.frame]
}

func (anim *Anim) Update(deltaTime float64) {
	anim.timer += deltaTime
	if anim.timer > anim.speed {
		anim.timer = 0.0
		anim.frame += 1
		if anim.frame >= len(anim.frames) {
			if anim.loop {
				anim.frame = 0
			} else {
				anim.frame -= 1
				anim.finished = true
			}
		}
		if anim.callback != nil {
			anim.callback(anim)
		}
	}
}
