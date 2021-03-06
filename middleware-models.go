package main

import (
	"errors"
	"fmt"
	"strconv"
)

type middlewareType int

const (
	middlewareTypeClass middlewareType = 0
	middlewareTypeEnum  middlewareType = 2
)

/**
Represent a type in a language. For example: class or enum.
Different languages can implement that in different ways.
*/
type middleware interface {
	getType() middlewareType
	addValue(name string, value string) error
}

type dataMember struct {
	memberType string
	name       string
}

type class struct {
	name        string
	dataMembers []*dataMember
}

func newClass(name string) *class {
	return &class{
		name:        name,
		dataMembers: make([]*dataMember, 0),
	}
}

/**
Add new data member to the class.
The value parameter is the data member type
*/
func (c *class) addValue(name string, value string) error {
	member := &dataMember{
		memberType: value,
		name:       toCamelCase(name),
	}

	if !memberUnique(c.dataMembers, member) {
		return errors.New(fmt.Sprintf(
			"tried to add member %s to class %s, but it is already exists", name, c.name))
	}

	c.dataMembers = append(c.dataMembers, member)

	return nil
}

func (c *class) getType() middlewareType {
	return middlewareTypeClass
}

type enumValue struct {
	name  string
	value int
}

type enum struct {
	name       string
	enumValues []*enumValue
}

func newEnum(name string) *enum {
	return &enum{
		name:       name,
		enumValues: make([]*enumValue, 0),
	}
}

func (e *enum) addValue(name string, value string) error {
	v, err := strconv.Atoi(value)
	if err != nil {
		return err
	}

	newEnumValue := &enumValue{
		name:  toCamelCase(name),
		value: v,
	}

	if !enumValueUnique(e.enumValues, newEnumValue) {
		return errors.New(fmt.Sprintf(
			"tried to add enum value %s to enum %s, but it's already exists", name, e.name))
	}
	e.enumValues = append(e.enumValues, newEnumValue)

	return nil
}

func (e *enum) getType() middlewareType {
	return middlewareTypeEnum
}
