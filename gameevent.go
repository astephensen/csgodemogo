package csgodemogo

import (
	"fmt"
	"reflect"

	"github.com/astephensen/csgodemogo/cstrikeproto"
	"github.com/golang/protobuf/proto"
)

func ParseGameEvent(gameEventList *cstrikeproto.CSVCMsg_GameEventList, buffer []byte) interface{} {
	message := cstrikeproto.CSVCMsg_GameEvent{}
	err := proto.Unmarshal(buffer, &message)
	if err != nil {
		panic(err)
	}

	eventDescriptor := gameEventList.GetEventDescriptor(message.GetEventid())
	eventName := eventDescriptor.GetName()
	gameEvent := GameEventFromName(eventName)

	// Unknown event. Log info
	if gameEvent == nil {
		fmt.Println(eventName)
		for eventKeyIndex, eventKey := range message.Keys {
			fmt.Printf("- %s: %s\n", eventDescriptor.Keys[eventKeyIndex].GetName(), eventKey.String())
		}
		return nil
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
