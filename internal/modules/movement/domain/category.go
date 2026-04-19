package domain

type Category struct {
	id   CategoryID
	name string
}

func RehydrateCategory(id CategoryID, name string) (Category, error) {
	if name == "" {
		return Category{}, ErrInvalidCategory
	}

	return Category{
		id:   id,
		name: name,
	}, nil
}

func (c Category) ID() CategoryID { return c.id }
func (c Category) Name() string   { return c.name }
