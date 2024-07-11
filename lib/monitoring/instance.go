package monitoring

import (
	"expvar"
	"strconv"
	"sync"
)

const defaultInstCapacity = 10000

func NewInstanceTracker(name string) *InstanceTracker {
	v := &InstanceTracker{ids: make(map[int]struct{}, defaultInstCapacity)}
	expvar.Publish(name, v)
	return v
}

type InstanceTracker struct {
	mu  sync.Mutex
	ids map[int]struct{}
	max int
}

func (u *InstanceTracker) String() string {
	u.mu.Lock()
	defer u.mu.Unlock()
	return strconv.Itoa(len(u.ids))
}

func (u *InstanceTracker) OnStart(id int) {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.ids[id] = struct{}{}
	u.max = max(u.max, len(u.ids))
}

func (u *InstanceTracker) OnFinish(id int) {
	u.mu.Lock()
	defer u.mu.Unlock()
	delete(u.ids, id)
}

func (u *InstanceTracker) Flush() int {
	u.mu.Lock()
	defer u.mu.Unlock()
	res := u.max
	u.max = len(u.ids)
	return res
}
