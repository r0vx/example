package models

import (
	"time"

	"github.com/r0vx/x/login"
	"github.com/r0vx/x/perm"
	"gorm.io/gorm"
)

const (
	RoleAdmin   = "Admin"
	RoleManager = "Manager"
	RoleEditor  = "Editor"
	RoleViewer  = "Viewer"

	OAuthProviderGoogle          = "google"
	OAuthProviderMicrosoftOnline = "microsoftonline"
	OAuthProviderGithub          = "github"

	StatusActive   = "active"
	StatusInactive = "inactive"
)

var DefaultRoles = []string{
	RoleAdmin,
	RoleManager,
	RoleEditor,
	RoleViewer,
}

var OAuthProviders = []string{
	OAuthProviderGoogle,
	OAuthProviderMicrosoftOnline,
	OAuthProviderGithub,
}

type User struct {
	gorm.Model

	Name             string
	Company          string
	Avatar           string // 持久化头像 URL（OAuthAvatar 是 gorm:"-" 瞬态，不入库，故单列存）
	Roles            []perm.Role `gorm:"many2many:user_role_join;"`
	Status           string
	UpdatedAt        time.Time
	CreatedAt        time.Time
	FavorPostID      uint
	RegistrationDate time.Time `gorm:"type:date"`

	// Username is email
	login.UserPass
	login.OAuthInfo
	login.SessionSecure
}

func (u User) GetName() string {
	return u.Name
}

// SetName 供 signup 流程写入用户显示名（实现 login.NameSetter 接口）
func (u *User) SetName(v string) { u.Name = v }

func (u User) GetID() uint {
	return u.ID
}

func (u User) GetRoles() (rs []string) {
	for _, r := range u.Roles {
		rs = append(rs, r.Name)
	}
	if len(rs) == 0 {
		rs = []string{RoleViewer}
	}
	return
}

func (u User) IsOAuthUser() bool {
	return u.OAuthProvider != "" && u.OAuthIdentifier != ""
}
