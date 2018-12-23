package limitware

import (
	"sync"
)

type LimitInterface interface {
	update(value interface{})
	read() int
}

type Limit struct {
	prop     LimitInterface
	maxvalue int
	sync.RWMutex
}

type Limitware struct {
	limits []Limit
}
