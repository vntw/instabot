package cache

type Cache interface {
	// Return -1 if no entry exists
	ReadLastDate() (int64, error)
	WriteLastDate(date int64) error
}
