package main

import "container/list"

type Signal int

const (
	SIGNAL_PLAYER_MOVED  Signal = iota //Fires once after the player has moved around a bit. For the tutorial mission.
	SIGNAL_PLAYER_SHOT                 //Fires once after the player has shot a few times. For the tutorial mission.
	SIGNAL_PLAYER_ASCEND               //Fires to indicate when the player ascends
	SIGNAL_CAT_RULE                    //Fires when the player tries to shoot the cat without being ascended
	SIGNAL_GAME_START                  //Fires at the start of the game after the intro transition
	SIGNAL_GAME_INIT                   //Fires before the game level is generated
)

type Observer interface {
	HandleSignal(kind Signal, src interface{}, params map[string]interface{})
}

var __observers map[Signal]*list.List

//Retrieve observer map and initialize if needed
func observers() map[Signal]*list.List {
	if __observers == nil {
		__observers = make(map[Signal]*list.List)
	}
	return __observers
}

func Listen_Signal(kind Signal, obs Observer) {
	//Initialize list if haven't already
	lst, ok := observers()[kind]
	if !ok {
		lst = list.New()
		observers()[kind] = lst
	}
	//Exit if observer is already registered
	for itr := lst.Front(); itr != nil; itr = itr.Next() {
		c_obs := itr.Value.(Observer)
		if c_obs == obs {
			println("Observer already added.")
			return
		}
	}
	//Add observer to list
	lst.PushBack(obs)
}

func Emit_Signal(kind Signal, src interface{}, params map[string]interface{}) {
	//Exit when no observers
	lst, ok := observers()[kind]
	if !ok {
		return
	}
	//Make empty map if passed nil
	if params == nil {
		params = make(map[string]interface{})
	}
	//Callback on all listening observers
	for itr := lst.Front(); itr != nil; itr = itr.Next() {
		obs := itr.Value.(Observer)
		obs.HandleSignal(kind, src, params)
	}
}
