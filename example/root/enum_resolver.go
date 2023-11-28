package root

import (
	"github.com/spf13/cobra"
	"github.com/walteh/snake/example/root/sample"
)

type EnumResolver struct {
	Myenum sample.SampleEnum `usage:"Enum" default:"z"`
}

func (me *EnumResolver) Run(cmd *cobra.Command) (*sample.SampleEnum, error) {

	return &me.Myenum, nil
}
