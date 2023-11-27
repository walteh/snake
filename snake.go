package snake

type SnakeFormatter interface {
	FormatJSON([]byte) ([]byte, error)
	FormatTable([]byte) ([]byte, error)
	FormatRaw([]byte) ([]byte, error)
}

// type SnakeOutput struct {
// 	Format SnakeFormat
// 	Data   []byte
// }

// type SnakeImplementation interface {

// }
