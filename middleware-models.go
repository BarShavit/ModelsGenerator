package ModelsGenerator

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
}

type dataMember struct {
	memberType string
	name       string
}

type class struct {
	dataMembers []*dataMember
}

func (c *class) getType() middlewareType {
	return middlewareTypeClass
}

type enumValue struct {
	name  string
	value int
}

type enum struct {
	enumValues []*enumValue
}

func (e *enum) getType() middlewareType {
	return middlewareTypeEnum
}
