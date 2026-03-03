package redis_st

import (
	"errors"
	"testing"
)

func TestRedisSt(t *testing.T) {
	t.Run("login", func(t *testing.T) {
		rdb := NewRedis()
		if err := rdb.Connect(); err != nil {
			t.Error(err)
			return
		}
		if err := rdb.Disconnect(); err != nil {
			t.Error(err)
		}
	})

	t.Run("set-get", func(t *testing.T) {
		rdb := NewRedis()
		if err := rdb.Connect(); err != nil {
			t.Error(err)
		}
		rdb.Set("key", 1024)
		val := rdb.GetInt("key")
		if val != 1024 {
			t.Error(errors.New("wrong stored value"))
		}
		if err := rdb.Disconnect(); err != nil {
			t.Error(err)
		}
	})

}
