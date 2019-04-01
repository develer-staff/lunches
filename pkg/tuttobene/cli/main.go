package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/develersrl/lunches/pkg/tuttobene"
)

const (
	usage = `
____  ______________ _____  _____
 | |  | |  | |  ||__]|___|\ ||___
 | |__| |  | |__||__]|___| \||___

A tool for parsing TuttoBene's menus

Usage: tuttobene <xlsx file> <output format>

Format can be:
- json
- tina

`
)

var TinaFormatTitles = map[tuttobene.MenuRowType]string{
	tuttobene.Primo:       "Primi Piatti",
	tuttobene.Secondo:     "Secondi Piatti",
	tuttobene.Contorno:    "Contorni",
	tuttobene.Vegetariano: "Piatti Vegetariano",
	tuttobene.Frutta:      "Frutta",
	tuttobene.Panino:      "Panini Espressi",
}

func main() {
	if len(os.Args) < 3 {
		fmt.Print(usage)
		os.Exit(1)
	}

	bs, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not open file: %v\n", err)
		os.Exit(1)
	}

	menu, err := tuttobene.ParseMenuBytes(bs)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not parse file: %v\n", err)
		os.Exit(1)
	}

	if menu == nil {
		fmt.Fprintf(os.Stderr, "Unexpected nil menu: %v\n", err)
		os.Exit(1)
	}

	switch os.Args[2] {
	case "json":
		out, err := json.MarshalIndent(menu, "", "		")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not marshal menu: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(string(out))
	case "tina":
		var currentSection tuttobene.MenuRowType
		for _, m := range menu.Rows {
			if currentSection != m.Type {
				currentSection = m.Type
				fmt.Println("\n" + TinaFormatTitles[currentSection])
			}

			if m.IsDailyProposal {
				fmt.Print("Proposta del giorno: ")
			}
			fmt.Println(m.Content)
		}
		fmt.Println("")
	default:
		fmt.Fprintf(os.Stderr, "Invalid format (json|tina): %v\n", os.Args[2])
		os.Exit(1)
	}
}
