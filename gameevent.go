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
		case 1: // string
			valueField.SetString(eventValue.GetValString())
		case 2: // float (float32)
			valueField.SetFloat(float64(eventValue.GetValFloat()))
		case 3: // long (int32)
			valueField.SetInt(int64(eventValue.GetValLong()))
		case 4: // short (int32)
			valueField.SetInt(int64(eventValue.GetValShort()))
		case 5: // byte (int32)
			valueField.SetInt(int64(eventValue.GetValByte()))
		case 6: // bool
			valueField.SetBool(eventValue.GetValBool())
		case 7: // binary data - not sure about this one
			valueField.SetBytes(eventValue.GetValWstring())
		}
	}

	return reflectGameEvent.Interface()
}
