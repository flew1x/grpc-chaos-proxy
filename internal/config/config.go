package config

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"sync/atomic"
)

const (
	yamlType = "yaml"
)

type Loader struct {
	v   *viper.Viper
	cfg atomic.Value
	ch  chan struct{} // broadcast on change
}

func NewLoader(path string, log *zap.Logger) (*Loader, error) {
	v := viper.New()
	v.SetConfigFile(path)
	v.SetConfigType(yamlType)

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	var c Config

	if err := v.Unmarshal(&c); err != nil {
		return nil, err
	}

	l := &Loader{
		v:  v,
		ch: make(chan struct{}, 1),
	}

	l.cfg.Store(&c)
	v.WatchConfig()

	v.OnConfigChange(func(e fsnotify.Event) {
		log.Info("[config] file changed, reloading…", zap.String("file", e.Name))

		var nc Config

		if err := v.Unmarshal(&nc); err == nil {
			l.cfg.Store(&nc)
			select {
			case l.ch <- struct{}{}:
			default:
			}
		} else {
			log.Warn("failed to unmarshal updated config", zap.Error(err))
		}
	})

	setGlobalLoader(l)

	return l, nil
}

func (l *Loader) Notify() <-chan struct{} { return l.ch }

func (l *Loader) Current() *Config {
	cfg := l.cfg.Load()
	if cfg == nil {
		return nil
	}

	return cfg.(*Config)
}

var globalLoader *Loader

func setGlobalLoader(l *Loader) {
	globalLoader = l
}

func GetCurrentConfig() *Config {
	if globalLoader == nil {
		return nil
	}

	return globalLoader.Current()
}
