package main

import (
	"bytes"
	"container/ring"
	"io/ioutil"
	"log"

	"github.com/hajimehoshi/ebiten/audio"
	"github.com/hajimehoshi/ebiten/audio/wav"
	"github.com/thetophatdemon/Feta-Feles-Remastered/assets"
)

var audioContext *audio.Context
var players map[string]*ring.Ring //Contains ring buffers of audio players for each sound effect that is loaded
var audioFiles map[string]string

func init() {
	audioContext = audio.NewContext(44100)
	players = make(map[string]*ring.Ring)
	audioFiles = map[string]string{
		"enemy_die":   assets.WAV_ENEMY_DIE,
		"enemy_hurt":  assets.WAV_ENEMY_HURT,
		"love_get":    assets.WAV_LOVE_GET,
		"player_hurt": assets.WAV_PLAYER_HURT,
		"player_shot": assets.WAV_PLAYER_SHOT,
		"voice":       assets.WAV_VOICE,
		"intro_chime": assets.WAV_INTRO_CHIME,
		"outro_chime": assets.WAV_OUTRO_CHIME,
	}
}

const PLAYERS_PER_SOUND = 8

func PlaySound(name string) {
	buffer, loaded := players[name]
	if !loaded {
		stream, err := wav.Decode(audioContext, bytes.NewReader(assets.Parse(audioFiles[name])))
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
			player.SetVolume(0.5)
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
