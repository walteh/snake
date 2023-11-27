package root

import (
	"github.com/spf13/pflag"
)

type CustomInterface interface {
}

type CustomInterfaceStruct struct {
}

type CustomResolver struct {
}

func (me *CustomResolver) Flags(flgs *pflag.FlagSet) {
}

func (me *CustomResolver) Run() (CustomInterface, error) {
	return &CustomInterfaceStruct{}, nil
}
