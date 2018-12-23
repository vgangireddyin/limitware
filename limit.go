package limitware

func (l *Limit) Read() int {
	l.RLock()
	val := l.prop.read()
	l.RUnlock()
	return val
}

func (l *Limit) Update(input int) {
	l.Lock()
	l.prop.update(input)
	l.Unlock()
}
