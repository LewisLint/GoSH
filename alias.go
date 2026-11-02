package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func laliases(path string) (map[string][]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	aliases := make(map[string][]string)
	scanner := bufio.NewScanner(file)

	for ln := 1; scanner.Scan(); ln++ {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		name, cmd, found := strings.Cut(line, "=")
		name = strings.TrimSpace(name)
		args := strings.Fields(cmd)

		if !found || len(strings.Fields(name)) != 1 || len(args) == 0 {
			return nil, fmt.Errorf("%s:%d: expected alias=command", path, ln)
		}
		aliases[name] = args
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("reading %s: %w", path, err)
	}
	return aliases, nil
}

func expandAlias(fields []string, aliases map[string][]string) []string {
	if len(fields) == 0 {
		return fields
	}

	replacement, ok := aliases[fields[0]]
	if !ok {
		return fields
	}

	expanded := append([]string{}, replacement...)
	return append(expanded, fields[1:]...)
}
