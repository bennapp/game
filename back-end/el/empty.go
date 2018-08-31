package el

const EMPTY = "empty"

type Empty struct {
	Location
}

func (empty *Empty) String() string {
	return " "
}

func (empty *Empty) Type() string {
	return EMPTY
}

func (empty *Empty) Id() int {
	return -1
}

func (empty *Empty) Key() string {
	return ""
}

func (empty *Empty) Serialize() string {
	return ""
}

func (empty *Empty) Deserialize(key string, values string) {
}
