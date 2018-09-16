package velocity

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

var VELOCITY_MAPPING map[string]float64

func Constants(movableType string) float64 {
	if VELOCITY_MAPPING == nil {
		initializeVelocityMappings()
	}

	velocity, ok := VELOCITY_MAPPING[movableType]

	if !ok {
		panic(fmt.Sprintf("Movable velocities should be registered. Velocity asked for: %v, registered velocities: %v.", movableType, VELOCITY_MAPPING))
	}

	return velocity
}

func initializeVelocityMappings() {
	VELOCITY_MAPPING = make(map[string]float64)

	absPath, _ := filepath.Abs("../../game-config/velocities.json")
	jsonFile, err := os.Open(absPath)
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	json.Unmarshal(byteValue, &VELOCITY_MAPPING)
}
