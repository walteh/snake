package root

import (
	"github.com/spf13/pflag"
	"github.com/walteh/snake/scobra"
)

type CustomInterface interface {
}

type CustomInterfaceStruct struct {
}

var _ scobra.Flagged = (*CustomResolver)(nil)

type CustomResolver struct {
}

func (me *CustomResolver) Flags(flgs *pflag.FlagSet) {
}

func (me *CustomResolver) Run() (CustomInterface, error) {
	return &CustomInterfaceStruct{}, nil
}
