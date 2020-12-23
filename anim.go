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
		if anim.callback != nil {
			anim.callback(anim)
		}
		if anim.frame >= len(anim.frames) {
			if anim.loop {
				anim.frame = 0
			} else {
				anim.frame -= 1
				anim.finished = true
			}
		}
	}
}
