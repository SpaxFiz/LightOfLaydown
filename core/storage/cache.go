// @program: unjuanable
// @author: Fizzy
// @created: 2021-11-25

package storage

import (
	"github.com/patrickmn/go-cache"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spaxfiz/unjuanable/config"
	"reflect"
	"sync"
	"time"
)

var (
	defaultIns *Cache
	lock       sync.Mutex
)

type Cache struct {
	c    *cache.Cache
	lock sync.Mutex
}

func GetCache() *Cache {
	lock.Lock()
	defer lock.Unlock()

	if defaultIns == nil {
		duration := time.Duration(config.GlobalConfig.Server.Cache.CacheDurationSecond) * time.Second
		defaultCache := cache.New(duration, duration)
		if err := defaultCache.LoadFile(config.GlobalConfig.Server.Cache.CachePath); err != nil {
			logrus.Warnf("load cache file failed. cache will be empty. err=%v", err)
		}
		defaultIns = &Cache{c: defaultCache}
	}
	return defaultIns
}

// LoadOrDo
/**
@param restore pointer of data receiver
*/
func (c *Cache) LoadOrDo(key string, restore interface{}, fn func() (shouldCache interface{}, err error)) error {
	restoreVal := reflect.Indirect(reflect.ValueOf(restore))

ReadCache:
	if val, hit := c.c.Get(key); hit {
		cacheType := reflect.TypeOf(val).Kind().String()

		if !restoreVal.CanSet() {
			return errors.New("restore cannot be set")
		}
		if cacheType == restoreVal.Kind().String() {
			restoreVal.Set(reflect.ValueOf(val))
			return nil
		} else {
			return errors.Errorf("incompatible cache type, got=%s, want=%s", cacheType, restoreVal.Kind().String())
		}
	}
	c.lock.Lock()
	defer c.lock.Unlock()
	if _, hit := c.c.Get(key); hit {
		goto ReadCache
	}
	if shouldCache, err := fn(); err != nil {
		return err
	} else if v := reflect.ValueOf(shouldCache); v.Kind() == reflect.Slice && v.Len() == 0 {
		return nil
	} else {
		restoreVal.Set(reflect.ValueOf(shouldCache))
		//c.c.Set(key, shouldCache, beforeTomorrowDuration())
		c.c.Set(key, shouldCache, time.Hour*2)
	}
	return nil
}

func (c *Cache) Save(key string, val interface{}, ttl time.Duration) {
	c.c.Set(key, val, ttl)
}

// SaveBeforeTomorrow
// @description save cache and set cache expire at the beginning of tomorrow
func (c *Cache) SaveBeforeTomorrow(key string, val interface{}) {
	c.c.Set(key, val, beforeTomorrowDuration())
}

func (c *Cache) Load(key string) (val interface{}, ok bool) {
	return c.c.Get(key)
}

// PersistCache
// @description save current cache to local file, wait for next program startup then load
func (c *Cache) PersistCache() error {
	return c.c.SaveFile(config.GlobalConfig.Server.Cache.CachePath)
}

func beforeTomorrowDuration() time.Duration {
	nextDay := time.Now().AddDate(0, 0, 1)
	y, m, d := nextDay.Date()
	nextDayTS := time.Date(y, m, d, 0, 0, 0, 0, time.Local).Second()
	return time.Duration(time.Duration(nextDayTS-time.Now().Second()) * time.Second)
}
