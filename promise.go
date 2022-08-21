package promise

import (
	"sync"
)

// Promise is a object that can be used to get the
// result of an async operation.
type Promise[T any] interface {
	// Await returns the result of the async operation.
	// If the operation has not completed it will block
	// until it does complete.
	// Otherwise it will return the result of the operation
	// immediately.
	Await() (T, error)

	// AwaitOr returns the result of the async operation
	// or a default value if the operation was not successful.
	AwaitOr(defaultValue T) T

	// Then executes the given functions when the promise
	// is either fulfilled or rejected respectively.
	// The functions are executed in a goroutine.
	Then(func(T), func(error))

	// OnSuccess executes the given function if the
	// promise is fulfilled.
	// The functions is executed in a goroutine.
	OnSuccess(func(T))

	// OnFailure executes the given function if the
	// promise is rejected.
	// The functions is executed in a goroutine.
	OnFailure(func(error))

	// Done returns a channel that's closed when the work done on behalf of this
	// promise is finished. Successive calls to Done return the same value.
	// The close of the Done channel may happen asynchronously,
	// after the function returns.
	//
	// Done is provided for use in select statements:
	//
	// // Produce generates a value by awaiting the promise and sending it to the out
	// // channel. It returns an error if the promise returned an error or ctx.Done is closed.
	// func Produce(ctx context.Context, p Promise[Value], out chan<- Value) error {
	// 	for {
	// 		select {
	// 		case <-ctx.Done():
	// 			return ctx.Err()
	// 		case <-p.Done():
	// 			v, err := p.Await() // returns immediately since the promise is already resolved
	// 			if err != nil {
	// 				return err
	// 			}
	// 			out <- v
	// 		}
	// 	}
	// }
	Done() <-chan struct{}
}

// promise is an implementation of the Promise interface.
type promise[T any] struct {
	wg    sync.WaitGroup
	once  sync.Once
	done  chan struct{}
	value T
	err   error
	fn    func() (T, error)
}

func (p *promise[T]) get() {
	p.once.Do(func() {
		p.wg.Add(1)
		go func() {
			defer p.wg.Done()
			p.value, p.err = p.fn()
			close(p.done)
		}()
	})
	p.wg.Wait()
}

func (p *promise[T]) Await() (T, error) {
	p.get()
	return p.value, p.err
}

func (p *promise[T]) AwaitOr(defaultValue T) T {
	if value, err := p.Await(); err == nil {
		return value
	}
	return defaultValue
}

func (p *promise[T]) Then(onSuccess func(T), onFailure func(error)) {
	go func() {
		if p.get(); p.err == nil {
			onSuccess(p.value)
		} else {
			onFailure(p.err)
		}
	}()
}

func (p *promise[T]) OnSuccess(fn func(T)) {
	go func() {
		if p.get(); p.err == nil {
			fn(p.value)
		}
	}()
}

func (p *promise[T]) OnFailure(fn func(error)) {
	go func() {
		if p.get(); p.err != nil {
			fn(p.err)
		}
	}()
}

func (p *promise[T]) Done() <-chan struct{} {
	go p.get()
	return p.done
}

// New returns a new Promise of type T.
// The given function will be executed in a goroutine.
// The function should return the result of an async operation
// or an error if the operation failed.
func New[T any](fn func() (T, error)) Promise[T] {
	return &promise[T]{
		fn:   fn,
		done: make(chan struct{}),
	}
}
