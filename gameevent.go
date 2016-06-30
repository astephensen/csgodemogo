/*
 * Notes on event key types:
 * 1 = string
 * 2 =
 * 3 = long (int32)
 * 4 = short (int32)
 * 5 = byte (int32)
 * 6 =
 * 7 =
 */

package csgodemogo

import (
	"fmt"
	"reflect"

	"github.com/astephensen/csgodemogo/cstrikeproto"
	"github.com/golang/protobuf/proto"
)

type GameEvent struct {
}

type GameEventRoundStart struct {
	TimeLimit int
	FragLimit int
	Objective string
}

type GameEventRoundEnd struct {
	Winner      int
	Reason      int
	Message     string
	Legacy      int
	PlayerCount int
}

type GameEventPlayerDeath struct {
	UserID                  int
	Attacker                int
	Assister                int
	Weapon                  string
	WeaponItemID            string
	WeaponFauxItemID        string
	WeaponOriginalOwnerXUID string
	Headshot                bool
	Dominated               int
	Revenge                 int
	Penetrated              int
	NoReplay                bool
}

func ParseGameEvent(gameEventList *cstrikeproto.CSVCMsg_GameEventList, buffer []byte) interface{} {
	message := cstrikeproto.CSVCMsg_GameEvent{}
	err := proto.Unmarshal(buffer, &message)
	if err != nil {
		panic(err)
	}

	var gameEvent interface{}
	eventDescriptor := gameEventList.GetEventDescriptor(message.GetEventid())
	eventName := eventDescriptor.GetName()
	switch eventName {
	case "round_start":
		gameEvent = GameEventRoundStart{}
	case "round_end":
		gameEvent = GameEventRoundEnd{}
	case "player_death":
		gameEvent = GameEventPlayerDeath{}
	default:
		return nil
		// Unknown event. Log info
		fmt.Println(eventName)
		for eventKeyIndex, eventKey := range message.Keys {
			fmt.Printf("- %s: %s\n", eventDescriptor.Keys[eventKeyIndex].GetName(), eventKey.String())
		}
	}

	// Use reflection to load the keys into the struct.
	reflectGameEvent := reflect.New(reflect.TypeOf(gameEvent)).Elem()
	for fieldIndex := 0; fieldIndex < reflectGameEvent.NumField(); fieldIndex++ {
		valueField := reflectGameEvent.Field(fieldIndex)

		// Get the key from the event.
		eventValue := message.Keys[fieldIndex]
		switch eventValue.GetType() {
		// 1 = String
		case 1:
			valueField.SetString(eventValue.GetValString())
		// 2 = float (float32)
		case 2:
			valueField.SetFloat(float64(eventValue.GetValFloat()))
		// 3 = long (int32)
		case 3:
			valueField.SetInt(int64(eventValue.GetValLong()))
		// 4 = short (int32)
		case 4:
			valueField.SetInt(int64(eventValue.GetValShort()))
		// 5 = byte (int32)
		case 5:
			valueField.SetInt(int64(eventValue.GetValByte()))
		// 6 = bool
		case 6:
			valueField.SetBool(eventValue.GetValBool())
		// 7 = Binary Data? (???)
		case 7:
		}
	}

	return reflectGameEvent.Interface()
}
