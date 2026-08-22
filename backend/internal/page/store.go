package page

import (
	"context"

	"github.com/google/uuid"
)

type PageStore interface {
	CreatePage(ctx context.Context, p *Page) error
	UpdatePage(ctx context.Context, p *Page) error
	GetPageByID(ctx context.Context, id uuid.UUID) (*Page, error)
	GetPageByUserID(ctx context.Context, userID uuid.UUID) (*Page, error)
	GetPageBySlug(ctx context.Context, slug string) (*Page, error)
	ListPages(ctx context.Context) ([]*Page, error)
}

type LinkStore interface {
	CreateLink(ctx context.Context, l *Link) error
	UpdateLink(ctx context.Context, l *Link) error
	DeleteLink(ctx context.Context, id uuid.UUID) error
	GetLinkByID(ctx context.Context, id uuid.UUID) (*Link, error)
	ListLinks(ctx context.Context, pageID uuid.UUID) ([]*Link, error)
	MaxOrder(ctx context.Context, pageID uuid.UUID) (int, error)
	Reorder(ctx context.Context, pageID uuid.UUID, ids []uuid.UUID) error
	IncrementClicks(ctx context.Context, pageID, linkID uuid.UUID) (int, error)
	SumClicks(ctx context.Context) (int, error)
}

type Store interface {
	PageStore
	LinkStore
}

func (m *MemoryStore) ListPages(ctx context.Context) ([]*Page, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	out := make([]*Page, 0, len(m.pages))
	for _, p := range m.pages {
		cp := *p
		out = append(out, &cp)
	}
	return out, nil
}
