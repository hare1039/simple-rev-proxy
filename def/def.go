package def

const (
	EXIT_ERROR          = 1
	EXIT_ERROR_INTERNET = 2

	BUF_SIZE           = 8
	CHANNEL_BUF_AMOUNT = 100
)

type Config struct {
	Port string
}

type TCPstream struct {
	Id   int    `json:"id"`
	Data []byte `json:"data"`
}
