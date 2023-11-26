package root

import (
	"context"
	"os"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestNewCommand(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		want    *cobra.Command
		wantErr bool
	}{
		{
			name: "test",
			args: args{
				ctx: context.Background(),
			},
			want: &cobra.Command{
				Use: "retab",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			cmd, hndl, err := NewCommand(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewCommand() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			err = os.Setenv("RETAB_COOL", "true")
			if err != nil {
				t.Errorf("Setenv() error = %v", err)
				return
			}

			os.Args = []string{"retab", "sample", "--value", "test"}

			err = cmd.Execute()
			if err != nil {
				t.Errorf("Execute() error = %v", err)
				return
			}

			assert.True(t, hndl.Cool)

			// assert.Equal(t, tt.want, got)
		})
	}
}
