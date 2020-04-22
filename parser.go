package main

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"
)

var spaceRegex = regexp.MustCompile(`\s+`)

func parseFile(path string) ([]middleware, error) {
	fileContent, err := readFile(path)
	if err != nil {
		return nil, err
	}

	return parse(string(fileContent))
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

/**
Get content and parse it to the middleware language.
Scanning row by row ignoring spaces, and do validation checks do (with the help methods).
Content can contain few classes and enums.
Classes and enums declaration should be like:

class className
{
	dataMember int
	another string
}

enum enumName
{
	value 5
	anotherValue 8
}
*/
func parse(fileContent string) ([]middleware, error) {
	scanner := bufio.NewScanner(strings.NewReader(fileContent))

	result := make([]middleware, 0)
	var currentMiddleware middleware = nil
	startAddMembers := false

	currentLine := 0

	for scanner.Scan() {
		line := trimContent(scanner.Text())
		currentLine++

		// Receiving { while reading another object
		if line == "{" {
			if startAddMembers {
				return nil, errors.New(fmt.Sprintf(
					"failed to parse. Unexpected { token in line %v", currentLine))
			}

			startAddMembers = true
			continue
		}

		// Finished to read the current object. Close him and add to result
		if line == "}" {
			result = append(result, currentMiddleware)
			currentMiddleware = nil
			startAddMembers = false
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

		// We got to a data member or enum value before getting {
		// after the class/enum declare
		if !startAddMembers {
			return nil, errors.New(fmt.Sprintf(
				"expected for { token before starting get values in row %v", currentLine))
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

	return middleware.addValue(splittedLine[0], splittedLine[1])
}

func trimContent(content string) string {
	// Delete all the double whitespaces before the final trim
	s := spaceRegex.ReplaceAllString(content, " ")
	return strings.TrimSpace(s)
}
