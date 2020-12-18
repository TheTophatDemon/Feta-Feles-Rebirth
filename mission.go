package main

type Mission struct {
	loveQuota  int
	knightFreq float64
	blarghFreq float64
	gopnikFreq float64
	wormFreq   float64
	barrelFreq float64
	crossFreq  float64
}

var missions []Mission

func init() {
	missions = []Mission{
		Mission{
			loveQuota:  100,
			knightFreq: 0.5,
			blarghFreq: 0.5,
			gopnikFreq: 0.5,
			wormFreq:   0.5,
			barrelFreq: 0.5,
			crossFreq:  0.0,
		},
	}
}
