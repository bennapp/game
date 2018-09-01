package rc

import (
	"fmt"
	"strings"
)

type Dbo interface {
	Type() string
	Key() string
	String() string
}

func GenerateKey(t string, id int) string {
	return fmt.Sprintf("%s:%v", t, id)
}

func SplitKey(key string) (string, string) {
	split := strings.Split(key, ":")

	return split[0], split[1]
}
