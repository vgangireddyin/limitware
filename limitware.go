package limitware

import (
	"net/http"
	"sync"
)

func New() *Limitware {
	return &Limitware{limits: make([]Limit, 0)}
}

func (lw *Limitware) Add(l Limit) {
	lw.limits = append(lw.limits, l)
}

func (lw *Limitware) Handler(next http.Handler, fail http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		failure := make(chan bool)
		success := make(chan bool)

		go func() {

			count := len(lw.limits)
			var wg sync.WaitGroup
			wg.Add(count)

			for i := 0; i < count; i++ {
				go func(index int, w *sync.WaitGroup) {
					res := lw.limits[index].Read()
					if res > lw.limits[index].maxvalue {
						failure <- true
					} else {
						w.Done()
					}
				}(i, &wg)
			}

			wg.Wait()
			success <- true
		}()

		select {
		case <-failure:
			fail.ServeHTTP(w, r)
			break
		case <-success:
			next.ServeHTTP(w, r)
			break
		}
	})

}
