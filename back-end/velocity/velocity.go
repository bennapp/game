package velocity

import (
	"../os_util"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
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

	jsonFile, err := os.Open(os_util.RelativeProjectPath("game-config/velocities.json"))
	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	json.Unmarshal(byteValue, &VELOCITY_MAPPING)
}
