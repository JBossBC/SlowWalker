package test

import (
	"replite_web/internal/app/utils"
	"testing"
	"time"
)

const DEFUALT_TEST_TIMES = 10

func TestRandom(t *testing.T) {
	var mapRand map[string]any = make(map[string]any)
	for i := 0; i < DEFUALT_TEST_TIMES; i++ {
		str := utils.NewRandomCode()
		if _, ok := mapRand[str]; ok {
			t.Fatalf("出现重复的验证码")
		}
		mapRand[str] = nil
		time.Sleep(1 * time.Nanosecond)
	}

}
