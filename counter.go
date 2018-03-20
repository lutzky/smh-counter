package counter

import (
	"time"
)

type Counter struct {
	buckets   [3600]uint64
	latest    time.Time
	anyCounts bool
}

func hourSecond(t time.Time) int {
	return t.Minute()*60 + t.Second()
}

func (c *Counter) clear(now time.Time) {
	if !c.anyCounts {
		return
	}

	offset := (hourSecond(c.latest) + 1) % 3600
	var secsToClear int
	if now.Sub(c.latest).Hours() > 1 {
		secsToClear = 3600
	} else {
		secsToClear = int(now.Sub(c.latest).Seconds())
	}

	for i := 0; i < secsToClear; i++ {
		c.buckets[(i+offset)%3600] = 0
	}
}

func (c *Counter) GetSecond(now time.Time) uint64 {
	c.clear(now)

	return c.buckets[hourSecond(now)]
}

func (c *Counter) GetMinute(now time.Time) uint64 {
	c.clear(now)

	sum := uint64(0)
	offset := hourSecond(now) + 3600 - 60 + 1 /* include "current" second */

	for i := 0; i < 60; i++ {
		sum += c.buckets[(offset+i)%3600]
	}

	return sum
}

func (c *Counter) GetHour(now time.Time) uint64 {
	c.clear(now)

	sum := uint64(0)
	for _, b := range c.buckets {
		sum += b
	}
	return sum
}

func (c *Counter) Count(now time.Time) {
	c.clear(now)
	c.buckets[hourSecond(now)]++
	c.anyCounts = true
	c.latest = now
}

func New() Counter {
	return Counter{}
}
