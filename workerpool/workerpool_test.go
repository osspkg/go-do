package workerpool

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"
)

// простой обработчик: удваивает число, для Data == 0 возвращает ошибку
func simpleHandler(ctx context.Context, task Task[int, int]) (int, error) {
	if task.Data == 0 {
		return 0, errors.New("zero data")
	}
	// небольшая имитация работы, чтобы проверить параллелизм
	time.Sleep(10 * time.Millisecond)
	return task.Data * 2, nil
}

func TestUnit_Pool_SuccessfulProcessing(t *testing.T) {
	p := New(3, simpleHandler)
	p.Start()

	done := make(chan struct{})

	n := 10
	expected := make(map[int]int, n)
	go func() {
		for i := 1; i <= n; i++ {
			expected[i] = i * 2
			p.Send(Task[int, int]{ID: i, Data: i})
		}
		close(done)
	}()

	results := make(map[int]int)
	for i := 0; i < n; i++ {
		res := <-p.Receive()
		if res.Err != nil {
			t.Errorf("unexpected error for task %d: %v", res.TaskID, res.Err)
		}
		results[res.TaskID] = res.Result
	}

	<-done

	for id, want := range expected {
		if got, ok := results[id]; !ok || got != want {
			t.Errorf("task %d: got %d, want %d", id, got, want)
		}
	}

	p.Close()
	// после Close канал результатов должен быть закрыт
	if p.Send(Task[int, int]{ID: 0, Data: 0}) {
		t.Error("send channel should be closed after Close()")
	}
	_, ok := <-p.Receive()
	if ok {
		t.Error("results channel should be closed after Close()")
	}
}

func TestUnit_Pool_HandlerError(t *testing.T) {
	p := New(2, simpleHandler)
	p.Start()

	p.Send(Task[int, int]{ID: 1, Data: 10})
	p.Send(Task[int, int]{ID: 2, Data: 0}) // вызовет ошибку

	res1 := <-p.Receive()
	res2 := <-p.Receive()

	// Порядок не гарантирован, проверяем по ID
	for _, r := range []Result[int, int]{res1, res2} {
		switch r.TaskID {
		case 1:
			if r.Err != nil {
				t.Errorf("task 1 should not have error: %v", r.Err)
			}
			if r.Result != 20 {
				t.Errorf("task 1 result: got %d, want 20", r.Result)
			}
		case 2:
			if r.Err == nil {
				t.Error("task 2 should have error")
			}
		default:
			t.Errorf("unexpected task ID %d", r.TaskID)
		}
	}

	p.Close()
}

func TestUnit_Pool_CloseCancelsWorkers(t *testing.T) {
	p := New(2, simpleHandler)
	p.Start()

	var wg sync.WaitGroup
	wg.Add(1)
	processed := make(chan int, 10)

	// читаем результаты в фоне
	go func() {
		defer wg.Done()
		for res := range p.Receive() {
			processed <- res.TaskID
		}
	}()

	// отправляем несколько задач
	for i := 1; i <= 5; i++ {
		p.Send(Task[int, int]{ID: i, Data: i})
	}

	// даём время на старт обработки
	time.Sleep(50 * time.Millisecond)

	// закрываем пул (отмена контекста)
	p.Close()

	// ждём завершения горутины-читателя
	wg.Wait()
	close(processed)

	// все задачи, которые были отправлены до Close, либо обработаны,
	// либо отброшены; ошибок быть не должно.
	ids := make(map[int]bool)
	for id := range processed {
		if id < 1 || id > 5 {
			t.Errorf("unexpected task id %d", id)
		}
		ids[id] = true
	}
	// Хотя бы часть задач должна была успеть обработаться
	if len(ids) == 0 {
		t.Error("no tasks were processed before close")
	}

	// Отправка после Close не должна блокироваться и не должна приводить к панике
	p.Send(Task[int, int]{ID: 100, Data: 100})
}

func TestUnit_Pool_SendAfterClose(t *testing.T) {
	p := New(1, simpleHandler)
	p.Start()
	p.Close()

	// Send не должен паниковать или зависать
	p.Send(Task[int, int]{ID: 1, Data: 42})
	// канал результатов закрыт, проверяем
	_, ok := <-p.Receive()
	if ok {
		t.Error("expected closed channel after Close")
	}
}

func TestUnit_Pool_ZeroWorkers(t *testing.T) {
	// New с 0 воркерами должен установить минимум 1
	p := New(0, simpleHandler)
	if p.workers != 1 {
		t.Errorf("expected 1 worker, got %d", p.workers)
	}
	p.Start()

	p.Send(Task[int, int]{ID: 1, Data: 5})
	res := <-p.Receive()
	if res.Err != nil || res.Result != 10 {
		t.Errorf("unexpected result: %+v", res)
	}

	p.Close()
}

func TestUnit_Pool_Concurrency(t *testing.T) {
	workers := 10
	tasksCount := 100
	p := New(workers, func(ctx context.Context, task Task[int, int]) (int, error) {
		time.Sleep(1 * time.Millisecond)
		return task.Data, nil
	})
	p.Start()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < tasksCount; i++ {
			p.Send(Task[int, int]{ID: i, Data: i})
		}
	}()

	received := 0
	for res := range p.Receive() {
		if res.Err != nil {
			t.Errorf("unexpected error: %v", res.Err)
		}
		received++
		if received == tasksCount {
			break
		}
	}
	wg.Wait()

	if received != tasksCount {
		t.Errorf("expected %d results, got %d", tasksCount, received)
	}

	p.Close()
}
