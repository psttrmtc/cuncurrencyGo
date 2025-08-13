package memo

// Func is the type of the function to memoize
type Func func(key string) (any, error)

// A result is the result of calling a Func
type result struct {
	value any
	err   error
}

type entry struct {
	res   result
	ready chan struct{} // closed when res is ready
}

// A request is a message requesting that the Func be applied to key
type request struct {
	key      string
	response chan<- result // the client wants a single result
}

type Memo struct{ requests chan request }

// New returns a memoize of f. Client must subsequently call Close()
func New(f Func) *Memo {
	memo := &Memo{requests: make(chan request)}
	go memo.server(f)
	return memo
}

func (memo *Memo) Get(key string) (any, error) {
	response := make(chan result)
	memo.requests <- request{key, response}
	res := <-response
	return res.value, res.err
}

func (memo *Memo) Close() { close(memo.requests) }

func (memo *Memo) server(f Func) {
	cache := make(map[string]*entry)
	for req := range memo.requests {
		e := cache[req.key]
		if e == nil {
			e = &entry{ready: make(chan struct{})}
			cache[req.key] = e
			go e.call(f, req.key)
		}
		go e.deliver(req.response)
	}
}
func (e *entry) call(f Func, key string) {
	// evaluate the function
	e.res.value, e.res.err = f(key)
	//broadcast the ready condition
	close(e.ready)
}
func (e *entry) deliver(response chan<- result) {
	// wait for the ready condition
	<-e.ready
	// Send the result to the client
	response <- e.res
}


