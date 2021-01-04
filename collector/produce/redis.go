package produce

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"time"
)

type Redis struct {
	Span *Span
}

var ctx = context.Background()

var client = redis.NewClient(&redis.Options{
	Addr:     "127.0.0.1:6379",
	PoolSize: 10,
})

func (span Redis) SendMessage() error {
	tm := time.Now()
	marshal, err := json.Marshal(span)

	if err != nil {
		return err
	}
	client.ZAdd(ctx, tm.Format("200601021504"), &redis.Z{
		Score:  float64(unixMs(tm)),
		Member: marshal,
	})
	return nil
}

func unixMs(tm time.Time) int64 {
	return tm.UnixNano() / int64(time.Millisecond)
}
