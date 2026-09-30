package shop

import (
	"fmt"

	"example/models"

	"github.com/r0vx/x/login"
	"github.com/r0vx/x/perm"
	"gorm.io/gorm"
)

// DemoAccounts 三种角色的演示账号（spec §9：运营管商品、客服管订单、译者只翻译指定语言）。
// 密码与种子管理员相同（LOGIN_INITIAL_USER_PASSWORD），权限策略见 admin/perm.go，译者的语言限定见 config.go。
var DemoAccounts = map[string]string{
	"shop-operator":   models.RoleShopOperator,
	"shop-support":    models.RoleShopSupport,
	"shop-translator": models.RoleShopTranslator,
}

// seedDemoUsers 建三种角色与各自的演示账号；已存在的跳过。密码为空时不建账号。
func seedDemoUsers(db *gorm.DB, password string) error {
	for account, roleName := range DemoAccounts {
		role := perm.Role{Name: roleName}
		if err := db.Where("name = ?", roleName).FirstOrCreate(&role).Error; err != nil {
			return fmt.Errorf("shop: 建角色 %s: %w", roleName, err)
		}
		if password == "" {
			continue
		}
		var n int64
		if err := db.Model(&models.User{}).Where("account = ?", account).Count(&n).Error; err != nil {
			return fmt.Errorf("shop: 查账号 %s: %w", account, err)
		}
		if n > 0 {
			continue
		}
		u := &models.User{Name: account, Status: models.StatusActive,
			UserPass: login.UserPass{Account: account, Password: password, Confirmed: true}}
		u.EncryptPassword()
		u.Roles = []perm.Role{role}
		if err := db.Create(u).Error; err != nil {
			return fmt.Errorf("shop: 建账号 %s: %w", account, err)
		}
	}
	return nil
}
