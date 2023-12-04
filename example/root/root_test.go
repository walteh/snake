package root

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/walteh/snake/example/resolvers"
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
			name: "simple test 1",
			args: args{
				ctx: context.Background(),
			},
			want: &cobra.Command{
				Use: "root",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctx := context.Background()

			_, cmd, hndl, err := NewCommand(ctx)
			if tt.wantErr {
				assert.Error(t, err)
				return
			} else {
				require.NoError(t, err)
			}

			err = os.Setenv("ROOT_COOL", "true")
			require.NoError(t, err)

			os.Args = []string{"root", "sample", "--value", "test123", "--sample-enum", "select", "--number-of-runs", "1", "--interval", "500ms", "--debug"}

			err = cmd.RootCommand.Execute()
			require.NoError(t, err)

			assert.True(t, hndl.Cool)
			assert.Equal(t, "test123", hndl.Value)

			args := hndl.Args()

			assert.Equal(t, resolvers.SampleEnumY, args.Enum)
		})
	}

}

func TestDocs(t *testing.T) {

	ctx := context.Background()

	_, cmd, _, err := NewCommand(ctx)
	require.NoError(t, err)

	tmp := os.TempDir()

	ref := filepath.Join(tmp, "root-docs")

	t.Cleanup(func() {
		os.RemoveAll(ref)
	})

	mdpath := filepath.Join(ref, "md")

	err = os.MkdirAll(mdpath, 0755)
	require.NoError(t, err)

	err = doc.GenMarkdownTree(cmd.RootCommand, mdpath)
	require.NoError(t, err)

	fle, err := os.Open(filepath.Join(mdpath, "root_sample.md"))
	require.NoError(t, err)

	defer fle.Close()

	stat, err := fle.Stat()
	require.NoError(t, err)

	dat := make([]byte, stat.Size())

	_, err = fle.Read(dat)
	require.NoError(t, err)

	assert.True(t, stat.Size() > 0)

	fmt.Println(string(dat))

}
