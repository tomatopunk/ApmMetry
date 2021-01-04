package produce

import (
	"testing"
	"time"
)

func TestRedisPing(t *testing.T) {
	if pingRes := client.Ping(ctx); pingRes.Err() != nil {
		t.Errorf(pingRes.Val())
	}
}
func TestRedisWrite(t *testing.T) {
	if err := client.Set(ctx, "test", "aaabbbccc", 0).Err(); err != nil {
		t.Errorf(err.Error())
	}

	if val, err2 := client.Get(ctx, "test").Result(); err2 != nil {
		t.Errorf(err2.Error())
	} else {
		if val != "aaabbbccc" {
			t.Errorf("Val is Error!")
		}
	}
}

func TestSendMessage(t *testing.T) {
	redisSpan := Redis{
		&Span{
			SpanId: "asdasdas",
		},
	}
	if err := redisSpan.SendMessage(); err != nil {
		t.Error(err)
	}
	val, err2 := client.ZRange(ctx, time.Now().Format("200601021504"), 0, 1).Result()
	if err2 != nil {
		t.Log(val)
		t.Error(err2)
	}
	t.Log(val)
}
