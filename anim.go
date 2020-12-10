package main

type Anim struct {
	frames   []*Sprite
	rate     float64 //Frames per second
	loop     bool
	timer    float64
	frame    int
	callback func(*Anim) //Called when the animation is finished
}

func (anim *Anim) GetSprite() *Sprite {
	return anim.frames[anim.frame]
}

func (anim *Anim) Update(deltaTime float64) {
	anim.timer += deltaTime
	if anim.timer > anim.rate {
		anim.timer = 0.0
		anim.frame += 1
		if anim.frame >= len(anim.frames) {
			if anim.callback != nil {
				anim.callback(anim)
			}
			if anim.loop {
				anim.frame = 0
			} else {
				anim.frame -= 1
			}
		}
	}
}
