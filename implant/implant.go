package implant

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/eiannone/keyboard"
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

// SelectImplant prompts the user to select an implant from the available options.
func SelectImplant(implants []Implant) Implant {
	fmt.Println("Choose your cybernetic implant:")
	for i, imp := range implants {
		fmt.Printf("%d. %s - %s\n", i+1, imp.Name, imp.Desc)
	}

	for {
		// Initialize keyboard listener
		err := keyboard.Open()
		if err != nil {
			fmt.Println("Error opening keyboard:", err)
			return Implant{} // Return a default value
		}
		defer keyboard.Close()

		// Listen for single key press
		char, _, err := keyboard.GetSingleKey()
		if err != nil {
			fmt.Println("Error reading keyboard input:", err)
			continue
		}

		// Convert the pressed key to index
		index, err := strconv.Atoi(string(char))
		if err == nil && index >= 1 && index <= len(implants) {
			return implants[index-1] // Return the selected implant
		}

		fmt.Println("Invalid choice, please select a valid implant.")
	}
}
