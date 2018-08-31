package el

import (
	"../gs"
	"../rc"
	"strconv"
	"strings"
)

const ELEMENT = "element"

type Element struct {
	Location
	DboKey string
}

func (element *Element) Type() string {
	return ELEMENT
}

func (element *Element) Key() string {
	return element.LocationKey()
}

func (element *Element) Serialize() string {
	return element.DboKey
}

func (element *Element) Deserialize(key string, dboKey string) {
	//fmt.Printf("element.go: Attempt to Deserialize. key: %s, val: %s\n", key, dboKey)

	subWorld, coord := rc.SplitKey(key)
	//fmt.Printf("element.go: Spliting key: %s, %s", subWorld, coord)

	subWorldCoordStringX := strings.Split(subWorld, ",")[0]
	subWorldCoordX, _ := strconv.Atoi(subWorldCoordStringX)
	subWorldCoordStringY := strings.Split(subWorld, ",")[1]
	subWorldCoordY, _ := strconv.Atoi(subWorldCoordStringY)

	gridCoordStringX := strings.Split(coord, ",")[0]
	gridCoordX, _ := strconv.Atoi(gridCoordStringX)
	gridCoordStringY := strings.Split(coord, ",")[1]
	gridCoordY, _ := strconv.Atoi(gridCoordStringY)

	element.SubWorldCoord = gs.NewCoord(subWorldCoordX, subWorldCoordY)
	element.GridCoord = gs.NewCoord(gridCoordX, gridCoordY)

	element.DboKey = dboKey
}

func (element *Element) String() string {
	return " "
}

func (element *Element) IsEmpty() bool {
	//fmt.Printf("element.go: DboKey: %s\n", element.DboKey)
	return element.DboKey == ""
}

func newElementDbo(id int) rc.Dbo {
	return &Element{}
}
