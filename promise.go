package go_promise

import (
	"errors"
	"sort"
	"time"
)

var (
	ErrRange error = errors.New("Range error: count, count should be less than or equal to the length of iterable")
)

var TimeoutWrapper = func(workload Workload, timeout time.Duration) Command {
	return func(args ...interface{}) interface{} {
		c := make(chan interface{})
		go func() {
			r := workload.Command()
			c <- r
		}()
		select {
		case data := <-c:
			return data
		case <-time.After(timeout):
			return workload.FallbackVal
		}
	}
}

type Promiser interface {
	AllSettled(iterable []Workload) []Result
	All(iterable []Workload) []Result
	Race(iterable []Workload) *Result
	RaceAll(iterable []Workload) []Result
}

type promise struct {
	Options Options
}

type Options struct {
	Timeout time.Duration
}

type Command func(args ...interface{}) interface{}

type Workload struct {
	Command     Command
	FallbackVal interface{}
}

type Result struct {
	Idx int
	R   interface{}
}

func NewPromise(options Options) *promise {
	return &promise{Options: options}
}

func (p *promise) AllSettled(iterable []Workload) []Result {
	if len(iterable) == 0 {
		return nil
	}
	c := make(chan struct{})
	r := make([]Result, len(iterable))
	for idx, iter := range iterable {
		iter := iter
		idx := idx
		go func() {
			r[idx].Idx = idx
			r[idx].R = iter.Command()
			c <- struct{}{}
		}()
	}
	for i := 0; i < len(iterable); i++ {
		<-c
	}
	sort.Slice(r, func(i, j int) bool {
		return r[i].Idx < r[j].Idx
	})
	close(c)
	return r
}

func (p *promise) All(iterable []Workload) []Result {
	if len(iterable) == 0 {
		return nil
	}
	c := make(chan Result, len(iterable))
	errc := make(chan Result, len(iterable))
	r := []Result{}
	for idx, iter := range iterable {
		iter := iter
		idx := idx
		go func() {
			result := iter.Command()
			if err, ok := result.(error); ok {
				errc <- Result{Idx: idx, R: err}
			} else {
				c <- Result{Idx: idx, R: result}
			}
		}()
	}
	for {
		select {
		case result := <-errc:
			return []Result{result}
		case result := <-c:
			r = append(r, result)
			if len(r) == len(iterable) {
				sort.Slice(r, func(i, j int) bool {
					return r[i].Idx < r[j].Idx
				})
				return r
			}
		}
	}
}

// Race
// The `Idx` field of `Result` struct will be always zero value 0
func (p *promise) Race(iterable []Workload) *Result {
	if len(iterable) == 0 {
		return nil
	}
	c := make(chan *Result, len(iterable))
	for _, iter := range iterable {
		iter := iter
		go func() {
			r := Result{R: iter.Command()}
			c <- &r
		}()
	}
	return <-c
}

// RaceAll
// https://stackoverflow.com/a/48578424/6463558
func (p *promise) RaceAll(iterable []Workload) []Result {
	if len(iterable) == 0 {
		return nil
	}
	if p.Options.Timeout == 0 {
		return p.All(iterable)
	}
	ws := []Workload{}
	rc := make(chan *Result)
	rs := []*Result{}
	for idx, iter := range iterable {
		iter := iter
		idx := idx
		go func() {
			r := p.Race([]Workload{
				iter,
				{Command: TimeoutWrapper(iter, p.Options.Timeout)},
			})
			r.Idx = idx
			rc <- r
		}()
	}
	for i := 0; i < len(iterable); i++ {
		rs = append(rs, <-rc)
	}
	sort.Slice(rs, func(i, j int) bool {
		return rs[i].Idx < rs[j].Idx
	})
	for _, r := range rs {
		result := r.R
		ws = append(ws, Workload{Command: func(args ...interface{}) interface{} {
			return result
		}})
	}

	return p.All(ws)
}

func (p *promise) Any(iterable []Workload) []Result {
	return p.Some(iterable, 1)
}

func (p *promise) Some(iterable []Workload, count uint) []Result {
	if len(iterable) < int(count) {
		return []Result{{R: ErrRange}}
	}
	if len(iterable) == 0 || count == 0 {
		return nil
	}
	errc := make(chan Result, len(iterable))
	c := make(chan Result, len(iterable))
	r := []Result{}
	errs := []Result{}
	for idx, iter := range iterable {
		idx := idx
		iter := iter
		go func() {
			result := iter.Command()
			if err, ok := result.(error); ok {
				errc <- Result{Idx: idx, R: err}
			} else {
				c <- Result{Idx: idx, R: result}
			}
		}()
	}
	for {
		select {
		case result := <-c:
			r = append(r, result)
			if len(r) == int(count) {
				sort.Slice(r, func(i, j int) bool {
					return r[i].Idx < r[j].Idx
				})
				return r
			}
		case result := <-errc:
			errs = append(errs, result)
			if (len(errs)+len(r)) == len(iterable) && len(r) < int(count) {
				sort.Slice(errs, func(i, j int) bool {
					return errs[i].Idx < errs[j].Idx
				})
				return errs
			}
		}
	}

}
