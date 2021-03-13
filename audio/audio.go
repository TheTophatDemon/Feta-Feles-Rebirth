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

var MuteSfx bool
var MuteMusic bool

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
		"descend":     assets.WAV_DESCEND,
		"menu":        assets.WAV_MENU,
		"button":      assets.WAV_BUTTON,
	}
	musFiles = map[string]string{
		"mystery":        assets.OGG_MYSTERY,
		"mystery_ingame": assets.OGG_MYSTERY_INGAME,
		"hope":           assets.OGG_HOPE,
		"hope_ingame":    assets.OGG_HOPE_INGAME,
		"malform":        assets.OGG_MALFORM,
		"malform_ingame": assets.OGG_MALFORM_INGAME,
	}
}

const MUS_VOL_SCALE = 0.6

func PlayMusic(name string) {
	nextSong = name
	musFadeTimer = 0.0
	if currSong != name {
		musFade = FADE_OUT
	}
}

func switchSongTo(name string) {
	musData, valid := musFiles[name]
	if valid {
		stream, err := vorbis.Decode(audioContext, assets.ReadCompressedString(musData))
		if err != nil {
			log.Fatal("Cannot decode music file: ", name)
		}
		musPlayer, err = audio.NewPlayer(audioContext, stream)
		if err != nil {
			log.Fatal("Failed to create stream for song: ", name)
		}
		musPlayer.Play()
		musFade = FADE_IN
	} else if musPlayer != nil {
		musPlayer.Close()
		musPlayer = nil
		musFade = FADE_NONE
	}
}

func Update(deltaTime float64) {
	if MuteMusic {
		if musPlayer != nil {
			musPlayer.SetVolume(0.0)
		}
	} else {
		//Handle song transition
		if musFade != FADE_NONE {
			musFadeTimer += deltaTime
			if musPlayer != nil {
				if musFade == FADE_IN {
					musPlayer.SetVolume(math.Min(MUS_VOL_SCALE, musFadeTimer*MUS_VOL_SCALE))
				} else if musFade == FADE_OUT {
					musPlayer.SetVolume(math.Max(0.0, MUS_VOL_SCALE*(FADE_TIME-musFadeTimer)))
				}
			}
			if musFadeTimer > FADE_TIME {
				musFadeTimer = 0.0
				if musFade == FADE_OUT {
					currSong = nextSong
					switchSongTo(currSong)
				} else if musFade == FADE_IN && musPlayer != nil {
					musPlayer.SetVolume(MUS_VOL_SCALE)
					musFade = FADE_NONE
				}
			}
		} else if musPlayer != nil {
			musPlayer.SetVolume(MUS_VOL_SCALE)
		}
		//Perform looping
		if musPlayer != nil && !musPlayer.IsPlaying() {
			musPlayer.Rewind()
			musPlayer.Play()
		}
	}
}

const PLAYERS_PER_SOUND = 8

func PlaySound(name string) {
	PlaySoundVolume(name, 0.5)
}

func PlaySoundVolume(name string, volume float64) {
	if MuteSfx {
		return
	}
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
