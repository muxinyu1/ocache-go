package ocache

import "sync"

type status struct {
	wg sync.WaitGroup
	res Value
	err error
}

type Table struct {
	mtx sync.Mutex
	m map[string]*status
}

func (t *Table) Do(group string, key string, fn func(string, string)(Value, error)) (Value, error) {
	t.mtx.Lock()
	if t.m == nil {
		t.m = make(map[string]*status)
	}
	s, ok := t.m[key]
	if ok {
		t.mtx.Unlock()
		s.wg.Wait()
		return s.res, s.err
	}
	s = new(status)
	s.wg.Add(1)
	t.m[key] = s
	t.mtx.Unlock()
	// 调用fn去拿
	s.res, s.err = fn(group, key)
	s.wg.Done()

	t.mtx.Lock()
	delete(t.m, key)
	t.mtx.Unlock()

	return s.res, s.err
}