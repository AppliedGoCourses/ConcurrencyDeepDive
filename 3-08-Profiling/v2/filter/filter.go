package filter

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"
)

// Grep reads lines from r until EOF and returns a slice of the lines that match the
// regular expression pattern, or an error if the pattern is invalid
// or the scanner fails to read a line.
func Grep(input io.Reader, pattern string) ([]string, error) {
	var filteredLines []string
	if input == nil {
		return nil, fmt.Errorf("input must not be nil")
	}
	scanner := bufio.NewScanner(input)
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("cannot compile regexp '%s': %w", pattern, err)
	}
	for scanner.Scan() {
		line := scanner.Text()
		if re.MatchString(line) {
			filteredLines = append(filteredLines, line)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return filteredLines, nil
}

func Match(input io.Reader, searchStr string) ([]string, error) {
	var filteredLines []string
	if input == nil {
		return nil, fmt.Errorf("input must not be nil")
	}
	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, searchStr) {
			filteredLines = append(filteredLines, line)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return filteredLines, nil
}
