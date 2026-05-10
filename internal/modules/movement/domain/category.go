package domain

type Category struct {
	id        CategoryID
	name      string
	isSystem  bool
	systemKey string
}

func RehydrateCategory(id CategoryID, name string, isSystem bool, systemKey string) (Category, error) {
	if name == "" {
		return Category{}, ErrInvalidCategory
	}
	if isSystem && systemKey == "" {
		return Category{}, ErrInvalidCategory
	}

	return Category{
		id:        id,
		name:      name,
		isSystem:  isSystem,
		systemKey: systemKey,
	}, nil
}

func (c Category) ID() CategoryID    { return c.id }
func (c Category) Name() string      { return c.name }
func (c Category) IsSystem() bool    { return c.isSystem }
func (c Category) SystemKey() string { return c.systemKey }
