package el

type Empty struct {
}

func (empty *Empty) String() string {
	return " "
}

func (empty *Empty) Id() int {
	return -1
}

func (empty *Empty) Key() string {
	return ""
}
