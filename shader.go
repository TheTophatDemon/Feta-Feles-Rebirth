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

import (
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
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
	noiseImg.WritePixels(noisePixels)

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
