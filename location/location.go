package location

import (
	"encoding/json"
	"fmt"
	"os"
)

// Location represents a game location.
type Location struct {
	Name string `json:"name"` // Name of the location
	Desc string `json:"desc"` // Additional information about the location (optional)
}

// LoadLocations loads locations from the provided JSON file.
func LoadLocations(filepath string) ([]Location, error) {
	var locations []Location

	// Read JSON file
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read locations file: %v", err)
	}

	// Unmarshal JSON data
	if err := json.Unmarshal(data, &locations); err != nil {
		return nil, fmt.Errorf("failed to unmarshal locations data: %v", err)
	}

	return locations, nil
}
