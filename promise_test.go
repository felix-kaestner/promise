package promise

import (
	"net/http"
	"reflect"
	"testing"
	"time"
)

func isNil(i any) bool {
	if i == nil {
		return true
	}

	v := reflect.ValueOf(i)
	switch v.Kind() {
	case reflect.Chan,
		reflect.Func,
		reflect.Map,
		reflect.Ptr,
		reflect.UnsafePointer,
		reflect.Interface,
		reflect.Slice:
		return v.IsNil()
	}

	return false
}

func assertNil(t *testing.T, actual any) {
	if !isNil(actual) {
		t.Errorf("Test %s: Expected value to be nil, Received `%v` (type %v)", t.Name(), actual, reflect.TypeOf(actual))
	}
}

func assertNotNil(t *testing.T, actual any) {
	if isNil(actual) {
		t.Errorf("Test %s: Expected value to not be nil, Received `%v` (type %v)", t.Name(), actual, reflect.TypeOf(actual))
	}
}

func assertEqual(t *testing.T, expected, actual any) {
	if (isNil(expected) && isNil(actual)) || reflect.DeepEqual(expected, actual) {
		return
	}

	t.Errorf("Test %s: Expected `%v` (type %v), Received `%v` (type %v)", t.Name(), expected, reflect.TypeOf(expected), actual, reflect.TypeOf(actual))
}

func TestPromise(t *testing.T) {
	{
		p := New(func() (*http.Response, error) {
			return http.Get("https://jsonplaceholder.typicode.com/posts/1")
		})

		p.Then(func(r *http.Response) {
			assertEqual(t, 200, r.StatusCode)
		}, func(err error) {
			assertNil(t, err)
		})

		p.OnSuccess(func(r *http.Response) {
			assertEqual(t, 200, r.StatusCode)
		})

		p.OnFailure(func(err error) {
			assertNil(t, err)
		})

		res, err := p.Await()
		assertEqual(t, 200, res.StatusCode)
		assertNil(t, err)

		assertEqual(t, res, p.AwaitOr(&http.Response{}))
	}
	{
		p := New(func() (*http.Response, error) {
			return http.Get("https://abc.def")
		})

		p.Then(func(r *http.Response) {
			assertNil(t, r)
		}, func(err error) {
			assertNotNil(t, err)
		})

		p.OnSuccess(func(r *http.Response) {
			assertNil(t, r)
		})

		p.OnFailure(func(err error) {
			assertNotNil(t, err)
		})

		res, err := p.Await()
		assertNil(t, res)
		assertNotNil(t, err)

		defaultRes := &http.Response{}
		res = p.AwaitOr(defaultRes)
		assertEqual(t, defaultRes, res)
	}
}

func TestDone(t *testing.T) {
	p := New(func() (bool, error) {
		<-time.After(100 * time.Millisecond)
		return true, nil
	})

	select {
	case <-p.Done():
		ok, err := p.Await()
		assertEqual(t, true, ok)
		assertNil(t, err)
	case <-time.After(5000 * time.Millisecond):
		t.Errorf("Test %s: Expected promise to be finished", t.Name())
	}
}
