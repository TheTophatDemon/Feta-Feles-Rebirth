package main

import (
	"container/ring"
	"io/ioutil"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/audio"
	"github.com/hajimehoshi/ebiten/audio/wav"
)

var audioContext *audio.Context
var players map[string]*ring.Ring //Contains ring buffers of audio players for each sound effect that is loaded

func init() {
	audioContext = audio.NewContext(44100)
	players = make(map[string]*ring.Ring)
}

const PLAYERS_PER_SOUND = 8

func PlaySound(name string) {
	buffer, loaded := players[name]
	if !loaded {
		//Load sound file if not cached yet
		reader, err := os.Open("assets/" + name + ".wav")
		defer reader.Close()
		if err != nil {
			log.Fatal(err)
			return
		}

		stream, err := wav.Decode(audioContext, reader)
		if err != nil {
			log.Fatal(err)
			return
		}
		bytes, err := ioutil.ReadAll(stream)
		if err != nil {
			log.Fatal(err)
			return
		}

		//Initialize audio players in the ring buffer
		buffer = ring.New(PLAYERS_PER_SOUND)
		for i := 0; i < PLAYERS_PER_SOUND; i++ {
			player := audio.NewPlayerFromBytes(audioContext, bytes)
			buffer.Value = player
			buffer = buffer.Next()
		}
		players[name] = buffer
	}
	for i := 0; i < PLAYERS_PER_SOUND; i++ {
		player := buffer.Value.(*audio.Player)
		if !player.IsPlaying() {
			player.Rewind()
			player.Play()
			break
		}
		buffer = buffer.Next()
	}
}
