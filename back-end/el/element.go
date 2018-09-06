package el

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
	return element.Key() == ""
}
