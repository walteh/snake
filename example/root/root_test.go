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
	"github.com/walteh/snake/example/root/sample"
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

			_, cmd, hndl, err := NewCommand(tt.args.ctx)
			if tt.wantErr {
				assert.Error(t, err)
				return
			} else {
				require.NoError(t, err)
			}

			err = os.Setenv("ROOT_COOL", "true")
			require.NoError(t, err)

			os.Args = []string{"root", "sample", "--value", "test123", "--sample-enum", "select"}

			err = cmd.RootCommand.Execute()
			require.NoError(t, err)

			assert.True(t, hndl.Cool)
			assert.Equal(t, "test123", hndl.Value)

			args := hndl.Args()

			assert.Equal(t, sample.SampleEnumY, args.Enum)
		})
	}

}

func TestDocs(t *testing.T) {

	ctx := context.Background()

	_, cmd, _, err := NewCommand(ctx)
	if err != nil {
		t.Errorf("NewCommand() error = %v", err)
		return
	}

	tmp := os.TempDir()

	ref := filepath.Join(tmp, "root-docs")

	t.Cleanup(func() {
		os.RemoveAll(ref)
	})

	mdpath := filepath.Join(ref, "md")

	if err := os.MkdirAll(mdpath, 0755); err != nil {
		t.Errorf("MkdirAll() error = %v", err)
		return
	}

	err = doc.GenMarkdownTree(cmd.RootCommand, mdpath)
	if err != nil {
		t.Errorf("GenMarkdownTree() error = %v", err)
		return
	}

	fle, err := os.Open(filepath.Join(mdpath, "root_sample.md"))
	if err != nil {
		t.Errorf("Open() error = %v", err)
		return
	}

	defer fle.Close()

	stat, err := fle.Stat()
	if err != nil {
		t.Errorf("Stat() error = %v", err)
		return
	}

	dat := make([]byte, stat.Size())

	_, err = fle.Read(dat)
	if err != nil {
		t.Errorf("Read() error = %v", err)
		return
	}

	assert.True(t, stat.Size() > 0)

	fmt.Println(string(dat))

}
