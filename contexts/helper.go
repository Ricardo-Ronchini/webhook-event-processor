package contexts

import (
	"time"

	"github.com/Ricardo-Ronchini/webhook-event-processor/common"
)

type Helper struct{}

func (h *Helper) TimeNow() time.Time {
	return time.Now().UTC()
}

func (h *Helper) TimeNowUTC(utc ...int) time.Time {
	now := time.Now()

	if len(utc) == 0 {
		return now.In(time.Local)
	}

	timezone := time.FixedZone("custom", utc[0]*3600)
	return now.In(timezone)
}

func (h *Helper) GenerateRandomID() string {
	return common.GenerateID()
}
