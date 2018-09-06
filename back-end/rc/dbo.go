package rc

import (
	"fmt"
	"github.com/google/uuid"
)

type Dbo interface {
	Key() string
}

func GenerateKey(id uuid.UUID) string {
	return fmt.Sprintf("%v", id)
}
