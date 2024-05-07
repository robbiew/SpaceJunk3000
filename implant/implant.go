package implant

import (
	"encoding/json"
	"fmt"
	"os"
	"spacejunk3000/door"
	"strconv"
)

// Implant represents the characteristics of a cybernetic implant.
type Implant struct {
	Name string `json:"name"`
	Desc string `json:"desc"`
}

func NewImplant(name, desc string) *Implant {
	return &Implant{
		Name: name,
		Desc: desc,
	}
}

// LoadImplants loads implants from a specified JSON file.
func LoadImplants(filename string) ([]Implant, error) {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var implants []Implant
	err = json.Unmarshal(bytes, &implants)
	if err != nil {
		return nil, err
	}
	return implants, nil
}

func SelectImplant(implants []Implant) Implant {
	door.ClearScreenAndDisplay("assets/selectImplant.ans")

	for {
		input, err := door.GetKeyboardInput()
		if err != nil {
			fmt.Println("Error reading keyboard input:", err)
			continue
		}

		index, err := strconv.Atoi(input)
		if err == nil && index >= 1 && index <= len(implants) {
			return implants[index-1]
		}

		door.HandleInvalidInput()
	}
}
