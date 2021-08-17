package wrapper

import (
	"time"
)

func NewTimeService() *TimeService {
	return &TimeService{}
}

type TimeService struct{}

func (s *TimeService) Now() time.Time {
	return time.Now()
}
