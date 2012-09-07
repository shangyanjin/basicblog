package main

import (
	"blog"
	"fmt"
)

func main() {
	b, err := blog.NewFromFile("data/entries.json")
	if err != nil {
		fmt.Printf("Error: unable to load blog entries: %s\n", err)
		return
	}

	mainMenu(b)
}

func mainMenu(b *blog.Blog) {
	for {
		fmt.Println("(1) List entries")
		fmt.Println("(2) Remove entry by ID\n")
		fmt.Println("(0) Quit\n")

		var n int
		fmt.Printf("> ")
		fmt.Scanf("%d", &n)
		fmt.Println("")

		switch n {
		case 0:
			var save, dmp rune
			fmt.Printf("Save? [Y/n]: ")
			fmt.Scanf("%c", &save)

			if save == 10 || save == 'Y' || save == 'y' {
				fmt.Println("Saving...")
				b.Save("data/entries.json")
			}

			// So for some reason, the left over '\n' in stdin will be
			// passed to bash causing an empty command to be ran. Not harmful
			// but confusing, so we catch it here.
			if save != 10 {
				fmt.Scanf("%c", &dmp)
			}

			return

		case 1:
			for _, e := range b.Entries {
				fmt.Printf("(%d) %q\n", e.ID, e.Title)
			}

		case 2:
			var id int
			fmt.Printf("ID: ")
			fmt.Scanf("%d", &id)

			var i int = -1
			for index, e := range b.Entries {
				if e.ID == id {
					i = index
					break
				}
			}

			if i >= 0 {
				b.Entries = append(b.Entries[:i], b.Entries[i+1:]...)
			} else {
				fmt.Printf("No entry found for ID %d\n", id)
			}

		default:
			fmt.Printf("Invalid option %d\n", n)
		}

		fmt.Println("")
	}
}
