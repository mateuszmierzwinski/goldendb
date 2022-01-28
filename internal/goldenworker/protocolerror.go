package goldenworker

import "net"

const (
	statusOK            = byte(0)
	statusNotFound      = byte(1)
	statusNotFit        = byte(2)
	statusProtocolError = byte(3)
	statusNoSpaceLeft   = byte(4)
	statusUnknownError  = byte(255)
)

func protocolStatusWrite(conn net.Conn, errorCode byte, text []byte) {
	conn.Write([]byte{errorCode})
	conn.Write(text)
}
