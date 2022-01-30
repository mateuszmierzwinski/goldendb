package goldenworker

import (
	"fmt"
	"golden/internal/goldendb"
	"log"
	"net"
)

type GoldenWorker struct {
	goldenDB *goldendb.GoldenDB
}

func (w *GoldenWorker) DispatchAndServe(conn net.Conn, workers chan *GoldenWorker) {
	// When this worker ends the job just send it back to queue
	defer func() { workers <- w }()
	defer conn.Close()

	command, err := readCommand(conn)
	if err != nil {
		log.Println(err.Error())
		return
	}

	switch command {
	case ping:
		w.cmdPong(conn)
		break
	case create:
		w.cmdCreate(conn)
		break
	case update:
		w.cmdCreate(conn)
		break
	case delete:
		w.cmdDelete(conn)
		break
	case read:
		w.cmdRead(conn)
		break
	}
}

func (w *GoldenWorker) cmdPong(conn net.Conn) {
	protocolStatusWrite(conn, statusOK, []byte(pongResponse))
}

func (w *GoldenWorker) cmdCreate(conn net.Conn) {
	// read key size
	sizeBuff := make([]byte, 8)
	if _, e := conn.Read(sizeBuff); e != nil {
		msg := fmt.Sprintf("cannot create object - %s", e.Error())
		log.Println(msg)
		protocolStatusWrite(conn, statusProtocolError, []byte(msg))
		return
	}
	keySize := bytesArrayToUint64(sizeBuff)

	// read key
	key := make([]byte, keySize)
	if _, e := conn.Read(key); e != nil {
		msg := fmt.Sprintf("cannot read object key - %s", e.Error())
		log.Println(msg)
		protocolStatusWrite(conn, statusProtocolError, []byte(msg))
		return
	}

	// read objSize
	if _, e := conn.Read(sizeBuff); e != nil {
		msg := fmt.Sprintf("cannot read object size - %s", e.Error())
		log.Println(msg)
		protocolStatusWrite(conn, statusProtocolError, []byte(msg))
		return
	}

	// memory checks
	objSize := bytesArrayToUint64(sizeBuff)
	if uint64(objSize) > w.goldenDB.MemoryLimit() {
		msg := fmt.Sprintf("object (%d bytes) does not fit in maximum memory allowed (%d bytes)", uint64(objSize), w.goldenDB.MemoryLimit())
		log.Println(msg)
		protocolStatusWrite(conn, statusNotFit, []byte(msg))
		return
	}
	log.Printf("Incomming label %s with size of %d bytes", string(key), objSize)

	if w.goldenDB.ObjectLimit() <= w.goldenDB.ObjectsCount() {
		msg := fmt.Sprintf("No space left")
		log.Println(msg)
		protocolStatusWrite(conn, statusNoSpaceLeft, []byte(msg))
		return
	}

	// safe to write space

	// make descriptor space
	w.goldenDB.RemakeKey(key)

	valBuff := make([]byte, defaultBuffSize)
	copied := 0
	for {
		sz, e := conn.Read(valBuff)
		copied = copied + sz
		w.goldenDB.Write(key, valBuff[:sz])
		if int64(copied) >= objSize || e != nil {
			break
		}
	}

	msg := fmt.Sprintf("object stored (%d bytes)", uint64(objSize))
	log.Println(msg)
	protocolStatusWrite(conn, statusOK, []byte(msg))
	return
}

func (w *GoldenWorker) cmdRead(conn net.Conn) {
	// read key size
	sizeBuff := make([]byte, 8)
	if _, e := conn.Read(sizeBuff); e != nil {
		msg := fmt.Sprintf("cannot create object - %s", e.Error())
		log.Println(msg)
		protocolStatusWrite(conn, statusProtocolError, []byte(msg))
		return
	}
	keySize := bytesArrayToUint64(sizeBuff)

	// read key
	key := make([]byte, keySize)
	if _, e := conn.Read(key); e != nil {
		msg := fmt.Sprintf("cannot read object key - %s", e.Error())
		log.Println(msg)
		protocolStatusWrite(conn, statusProtocolError, []byte(msg))
		return
	}

	obj, err := w.goldenDB.GetObject(key)
	if err != nil {
		msg := fmt.Sprintf("cannot read object key - %s", err.Error())
		log.Println(msg)
		protocolStatusWrite(conn, statusNotFound, []byte(msg))
		return
	}

	conn.Write([]byte{statusOK})
	conn.Write(intToBytes(len(obj)))
	conn.Write(obj)
}

func (w *GoldenWorker) cmdDelete(conn net.Conn) {
	// read key size
	sizeBuff := make([]byte, 8)
	if _, e := conn.Read(sizeBuff); e != nil {
		msg := fmt.Sprintf("cannot create object - %s", e.Error())
		log.Println(msg)
		protocolStatusWrite(conn, statusProtocolError, []byte(msg))
		return
	}
	keySize := bytesArrayToUint64(sizeBuff)

	// read key
	key := make([]byte, keySize)
	if _, e := conn.Read(key); e != nil {
		msg := fmt.Sprintf("cannot read object key - %s", e.Error())
		log.Println(msg)
		protocolStatusWrite(conn, statusProtocolError, []byte(msg))
		return
	}

	w.goldenDB.Delete(key)
	protocolStatusWrite(conn, statusOK, []byte("ok"))
}

func New(goldenDB *goldendb.GoldenDB) *GoldenWorker {
	return &GoldenWorker{
		goldenDB: goldenDB,
	}
}
