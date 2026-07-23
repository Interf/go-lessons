package tasks

import (
	"context"
	"sync/atomic"
	"testing"
	"time"
)

func TestScheduleAfter(t *testing.T) {
	s := NewScheduler()
	defer s.Stop()

	var counter int32
	task := NewTask(func(ctx context.Context) error {
		atomic.AddInt32(&counter, 1)
		return nil
	})

	s.ScheduleAfter(task, 50*time.Millisecond)

	time.Sleep(150 * time.Millisecond)

	if atomic.LoadInt32(&counter) != 1 {
		t.Errorf("expected task to run once, got %d", counter)
	}
}

func TestCancel(t *testing.T) {
	s := NewScheduler()
	defer s.Stop()

	var counter int32
	task := NewTask(func(ctx context.Context) error {
		atomic.AddInt32(&counter, 1)
		return nil
	})

	s.ScheduleAfter(task, 50*time.Millisecond)
	s.Cancel(task.ID)

	time.Sleep(150 * time.Millisecond)

	if atomic.LoadInt32(&counter) != 0 {
		t.Errorf("expected task to be cancelled, but it ran %d times", counter)
	}
}

func TestGetPendingTasks(t *testing.T) {
	s := NewScheduler()
	defer s.Stop()

	task1 := NewTask(func(ctx context.Context) error { return nil })
	task2 := NewTask(func(ctx context.Context) error { return nil })

	s.ScheduleAfter(task1, time.Hour)
	s.ScheduleAfter(task2, time.Hour)

	pending := s.GetPendingTasks()
	if len(pending) != 2 {
		t.Errorf("expected 2 pending tasks, got %d", len(pending))
	}
}

func TestRecurring(t *testing.T) {
	s := NewScheduler()
	defer s.Stop()

	var counter int32
	task := NewTask(func(ctx context.Context) error {
		atomic.AddInt32(&counter, 1)
		return nil
	})

	s.ScheduleEvery(task, 50*time.Millisecond)

	time.Sleep(250 * time.Millisecond)

	count := atomic.LoadInt32(&counter)
	if count < 3 {
		t.Errorf("expected recurring task to run at least 3 times, got %d", count)
	}

	pending := s.GetPendingTasks()
	found := false
	for _, p := range pending {
		if p.ID == task.ID {
			found = true
			break
		}
	}
	if !found {
		t.Error("recurring task should still be pending")
	}
}

func TestScheduleAt(t *testing.T) {
	s := NewScheduler()
	defer s.Stop()

	var counter int32
	task := NewTask(func(ctx context.Context) error {
		atomic.AddInt32(&counter, 1)
		return nil
	})

	future := time.Now().Add(50 * time.Millisecond)
	s.ScheduleAt(task, future)

	time.Sleep(150 * time.Millisecond)

	if atomic.LoadInt32(&counter) != 1 {
		t.Errorf("expected task to run once, got %d", counter)
	}
}

func TestStopWaitsForCurrentTask(t *testing.T) {
	s := NewScheduler()

	var started int32
	var finished int32
	startedCh := make(chan struct{})
	task := NewTask(func(ctx context.Context) error {
		atomic.StoreInt32(&started, 1)
		close(startedCh)
		time.Sleep(100 * time.Millisecond)
		atomic.StoreInt32(&finished, 1)
		return nil
	})

	s.ScheduleAfter(task, 0)

	<-startedCh

	s.Stop()

	if atomic.LoadInt32(&finished) != 1 {
		t.Error("Stop should wait for current task to finish")
	}
}
