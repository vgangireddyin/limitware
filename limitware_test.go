package limitware

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// defining required limits
type a struct {
	fast int
}

func (aa *a) update(value interface{}) {
	aa.fast = value.(int)
}

func (aa *a) read() int {
	return aa.fast
}

type b struct {
	slow int
}

func (bb *b) update(value interface{}) {
	bb.slow = value.(int)
}

func (bb *b) read() int {
	time.Sleep(1 * time.Second)
	return bb.slow
}

type c struct {
	veryslow int
}

func (cc *c) update(value interface{}) {
	cc.veryslow = value.(int)
}

func (cc *c) read() int {
	time.Sleep(11 * time.Second)
	return cc.veryslow
}

// Http handlers using in test
func next(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "from next")
}

func fail(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	io.WriteString(w, "from fail")
}

func TestLimit(t *testing.T) {
	fastlimit := Limit{prop: &a{fast: 0}, maxvalue: 100}
	fastlimit.Update(100)
	if fastlimit.Read() != 100 {
		t.Error("Limit object Update and Read operations not working.", fastlimit.Read())
	}
}

func TestLimitware(t *testing.T) {
	lw := New()
	fastlimit := Limit{prop: &a{fast: 0}, maxvalue: 100}
	lw.Add(fastlimit)
	if len(lw.limits) != 1 {
		t.Error("Adding to limitware is not working.", len(lw.limits))
	}
}

func TestHandler1(t *testing.T) {

	//create limitware
	lw := New()
	fastlimit := Limit{prop: &a{fast: 0}, maxvalue: 100}
	lw.Add(fastlimit)
	if len(lw.limits) != 1 {
		t.Error("Adding to limitware is not working.", len(lw.limits))
	}

	//create a request
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	nextHandler := http.HandlerFunc(next)
	failHandler := http.HandlerFunc(fail)

	handler := lw.Handler(nextHandler, failHandler)

	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler Test 1: handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestHandler2(t *testing.T) {

	//create limitware
	lw := New()
	fastlimit := Limit{prop: &a{fast: 0}, maxvalue: 100}
	lw.Add(fastlimit)

	//create a request
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	nextHandler := http.HandlerFunc(next)
	failHandler := http.HandlerFunc(fail)

	handler := lw.Handler(nextHandler, failHandler)

	fastlimit.Update(101)

	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("Handler Test 2: handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}
}

func TestHandler3(t *testing.T) {

	//create limitware
	lw := New()
	fastlimit := Limit{prop: &a{fast: 0}, maxvalue: 100}
	slowlimit := Limit{prop: &b{slow: 0}, maxvalue: 100}
	lw.Add(fastlimit)
	lw.Add(slowlimit)

	//create a request
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	nextHandler := http.HandlerFunc(next)
	failHandler := http.HandlerFunc(fail)

	handler := lw.Handler(nextHandler, failHandler)

	fastlimit.Update(101)
	slowlimit.Update(1)

	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("Handler Test 3: handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}
}

func TestHandler4(t *testing.T) {

	//create limitware
	lw := New()
	fastlimit := Limit{prop: &a{fast: 0}, maxvalue: 100}
	slowlimit := Limit{prop: &b{slow: 0}, maxvalue: 100}
	veryslowlimit := Limit{prop: &c{veryslow: 0}, maxvalue: 100}
	lw.Add(fastlimit)
	lw.Add(slowlimit)
	lw.Add(veryslowlimit)

	//create a request
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	nextHandler := http.HandlerFunc(next)
	failHandler := http.HandlerFunc(fail)

	handler := lw.Handler(nextHandler, failHandler)

	fastlimit.Update(101)

	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("Handler Test 3: handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}
}
