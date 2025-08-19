package editor

import "sync"

// ReadyState manages the "ready" state of the editor safely using sync.Cond.
type ReadyState struct {
	ready bool
	cond  *sync.Cond
}

// NewEditorReadyState creates and returns a new editorReadyState
// instance. The editor starts in a "not ready" state.
func NewEditorReadyState() *ReadyState {
	return &ReadyState{
		cond: sync.NewCond(&sync.Mutex{}),
	}
}

// Get returns the current "ready" state of the editor.
func (s *ReadyState) Get() bool {
	s.cond.L.Lock()
	defer s.cond.L.Unlock()
	return s.ready
}

// Set updates the "ready" state of the editor.
func (s *ReadyState) Set(ready bool) {
	s.cond.L.Lock()
	defer s.cond.L.Unlock()
	if s.ready == ready {
		return
	}
	s.ready = ready
	s.cond.Broadcast()
}

// Wait returns a channel that will be closed once the editor becomes ready.
func (s *ReadyState) Wait() <-chan struct{} {
	ch := make(chan struct{})
	go func() {
		s.cond.L.Lock()
		for !s.ready {
			s.cond.Wait()
		}
		s.cond.L.Unlock()
		close(ch)
	}()
	return ch
}
