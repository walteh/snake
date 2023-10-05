package snake

import (
	"context"

	"github.com/spf13/cobra"
)

func Apply(ctx context.Context, me *Ctx, root *cobra.Command) error {

	for _, exer := range me.resolvers {

		cmd := exer.Command().Cobra()

		if exer.Command() == nil {
			continue
		}

		if flgs, err := me.FlagsFor(exer); err != nil {
			return err
		} else {
			cmd.Flags().AddFlagSet(flgs)
		}

		err := exer.ValidateResponse()
		if err != nil {
			return err
		}

		oldRunE := cmd.RunE

		cmd.RunE = func(cmd *cobra.Command, args []string) error {
			defer setBindingWithLock(me, cmd)()
			err := me.Run(exer)
			if err != nil {
				return err
			}
			if oldRunE != nil {
				return oldRunE(cmd, args)
			}
			return nil
		}

	}

	return nil
}

func Build(ctx context.Context, me *Ctx) (*cobra.Command, error) {

	cmd := &cobra.Command{}

	if err := Apply(ctx, me, cmd); err != nil {
		return nil, err
	}

	for nme, sub := range me.cmds {
		cmdn := sub.Cobra()
		cmd.AddCommand(cmdn)

		if cmdn.Use == "" {
			cmdn.Use = nme
		}
	}

	return cmd, nil

}
