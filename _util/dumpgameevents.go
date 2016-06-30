package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/astephensen/csgodemogo"
	"github.com/astephensen/csgodemogo/cstrikeproto"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage: %s [demo.dem]\n", os.Args[0])
		os.Exit(2)
	}

	fmt.Println("Dumping game events...")

	demofile := csgodemogo.Open(os.Args[1])

	demofile.GameEventListEmitter = func(gameEventList *cstrikeproto.CSVCMsg_GameEventList) {
		fileOutput := "// Note: This file is automatically generated. Do not edit directly.\n\n"
		fileOutput += "package csgodemogo\n\n"

		// Create the interface function.
		fileOutput += CreateInterfaceFunction(gameEventList)

		// Create the game event structs.
		for _, gameEventListDescriptor := range gameEventList.Descriptors {
			fileOutput += CreateGameEventStruct(gameEventListDescriptor)
		}

		fmt.Println(fileOutput)

		// The event table has been dumped. Exit cleanly.
		os.Exit(0)
	}

	for demofile.Finished == false {
		demofile.GetFrame()
	}
}

func CreateInterfaceFunction(gameEventList *cstrikeproto.CSVCMsg_GameEventList) string {
	interfaceFunction := "func GameEventFromName(eventName string) interface{} {\n"
	interfaceFunction += "\tswitch eventName {\n"

	// Add a case statement for each game event.
	for _, gameEventListDescriptor := range gameEventList.Descriptors {
		interfaceFunction += fmt.Sprintf("\tcase \"%s\":\n", gameEventListDescriptor.GetName())
		interfaceFunction += fmt.Sprintf("\t\treturn GameEvent%s{}\n", TitleCaseUnderscoredString(gameEventListDescriptor.GetName()))
	}

	interfaceFunction += "\t}\n\treturn nil\n}\n"
	return interfaceFunction
}

func CreateGameEventStruct(gameEventListDescriptor *cstrikeproto.CSVCMsg_GameEventListDescriptorT) string {
	gameEventStruct := ""

	// Add the event name.
	eventName := fmt.Sprintf("GameEvent%s", TitleCaseUnderscoredString(gameEventListDescriptor.GetName()))
	gameEventStruct += fmt.Sprintf("type %s struct {\n", eventName)

	// Loop through each of the keys.
	for _, descriptorKey := range gameEventListDescriptor.Keys {
		descriptorKeyName := TitleCaseUnderscoredString(descriptorKey.GetName())
		descriptorKeyType := KeyTypeFromValue(descriptorKey.GetType())
		gameEventStruct += fmt.Sprintf("\t%s\t%s\n", descriptorKeyName, descriptorKeyType)
	}

	gameEventStruct += "}\n\n"

	return gameEventStruct
}

func TitleCaseUnderscoredString(underscoredString string) string {
	splitString := strings.Split(underscoredString, "_")
	for stringComponentIndex, stringComponent := range splitString {
		splitString[stringComponentIndex] = strings.Title(stringComponent)
	}
	joinedString := strings.Join(splitString, "")
	// Just for niceties convert ID and XUID to uppercase.
	if strings.HasSuffix(joinedString, "xuid") {
		joinedString = strings.TrimSuffix(joinedString, "xuid")
		joinedString += "XUID"
	} else if strings.HasSuffix(joinedString, "id") {
		joinedString = strings.TrimSuffix(joinedString, "id")
		joinedString += "ID"
	}
	return joinedString
}

func KeyTypeFromValue(value int32) string {
	switch value {
	case 1:
		return "string"
	case 2:
		return "float64"
	case 3:
		return "int"
	case 4:
		return "int"
	case 5:
		return "int"
	case 6:
		return "bool"
	case 7:
		return "[]byte"
	}
	panic("Encountered unknown value")
}
