package rc

import (
	"fmt"
	"strings"
	"github.com/google/uuid"
)

type Dbo interface {
	Key() string
	String() string
	Load() Dbo
}

func GenerateKey(id uuid.UUID) string {
	return fmt.Sprintf("%v", id)
}

func SplitKey(key string) (string, string) {
	split := strings.Split(key, ":")

	return split[0], split[1]
}
