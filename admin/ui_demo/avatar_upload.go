package ui_demo

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/r0vx/admin/avatarupload"
	"github.com/r0vx/admin/presets"
	"gorm.io/gorm"
)

// AvatarUploadPath 头像上传端点
const AvatarUploadPath = "/avatar-demo/upload"

// avatarUploadDir 本地存储目录（demo 用；生产换 OSS）
const avatarUploadDir = "./public/uploads/avatars"

// avatarPublicPrefix 对外访问前缀（router 挂 /uploads/ → ./public/uploads）
const avatarPublicPrefix = "/uploads/avatars"

// MemberProfile 头像上传演示模型
type MemberProfile struct {
	ID        uint `gorm:"primarykey"`
	Name      string
	Avatar    string // 存头像 URL
	CreatedAt time.Time
	UpdatedAt time.Time
}

// ConfigAvatarUploadDemo 注册头像上传演示
func ConfigAvatarUploadDemo(b *presets.Builder, db *gorm.DB) {
	if err := db.AutoMigrate(&MemberProfile{}); err != nil {
		panic(err)
	}
	seedMemberProfiles(db)

	mb := b.Model(&MemberProfile{}).URIName("avatar-upload-demo")
	mb.Listing("ID", "Name", "Avatar")
	mb.Editing("Name", "Avatar")
	avatarupload.Configure(mb, "Avatar", avatarupload.Config{
		UploadURL: AvatarUploadPath,
		Shape:     "circle",
		Size:      96,
	})
}

// AvatarUploadHandler 本地磁盘存储 handler（demo）
func AvatarUploadHandler() http.Handler {
	return avatarupload.New(func(r *http.Request, f multipart.File, hdr *multipart.FileHeader) (string, error) {
		if err := os.MkdirAll(avatarUploadDir, 0o755); err != nil {
			return "", err
		}
		name := fmt.Sprintf("%d.png", time.Now().UnixNano())
		dst, err := os.Create(filepath.Join(avatarUploadDir, name))
		if err != nil {
			return "", err
		}
		defer dst.Close()
		if _, err := io.Copy(dst, f); err != nil {
			return "", err
		}
		return avatarPublicPrefix + "/" + name, nil
	})
}

// seedMemberProfiles 首次插入一条演示数据
func seedMemberProfiles(db *gorm.DB) {
	var count int64
	db.Model(&MemberProfile{}).Count(&count)
	if count > 0 {
		return
	}
	db.Create(&MemberProfile{Name: "示例成员", Avatar: ""})
}
