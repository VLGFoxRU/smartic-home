package domain

import "time"

type Home struct {
    id        string
    name      string
    ownerID   string
    createdAt time.Time
}

func NewHome(id, name, ownerID string) *Home {
    return &Home{
        id:      id,
        name:    name,
        ownerID: ownerID,
        createdAt: time.Now(),
    }
}

func (h *Home) ID() string      { return h.id }
func (h *Home) Name() string    { return h.name }
func (h *Home) OwnerID() string { return h.ownerID }
func (h *Home) CreatedAt() time.Time { return h.createdAt }

func (h *Home) SetName(name string) { h.name = name }