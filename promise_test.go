package go_promise_test

import (
	"errors"
	"reflect"
	"runtime"
	"testing"
	"time"

	go_promise "github.com/mrdulin/go-promise"
)

func workloadGen(i interface{}, options ...interface{}) go_promise.Workload {
	var sec time.Duration = 1
	var fallbackVal interface{}
	if len(options) != 0 {
		sec = options[0].(time.Duration)
	}
	if len(options) > 1 {
		fallbackVal = options[1]
	}
	command := func(args ...interface{}) interface{} {
		time.Sleep(time.Second * sec)
		return i
	}

	return go_promise.Workload{Command: command, FallbackVal: fallbackVal}
}

func TestPromise_AllSettled(t *testing.T) {
	p := go_promise.NewPromise(go_promise.Options{})
	t.Run("should settle all workloads and return results", func(t *testing.T) {
		var d time.Duration = 3
		w1 := workloadGen(1, d)
		w2 := workloadGen(2, d)
		w3 := workloadGen(3, d)
		ws := []go_promise.Workload{w1, w2, w3}
		got := p.AllSettled(ws)
		want := []go_promise.Result{{Idx: 0, R: 1}, {Idx: 1, R: 2}, {Idx: 2, R: 3}}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("got: %+v, want: %+v", got, want)
		}
	})
	t.Run("should settle empty workloads and return nil", func(t *testing.T) {
		ws := []go_promise.Workload{}
		got := p.AllSettled(ws)
		var want []go_promise.Result
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("got: %#v, want: %#v", got, want)
		}
	})

	t.Run("should handle many workloads concurrency", func(t *testing.T) {
		ws := []go_promise.Workload{}
		c := 100
		for i := 0; i < c; i++ {
			ws = append(ws, workloadGen(i))
		}
		t.Logf("goroutine num: %d", runtime.NumGoroutine())
		got := p.AllSettled(ws)
		time.Sleep(time.Second * 5)
		t.Logf("goroutine num: %d", runtime.NumGoroutine())
		if len(got) != c {
			t.Fatalf("got: %+v, want: %+v", got, c)
		}
	})
}

func TestPromise_All(t *testing.T) {
	p := go_promise.NewPromise(go_promise.Options{})
	t.Run("should handle all workloads", func(t *testing.T) {
		w1 := workloadGen(1)
		w2 := workloadGen(2)
		w3 := workloadGen(3)
		ws := []go_promise.Workload{w1, w2, w3}
		got := p.All(ws)
		want := []go_promise.Result{{Idx: 0, R: 1}, {Idx: 1, R: 2}, {Idx: 2, R: 3}}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("got: %#v, want: %#v", got, want)
		}
	})

	t.Run("should reject the promise when any workload is rejected", func(t *testing.T) {
		err := errors.New("network")
		ws := []go_promise.Workload{}
		for i := 0; i < 1000; i++ {
			ws = append(ws, workloadGen(err))
		}
		//w1 := workloadGen(err)
		//w2 := workloadGen(err)
		//w3 := workloadGen(err)
		//ws := []go_promise.Workload{w1, w2, w3}
		t.Logf("goroutine num: %d", runtime.NumGoroutine())
		got := p.All(ws)
		time.Sleep(time.Second * 5)
		t.Logf("goroutine num: %d", runtime.NumGoroutine())
		if len(got) != 1 {
			t.Fatalf("should got one result, got: %+v", len(got))
		}
		if !reflect.DeepEqual(got[0].R, err) {
			t.Fatalf("got: %+v, want: %+v", got, err)
		}
	})
}

func TestPromise_Race(t *testing.T) {
	p := go_promise.NewPromise(go_promise.Options{})
	t.Run("should race all resolved promise", func(t *testing.T) {
		var t1 time.Duration = 2
		var t2 time.Duration = 1
		w1 := workloadGen(1, t1)
		w2 := workloadGen(2, t1)
		w3 := workloadGen(3, t2)
		ws := []go_promise.Workload{w1, w2, w3}
		got := p.Race(ws)
		want := &go_promise.Result{Idx: 0, R: 3}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("got: %+v, want: %+v", got, want)
		}
	})

	t.Run("should race promises with a rejected promise", func(t *testing.T) {
		err := errors.New("network")
		var t1 time.Duration = 2
		var t2 time.Duration = 1
		w1 := workloadGen(1, t1)
		w2 := workloadGen(2, t1)
		w3 := workloadGen(err, t2)
		ws := []go_promise.Workload{w1, w2, w3}
		got := p.Race(ws)
		want := &go_promise.Result{Idx: 0, R: err}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("got: %+v, want: %+v", got, want)
		}
	})

	t.Run("should not goroutine memory leak", func(t *testing.T) {
		ws := []go_promise.Workload{}
		var d time.Duration = 1
		for i := 0; i < 100; i++ {
			ws = append(ws, workloadGen(i, d))
		}
		t.Logf("goroutine num: %d", runtime.NumGoroutine())
		p.Race(ws)
		time.Sleep(time.Second * 5)
		t.Logf("goroutine num: %d", runtime.NumGoroutine())
	})

	t.Run("should work fine with TimeoutWrapper", func(t *testing.T) {
		ws := []go_promise.Workload{
			workloadGen(1, time.Duration(2)),
			{Command: go_promise.TimeoutWrapper(workloadGen(2, time.Duration(2), 'a'), 1*time.Second)},
		}
		got := p.Race(ws)
		want := &go_promise.Result{Idx: 0, R: 'a'}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("got: %#v, want: %+v", got, want)
		}
	})
}

func TestPromise_RaceAll(t *testing.T) {
	t.Run("should resolve all promises", func(t *testing.T) {
		p := go_promise.NewPromise(go_promise.Options{})
		w1 := workloadGen(1)
		w2 := workloadGen(2)
		w3 := workloadGen(3)
		ws := []go_promise.Workload{w1, w2, w3}
		got := p.RaceAll(ws)
		want := []go_promise.Result{{Idx: 0, R: 1}, {Idx: 1, R: 2}, {Idx: 2, R: 3}}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("got: %+v, want: %+v", got, want)
		}
	})

	t.Run("should rejected if any promise is rejected", func(t *testing.T) {
		p := go_promise.NewPromise(go_promise.Options{})
		err := errors.New("network")
		w1 := workloadGen(1)
		w2 := workloadGen(2)
		w3 := workloadGen(err)
		ws := []go_promise.Workload{w1, w2, w3}
		got := p.RaceAll(ws)
		want := []go_promise.Result{{Idx: 2, R: err}}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("got: %+v, want: %+v", got, want)
		}
	})

	t.Run("should return fallback value of the workload if it's timeout", func(t *testing.T) {
		p := go_promise.NewPromise(go_promise.Options{Timeout: 2 * time.Second})
		w1 := workloadGen(1, time.Duration(3), 'a')
		w2 := workloadGen(2, time.Duration(1))
		w3 := workloadGen(3, time.Duration(1))
		ws := []go_promise.Workload{w1, w2, w3}
		got := p.RaceAll(ws)
		want := []go_promise.Result{{Idx: 0, R: 'a'}, {Idx: 1, R: 2}, {Idx: 2, R: 3}}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("got: %+v, want: %+v", got, want)
		}
	})

	t.Run("should not goroutine memory leak", func(t *testing.T) {
		p := go_promise.NewPromise(go_promise.Options{Timeout: 2 * time.Second})
		ws := []go_promise.Workload{}
		c := 100
		for i := 0; i < c; i++ {
			ws = append(ws, workloadGen(i))
		}
		t.Logf("goroutine num: %d", runtime.NumGoroutine())
		got := p.RaceAll(ws)
		time.Sleep(time.Second * 5)
		t.Logf("goroutine num: %d", runtime.NumGoroutine())
		if len(got) != c {
			t.Fatalf("got: %+v, want: %+v", got, c)
		}
	})
}