package main

import (
	"math/rand"

	"github.com/hajimehoshi/ebiten"
)

var whiteFadeShader *ebiten.Shader
var noiseImg *ebiten.Image

func init() {
	noiseImg = ebiten.NewImage(SCR_WIDTH, SCR_HEIGHT)
	noisePixels := make([]byte, SCR_WIDTH*SCR_HEIGHT*4)
	for i := 0; i < SCR_HEIGHT*SCR_WIDTH; i++ {
		noisePixels[i*4+0] = byte(rand.Intn(255))
		noisePixels[i*4+1] = byte(rand.Intn(255))
		noisePixels[i*4+2] = byte(rand.Intn(255))
		noisePixels[i*4+3] = 255
	}
	noiseImg.ReplacePixels(noisePixels)

	var err error
	whiteFadeShader, err = ebiten.NewShader([]byte(`
		package main

		var Coverage float

		func Fragment(position vec4, texCoord vec2, color vec4) vec4 {
			diffuse := imageSrc0UnsafeAt(texCoord)
			noise := imageSrc1UnsafeAt(texCoord)
			mask := step(1.0 - Coverage, noise.r)
			return min(diffuse + mask, vec4(0.9, 0.9, 0.9, 1.0))
		}
	`))
	if err != nil {
		println(err)
		panic(err)
	}
}
