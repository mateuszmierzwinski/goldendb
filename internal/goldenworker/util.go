package goldenworker

import (
	"encoding/binary"
	"fmt"
	"net"
)

func readCommand(conn net.Conn) (byte, error) {
	// read command
	cmd := make([]byte, 1)
	n, e := conn.Read(cmd)
	if e != nil {
		return 0, e
	}
	// did we pulled the content?
	if n < 1 {
		return 0, fmt.Errorf("cannot process event - cannot read package")
	}

	return cmd[0], nil
}

func intToBytes(value int) []byte {
	return int64toBytes(int64(value))
}

func int64toBytes(value int64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(value))
	return b
}

func bytesArrayToUint64(bytes []byte) int64 {
	if len(bytes) < 8 {
		return 0
	}
	return int64(binary.LittleEndian.Uint64(bytes))
}
