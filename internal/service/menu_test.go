package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gogo/internal/dto"
	"gogo/internal/model"
)

// menuRepoStub is a stub for MenuRepository.
type menuRepoStub struct {
	menus map[int64]*model.Menu
	nextID int64
}

func newMenuRepoStub() *menuRepoStub {
	return &menuRepoStub{menus: make(map[int64]*model.Menu), nextID: 1}
}

func (s *menuRepoStub) Create(ctx context.Context, menu *model.Menu) error {
	menu.ID = s.nextID
	s.nextID++
	s.menus[menu.ID] = menu
	return nil
}

func (s *menuRepoStub) GetByID(ctx context.Context, id int64) (*model.Menu, error) {
	m, ok := s.menus[id]
	if !ok {
		return nil, nil
	}
	return m, nil
}

func (s *menuRepoStub) List(ctx context.Context) ([]model.Menu, error) {
	var result []model.Menu
	for _, m := range s.menus {
		result = append(result, *m)
	}
	return result, nil
}

func (s *menuRepoStub) Update(ctx context.Context, menu *model.Menu) error {
	s.menus[menu.ID] = menu
	return nil
}

func (s *menuRepoStub) Delete(ctx context.Context, id int64) error {
	delete(s.menus, id)
	return nil
}

func (s *menuRepoStub) HasChildren(ctx context.Context, parentID int64) (bool, error) {
	for _, m := range s.menus {
		if m.ParentID == parentID {
			return true, nil
		}
	}
	return false, nil
}

func (s *menuRepoStub) GetMenusByRoleID(ctx context.Context, roleID int64) ([]int64, error) {
	return nil, nil
}

func TestMenuService_Tree(t *testing.T) {
	repo := newMenuRepoStub()
	svc := NewMenuService(repo)

	// Create root
	root, _ := svc.Create(context.Background(), dto.CreateMenuReq{
		Name: "系统管理", Type: 1, SortOrder: 1,
	})

	// Create children
	svc.Create(context.Background(), dto.CreateMenuReq{
		ParentID: root.ID, Name: "用户管理", Type: 2, SortOrder: 1,
	})
	svc.Create(context.Background(), dto.CreateMenuReq{
		ParentID: root.ID, Name: "角色管理", Type: 2, SortOrder: 2,
	})

	// Create a button under user management
	child, _ := svc.Create(context.Background(), dto.CreateMenuReq{
		ParentID: root.ID, Name: "用户管理", Type: 2, SortOrder: 1,
	})
	svc.Create(context.Background(), dto.CreateMenuReq{
		ParentID: child.ID, Name: "创建用户", Type: 3, Perms: "sys:user:add",
	})

	tree, err := svc.Tree(context.Background())
	require.NoError(t, err)
	assert.Len(t, tree, 1) // root node
	assert.Equal(t, "系统管理", tree[0].Name)
	assert.Len(t, tree[0].Children, 3)
}

func TestMenuService_DeleteWithChildren(t *testing.T) {
	repo := newMenuRepoStub()
	svc := NewMenuService(repo)

	root, _ := svc.Create(context.Background(), dto.CreateMenuReq{Name: "系统管理", Type: 1, SortOrder: 1})
	svc.Create(context.Background(), dto.CreateMenuReq{ParentID: root.ID, Name: "用户管理", Type: 2, SortOrder: 1})

	err := svc.Delete(context.Background(), root.ID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "存在子菜单")
}

func TestMenuService_DeleteLeaf(t *testing.T) {
	repo := newMenuRepoStub()
	svc := NewMenuService(repo)

	leaf, _ := svc.Create(context.Background(), dto.CreateMenuReq{Name: "创建用户", Type: 3, Perms: "sys:user:add"})

	err := svc.Delete(context.Background(), leaf.ID)
	assert.NoError(t, err)

	_, err = svc.GetByID(context.Background(), leaf.ID)
	assert.Error(t, err)
}

func TestBuildTree_Empty(t *testing.T) {
	tree := buildTree(nil, 0)
	assert.Nil(t, tree)
}

func TestBuildTree_SingleNode(t *testing.T) {
	menus := []model.Menu{
		{ID: 1, ParentID: 0, Name: "Root", Type: 1},
	}
	tree := buildTree(menus, 0)
	assert.Len(t, tree, 1)
	assert.Equal(t, "Root", tree[0].Name)
	assert.Empty(t, tree[0].Children)
}

func TestBuildTree_Nested(t *testing.T) {
	menus := []model.Menu{
		{ID: 1, ParentID: 0, Name: "A", Type: 1},
		{ID: 2, ParentID: 1, Name: "A1", Type: 2},
		{ID: 3, ParentID: 1, Name: "A2", Type: 2},
		{ID: 4, ParentID: 0, Name: "B", Type: 1},
	}

	tree := buildTree(menus, 0)
	assert.Len(t, tree, 2)
	assert.Equal(t, "A", tree[0].Name)
	assert.Len(t, tree[0].Children, 2)
	assert.Equal(t, "B", tree[1].Name)
	assert.Empty(t, tree[1].Children)
}
