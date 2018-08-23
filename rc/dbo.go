package rc

import "fmt"

type Dbo interface {
	Type() string
	Id() int
	Serialize() string
	Deserialize(v string) Dbo
	DboManager() RedisManager
}

func generateKey(o Dbo) string {
	return fmt.Sprintf("%s:%v", o.Type(), o.Id())
}

func Save(o Dbo) {
	err := o.DboManager().Client(o).Set(generateKey(o), o.Serialize(), 0).Err()
	if err != nil {
		panic(err)
	}
}

func Delete(o Dbo) {
	err := o.DboManager().Client(o).Set(generateKey(o), nil, 0).Err()
	if err != nil {
		panic(err)
	}
}
