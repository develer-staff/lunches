package main

import (
	"encoding/json"
	"fmt"
	"github.com/develersrl/lunches/pkg/tuttobene"
	"io/ioutil"
	"os"
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

func main() {
	if len(os.Args) < 3 {
		fmt.Print(usage)
		os.Exit(1)
	}

	bs, err := ioutil.ReadFile(os.Args[1]);
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not open file: %v", err)
		os.Exit(1)
	}

	menu, err := tuttobene.ParseMenuBytes(bs)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not parse file: %v", err)
		os.Exit(1)
	}

	if menu == nil {
		fmt.Fprintf(os.Stderr, "Unexpected nil menu: %v", err)
		os.Exit(1)
	}

	switch (os.Args[2]) {
	case "json":
		out, err := json.MarshalIndent(menu, "", "		")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not marshal menu: %v", err)
			os.Exit(1)
		}
		fmt.Println(string(out))
	case "tina":
		var currentSection tuttobene.MenuRowType
		for _, m := range *menu {
			if currentSection != m.Type {
				currentSection = m.Type
				fmt.Println("\n" + tuttobene.Titles[currentSection])
			}
			fmt.Println(m.Content)
		}
		fmt.Println("")
	}
}
