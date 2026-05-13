package casbin

import (
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
)

// NewEnforcer creates a Casbin enforcer backed by a GORM adapter.
// Uses a simple RBAC-like model for (role, path, method) enforcement.
func NewEnforcer(db *gorm.DB) (*casbin.Enforcer, error) {
	adapter, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		return nil, err
	}

	m := model.NewModel()
	m.AddDef("r", "r", "sub, obj, act")
	m.AddDef("p", "p", "sub, obj, act")
	m.AddDef("e", "e", "some(where (p.eft == allow))")
	m.AddDef("m", "m", "r.sub == p.sub && keyMatch(r.obj, p.obj) && regexMatch(r.act, p.act)")

	enforcer, err := casbin.NewEnforcer(m, adapter)
	if err != nil {
		return nil, err
	}

	if err := enforcer.LoadPolicy(); err != nil {
		return nil, err
	}

	return enforcer, nil
}

// AddPolicy adds a policy rule for a role: (role_code, path, method).
func AddPolicy(enforcer *casbin.Enforcer, roleCode, path, method string) (bool, error) {
	return enforcer.AddPolicy(roleCode, path, method)
}

// RemovePolicy removes a policy rule.
func RemovePolicy(enforcer *casbin.Enforcer, roleCode, path, method string) (bool, error) {
	return enforcer.RemovePolicy(roleCode, path, method)
}
