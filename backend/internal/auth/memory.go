package auth

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
)

type MemoryStore struct {
	mu       sync.Mutex
	users    map[string]*User
	byID     map[uuid.UUID]*User
	reset    map[string]resetEntry
	settings Settings
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		users:    map[string]*User{},
		byID:     map[uuid.UUID]*User{},
		reset:    map[string]resetEntry{},
		settings: DefaultSettings(),
	}
}

func (m *MemoryStore) UserCount(ctx context.Context) (int, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.users), nil
}

func (m *MemoryStore) CreateUser(ctx context.Context, u *User) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.users[u.Email]; ok {
		return ErrExists
	}
	cp := *u
	m.users[u.Email] = &cp
	m.byID[u.ID] = &cp
	return nil
}

func (m *MemoryStore) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	u, ok := m.users[email]
	if !ok {
		return nil, ErrNotFound
	}
	cp := *u
	return &cp, nil
}

func (m *MemoryStore) GetUserByID(ctx context.Context, id uuid.UUID) (*User, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	u, ok := m.byID[id]
	if !ok {
		return nil, ErrNotFound
	}
	cp := *u
	return &cp, nil
}

func (m *MemoryStore) ListUsers(ctx context.Context) ([]*User, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	out := make([]*User, 0, len(m.byID))
	for _, u := range m.byID {
		cp := *u
		out = append(out, &cp)
	}
	return out, nil
}

func (m *MemoryStore) UpdateUser(ctx context.Context, u *User) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	old, ok := m.byID[u.ID]
	if !ok {
		return ErrNotFound
	}
	if old.Email != u.Email {
		delete(m.users, old.Email)
	}
	cp := *u
	m.byID[u.ID] = &cp
	m.users[u.Email] = &cp
	return nil
}

func (m *MemoryStore) DeleteUser(ctx context.Context, id uuid.UUID) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	u, ok := m.byID[id]
	if !ok {
		return ErrNotFound
	}
	delete(m.byID, id)
	delete(m.users, u.Email)
	return nil
}

func (m *MemoryStore) GetSettings(ctx context.Context) (Settings, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	st := m.settings
	st.Nav = navOrDefault(m.settings.Nav)
	if m.settings.Nav != nil {
		st.Nav = append([]NavItem(nil), m.settings.Nav...)
	}
	return st, nil
}

func (m *MemoryStore) PutSettings(ctx context.Context, st Settings) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.settings = st
	if st.Nav != nil {
		m.settings.Nav = append([]NavItem(nil), st.Nav...)
	}
	return nil
}

func (m *MemoryStore) PutResetToken(ctx context.Context, token string, userID uuid.UUID, exp time.Time) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.reset == nil {
		m.reset = map[string]resetEntry{}
	}
	m.reset[token] = resetEntry{userID: userID, exp: exp}
	return nil
}

func (m *MemoryStore) GetResetToken(ctx context.Context, token string) (uuid.UUID, time.Time, bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.reset == nil {
		return uuid.Nil, time.Time{}, false, nil
	}
	e, ok := m.reset[token]
	if !ok {
		return uuid.Nil, time.Time{}, false, nil
	}
	return e.userID, e.exp, true, nil
}

func (m *MemoryStore) DeleteResetToken(ctx context.Context, token string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.reset, token)
	return nil
}
