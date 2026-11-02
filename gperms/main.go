package main

import (
	"fmt"
	"os"
)

func main() {
	inf, e := os.Stat(os.Args[1])
	if e != nil {
		fmt.Printf("stat: %t\n", e)
		return
	}
	fmt.Printf("%v:\nsize: %v\nperms: %v\nmodified: %v\n", inf.Name(), inf.Size(), inf.Mode(), inf.ModTime())

	//	fmt.Printf("\n\n\n")
	//	//	spew.Dump(inf)
	//permss := inf.Mode().Perm().String()
	permss := inf.Mode().String()
	//	fmt.Printf("\n\n\n\n%v\n%o", perms, m.Perm())
	//	ow := permss[0:3]
	//	gr := permss[3:6]
	//	ot := permss[6:9]
	fmt.Printf("\n\n")
	fmt.Printf("\nDEBUG: %q len=%d\n", permss, len(permss))

	for i, c := range permss {
		if i == 0 {
			fmt.Printf("\nType:\n    ")
			switch c {
			case 'd':
				fmt.Println("dir")
			case '-':
				fmt.Println("file")
			case 'l':
				fmt.Println("symlink")
			case 'c':
				fmt.Println("chardev")
			case 'b':
				fmt.Println("blockdev")
			case 'p':
				fmt.Println("pipe")
			case 's':
				fmt.Println("domain socket")
			}
			continue
		}

		switch i {
		case 1:
			fmt.Printf("\nOwner:\n")
		case 4:
			fmt.Printf("\nGroup:\n")
		case 7:
			fmt.Printf("\nUniverse:\n")
		}

		switch c {
		case 'r':
			fmt.Printf("    Read\n")
		case 'w':
			fmt.Printf("    Write\n")
		case 'x':
			fmt.Printf("    Execute\n")
		case '-':
			fmt.Printf("    No permission\n")
		}
	}
}

func ln(s string, n int) string {
	runes := []rune(s)
	if n < 0 {
		n = 0
	}
	if n > len(runes) {
		n = len(runes)
	}
	return string(runes[:n])
}
