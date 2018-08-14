package main

import "fmt"

type Cell struct {
	element interface{}
	code    int
}

type Element struct {
}
type Rock struct {
	Element
}
type Snake struct {
	data int
	Element
}

func (element *Element) Work() {
	fmt.Println("element work")
}

//func (rock *Rock) Work() {
//	fmt.Println("rock work")
//}

func (element *Element) Interact(otherElement *Element) {

}

func (rock *Rock) Interact(snake *Snake) {
	fmt.Println("rock smash snake!")
	fmt.Println(snake.data)
}

func main() {
	rock := Rock{}
	rock.Work()

	snake := Snake{data: 1}

	cell := Cell{element: &snake}

	// This does not work
	// rock.Interact(cell.element)

	// This does work
	switch v := cell.element.(type) {
	case *Snake:
		snakePointer := cell.element.(*Snake)
		rock.Interact(snakePointer)
	default:
		fmt.Printf("I don't know about type %T!\n", v)
	}
}
