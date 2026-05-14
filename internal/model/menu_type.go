package model

// MenuType represents the type of menu item.
type MenuType int16

const (
	MenuTypeDirectory MenuType = 1
	MenuTypePage      MenuType = 2
	MenuTypeButton    MenuType = 3
)

func (t MenuType) String() string {
	switch t {
	case MenuTypeDirectory:
		return "directory"
	case MenuTypePage:
		return "page"
	case MenuTypeButton:
		return "button"
	default:
		return "unknown"
	}
}
