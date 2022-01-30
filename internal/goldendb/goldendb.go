package goldendb

import (
	"fmt"
	"sync"
)

type GoldenDB struct {
	objLimit uint64
	memLimit uint64
	objCount uint64

	storage    map[string][]byte
	storageMtx sync.Mutex
}

func (d *GoldenDB) ObjectLimit() uint64 {
	return d.objLimit
}

func (d *GoldenDB) MemoryLimit() uint64 {
	return d.memLimit
}

func (d *GoldenDB) ObjectsCount() uint64 {
	return d.objCount
}

func (d *GoldenDB) InitDB(objLimit uint64, memLimit uint64) {
	d.objCount = 0
	d.objLimit = objLimit
	d.memLimit = memLimit

	d.storage = make(map[string][]byte)
}

func (d *GoldenDB) RemakeKey(key []byte) {
	d.storageMtx.Lock()
	defer d.storageMtx.Unlock()

	d.storage[string(key)] = []byte{}
	d.objCount++
}

func (d *GoldenDB) Write(key []byte, bytes []byte) {
	d.storageMtx.Lock()
	defer d.storageMtx.Unlock()
	obj := append(d.storage[string(key)], bytes...)
	d.storage[string(key)] = obj
}

func (d *GoldenDB) GetObject(key []byte) ([]byte, error) {
	d.storageMtx.Lock()
	defer d.storageMtx.Unlock()

	if v, ok := d.storage[string(key)]; ok {
		return v, nil
	} else {
		return []byte{}, fmt.Errorf("object not found")
	}
}

func (d *GoldenDB) Delete(key []byte) {
	d.storageMtx.Lock()
	defer d.storageMtx.Unlock()

	d.storage[string(key)] = nil
	d.objCount--
}
