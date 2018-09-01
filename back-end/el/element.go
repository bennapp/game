package el

import (
	"../rc"
)

const ELEMENT = "element"

type Element struct {
	Location
	DboKey string
}

func (element *Element) Key() string {
	return element.LocationKey()
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
