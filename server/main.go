package main

import (
	"fmt"
)

func main() {
	p1 := NewParticipant("fulano", []Profile{BackEnd}, ContactInfo{"fulano@foo.com"})
	p2 := NewParticipant("altrano", []Profile{Web}, ContactInfo{"altrano@foo.com"})
	t := NewTeam("Team", p1)
	t.AddParticipant(p2) // POST
	t.UpdateWishList([]Profile{Mobile,UX}) // PUT
	fmt.Println(t)
}
