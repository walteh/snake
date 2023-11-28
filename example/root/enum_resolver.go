package root

import (
	"github.com/spf13/cobra"
	"github.com/walteh/snake/example/root/sample"
)

type EnumResolver struct {
	Myenum sample.SampleEnum `usage:"Enum" default:"r"`
}

func (me *EnumResolver) Run(cmd *cobra.Command) (*sample.SampleEnum, error) {

	return &me.Myenum, nil
}
