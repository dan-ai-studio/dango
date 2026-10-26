package server

import (
	"sync"
	"time"
)

// Task 任务模型
type Task struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Done      bool      `json:"done"`
	CreatedAt time.Time `json:"created_at"`
}

// Store 内存任务存储（demo 用，无需数据库）
type Store struct {
	mu     sync.RWMutex
	tasks  map[int64]*Task
	nextID int64
}

func NewStore() *Store {
	s := &Store{tasks: make(map[int64]*Task), nextID: 1}
	// 种子数据
	s.Create("熟悉 Gin + Vue 项目结构")
	s.Create("试试新增一条任务")
	t := s.Create("勾选完成一条任务")
	s.Toggle(t.ID)
	return s
}

func (s *Store) List(q string, done *bool) []*Task {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]*Task, 0, len(s.tasks))
	for _, t := range s.tasks {
		if q != "" && !containsFold(t.Title, q) {
			continue
		}
		if done != nil && t.Done != *done {
			continue
		}
		out = append(out, &Task{ID: t.ID, Title: t.Title, Done: t.Done, CreatedAt: t.CreatedAt})
	}
	// 按 ID 排序（简单冒泡即可，数据量小）
	for i := 0; i < len(out); i++ {
		for j := i + 1; j < len(out); j++ {
			if out[j].ID < out[i].ID {
				out[i], out[j] = out[j], out[i]
			}
		}
	}
	return out
}

func (s *Store) Create(title string) *Task {
	s.mu.Lock()
	defer s.mu.Unlock()
	t := &Task{ID: s.nextID, Title: title, CreatedAt: time.Now()}
	s.nextID++
	s.tasks[t.ID] = t
	cp := *t
	return &cp
}

func (s *Store) Get(id int64) (*Task, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	t, ok := s.tasks[id]
	if !ok {
		return nil, false
	}
	cp := *t
	return &cp, true
}

func (s *Store) Update(id int64, title *string, done *bool) (*Task, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	t, ok := s.tasks[id]
	if !ok {
		return nil, false
	}
	if title != nil {
		t.Title = *title
	}
	if done != nil {
		t.Done = *done
	}
	cp := *t
	return &cp, true
}

func (s *Store) Toggle(id int64) (*Task, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	t, ok := s.tasks[id]
	if !ok {
		return nil, false
	}
	t.Done = !t.Done
	cp := *t
	return &cp, true
}

func (s *Store) Delete(id int64) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.tasks[id]; !ok {
		return false
	}
	delete(s.tasks, id)
	return true
}

func (s *Store) Stats() (total, done, pending int) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	total = len(s.tasks)
	for _, t := range s.tasks {
		if t.Done {
			done++
		}
	}
	pending = total - done
	return
}

func containsFold(s, sub string) bool {
	return len(sub) == 0 || searchFold(s, sub)
}

func searchFold(s, sub string) bool {
	// 简单不区分大小写的包含（ASCII 够用，中文直接比较）
	ls, lsub := []rune(s), []rune(sub)
	if len(lsub) > len(ls) {
		return false
	}
	for i := 0; i <= len(ls)-len(lsub); i++ {
		match := true
		for j := 0; j < len(lsub); j++ {
			a, b := ls[i+j], lsub[j]
			if a >= 'A' && a <= 'Z' {
				a += 'a' - 'A'
			}
			if b >= 'A' && b <= 'Z' {
				b += 'a' - 'A'
			}
			if a != b {
				match = false
				break
			}
		}
		if match {
			return true
		}
	}
	return false
}
