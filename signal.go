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

import "container/list"

type Signal int

const (
	SIGNAL_PLAYER_MOVED  Signal = iota //Fires continuously as long as the player is moving. For the tutorial mission.
	SIGNAL_PLAYER_SHOT                 //Fires every time the player shoots. For the tutorial mission.
	SIGNAL_PLAYER_EDGE				   //Fires when the camera reaches the edge of the map. For the tutorial mission.
	SIGNAL_PLAYER_ASCEND               //Fires to indicate when the player ascends
	SIGNAL_CAT_RULE                    //Fires when the player tries to shoot the cat without being ascended
	SIGNAL_CAT_DIE                     //Fires when the cat is killed
	SIGNAL_GAME_START                  //Fires at the start of the game after the intro transition
	SIGNAL_GAME_INIT                   //Fires before the game level is generated
)

//Represents the number of times a signal has been emitted
var __signal_counts map[Signal]int

//Returns the number of times the given signal has been emitted
func Get_Signal_Count(sig Signal) int {
	return __signal_counts[sig]
}

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
	//Update signal count
	if __signal_counts == nil {
		__signal_counts = make(map[Signal]int)
	}
	_, ok := __signal_counts[kind]
	if !ok {
		__signal_counts[kind] = 1
	} else {
		__signal_counts[kind] += 1
	}
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
