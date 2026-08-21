package page

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
)

type MemoryStore struct {
	mu     sync.Mutex
	pages  map[uuid.UUID]*Page
	byUser map[uuid.UUID]uuid.UUID
	bySlug map[string]uuid.UUID
	links  map[uuid.UUID]*Link
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		pages:  map[uuid.UUID]*Page{},
		byUser: map[uuid.UUID]uuid.UUID{},
		bySlug: map[string]uuid.UUID{},
		links:  map[uuid.UUID]*Link{},
	}
}

func (m *MemoryStore) CreatePage(ctx context.Context, p *Page) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.bySlug[p.Slug]; ok {
		return ErrConflict
	}
	now := time.Now().UTC()
	p.CreatedAt, p.UpdatedAt = now, now
	cp := *p
	m.pages[p.ID] = &cp
	m.byUser[p.UserID] = p.ID
	m.bySlug[p.Slug] = p.ID
	return nil
}

func (m *MemoryStore) UpdatePage(ctx context.Context, p *Page) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if other, ok := m.bySlug[p.Slug]; ok && other != p.ID {
		return ErrConflict
	}
	old := m.pages[p.ID]
	if old != nil {
		delete(m.bySlug, old.Slug)
	}
	p.UpdatedAt = time.Now().UTC()
	cp := *p
	m.pages[p.ID] = &cp
	m.bySlug[p.Slug] = p.ID
	m.byUser[p.UserID] = p.ID
	return nil
}

func (m *MemoryStore) GetPageByID(ctx context.Context, id uuid.UUID) (*Page, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	p, ok := m.pages[id]
	if !ok {
		return nil, ErrNotFound
	}
	cp := *p
	return &cp, nil
}

func (m *MemoryStore) GetPageByUserID(ctx context.Context, userID uuid.UUID) (*Page, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	id, ok := m.byUser[userID]
	if !ok {
		return nil, ErrNotFound
	}
	cp := *m.pages[id]
	return &cp, nil
}

func (m *MemoryStore) GetPageBySlug(ctx context.Context, slug string) (*Page, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	id, ok := m.bySlug[slug]
	if !ok {
		return nil, ErrNotFound
	}
	cp := *m.pages[id]
	return &cp, nil
}

func (m *MemoryStore) CreateLink(ctx context.Context, l *Link) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	now := time.Now().UTC()
	l.CreatedAt, l.UpdatedAt = now, now
	cp := *l
	m.links[l.ID] = &cp
	return nil
}

func (m *MemoryStore) UpdateLink(ctx context.Context, l *Link) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	l.UpdatedAt = time.Now().UTC()
	cp := *l
	m.links[l.ID] = &cp
	return nil
}

func (m *MemoryStore) DeleteLink(ctx context.Context, id uuid.UUID) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.links, id)
	return nil
}

func (m *MemoryStore) GetLinkByID(ctx context.Context, id uuid.UUID) (*Link, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	l, ok := m.links[id]
	if !ok {
		return nil, ErrNotFound
	}
	cp := *l
	return &cp, nil
}

func (m *MemoryStore) ListLinks(ctx context.Context, pageID uuid.UUID) ([]*Link, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	var out []*Link
	for _, l := range m.links {
		if l.PageID == pageID {
			cp := *l
			out = append(out, &cp)
		}
	}
	for i := 0; i < len(out); i++ {
		for j := i + 1; j < len(out); j++ {
			if out[j].Order < out[i].Order {
				out[i], out[j] = out[j], out[i]
			}
		}
	}
	return out, nil
}

func (m *MemoryStore) MaxOrder(ctx context.Context, pageID uuid.UUID) (int, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	max := -1
	for _, l := range m.links {
		if l.PageID == pageID && l.Order > max {
			max = l.Order
		}
	}
	return max, nil
}

func (m *MemoryStore) Reorder(ctx context.Context, pageID uuid.UUID, ids []uuid.UUID) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	for i, id := range ids {
		l, ok := m.links[id]
		if !ok || l.PageID != pageID {
			return ErrNotFound
		}
		l.Order = i
	}
	return nil
}
