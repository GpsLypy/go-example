
```go
// configFromFile loads config from file.
func (c *Config) configFromFile(path string) error {
	cfg, err := c.parseFileContent(path)
	if err != nil {
		return errors.Trace(err)
	}
	if cfg[0] == '{' {
		// Parse as a json object
		if err := json.Unmarshal([]byte(cfg), c); nil != err {
			return errors.Trace(err)
		}
	} else {
		// Parse as a toml object
		meta, err := toml.Decode(cfg, c)
		//meta, err := toml.DecodeFile(path, c)
		if err != nil {
			return errors.Trace(err)
		}
		if len(meta.Undecoded()) > 0 {
			return errors.Errorf("unknown keys in config file %s: %v", path, meta.Undecoded())
		}
	}
	return nil
}
```