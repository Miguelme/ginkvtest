package service

// DatabaseService defines an interface for database operations
type DatabaseService interface {
	GetValueByKey(key string) (string, error)
}
