package config

type Provider interface {
	NotifyChange() chan bool
	Get(key string) (interface{}, error)
}
