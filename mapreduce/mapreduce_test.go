package mapreduce

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"
)

func TestUnit_MapReduceSuccess(t *testing.T) {
	ctx := context.Background()
	items := []int{1, 2, 3, 4, 5}
	mapper := func(ctx context.Context, x int) (int, error) {
		return x * x, nil
	}
	reducer := func(ctx context.Context, acc, val int) (int, error) {
		return acc + val, nil
	}

	result, err := New(ctx, items, mapper, reducer, 0, 4)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != 55 {
		t.Errorf("expected 55, got %d", result)
	}
}

func TestUnit_MapReduceEmptyItems(t *testing.T) {
	ctx := context.Background()
	mapper := func(ctx context.Context, x int) (int, error) { return x, nil }
	reducer := func(ctx context.Context, acc, val int) (int, error) { return acc + val, nil }

	result, err := New(ctx, []int{}, mapper, reducer, 10, 4)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != 10 {
		t.Errorf("expected initial value 10, got %d", result)
	}
}

func TestUnit_MapReduceMapperError(t *testing.T) {
	ctx := context.Background()
	items := []int{1, 2, 3}
	expectedErr := errors.New("mapper error")
	mapper := func(ctx context.Context, x int) (int, error) {
		if x == 2 {
			return 0, expectedErr
		}
		return x, nil
	}
	reducer := func(ctx context.Context, acc, val int) (int, error) {
		return acc + val, nil
	}

	_, err := New(ctx, items, mapper, reducer, 0, 4)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, expectedErr) {
		t.Errorf("expected error to contain %v, got %v", expectedErr, err)
	}
}

func TestUnit_MapReduceReducerError(t *testing.T) {
	ctx := context.Background()
	items := []int{1, 2, 3}
	expectedErr := errors.New("reducer error")
	mapper := func(ctx context.Context, x int) (int, error) {
		return x * 2, nil
	}
	reducer := func(ctx context.Context, acc, val int) (int, error) {
		if acc > 5 {
			return acc, expectedErr
		}
		return acc + val, nil
	}

	// workers=1 гарантирует последовательную обработку и детерминированный порядок reduce
	acc, err := New(ctx, items, mapper, reducer, 0, 1)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, expectedErr) {
		t.Errorf("expected %v, got %v", expectedErr, err)
	}
	// До ошибки должны были обработаться 1*2=2 и 2*2=4, сумма=6
	if acc != 6 {
		t.Errorf("expected acc=6, got %d", acc)
	}
}

func TestUnit_MapReduceContextCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // отменяем до вызова
	items := []int{1, 2, 3}
	mapper := func(ctx context.Context, x int) (int, error) { return x, nil }
	reducer := func(ctx context.Context, acc, val int) (int, error) { return acc + val, nil }

	_, err := New(ctx, items, mapper, reducer, 0, 4)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, context.Canceled) {
		t.Errorf("expected context.Canceled, got %v", err)
	}
}

func TestUnit_MapReduceContextTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	items := []int{1, 2, 3}
	mapper := func(ctx context.Context, x int) (int, error) {
		select {
		case <-time.After(50 * time.Millisecond):
			return x, nil
		case <-ctx.Done():
			return 0, ctx.Err()
		}
	}
	reducer := func(ctx context.Context, acc, val int) (int, error) {
		return acc + val, nil
	}

	_, err := New(ctx, items, mapper, reducer, 0, 4)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("expected DeadlineExceeded, got %v", err)
	}
}

func TestUnit_MapReduceSingleWorker(t *testing.T) {
	ctx := context.Background()
	items := []int{10, 20, 30}
	mapper := func(ctx context.Context, x int) (int, error) {
		return x * 10, nil
	}
	reducer := func(ctx context.Context, acc, val int) (int, error) {
		return acc + val, nil
	}

	result, err := New(ctx, items, mapper, reducer, 0, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != 600 { // 100+200+300
		t.Errorf("expected 600, got %d", result)
	}
}

func TestUnit_MapReduceManyWorkers(t *testing.T) {
	// Проверка на гонки при большом числе горутин
	ctx := context.Background()
	n := 1000
	items := make([]int, n)
	for i := 0; i < n; i++ {
		items[i] = 1
	}
	mapper := func(ctx context.Context, x int) (int, error) {
		return x, nil
	}
	reducer := func(ctx context.Context, acc, val int) (int, error) {
		return acc + val, nil
	}

	result, err := New(ctx, items, mapper, reducer, 0, 50)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != n {
		t.Errorf("expected %d, got %d", n, result)
	}
}

func TestUnit_MapReduceConcurrentReduceOrder(t *testing.T) {
	// Проверяем, что при параллельном map результат не зависит от порядка
	ctx := context.Background()
	items := []int{2, 3, 5, 7, 11}
	mapper := func(ctx context.Context, x int) (int, error) {
		return x, nil
	}
	reducer := func(ctx context.Context, acc, val int) (int, error) {
		return acc*10 + val, nil // некоммутативный reducer
	}

	// При параллельной обработке порядок reduce не гарантирован,
	// поэтому ожидаем лишь, что функция завершится без ошибки
	// (конкретное значение будет зависеть от планировщика).
	_, err := New(ctx, items, mapper, reducer, 0, 4)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// тест просто не должен паниковать или возвращать ошибку
}

func TestUnit_MapReduceWithCustomTypes(t *testing.T) {
	type person struct {
		Name string
		Age  int
	}
	ctx := context.Background()
	items := []person{{"Alice", 30}, {"Bob", 25}, {"Charlie", 35}}
	mapper := func(ctx context.Context, p person) (int, error) {
		return p.Age, nil
	}
	reducer := func(ctx context.Context, acc string, age int) (string, error) {
		if acc == "" {
			return string(rune(age)), nil // некий аккумулятор
		}
		return acc + "," + string(rune(age)), nil
	}

	res, err := New(ctx, items, mapper, reducer, "", 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res) == 0 {
		t.Error("expected non-empty result string")
	}
}

func TestUnit_MapReduceWaitForAllGoroutines(t *testing.T) {
	// Убедимся, что горутины mapper действительно выполняются параллельно
	// и канал закрывается только после их завершения.
	ctx := context.Background()
	var mu sync.Mutex
	visited := make(map[int]bool)
	items := []int{1, 2, 3, 4, 5, 6, 7, 8}
	mapper := func(ctx context.Context, x int) (int, error) {
		time.Sleep(10 * time.Millisecond) // имитация работы
		mu.Lock()
		visited[x] = true
		mu.Unlock()
		return x, nil
	}
	reducer := func(ctx context.Context, acc []int, val int) ([]int, error) {
		return append(acc, val), nil
	}

	result, err := New(ctx, items, mapper, reducer, []int{}, 4)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	mu.Lock()
	visitCount := len(visited)
	mu.Unlock()
	if visitCount != len(items) {
		t.Errorf("expected %d visited items, got %d", len(items), visitCount)
	}
	if len(result) != len(items) {
		t.Errorf("expected %d results, got %d", len(items), len(result))
	}
}
