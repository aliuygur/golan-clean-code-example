package app

type DBWhere map[string]interface{}

type DBFilter struct {
	Limit   int
	Offset  int
	OrderBy string
	Reverse bool
	Preload []string
}

type DBCreator interface {
	Store(interface{}) error
}

type DBFinder interface {
	One(model interface{}, id interface{}) error
	OneBy(model interface{}, w DBWhere) error
	FindBy(models interface{}, w DBWhere, f *DBFilter) error
	FirstOrInit(m interface{}, w DBWhere) error
}

type DBExistser interface {
	ExistsBy(b interface{}, w DBWhere) (bool, error)
}

type DBCreatorExistser interface {
	DBCreator
	DBExistser
}

type DBUpdater interface {
	Save(model interface{}) error
	UpdateField(model interface{}, field string, value interface{}) error
	UpdateFields(model interface{}, kv map[string]interface{}) error
}

// Databaser is database interface
type Databaser interface {
	DBCreator
	DBExistser
	DBFinder
	DBUpdater
	IsNotFoundErr(error) bool
}
