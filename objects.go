package quickjs

import "sync"

type ObjectId int64

var refs struct {
	sync.RWMutex
	objs map[ObjectId]interface{}
	next ObjectId
}

func init() {
	refs.Lock()
	defer refs.Unlock()

	refs.objs = make(map[ObjectId]interface{})
	refs.next = 1000
}

func NewObjectId(obj interface{}) ObjectId {
	refs.Lock()
	defer refs.Unlock()

	id := refs.next
	refs.objs[id] = obj
	refs.next++

	return id
}

func (id ObjectId) Get() (interface{}, bool) {
	refs.RLock()
	defer refs.RUnlock()

	if id.IsNil() {
		return nil, false
	}

	obj, ok := refs.objs[id]
	return obj, ok
}

func (id ObjectId) IsNil() bool {
	return id == 0
}

func (id *ObjectId) Free() {
	refs.Lock()
	defer refs.Unlock()

	if id.IsNil() {
		return
	}

	_, ok := refs.objs[*id]
	if ok {
		delete(refs.objs, *id)
	}

	*id = 0
}
