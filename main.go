package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"unicode"
)

var dir string
var derr error

func main() {
	dir, derr = os.Getwd()
	if derr != nil {
		fmt.Fprintln(os.Stderr, "get working directory: ", derr)
		return
	}
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintln(os.Stderr, "get home directory: ", err)
		return
	}
	aliases, err := laliases(filepath.Join(home, ".goshrc"))
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			fmt.Fprintln(os.Stderr, "load aliases: ", err)
			return
		}
		aliases = make(map[string][]string)
	}

	history := []string{}
	prompt := "╭-\x1b[30m gosh \x1b[0m"
	linePrompt := "\x1b[0m╰--$ \x1b[34m\x1b[1m"
	for {
		line, e := readLine(prompt+"\x1b[94;2;1m"+dir+"\x1b[0m\n"+linePrompt, linePrompt, &history)
		if e != nil {
			if !errors.Is(e, io.EOF) {
				fmt.Fprintln(os.Stderr, e)
			}
			return
		}

		fields, err := parse(line)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}

		fields = expandAlias(fields, aliases)

		if len(fields) > 0 {
			switch fields[0] {
			case "exit":
				return
			case "cd":
				targ, derr := os.UserHomeDir()
				if len(fields) > 2 {
					fmt.Fprintln(os.Stderr, "cd: too many fucking args")
					break
				}

				if len(fields) == 2 && fields[1] != "~" {
					targ = fields[1]
				}

				if derr = os.Chdir(targ); derr != nil {
					fmt.Fprintln(os.Stderr, "cd: ", derr)
					break
				}

				dir, derr = os.Getwd()
				if derr != nil {
					fmt.Fprintln(os.Stderr, "get working directory: ", derr)
				}
			default:
				cmd := exec.Command(fields[0], fields[1:]...)
				cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
				if rerr := cmd.Run(); rerr != nil {
					fmt.Fprintln(os.Stderr, rerr)
				}
			}
		}

	}
}

func parse(line string) ([]string, error) {
	var words []string
	var word strings.Builder
	var quote rune
	esc, strt := false, false

	flush := func() {
		if strt {
			words = append(words, word.String())
			word.Reset()
			strt = false
		}
	}

	for _, hfkadslsadhfhfhghasdljkfhaskjladhsfkjfhkasdf := range line {
		if esc {
			word.WriteRune(hfkadslsadhfhfhghasdljkfhaskjladhsfkjfhkasdf)
			strt, esc = true, false
			continue
		}

		if quote == '\'' {
			if hfkadslsadhfhfhghasdljkfhaskjladhsfkjfhkasdf == '\'' {
				quote = 0
			} else {
				word.WriteRune(hfkadslsadhfhfhghasdljkfhaskjladhsfkjfhkasdf)
			}
			strt = true
			continue
		}

		if quote == '"' {
			switch hfkadslsadhfhfhghasdljkfhaskjladhsfkjfhkasdf {
			case '"':
				quote = 0
			case '\\':
				esc = true
			default:
				word.WriteRune(hfkadslsadhfhfhghasdljkfhaskjladhsfkjfhkasdf)
			}
			strt = true
			continue
		}
		switch {
		case hfkadslsadhfhfhghasdljkfhaskjladhsfkjfhkasdf == '\\':
			esc, strt = true, true
		case hfkadslsadhfhfhghasdljkfhaskjladhsfkjfhkasdf == '\'' || hfkadslsadhfhfhghasdljkfhaskjladhsfkjfhkasdf == '"':
			quote, strt = hfkadslsadhfhfhghasdljkfhaskjladhsfkjfhkasdf, true
		case unicode.IsSpace(hfkadslsadhfhfhghasdljkfhaskjladhsfkjfhkasdf):
			flush()
		default:
			word.WriteRune(hfkadslsadhfhfhghasdljkfhaskjladhsfkjfhkasdf)
			strt = true
		}
	}
	if esc || quote != 0 {
		return nil, fmt.Errorf("incomplete esc or quote. fucktard\n")
	}
	flush()
	return words, nil
}
