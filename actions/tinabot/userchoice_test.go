package tinabot

import (
	"fmt"
	"testing"

	"github.com/develersrl/lunches/pkg/tuttobene"
)

func assertEqual(t *testing.T, a interface{}, b interface{}, message string) {
	if a == b {
		return
	}
	if len(message) == 0 {
		message = fmt.Sprintf("%v != %v", a, b)
	}
	t.Fatal(message)
}

func assertNotEqual(t *testing.T, a interface{}, b interface{}, message string) {
	if a != b {
		return
	}
	if len(message) == 0 {
		message = fmt.Sprintf("%v == %v", a, b)
	}
	t.Fatal(message)
}

func TestUserChoice(t *testing.T) {
	var choice UserChoice
	p := tuttobene.MenuRow{
		Content: "primo",
		Type:    tuttobene.Primo,
	}
	s1 := tuttobene.MenuRow{
		Content: "secondo1",
		Type:    tuttobene.Secondo,
	}
	s2 := tuttobene.MenuRow{
		Content: "secondo2",
		Type:    tuttobene.Secondo,
	}
	c1 := tuttobene.MenuRow{
		Content: "contorno1",
		Type:    tuttobene.Contorno,
	}
	c2 := tuttobene.MenuRow{
		Content: "contorno2",
		Type:    tuttobene.Contorno,
	}
	v := tuttobene.MenuRow{
		Content: "vegetariano",
		Type:    tuttobene.Vegetariano,
	}
	pa := tuttobene.MenuRow{
		Content: "panino",
		Type:    tuttobene.Panino,
	}
	f := tuttobene.MenuRow{
		Content: "frutta",
		Type:    tuttobene.Frutta,
	}

	e := choice.Add(p)
	assertEqual(t, e, nil, "")

	assertEqual(t, choice.String(), "primo", "")
	assertEqual(t, choice.OrdString(), "0004-primo", "")

	e = choice.Add(s1)
	assertNotEqual(t, e, nil, "Non si può comporre un primo con un secondo")

	e = choice.Add(f)
	assertNotEqual(t, e, nil, "Non si può comporre un primo con una frutta")

	e = choice.Add(pa)
	assertNotEqual(t, e, nil, "Non si può comporre un primo con un panino")

	assertEqual(t, choice.Customized(), false, "")

	choice.Clear()
	assertEqual(t, choice.DishMask, uint(0), "")
	assertEqual(t, len(choice.Dishes), 0, "")

	e = choice.Add(v)
	assertEqual(t, e, nil, "")
	assertEqual(t, choice.String(), "vegetariano", "")
	assertEqual(t, choice.OrdString(), "0032-vegetariano", "")

	e = choice.Add(s1)
	assertEqual(t, e, nil, "")
	assertEqual(t, choice.String(), "secondo1 con vegetariano", "")
	assertEqual(t, choice.OrdString(), "0040-secondo1 con vegetariano", "")

	e = choice.Add(c2)
	assertEqual(t, e, nil, "")
	assertEqual(t, choice.String(), "secondo1 con contorno2, vegetariano", "")
	assertEqual(t, choice.OrdString(), "0056-secondo1 con contorno2, vegetariano", "")

	e = choice.Add(s2)
	assertNotEqual(t, e, nil, "Non si può comporre un secondo con un secondo")

	e = choice.Add(c1)
	assertEqual(t, e, nil, "")

	assertEqual(t, choice.String(), "secondo1 con contorno1, contorno2, vegetariano", "")
	assertEqual(t, choice.OrdString(), "0056-secondo1 con contorno1, contorno2, vegetariano", "")

	e = choice.Add(pa)
	assertNotEqual(t, e, nil, "Non si può comporre un secondo con un panino")

	assertEqual(t, choice.Customized(), true, "")
}
