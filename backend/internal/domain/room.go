package domain

type Room struct {
    id     string
    homeID string
    name   string
}

func NewRoom(id, homeID, name string) *Room {
    return &Room{
        id:     id,
        homeID: homeID,
        name:   name,
    }
}

func (r *Room) ID() string     { return r.id }
func (r *Room) HomeID() string { return r.homeID }
func (r *Room) Name() string   { return r.name }

func (r *Room) SetName(name string) { r.name = name }