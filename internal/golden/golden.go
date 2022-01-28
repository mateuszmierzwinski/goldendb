package golden

import (
	"golden/internal/goldendb"
	"golden/internal/goldenworker"
	"log"
	"net"
)

const defaultMaxWorkersCount = uint64(200)

type Golden struct {
	ObjLimit   uint64
	MemLimit   uint64
	MaxWorkers uint64
	Addr       string
	isBinded   bool

	workers  chan *goldenworker.GoldenWorker
	goldenDB *goldendb.GoldenDB
}

func (g *Golden) Bind(addr string) (err error) {
	var listener net.Listener

	g.Addr = addr

	g.loadWorkers()
	g.goldenDB = &goldendb.GoldenDB{}
	g.goldenDB.InitDB(g.ObjLimit, g.MemLimit)

	if listener, err = net.Listen("tcp", addr); err != nil {
		return
	}
	defer listener.Close()

	for {
		conn, errAccept := listener.Accept()
		if errAccept != nil {
			log.Println(errAccept)
		}
		log.Printf("Incomming: %s @ %s", conn.RemoteAddr().String(), conn.RemoteAddr().Network())
		worker := <-g.workers
		worker.DispatchAndServe(conn, g.workers)
	}

	return nil
}

func (g *Golden) loadWorkers() {
	if g.MaxWorkers == 0 {
		g.MaxWorkers = defaultMaxWorkersCount
	}

	g.workers = make(chan *goldenworker.GoldenWorker, g.MaxWorkers)
	for i := uint64(0); i < g.MaxWorkers; i++ {
		g.workers <- goldenworker.New(g.goldenDB)
	}
}
