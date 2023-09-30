package snake

import (
	"reflect"

	"github.com/spf13/cobra"
)

var (
	root = Ctx{
		bindings:  make(map[string]*reflect.Value),
		resolvers: make(map[string]Execable),
		cmds:      make(map[string]*cobra.Command),
	}
)

func NewArgument[I any](method Flagged) (Execable, error) {
	return NewArgContext[I](&root, method)
}

func NewCmd(cmd *cobra.Command, method Flagged) (Execable, error) {
	return NewCmdContext(&root, cmd.Name(), cmd, method)
}
