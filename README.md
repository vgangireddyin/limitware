# Limitware

Limitware is a generic middleware framework to limit available resources for Golang.

  - Simple
  - Concurrent verification
  - Concurrent execution
  - Thread safe
  - Easy to extend

## Usage

It is assumed that every resource property to be limited can be accessed and modified using `read()` and `update()` function calls respectively.

```golang
type LimitInterface interface {
	update(value interface{})
	read() int
}
```

You can define resource property by implementation these two functions.

```
type a struct {
	fast int
}

func (aa *a) update(value interface{}) {
	aa.fast = value.(int)
}

func (aa *a) read() int {
	return aa.fast
}
```

These type of resources can be covered into `Limit` object to access them in thread safe manner using `Limit.Read` and `Limit.Update` function calls.

```golang
type Limit struct {
	prop     LimitInterface
	maxvalue int
	sync.RWMutex
}

limit := limitware.Limit{prop: &a{fast: 0}, maxvalue: 100}
limit.Update(100)
curr := limit.Read()
```

`Limitware` is collection of such `Limit` objects which need to verified before accessing the resources. These `Limit` objects can be add using `Limitware.Add` function call. The verification is done by `Limitware.Handler`, and it can be replaced your regular `http.Handler`.

```golang
type Limitware struct {
	limits []Limit
}

lw := limitware.New()
lw.Add(limit)

func next(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "from next")
}

func fail(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	io.WriteString(w, "from fail")
}

nextHandler := http.HandlerFunc(next)
failHandler := http.HandlerFunc(fail)

http.Handle("/", lw.Handler(nextHandler, failHandler)

```

Feel free to open an issue in case any modification are required. 

### TODO !

  - Add typesafe mechanism

License
----

MIT
