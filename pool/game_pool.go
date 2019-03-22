package pool

import "go_resource_pool_demo/model"
import "fmt"

var objCounter = 0

type (
	// GamePool manages a pool of GamePiece objects to ease GC churn
	GamePool struct {
		store    chan *model.GamePiece
		highmark int // max objects in pool during life
	}
)

// Alloc returns a game piece unless your OS is out of memory or 'p' is nil
// 1. first it tries to get an object form the store
// 2. if the store is empty, then it allocs one an return it
func (p *GamePool) Alloc() *model.GamePiece {
	if p == nil {
		return nil
	}
	// non-blocking channel receive
	var obj *model.GamePiece
	select {
	case obj = <-p.store:
	default:
		obj = &model.GamePiece{}
	}

	objCounter++
	obj.Id = objCounter // do some init for a new obj
	return obj
}

// Release always returns the obj to the store and clean up so it might
// not be use by anyone who still holds a reference
func (p *GamePool) Release(obj *model.GamePiece) {
	if p == nil || obj == nil { // no-op essentially
		return
	}
	// don't block trying to add to a full store
	if len(p.store) <= cap(p.store) {
		// clean up obj
		obj.Id = 0
		p.store <- obj
		if p.highmark <= len(p.store) { // accounting stuff
			p.highmark = len(p.store)
		}
	} else {
		obj = nil // else GC will take care of it
	}
}

// New creates a new GamePool manager
func New(size int) (*GamePool, error) {
	if size < 1 {
		return nil, fmt.Errorf("Invalid size value %v, size must be > 0", size)
	}
	p := &GamePool{store: make(chan *model.GamePiece, size)}
	return p, nil
}

func (p *GamePool) String() string {
	if p == nil {
		return "GamePool<nil>"
	}
	return fmt.Sprintf("GamePool high water mark: %v", p.highmark)
}
