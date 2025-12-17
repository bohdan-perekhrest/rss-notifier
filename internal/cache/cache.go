package cache

type Cache interface {
	GetLastItemURL(url string) (string, error)
	SaveLastItemURL(urls map[string]string) error
	Close() error
}
