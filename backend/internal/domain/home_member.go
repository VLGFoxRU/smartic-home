package domain

type HomeRole string

const (
    HomeRoleOwner      HomeRole = "owner"
    HomeRoleAdmin      HomeRole = "admin"
    HomeRoleController HomeRole = "controller"
    HomeRoleViewer     HomeRole = "viewer"
)

type HomeMember struct {
    homeID string
    userID string
    role   HomeRole
}

func NewHomeMember(homeID, userID string, role HomeRole) *HomeMember {
    return &HomeMember{
        homeID: homeID,
        userID: userID,
        role:   role,
    }
}

func (m *HomeMember) HomeID() string  { return m.homeID }
func (m *HomeMember) UserID() string  { return m.userID }
func (m *HomeMember) Role() HomeRole  { return m.role }

func (m *HomeMember) SetRole(role HomeRole) { m.role = role }