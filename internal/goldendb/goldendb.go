package goldendb

import "fmt"

type GoldenDB struct {
	objLimit uint64
	memLimit uint64
	objCount uint64

	opStorage chan map[string][]byte
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
	
	d.opStorage = make(chan map[string][]byte, 1)
}

func (d *GoldenDB) RemakeKey(key []byte) {
	m := <-d.opStorage
	m[string(key)] = []byte{}
	d.objCount++
	d.opStorage <- m
}

func (d *GoldenDB) Write(key []byte, bytes []byte) {
	m := <-d.opStorage
	obj := append(m[string(key)], bytes...)
	m[string(key)] = obj
	d.opStorage <- m
}

func (d *GoldenDB) GetObject(key []byte) ([]byte, error) {
	m := <-d.opStorage
	if v, ok := m[string(key)]; ok {
		d.opStorage <- m
		return v, nil
	} else {
		d.opStorage <- m
		return []byte{}, fmt.Errorf("object not found")
	}
}

func (d *GoldenDB) Delete(key []byte) {
	m := <-d.opStorage
	m[string(key)] = nil
	d.opStorage <- m
	d.objCount--
}
