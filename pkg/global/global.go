package global

import (
	"sync"

	"go.uber.org/zap"
)

var (
	App     string = "emri"
	Version string = "v0.0-dev"
	Build   string = "ad-hoc"

	_gcMu sync.RWMutex
	_gC   *GlobalConfig
)

// Return the current global configuration at request time
func C() *GlobalConfig {
	_gcMu.Lock()
	if _gC == nil {
		// TODO: this is just preventing nil deref for now
		// should do more here in this situation maybe
		zap.S().Warn("encountered nil global config")
		_gC = DefaultConfig()
	}
	defer _gcMu.Unlock()
	return _gC
}
