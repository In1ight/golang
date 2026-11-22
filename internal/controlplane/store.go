package controlplane

import (
	"github.com/google/uuid"
	"sort"
	"sync"
	"time"
)

type ProjectStore struct {
	mu       sync.RWMutex
	projects map[string]Project
}

func NewProjectStore() *ProjectStore {
	return &ProjectStore{
		projects: make(map[string]Project),
	}
}

func (s *ProjectStore) Create(p Project) Project {
	p.ID = uuid.NewString()

	now := time.Now().UTC()
	p.CreatedAt = now
	p.UpdatedAt = now

	if p.Timezone == "" {
		p.Timezone = "UTC"
	}
	s.mu.Lock()
	defer s.mu.Unlock()

	s.projects[p.ID] = p
	return p
}

func (s *ProjectStore) Get(id string) (Project, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	p, ok := s.projects[id]
	return p, ok
}

func (s *ProjectStore) List() []Project {
	s.mu.RLock()
	defer s.mu.RUnlock()

	items := make([]Project, 0, len(s.projects))
	for _, p := range s.projects {
		items = append(items, p)
	}
	sort.Slice(items, func(i, j int) bool {
		if items[i].CreatedAt.Equal(items[j].CreatedAt) {
			return items[i].ID < items[j].ID
		}
		return items[i].CreatedAt.Before(items[j].CreatedAt)
	})
	return items
}

func (s *ProjectStore) Update(id string, input UpdateProjectInput) (Project, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	p, ok := s.projects[id]
	if !ok {
		return Project{}, false
	}

	if input.Name != nil {
		p.Name = *input.Name
	}

	if input.Timezone != nil {
		p.Timezone = *input.Timezone
	}

	p.UpdatedAt = time.Now().UTC()

	s.projects[id] = p

	return p, true
}

func (s *ProjectStore) Delete(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.projects[id]
	if !ok {
		return false
	}

	delete(s.projects, id)
	return true
}
