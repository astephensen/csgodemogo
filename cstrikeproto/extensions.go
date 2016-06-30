package cstrikeproto

import (
	"fmt"
)

func (gameEventList *CSVCMsg_GameEventList) PrintEventTable() {
	for _, eventDescriptor := range gameEventList.Descriptors {
		fmt.Printf("%s - %d\n", eventDescriptor.GetName(), eventDescriptor.GetEventid())
		for _, eventDescriptorKey := range eventDescriptor.Keys {
			fmt.Printf("- %s | %d\n", eventDescriptorKey.GetName(), eventDescriptorKey.GetType())
		}
	}
}

func (gameEventList *CSVCMsg_GameEventList) GetEventDescriptor(eventID int32) *CSVCMsg_GameEventListDescriptorT {
	// This should be some sort of lookup table instead.
	for _, eventDescriptor := range gameEventList.Descriptors {
		if eventDescriptor.GetEventid() == eventID {
			return eventDescriptor
		}
	}
	return nil
}
