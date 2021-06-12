package config

// Cfg return a singleton configuration instance.
func Cfg() Configuration {
	return cfg
}

// Load initialize this package. It will load config into cache and get ready to work.  However,
// if the config file was not found, `ErrFileNotFound` will be returned
func Load() error {
	return cfg.Load()
}

// LoadContent reads the val as content.
func LoadContent(content string) error {
	return cfg.LoadContent(content)
}

// FileName return config file name.  The default config file name is "app.yml"
func FileName() string {
	return cfg.FileName()
}

// SetFileName set a config file name to default configuration instance.  The default config file name is "app.yml"
func SetFileName(fileName string) {
	cfg.SetFileName(fileName)
}

// String returns a string type value which has the key.  If the value can't convert to string type,
func String(key string, defaultValue ...string) (string, error) {
	return cfg.String(key, defaultValue...)
}

// Set set a new value with key into config.  If the key doesn't exist, a new key will be created and no error be returned.
func Set(key string, val string) error {
	return cfg.Set(key, val)
}

// Int32 returns a int32 type value which has the key.  If the value can't convert to string type,
func Int32(key string, defaultValue ...int32) (int32, error) {
	return cfg.Int32(key, defaultValue...)
}
