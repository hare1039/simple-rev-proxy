package def

const (
	EXIT_ERROR          = 1
	EXIT_ERROR_INTERNET = 2

	BUF_SIZE = 8096
)

type Config struct {
	Port string
}

type TCPstream struct {
	Id   int
	Data []byte
}
