package resolvers

import "github.com/walteh/snake"

func CustomRunner() snake.Runner {
	return snake.GenRunResolver_In00_Out02(&CustomResolver{})
}

type CustomInterface interface {
}

type CustomResolver struct {
}

func (me *CustomResolver) Run() (CustomInterface, error) {
	return struct{}{}, nil
}
