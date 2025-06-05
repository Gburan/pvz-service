package nower

import "time"

type Nower struct{}

func (n Nower) Now() time.Time {
	return time.Now()
}
