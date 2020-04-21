package ModelsGenerator

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
)

func parseFile(path string) ([]*middleware, error) {
	fileContent, err := readFile(path)
	if err != nil {
		return nil, err
	}

	parse(string(fileContent))

	return nil, nil
}

/**
Get a file path and read all the content.
Return the content as byte slice.
The file must to be with .gen extension.
*/
func readFile(path string) ([]byte, error) {
	if filepath.Ext(path) != ".gen" {
		return nil, errors.New(
			"tried to parse a file with wrong type. the tool supports .gen files only")
	}

	return ioutil.ReadFile(path)
}

func parse(fileContent string) ([]middleware, error) {
	scanner := bufio.NewScanner(strings.NewReader(fileContent))

	result := make([]middleware, 0)
	var currentMiddleware middleware = nil

	currentLine := 0

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		currentLine++

		// Receiving { while reading another object
		if currentMiddleware != nil && line == "{" {
			return nil, errors.New(fmt.Sprintf(
				"failed to parse. Unexpected { token in line %v", currentLine))
		}

		if line == "{" {
			continue
		}

		// Finished to read the current object. Close him and add to result
		if line == "}" {
			result = append(result, currentMiddleware)
			currentMiddleware = nil
			continue
		}

		// If the line isn't { or }, and the current middleware is nil
		// we expect a new middleware declare, meaning a new class or enum title
		if currentMiddleware == nil {
			mw, err := readMiddlewareDeclare(line)
			if err != nil {
				return nil, err
			}

			currentMiddleware = mw
			continue
		}

		if err := readMiddlewareValue(currentMiddleware, line); err != nil {
			return nil, err
		}
	}

	return result, nil
}

func readMiddlewareDeclare(line string) (middleware, error) {
	splittedLine := strings.Split(line, " ")

	if len(splittedLine) != 2 {
		return nil, errors.New(fmt.Sprintf(
			"tried to declare middleware, but got string with the wrong length %s", line))
	}

	switch splittedLine[0] {
	case "class":
		return newClass(splittedLine[1]), nil
	case "enum":
		return newEnum(splittedLine[1]), nil
	}

	return nil, errors.New(
		fmt.Sprintf("tried to declare a new middleware, but got wrong declare %s", line))
}

func readMiddlewareValue(middleware middleware, line string) error {
	splittedLine := strings.Split(line, " ")

	if len(splittedLine) != 2 {
		return errors.New(fmt.Sprintf(
			"tried to read a new value, but got string with the wrong length %s", line))
	}

	middleware.addValue(splittedLine[0], splittedLine[1])

	return nil
}
