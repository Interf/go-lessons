package tasks

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Scheduler interface {
	ScheduleAt(Task, time.Time)
	ScheduleAfter(Task, time.Duration)
	ScheduleEvery(Task, time.Duration)
	Cancel(uuid.UUID)
	GetPendingTasks() []Task
	Stop()
}

type scheduler struct {
	mu       sync.Mutex
	tasks    map[uuid.UUID]Task
	ctx      context.Context
	cancel   context.CancelFunc
	stopChan chan struct{}
	doneChan chan struct{}
	wg       sync.WaitGroup
}

func NewScheduler() Scheduler {
	ctx, cancel := context.WithCancel(context.Background())
	s := &scheduler{
		tasks:    make(map[uuid.UUID]Task),
		ctx:      ctx,
		cancel:   cancel,
		stopChan: make(chan struct{}),
		doneChan: make(chan struct{}),
	}
	go s.loop()
	return s
}

func (s *scheduler) ScheduleAt(task Task, t time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()

	task.ExecuteAt = t
	s.tasks[task.ID] = task
}

func (s *scheduler) ScheduleAfter(task Task, d time.Duration) {
	s.ScheduleAt(task, time.Now().Add(d))
}

func (s *scheduler) ScheduleEvery(task Task, d time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()

	task.Recurring = true
	task.Interval = d
	task.ExecuteAt = time.Now().Add(d)
	s.tasks[task.ID] = task
}

func (s *scheduler) Cancel(id uuid.UUID) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.tasks, id)
}

func (s *scheduler) GetPendingTasks() []Task {
	s.mu.Lock()
	defer s.mu.Unlock()

	result := make([]Task, 0, len(s.tasks))
	for _, t := range s.tasks {
		result = append(result, t)
	}
	return result
}

func (s *scheduler) Stop() {
	close(s.stopChan)
	s.wg.Wait()
	<-s.doneChan
}

func (s *scheduler) loop() {
	defer close(s.doneChan)

	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.executeReadyTasks()
		case <-s.stopChan:
			s.executeReadyTasks()
			return
		case <-s.ctx.Done():
			return
		}
	}
}

func (s *scheduler) executeReadyTasks() {
	s.mu.Lock()
	now := time.Now()
	var ready []Task
	for id, t := range s.tasks {
		if !t.ExecuteAt.After(now) {
			ready = append(ready, t)
			if !t.Recurring {
				delete(s.tasks, id)
			}
		}
	}
	s.mu.Unlock()

	for _, task := range ready {
		s.wg.Add(1)
		go func() {
			defer s.wg.Done()
			if task.Command != nil {
				task.Command(s.ctx)
			}

			if task.Recurring && task.Interval > 0 {
				s.mu.Lock()
				updated := s.tasks[task.ID]
				updated.ExecuteAt = time.Now().Add(task.Interval)
				s.tasks[task.ID] = updated
				s.mu.Unlock()
			}
		}()
	}
}
