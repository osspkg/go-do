package pipline

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestUnit_New(t *testing.T) {
	p := New[string, int]()
	if p == nil {
		t.Fatal("New returned nil")
	}
	if p.pips == nil {
		t.Fatal("pips map is nil")
	}
	if len(p.pips) != 0 {
		t.Fatal("pips should be empty initially")
	}
}

func TestUnit_Set_NilHandler(t *testing.T) {
	p := New[string, int]()
	err := p.Set(Pipe[string, int]{
		Current: "start",
		Next:    "end",
		Handler: nil,
	})
	if err == nil || err.Error() != "nil handler" {
		t.Fatalf("expected 'nil handler' error, got %v", err)
	}
}

func TestUnit_Set_DuplicateCurrent(t *testing.T) {
	p := New[string, int]()
	handler := func(ctx context.Context, arg int) (int, error) { return arg, nil }
	p.Set(Pipe[string, int]{Current: "a", Next: "b", Handler: handler})
	err := p.Set(Pipe[string, int]{Current: "a", Next: "c", Handler: handler})
	if err == nil || err.Error() != "current state is already set" {
		t.Fatalf("expected 'current state is already set' error, got %v", err)
	}
}

func TestUnit_Set_SameCurrentNext(t *testing.T) {
	p := New[string, int]()
	err := p.Set(Pipe[string, int]{
		Current: "a",
		Next:    "a",
		Handler: func(ctx context.Context, arg int) (int, error) { return arg, nil },
	})
	if err == nil || err.Error() != "current and next state is identical" {
		t.Fatalf("expected 'current and next state is identical' error, got %v", err)
	}
}

func TestUnit_Set_Success(t *testing.T) {
	p := New[string, int]()
	h := func(ctx context.Context, arg int) (int, error) { return arg * 2, nil }
	err := p.Set(Pipe[string, int]{Current: "s1", Next: "s2", Handler: h})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(p.pips) != 1 {
		t.Fatal("expected one pipe in map")
	}
}

func TestUnit_Do_SimpleChain(t *testing.T) {
	p := New[string, int]()
	// Цепочка: "init" -> "a" -> "b" -> "c"
	// Удвоение на каждом шаге
	p.Set(Pipe[string, int]{
		Current: "init",
		Next:    "a",
		Handler: func(ctx context.Context, arg int) (int, error) { return arg * 2, nil },
	})
	p.Set(Pipe[string, int]{
		Current: "a",
		Next:    "b",
		Handler: func(ctx context.Context, arg int) (int, error) { return arg + 10, nil },
	})
	p.Set(Pipe[string, int]{
		Current: "b",
		Next:    "c",
		Handler: func(ctx context.Context, arg int) (int, error) { return arg - 3, nil },
	})

	state, val, err := p.Do(context.Background(), "init", 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if state != "c" {
		t.Errorf("expected final state 'c', got %q", state)
	}
	// 5*2=10, +10=20, -3=17
	if val != 17 {
		t.Errorf("expected val 17, got %d", val)
	}
}

func TestUnit_Do_UnknownState(t *testing.T) {
	p := New[string, int]()
	p.Set(Pipe[string, int]{
		Current: "start",
		Next:    "end",
		Handler: func(ctx context.Context, arg int) (int, error) { return arg, nil },
	})

	// Начинаем с состояния, которого нет в pips
	state, val, err := p.Do(context.Background(), "unknown", 42)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if state != "unknown" {
		t.Errorf("expected state 'unknown', got %q", state)
	}
	if val != 42 {
		t.Errorf("expected val 42, got %d", val)
	}
}

func TestUnit_Do_HandlerError(t *testing.T) {
	p := New[string, int]()
	expectedErr := errors.New("handler failed")
	p.Set(Pipe[string, int]{
		Current: "init",
		Next:    "a",
		Handler: func(ctx context.Context, arg int) (int, error) { return arg, nil },
	})
	p.Set(Pipe[string, int]{
		Current: "a",
		Next:    "b",
		Handler: func(ctx context.Context, arg int) (int, error) { return 0, expectedErr },
	})
	p.Set(Pipe[string, int]{
		Current: "b",
		Next:    "c",
		Handler: func(ctx context.Context, arg int) (int, error) { return arg, nil },
	})

	state, val, err := p.Do(context.Background(), "init", 10)
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected error %v, got %v", expectedErr, err)
	}
	// Состояние должно быть тем, на котором случилась ошибка – "a"
	if state != "a" {
		t.Errorf("expected state 'a', got %q", state)
	}
	// Значение – то, что вернул ошибочный обработчик (0)
	if val != 0 {
		t.Errorf("expected val 0, got %d", val)
	}
}

func TestUnit_Do_ContextCancellation(t *testing.T) {
	p := New[string, int]()
	// Создаём цепочку с блокировкой на втором шаге
	p.Set(Pipe[string, int]{
		Current: "start",
		Next:    "blocked",
		Handler: func(ctx context.Context, arg int) (int, error) {
			// Не блокируем, просто пропускаем
			return arg + 1, nil
		},
	})
	p.Set(Pipe[string, int]{
		Current: "blocked",
		Next:    "end",
		Handler: func(ctx context.Context, arg int) (int, error) {
			<-ctx.Done()
			return arg, nil
		},
	})

	ctx, cancel := context.WithCancel(context.Background())
	errCh := make(chan error, 1)
	var state string
	var val int
	go func() {
		var e error
		state, val, e = p.Do(ctx, "start", 100)
		errCh <- e
	}()

	// Даём горутине дойти до заблокированного handler'а
	time.Sleep(50 * time.Millisecond)
	cancel()

	select {
	case err := <-errCh:
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("expected context.Canceled, got %v", err)
		}
		if state != "end" {
			t.Errorf("expected state 'blocked', got %q", state)
		}
		if val != 101 { // значение после первого шага
			t.Errorf("expected val 101, got %d", val)
		}
	case <-time.After(time.Second):
		t.Fatal("test timed out")
	}
}

func TestUnit_Do_EmptyPipeline(t *testing.T) {
	p := New[string, int]()
	state, val, err := p.Do(context.Background(), "any", 123)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if state != "any" {
		t.Errorf("expected state 'any', got %q", state)
	}
	if val != 123 {
		t.Errorf("expected val 123, got %d", val)
	}
}
