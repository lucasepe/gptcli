package shortcuts

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

const (
	prefix = "@@"
)

func Expand(refs map[string]string, text string) string {
	if !strings.HasPrefix(text, prefix) {
		return ""
	}

	text = strings.TrimPrefix(text, prefix)
	parts := strings.SplitN(text, " ", 2)

	var prompt string
	for k, v := range refs {
		if k == parts[0] {
			prompt = v
			break
		}
	}

	if len(prompt) == 0 {
		return ""
	}

	return strings.Replace(text, parts[0], fmt.Sprintf("%s:", prompt), 1)
}

// FromReader read and parse an shortcuts file from
// an `io.Reader`, returning a map of keys and values.
func FromReader(r io.Reader) (res map[string]string, err error) {
	res = make(map[string]string)

	var lines []string
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err = scanner.Err(); err != nil {
		return
	}

	for _, fullLine := range lines {
		if !isIgnoredLine(fullLine) {
			var key, value string
			key, value, err = parseLine(fullLine, res)
			if err != nil {
				return
			}
			res[key] = value
		}
	}
	return
}

// FromFile read and parse a shortcurts from
// a file, returning a map of keys and values.
func FromFile(filename string) (res map[string]string, err error) {
	file, err := os.Open(filename)
	if err != nil {
		return
	}
	defer file.Close()

	return FromReader(file)
}

var (
	singleQuotesRegex  = regexp.MustCompile(`\A'(.*)'\z`)
	doubleQuotesRegex  = regexp.MustCompile(`\A"(.*)"\z`)
	escapeRegex        = regexp.MustCompile(`\\.`)
	unescapeCharsRegex = regexp.MustCompile(`\\([^$])`)
)

func parseLine(line string, envMap map[string]string) (string, string, error) {
	if len(line) == 0 {
		return "", "", errors.New("zero length string")
	}

	// ditch the comments (but keep quoted hashes)
	if strings.Contains(line, "#") {
		segmentsBetweenHashes := strings.Split(line, "#")
		quotesAreOpen := false
		var segmentsToKeep []string
		for _, segment := range segmentsBetweenHashes {
			if strings.Count(segment, "\"") == 1 || strings.Count(segment, "'") == 1 {
				if quotesAreOpen {
					quotesAreOpen = false
					segmentsToKeep = append(segmentsToKeep, segment)
				} else {
					quotesAreOpen = true
				}
			}

			if len(segmentsToKeep) == 0 || quotesAreOpen {
				segmentsToKeep = append(segmentsToKeep, segment)
			}
		}

		line = strings.Join(segmentsToKeep, "#")
	}

	splitString := strings.SplitN(line, "=", 2)
	if len(splitString) != 2 {
		return "", "", errors.New("can't separate key from value")
	}

	// Parse the key
	key := splitString[0]
	key = strings.TrimSpace(key)

	// Parse the value
	value, err := parseValue(splitString[1], envMap)

	return key, value, err
}

func parseValue(value string, envMap map[string]string) (res string, err error) {
	// trim
	res = strings.Trim(value, " ")
	if len(res) <= 1 {
		return res, nil
	}

	// check if we've got quoted values or possible escapes
	singleQuotes := singleQuotesRegex.FindStringSubmatch(res)
	doubleQuotes := doubleQuotesRegex.FindStringSubmatch(res)

	if singleQuotes != nil || doubleQuotes != nil {
		// pull the quotes off the edges
		res = res[1 : len(res)-1]
	}

	if doubleQuotes != nil {
		// expand newlines
		res = escapeRegex.ReplaceAllStringFunc(res, func(match string) string {
			c := strings.TrimPrefix(match, `\`)
			switch c {
			case "n":
				return "\n"
			case "r":
				return "\r"
			default:
				return match
			}
		})
		// unescape characters
		res = unescapeCharsRegex.ReplaceAllString(res, "$1")
	}

	return res, err
}

func isIgnoredLine(line string) bool {
	trimmedLine := strings.TrimSpace(line)
	return len(trimmedLine) == 0 || strings.HasPrefix(trimmedLine, "#")
}
