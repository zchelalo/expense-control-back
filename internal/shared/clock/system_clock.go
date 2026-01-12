package clock

import "time"

type SystemClock struct{}

func (sc *SystemClock) Now() time.Time { return time.Now().UTC() }