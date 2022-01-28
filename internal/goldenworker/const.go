package goldenworker

const (
	ping   = byte(0)
	create = byte(1)
	update = byte(2)
	delete = byte(3)
	read   = byte(4)

	pongResponse = "pong"

	defaultBuffSize = 65532
)
