package page

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
)

type MemoryStore struct {
	mu          sync.Mutex
	pages       map[uuid.UUID]*Page
	byUser      map[uuid.UUID]uuid.UUID
	bySlug      map[string]uuid.UUID
	links       map[uuid.UUID]*Link
	subscribers map[uuid.UUID]*Subscriber
	days        map[uuid.UUID]map[string]int
	referrers   map[uuid.UUID]map[string]int
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		pages:       map[uuid.UUID]*Page{},
		byUser:      map[uuid.UUID]uuid.UUID{},
		bySlug:      map[string]uuid.UUID{},
		links:       map[uuid.UUID]*Link{},
		subscribers: map[uuid.UUID]*Subscriber{},
		days:        map[uuid.UUID]map[string]int{},
		referrers:   map[uuid.UUID]map[string]int{},
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
	if p.PublishedAt == nil {
		p.PublishedAt = &now
	}
	m.pages[p.ID] = clonePage(p)
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
	if p.PublishedAt == nil {
		now := p.UpdatedAt
		p.PublishedAt = &now
	}
	m.pages[p.ID] = clonePage(p)
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
	return clonePage(p), nil
}

func (m *MemoryStore) GetPageByUserID(ctx context.Context, userID uuid.UUID) (*Page, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	id, ok := m.byUser[userID]
	if !ok {
		return nil, ErrNotFound
	}
	return clonePage(m.pages[id]), nil
}

func (m *MemoryStore) GetPageBySlug(ctx context.Context, slug string) (*Page, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	id, ok := m.bySlug[slug]
	if !ok {
		return nil, ErrNotFound
	}
	return clonePage(m.pages[id]), nil
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

func (m *MemoryStore) IncrementClicks(ctx context.Context, pageID, linkID uuid.UUID) (int, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	l, ok := m.links[linkID]
	if !ok || l.PageID != pageID || !l.Active {
		return 0, ErrNotFound
	}
	now := time.Now().UTC()
	l.Clicks++
	l.LastClickedAt = &now
	l.UpdatedAt = now
	day := now.Format("2006-01-02")
	if m.days[pageID] == nil {
		m.days[pageID] = map[string]int{}
	}
	m.days[pageID][day]++
	return l.Clicks, nil
}

func (m *MemoryStore) DailyClicks(ctx context.Context, pageID uuid.UUID) (map[string]int, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	src := m.days[pageID]
	out := map[string]int{}
	for k, v := range src {
		out[k] = v
	}
	return out, nil
}

func (m *MemoryStore) BumpReferrer(ctx context.Context, pageID uuid.UUID, host string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.referrers[pageID] == nil {
		m.referrers[pageID] = map[string]int{}
	}
	m.referrers[pageID][host]++
	return nil
}

func (m *MemoryStore) Referrers(ctx context.Context, pageID uuid.UUID) (map[string]int, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	src := m.referrers[pageID]
	out := map[string]int{}
	for k, v := range src {
		out[k] = v
	}
	return out, nil
}

func (m *MemoryStore) SumClicks(ctx context.Context) (int, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	n := 0
	for _, l := range m.links {
		n += l.Clicks
	}
	return n, nil
}

func (m *MemoryStore) ClicksByPage(ctx context.Context) (map[uuid.UUID]int, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	out := make(map[uuid.UUID]int, len(m.pages))
	for _, p := range m.pages {
		out[p.ID] = 0
	}
	for _, l := range m.links {
		out[l.PageID] += l.Clicks
	}
	return out, nil
}

func (m *MemoryStore) CreateSubscriber(ctx context.Context, s *Subscriber) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, existing := range m.subscribers {
		if existing.PageID == s.PageID && existing.Email == s.Email {
			return ErrConflict
		}
	}
	s.CreatedAt = time.Now().UTC()
	cp := *s
	m.subscribers[s.ID] = &cp
	return nil
}

func (m *MemoryStore) ListSubscribers(ctx context.Context, pageID uuid.UUID) ([]*Subscriber, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	var out []*Subscriber
	for _, s := range m.subscribers {
		if s.PageID == pageID {
			cp := *s
			out = append(out, &cp)
		}
	}
	for i := 0; i < len(out); i++ {
		for j := i + 1; j < len(out); j++ {
			if out[j].CreatedAt.After(out[i].CreatedAt) {
				out[i], out[j] = out[j], out[i]
			}
		}
	}
	if out == nil {
		out = []*Subscriber{}
	}
	return out, nil
}

func (m *MemoryStore) GetSubscriberByID(ctx context.Context, id uuid.UUID) (*Subscriber, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	s, ok := m.subscribers[id]
	if !ok {
		return nil, ErrNotFound
	}
	cp := *s
	return &cp, nil
}

func (m *MemoryStore) DeleteSubscriber(ctx context.Context, id uuid.UUID) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.subscribers[id]; !ok {
		return ErrNotFound
	}
	delete(m.subscribers, id)
	return nil
}
