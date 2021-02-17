package audio

import (
	"container/ring"
	"io/ioutil"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/audio"
	"github.com/hajimehoshi/ebiten/audio/vorbis"
	"github.com/hajimehoshi/ebiten/audio/wav"
	"github.com/thetophatdemon/Feta-Feles-Remastered/assets"
	"github.com/thetophatdemon/Feta-Feles-Remastered/vmath"
)

var audioContext *audio.Context
var sfxPlayers map[string]*ring.Ring //Contains ring buffers of audio players for each sound effect that is loaded
var sfxFiles map[string]string
var musFiles map[string]string
var musPlayer *audio.Player
var currSong string
var nextSong string

type FadeType int

const (
	FADE_NONE FadeType = iota
	FADE_IN
	FADE_OUT
)

const FADE_TIME float64 = 1.0

var musFade FadeType = FADE_NONE
var musFadeTimer float64 = 0.0

func init() {
	audioContext = audio.NewContext(44100)
	sfxPlayers = make(map[string]*ring.Ring)
	sfxFiles = map[string]string{
		"enemy_die":   assets.WAV_ENEMY_DIE,
		"enemy_hurt":  assets.WAV_ENEMY_HURT,
		"love_get":    assets.WAV_LOVE_GET,
		"player_hurt": assets.WAV_PLAYER_HURT,
		"player_shot": assets.WAV_PLAYER_SHOT,
		"voice":       assets.WAV_VOICE,
		"intro_chime": assets.WAV_INTRO_CHIME,
		"outro_chime": assets.WAV_OUTRO_CHIME,
		"explode":     assets.WAV_EXPLODE,
		"cat_die":     assets.WAV_CAT_DIE,
		"cat_meow":    assets.WAV_CAT_MEOW,
		"evil_voice":  assets.WAV_EVIL_VOICE,
		"ascend":      assets.WAV_ASCEND,
		"roar":        assets.WAV_ROAR,
	}
	musFiles = map[string]string{
		"mystery": assets.OGG_MYSTERY,
		"hope":    assets.OGG_HOPE,
	}
}

func PlayMusic(name string) {
	nextSong = name
	musData, valid := musFiles[name]
	musFadeTimer = 0.0
	if len(name) == 0 || !valid {
		musFade = FADE_OUT
	} else if currSong != name {
		stream, err := vorbis.Decode(audioContext, assets.ReadCompressedString(musData))
		if err != nil {
			log.Fatal("Cannot decode music file: ", name)
		}
		musPlayer, err = audio.NewPlayer(audioContext, stream)
		if err != nil {
			log.Fatal("Failed to create stream for song: ", name)
		}
		musPlayer.Play()
		musPlayer.SetVolume(1.0)
		musFade = FADE_NONE
	}
}

func Update(deltaTime float64) {
	if musFade != FADE_NONE {
		musFadeTimer += deltaTime
		if musFadeTimer > FADE_TIME {
			musFadeTimer = 0.0
			currSong = nextSong
			musFade = FADE_NONE
		}
		if musPlayer != nil {
			switch musFade {
			case FADE_IN:
				musPlayer.SetVolume(musFadeTimer)
			case FADE_OUT:
				musPlayer.SetVolume(FADE_TIME - musFadeTimer)
			}
		}
	}
	if musPlayer != nil && !musPlayer.IsPlaying() {
		musPlayer.Rewind()
		musPlayer.Play()
	}
}

const PLAYERS_PER_SOUND = 8

func PlaySound(name string) {
	PlaySoundVolume(name, 0.5)
}

func PlaySoundVolume(name string, volume float64) {
	buffer, loaded := sfxPlayers[name]
	//Load the sound in if it hasn't been already
	if !loaded {
		stream, err := wav.Decode(audioContext, assets.ReadCompressedString(sfxFiles[name]))
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
			player.SetVolume(volume)
			buffer.Value = player
			buffer = buffer.Next()
		}
		sfxPlayers[name] = buffer
	}
	//Play the sound in the first buffer that isn't already playing
	for i := 0; i < PLAYERS_PER_SOUND; i++ {
		player := buffer.Value.(*audio.Player)
		if !player.IsPlaying() {
			player.Rewind()
			player.SetVolume(volume)
			player.Play()
			break
		} else if player.Current().Seconds() < 0.1 {
			//Abort if sound has already been triggered around the same time
			//Prevents earrape
			return
		}
		buffer = buffer.Next()
	}
}

//Plays a sound that gets quieter the farther it is from the camera
func PlaySoundAttenuated(name string, factor float64, src *vmath.Vec2f, listenerMin, listenerMax *vmath.Vec2f) {
	closestCamPoint := vmath.VecMin(listenerMax, vmath.VecMax(listenerMin, src.Clone()))
	dist := closestCamPoint.Clone().Sub(src).Length()
	PlaySoundVolume(name, math.Max(0.0, math.Min(1.0, 0.5-(dist/factor))))
}
