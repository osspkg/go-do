/*
 *  Copyright (c) 2024-2025 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package do_test

import (
	"context"
	"errors"
	"io"
	"strings"
	"sync"
	"testing"

	"go.osspkg.com/do"
)

// Test types
type TestState string
type TestData struct {
	Value    int
	Messages []string
	Error    error
}

const (
	StateInit    TestState = "init"
	StateWorking TestState = "working"
	StateDone    TestState = "done"
)

func TestNewStateMachine(t *testing.T) {
	t.Run("should create new state machine", func(t *testing.T) {
		sm := do.NewStateMachine[TestState, TestData]()

		if sm == nil {
			t.Error("Expected state machine to be created, got nil")
		}
	})
}

func TestStateMachine_Add(t *testing.T) {
	t.Run("should add valid transition", func(t *testing.T) {
		sm := do.NewStateMachine[TestState, TestData]()
		nextState := StateWorking

		transition := &do.Transition[TestState, TestData]{
			Previous: StateInit,
			Next:     &nextState,
			Apply: []func(ctx context.Context, data TestData) (TestData, error){
				func(ctx context.Context, data TestData) (TestData, error) {
					data.Value += 1
					return data, nil
				},
			},
		}

		err := sm.Add(transition)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	t.Run("should reject nil transition", func(t *testing.T) {
		sm := do.NewStateMachine[TestState, TestData]()

		err := sm.Add(nil)
		if err == nil {
			t.Error("Expected error for nil transition, got nil")
		}

		expectedError := "transition cannot be nil"
		if err.Error() != expectedError {
			t.Errorf("Expected error '%s', got '%v'", expectedError, err)
		}
	})

	t.Run("should reject transition with no apply functions", func(t *testing.T) {
		sm := do.NewStateMachine[TestState, TestData]()
		nextState := StateWorking

		transition := &do.Transition[TestState, TestData]{
			Previous: StateInit,
			Next:     &nextState,
			Apply:    []func(ctx context.Context, data TestData) (TestData, error){},
		}

		err := sm.Add(transition)
		if err == nil {
			t.Error("Expected error for transition with no apply functions, got nil")
		}

		expectedError := "transition must have at least one apply"
		if err.Error() != expectedError {
			t.Errorf("Expected error '%s', got '%v'", expectedError, err)
		}
	})

	t.Run("should reject duplicate previous state", func(t *testing.T) {
		sm := do.NewStateMachine[TestState, TestData]()
		nextState := StateWorking

		transition1 := &do.Transition[TestState, TestData]{
			Previous: StateInit,
			Next:     &nextState,
			Apply: []func(ctx context.Context, data TestData) (TestData, error){
				func(ctx context.Context, data TestData) (TestData, error) {
					return data, nil
				},
			},
		}

		transition2 := &do.Transition[TestState, TestData]{
			Previous: StateInit, // Same previous state
			Next:     &nextState,
			Apply: []func(ctx context.Context, data TestData) (TestData, error){
				func(ctx context.Context, data TestData) (TestData, error) {
					return data, nil
				},
			},
		}

		err1 := sm.Add(transition1)
		if err1 != nil {
			t.Errorf("First add should succeed, got %v", err1)
		}

		err2 := sm.Add(transition2)
		if err2 == nil {
			t.Error("Expected error for duplicate previous state, got nil")
		}

		expectedError := "transition already has a previous state"
		if err2.Error() != expectedError {
			t.Errorf("Expected error '%s', got '%v'", expectedError, err2)
		}
	})

	t.Run("should reject identical previous and next state", func(t *testing.T) {
		sm := do.NewStateMachine[TestState, TestData]()
		nextState := StateInit // Same as Previous

		transition := &do.Transition[TestState, TestData]{
			Previous: StateInit,
			Next:     &nextState,
			Apply: []func(ctx context.Context, data TestData) (TestData, error){
				func(ctx context.Context, data TestData) (TestData, error) {
					return data, nil
				},
			},
		}

		err := sm.Add(transition)
		if err == nil {
			t.Error("Expected error for identical previous and next state, got nil")
		}

		expectedError := "transition has a identical previous and next state"
		if err.Error() != expectedError {
			t.Errorf("Expected error '%s', got '%v'", expectedError, err)
		}
	})

	t.Run("should reject direct state loop", func(t *testing.T) {
		sm := do.NewStateMachine[TestState, TestData]()

		// First transition: A -> B
		stateB := StateWorking
		transitionAB := &do.Transition[TestState, TestData]{
			Previous: StateInit,
			Next:     &stateB,
			Apply: []func(ctx context.Context, data TestData) (TestData, error){
				func(ctx context.Context, data TestData) (TestData, error) {
					return data, nil
				},
			},
		}

		// Second transition: B -> A (creates loop)
		stateA := StateInit
		transitionBA := &do.Transition[TestState, TestData]{
			Previous: StateWorking,
			Next:     &stateA,
			Apply: []func(ctx context.Context, data TestData) (TestData, error){
				func(ctx context.Context, data TestData) (TestData, error) {
					return data, nil
				},
			},
		}

		err1 := sm.Add(transitionAB)
		if err1 != nil {
			t.Fatalf("First transition should succeed, got %v", err1)
		}

		err2 := sm.Add(transitionBA)
		if err2 == nil {
			t.Error("Expected error for direct state loop, got nil")
		}

		expectedError := "transaction has a direct state loop"
		if err2.Error() != expectedError {
			t.Errorf("Expected error '%s', got '%v'", expectedError, err2)
		}
	})

	t.Run("should allow transition with nil next state", func(t *testing.T) {
		sm := do.NewStateMachine[TestState, TestData]()

		transition := &do.Transition[TestState, TestData]{
			Previous: StateInit,
			Next:     nil, // Terminal state
			Apply: []func(ctx context.Context, data TestData) (TestData, error){
				func(ctx context.Context, data TestData) (TestData, error) {
					return data, nil
				},
			},
		}

		err := sm.Add(transition)
		if err != nil {
			t.Errorf("Expected no error for nil next state, got %v", err)
		}
	})
}

func TestStateMachine_Apply(t *testing.T) {
	t.Run("should execute single transition", func(t *testing.T) {
		sm := do.NewStateMachine[TestState, TestData]()
		nextState := StateWorking

		called := false
		transition := &do.Transition[TestState, TestData]{
			Previous: StateInit,
			Next:     &nextState,
			Apply: []func(ctx context.Context, data TestData) (TestData, error){
				func(ctx context.Context, data TestData) (TestData, error) {
					called = true
					data.Value += 10
					return data, nil
				},
			},
		}

		err := sm.Add(transition)
		if err != nil {
			t.Fatalf("Failed to add transition: %v", err)
		}

		initialData := TestData{Value: 5}
		err = sm.Apply(context.Background(), StateInit, initialData)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if !called {
			t.Error("Apply function was not called")
		}
	})

	t.Run("should execute multiple apply functions", func(t *testing.T) {
		sm := do.NewStateMachine[TestState, TestData]()
		nextState := StateDone

		callCount := 0
		transition := &do.Transition[TestState, TestData]{
			Previous: StateInit,
			Next:     &nextState,
			Apply: []func(ctx context.Context, data TestData) (TestData, error){
				func(ctx context.Context, data TestData) (TestData, error) {
					callCount++
					data.Value += 1
					return data, nil
				},
				func(ctx context.Context, data TestData) (TestData, error) {
					callCount++
					data.Value *= 2
					return data, nil
				},
				func(ctx context.Context, data TestData) (TestData, error) {
					callCount++
					data.Value += 3
					return data, nil
				},
			},
		}

		err := sm.Add(transition)
		if err != nil {
			t.Fatalf("Failed to add transition: %v", err)
		}

		initialData := TestData{Value: 5}
		err = sm.Apply(context.Background(), StateInit, initialData)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if callCount != 3 {
			t.Errorf("Expected 3 apply calls, got %d", callCount)
		}
	})

	t.Run("should chain multiple transitions", func(t *testing.T) {
		sm := do.NewStateMachine[TestState, TestData]()

		// Transition from Init -> Working
		workingState := StateWorking
		initTransition := &do.Transition[TestState, TestData]{
			Previous: StateInit,
			Next:     &workingState,
			Apply: []func(ctx context.Context, data TestData) (TestData, error){
				func(ctx context.Context, data TestData) (TestData, error) {
					data.Messages = append(data.Messages, "init->working")
					return data, nil
				},
			},
		}

		// Transition from Working -> Done
		doneState := StateDone
		workingTransition := &do.Transition[TestState, TestData]{
			Previous: StateWorking,
			Next:     &doneState,
			Apply: []func(ctx context.Context, data TestData) (TestData, error){
				func(ctx context.Context, data TestData) (TestData, error) {
					data.Messages = append(data.Messages, "working->done")
					return data, nil
				},
			},
		}

		err := sm.Add(initTransition)
		if err != nil {
			t.Fatalf("Failed to add init transition: %v", err)
		}

		err = sm.Add(workingTransition)
		if err != nil {
			t.Fatalf("Failed to add working transition: %v", err)
		}

		initialData := TestData{Messages: []string{}}
		err = sm.Apply(context.Background(), StateInit, initialData)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	t.Run("should stop at transition with nil next state", func(t *testing.T) {
		sm := do.NewStateMachine[TestState, TestData]()

		callCount := 0
		transition := &do.Transition[TestState, TestData]{
			Previous: StateInit,
			Next:     nil, // No next state - should stop here
			Apply: []func(ctx context.Context, data TestData) (TestData, error){
				func(ctx context.Context, data TestData) (TestData, error) {
					callCount++
					return data, nil
				},
			},
		}

		err := sm.Add(transition)
		if err != nil {
			t.Fatalf("Failed to add transition: %v", err)
		}

		initialData := TestData{}
		err = sm.Apply(context.Background(), StateInit, initialData)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if callCount != 1 {
			t.Errorf("Expected 1 apply call, got %d", callCount)
		}
	})

	t.Run("should return error when apply function fails", func(t *testing.T) {
		sm := do.NewStateMachine[TestState, TestData]()
		nextState := StateWorking

		expectedError := errors.New("apply function failed")
		transition := &do.Transition[TestState, TestData]{
			Previous: StateInit,
			Next:     &nextState,
			Apply: []func(ctx context.Context, data TestData) (TestData, error){
				func(ctx context.Context, data TestData) (TestData, error) {
					return data, expectedError
				},
			},
		}

		err := sm.Add(transition)
		if err != nil {
			t.Fatalf("Failed to add transition: %v", err)
		}

		initialData := TestData{}
		err = sm.Apply(context.Background(), StateInit, initialData)

		if err == nil {
			t.Error("Expected error from apply function, got nil")
		}
		if !errors.Is(err, expectedError) {
			t.Errorf("Expected error '%v', got '%v'", expectedError, err)
		}
	})

	t.Run("should handle unknown state gracefully", func(t *testing.T) {
		sm := do.NewStateMachine[TestState, TestData]()

		// No transitions added, so any state should be unknown
		initialData := TestData{}
		err := sm.Apply(context.Background(), StateInit, initialData)

		if err != nil {
			t.Errorf("Expected no error for unknown state, got %v", err)
		}
	})

	t.Run("should handle panic in apply function", func(t *testing.T) {
		sm := do.NewStateMachine[TestState, TestData]()
		nextState := StateWorking

		transition := &do.Transition[TestState, TestData]{
			Previous: StateInit,
			Next:     &nextState,
			Apply: []func(ctx context.Context, data TestData) (TestData, error){
				func(ctx context.Context, data TestData) (TestData, error) {
					panic("test panic")
				},
			},
		}

		err := sm.Add(transition)
		if err != nil {
			t.Fatalf("Failed to add transition: %v", err)
		}

		initialData := TestData{}
		err = sm.Apply(context.Background(), StateInit, initialData)

		if err == nil {
			t.Error("Expected error from panic, got nil")
		}

		// Should recover from panic and return error
		if !strings.Contains(err.Error(), "test panic") {
			t.Errorf("Expected panic error, got '%v'", err)
		}
	})

	t.Run("should return nil on EOF error", func(t *testing.T) {
		sm := do.NewStateMachine[TestState, TestData]()
		nextState := StateWorking

		transition := &do.Transition[TestState, TestData]{
			Previous: StateInit,
			Next:     &nextState,
			Apply: []func(ctx context.Context, data TestData) (TestData, error){
				func(ctx context.Context, data TestData) (TestData, error) {
					return data, io.EOF
				},
			},
		}

		err := sm.Add(transition)
		if err != nil {
			t.Fatalf("Failed to add transition: %v", err)
		}

		initialData := TestData{}
		err = sm.Apply(context.Background(), StateInit, initialData)

		if err != nil {
			t.Errorf("Expected nil error for EOF, got %v", err)
		}
	})

	t.Run("should pass context to apply functions", func(t *testing.T) {
		sm := do.NewStateMachine[TestState, TestData]()
		nextState := StateDone

		var receivedContext context.Context
		transition := &do.Transition[TestState, TestData]{
			Previous: StateInit,
			Next:     &nextState,
			Apply: []func(ctx context.Context, data TestData) (TestData, error){
				func(ctx context.Context, data TestData) (TestData, error) {
					receivedContext = ctx
					return data, nil
				},
			},
		}

		err := sm.Add(transition)
		if err != nil {
			t.Fatalf("Failed to add transition: %v", err)
		}

		ctx := context.WithValue(context.Background(), "testKey", "testValue")
		initialData := TestData{}
		err = sm.Apply(ctx, StateInit, initialData)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if receivedContext == nil {
			t.Error("Context was not passed to apply function")
		} else if receivedContext.Value("testKey") != "testValue" {
			t.Error("Context values were not preserved")
		}
	})
}

func TestStateMachine_Concurrency(t *testing.T) {
	t.Run("should handle concurrent access", func(t *testing.T) {
		sm := do.NewStateMachine[TestState, TestData]()

		var wg sync.WaitGroup
		errC := make(chan error, 10)

		// Concurrent adds
		for i := 0; i < 5; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()

				state := TestState(string(StateInit) + string(rune('A'+i)))
				nextState := StateWorking

				transition := &do.Transition[TestState, TestData]{
					Previous: state,
					Next:     &nextState,
					Apply: []func(ctx context.Context, data TestData) (TestData, error){
						func(ctx context.Context, data TestData) (TestData, error) {
							return data, nil
						},
					},
				}

				if err := sm.Add(transition); err != nil {
					errC <- err
				}
			}(i)
		}

		// Concurrent applies
		for i := 0; i < 5; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()

				state := TestState(string(StateInit) + string(rune('A'+i)))
				initialData := TestData{Value: i}

				if err := sm.Apply(context.Background(), state, initialData); err != nil {
					errC <- err
				}
			}(i)
		}

		wg.Wait()
		close(errC)

		for err := range errC {
			t.Errorf("Unexpected error in concurrent test: %v", err)
		}
	})
}

func TestStateMachine_EdgeCases(t *testing.T) {
	t.Run("should handle transition to unknown state", func(t *testing.T) {
		sm := do.NewStateMachine[TestState, TestData]()

		// Transition from A -> B, but B has no transition defined
		unknownState := StateWorking
		transition := &do.Transition[TestState, TestData]{
			Previous: StateInit,
			Next:     &unknownState,
			Apply: []func(ctx context.Context, data TestData) (TestData, error){
				func(ctx context.Context, data TestData) (TestData, error) {
					data.Messages = append(data.Messages, "moved to unknown state")
					return data, nil
				},
			},
		}

		err := sm.Add(transition)
		if err != nil {
			t.Fatalf("Failed to add transition: %v", err)
		}

		initialData := TestData{Messages: []string{}}
		err = sm.Apply(context.Background(), StateInit, initialData)

		if err != nil {
			t.Errorf("Expected no error when transitioning to unknown state, got %v", err)
		}
	})
}
