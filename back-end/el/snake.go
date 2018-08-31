package el

type Snake struct {
	Location
}

func (s *Snake) String() string {
	return "S"
}
func (s *Snake) Attack(player *Player) {
	player.DecreaseHp(1)
}
func (snake *Snake) Interact(element interface{}) bool {
	switch element.(type) {
	case *Player:
		player := element.(*Player)
		snake.Attack(player)
		return true
	case *Empty:
		return true
	default:
		return false
	}

	return false
}
