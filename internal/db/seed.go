package db

import (
	"log/slog"

	"gogo/internal/model"
	"gogo/internal/pkg"

	"gorm.io/gorm"
)

// Seed inserts initial data if not already present.
func Seed(db *gorm.DB) error {
	// Check if SUPER_ADMIN role already exists
	var count int64
	db.Model(&model.Role{}).Where("code = ?", "SUPER_ADMIN").Count(&count)
	if count > 0 {
		slog.Info("seed data already exists, skipping")
		return nil
	}

	slog.Info("seeding initial data")

	if err := doSeed(db); err != nil {
		return err
	}

	slog.Info("seed data inserted successfully")
	return nil
}

func doSeed(db *gorm.DB) error {
	return db.Transaction(func(tx *gorm.DB) error {
		// 1. SUPER_ADMIN role
		superAdminRole := &model.Role{
			Name:        "超级管理员",
			Code:        "SUPER_ADMIN",
			Description: "系统内置超级管理员，拥有所有权限",
			Status:      1,
		}
		if err := tx.Create(superAdminRole).Error; err != nil {
			return err
		}

		// 2. Admin user (password: Admin123!)
		hash, err := pkg.HashPassword("Admin123!")
		if err != nil {
			return err
		}
		adminUser := &model.User{
			Username: "admin",
			Email:    "admin@system.local",
			Password: hash,
			RealName: "系统管理员",
			Status:   1,
		}
		if err := tx.Create(adminUser).Error; err != nil {
			return err
		}

		// 3. Assign SUPER_ADMIN role to admin
		if err := tx.Exec("INSERT INTO user_roles (user_id, role_id) VALUES (?, ?)", adminUser.ID, superAdminRole.ID).Error; err != nil {
			return err
		}

		// 4. Base menus
		rootMenu := &model.Menu{
			ParentID:  0,
			Name:      "系统管理",
			Path:      "/system",
			Component: "Layout",
			Icon:      "system",
			Type:      1,
			SortOrder: 1,
			Visible:   true,
			Status:    1,
		}
		if err := tx.Create(rootMenu).Error; err != nil {
			return err
		}

		childMenus := []*model.Menu{
			{ParentID: rootMenu.ID, Name: "用户管理", Path: "/system/user", Component: "system/user/index", Icon: "user", Type: 2, SortOrder: 1, Visible: true, Status: 1},
			{ParentID: rootMenu.ID, Name: "角色管理", Path: "/system/role", Component: "system/role/index", Icon: "role", Type: 2, SortOrder: 2, Visible: true, Status: 1},
			{ParentID: rootMenu.ID, Name: "菜单管理", Path: "/system/menu", Component: "system/menu/index", Icon: "menu", Type: 2, SortOrder: 3, Visible: true, Status: 1},
			{ParentID: rootMenu.ID, Name: "门店管理", Path: "/system/store", Component: "system/store/index", Icon: "store", Type: 2, SortOrder: 4, Visible: true, Status: 1},
			{ParentID: rootMenu.ID, Name: "终端管理", Path: "/system/terminal", Component: "system/terminal/index", Icon: "terminal", Type: 2, SortOrder: 5, Visible: true, Status: 1},
			{ParentID: rootMenu.ID, Name: "日志管理", Path: "/system/log", Component: "system/log/index", Icon: "log", Type: 2, SortOrder: 6, Visible: true, Status: 1},
		}
		buttonPerms := []struct {
			parentIdx int
			name      string
			perms     string
			sortOrder int
		}{
			{0, "创建用户", "sys:user:add", 1},
			{0, "编辑用户", "sys:user:edit", 2},
			{0, "删除用户", "sys:user:delete", 3},
			{0, "重置密码", "sys:user:reset-pwd", 4},
			{0, "分配角色", "sys:user:assign-role", 5},
			{0, "分配门店", "sys:user:assign-store", 6},
			{0, "用户列表", "sys:user:list", 7},
			{1, "创建角色", "sys:role:add", 1},
			{1, "编辑角色", "sys:role:edit", 2},
			{1, "删除角色", "sys:role:delete", 3},
			{1, "分配菜单", "sys:role:assign-menu", 4},
			{1, "角色列表", "sys:role:list", 5},
			{2, "创建菜单", "sys:menu:add", 1},
			{2, "编辑菜单", "sys:menu:edit", 2},
			{2, "删除菜单", "sys:menu:delete", 3},
			{2, "菜单列表", "sys:menu:list", 4},
			{3, "创建门店", "sys:store:add", 1},
			{3, "编辑门店", "sys:store:edit", 2},
			{3, "删除门店", "sys:store:delete", 3},
			{3, "门店列表", "sys:store:list", 4},
			{4, "创建终端", "sys:terminal:add", 1},
			{4, "编辑终端", "sys:terminal:edit", 2},
			{4, "删除终端", "sys:terminal:delete", 3},
			{4, "终端列表", "sys:terminal:list", 4},
			{5, "日志列表", "sys:log:list", 1},
		}

		for _, child := range childMenus {
			if err := tx.Create(child).Error; err != nil {
				return err
			}
		}

		for _, bp := range buttonPerms {
			button := &model.Menu{
				ParentID:  childMenus[bp.parentIdx].ID,
				Name:      bp.name,
				Type:      3,
				Perms:     bp.perms,
				SortOrder: bp.sortOrder,
				Visible:   true,
				Status:    1,
			}
			if err := tx.Create(button).Error; err != nil {
				return err
			}
		}

		// 5. Assign all menus to SUPER_ADMIN role
		var allMenuIDs []int64
		if err := tx.Model(&model.Menu{}).Pluck("id", &allMenuIDs).Error; err != nil {
			return err
		}
		for _, menuID := range allMenuIDs {
			if err := tx.Exec("INSERT INTO role_menus (role_id, menu_id) VALUES (?, ?)", superAdminRole.ID, menuID).Error; err != nil {
				return err
			}
		}

		return nil
	})
}
