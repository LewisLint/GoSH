package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"syscall"
	"unicode/utf8"
	"unsafe"
)

func readLine(prompt, linePrompt string, history *[]string) (result string, resultErr error) {
	fd := int(os.Stdin.Fd())
	var original syscall.Termios
	if err := ioctl(fd, syscall.TCGETS, &original); err != nil {
		if errors.Is(err, syscall.ENOTTY) {
			return readLineWithoutTerminal(prompt)
		}
		return "", fmt.Errorf("read terminal settings: %w", err)
	}

	raw := original
	raw.Lflag &^= syscall.ICANON | syscall.ECHO | syscall.ISIG
	raw.Cc[syscall.VMIN] = 0
	raw.Cc[syscall.VTIME] = 1
	if err := ioctl(fd, syscall.TCSETS, &raw); err != nil {
		return "", fmt.Errorf("enable character input: %w", err)
	}
	defer func() {
		if err := ioctl(fd, syscall.TCSETS, &original); err != nil {
			restoreErr := fmt.Errorf("restore terminal settings: %w", err)
			if resultErr == nil {
				resultErr = restoreErr
			} else {
				resultErr = errors.Join(resultErr, restoreErr)
			}
		}
	}()

	fmt.Print(prompt)
	var line []rune
	cursor := 0
	historyIndex := len(*history)
	draft := ""

	for {
		ch, ok, err := readInputByte(fd)
		if err != nil {
			return "", fmt.Errorf("read terminal input: %w", err)
		}
		if !ok {
			continue
		}

		switch ch {
		case '\r', '\n':
			value := string(line)
			fmt.Print("\x1b[0m\n")
			if value != "" && (len(*history) == 0 || (*history)[len(*history)-1] != value) {
				*history = append(*history, value)
			}
			return value, nil
		case 3:
			fmt.Print("\x1b[0m^C\n")
			return "", nil
		case 4:
			if len(line) == 0 {
				fmt.Print("\x1b[0m\n")
				return "", io.EOF
			}
		case 8, 127:
			if cursor > 0 {
				line = append(line[:cursor-1], line[cursor:]...)
				cursor--
				redrawLine(linePrompt, line, cursor)
			}
		case 27:
			sequence, ok, err := readEscapeSequence(fd)
			if err != nil {
				return "", fmt.Errorf("read terminal input: %w", err)
			}
			if ok {
				switch sequence {
				case "[A":
					if historyIndex > 0 {
						if historyIndex == len(*history) {
							draft = string(line)
						}
						historyIndex--
						line = []rune((*history)[historyIndex])
						cursor = len(line)
						redrawLine(linePrompt, line, cursor)
					}
				case "[B":
					if historyIndex < len(*history) {
						historyIndex++
						value := draft
						if historyIndex < len(*history) {
							value = (*history)[historyIndex]
						}
						line = []rune(value)
						cursor = len(line)
						redrawLine(linePrompt, line, cursor)
					}
				case "[C":
					if cursor < len(line) {
						cursor++
						redrawLine(linePrompt, line, cursor)
					}
				case "[D":
					if cursor > 0 {
						cursor--
						redrawLine(linePrompt, line, cursor)
					}
				case "[H", "OH", "[1~", "[7~":
					cursor = 0
					redrawLine(linePrompt, line, cursor)
				case "[F", "OF", "[4~", "[8~":
					cursor = len(line)
					redrawLine(linePrompt, line, cursor)
				case "[3~":
					if cursor < len(line) {
						line = append(line[:cursor], line[cursor+1:]...)
						redrawLine(linePrompt, line, cursor)
					}
				}
			}
		default:
			if ch >= 32 {
				r, err := readInputRune(fd, ch)
				if err != nil {
					return "", fmt.Errorf("read terminal input: %w", err)
				}
				line = append(line, 0)
				copy(line[cursor+1:], line[cursor:])
				line[cursor] = r
				cursor++
				redrawLine(linePrompt, line, cursor)
			}
		}
	}
}

func readInputByte(fd int) (byte, bool, error) {
	var buffer [1]byte
	n, err := syscall.Read(fd, buffer[:])
	if err != nil {
		return 0, false, err
	}
	if n == 0 {
		return 0, false, nil
	}
	return buffer[0], true, nil
}

func readInputRune(fd int, first byte) (rune, error) {
	size := 1
	switch {
	case first&0xE0 == 0xC0:
		size = 2
	case first&0xF0 == 0xE0:
		size = 3
	case first&0xF8 == 0xF0:
		size = 4
	}

	bytes := []byte{first}
	for len(bytes) < size {
		next, ok, err := readInputByte(fd)
		if err != nil {
			return 0, err
		}
		if !ok {
			return utf8.RuneError, nil
		}
		bytes = append(bytes, next)
	}
	r, _ := utf8.DecodeRune(bytes)
	return r, nil
}

func readEscapeSequence(fd int) (string, bool, error) {
	second, ok, err := readInputByte(fd)
	if err != nil || !ok {
		return "", false, err
	}
	if second != '[' && second != 'O' {
		return "", false, nil
	}

	sequence := []byte{second}
	for len(sequence) < 4 {
		next, ok, err := readInputByte(fd)
		if err != nil {
			return "", false, err
		}
		if !ok {
			return "", false, nil
		}
		sequence = append(sequence, next)
		if next >= 'A' && next <= 'Z' || next == '~' {
			break
		}
	}
	return string(sequence), true, nil
}

func redrawLine(prompt string, line []rune, cursor int) {
	fmt.Print("\r\x1b[2K", prompt, string(line))
	if remaining := len(line) - cursor; remaining > 0 {
		fmt.Printf("\x1b[%dD", remaining)
	}
}

func readLineWithoutTerminal(prompt string) (string, error) {
	fmt.Print(prompt)
	var line strings.Builder
	var buffer [1]byte
	for {
		n, err := os.Stdin.Read(buffer[:])
		if n > 0 {
			if buffer[0] == '\n' {
				fmt.Print("\x1b[0m\n")
				return strings.TrimSuffix(line.String(), "\r"), nil
			}
			line.WriteByte(buffer[0])
		}
		if err != nil {
			if errors.Is(err, io.EOF) && line.Len() > 0 {
				return line.String(), nil
			}
			return "", err
		}
	}
}

func ioctl(fd int, request uintptr, value *syscall.Termios) error {
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), request, uintptr(unsafe.Pointer(value)))
	if errno != 0 {
		return errno
	}
	return nil
}
