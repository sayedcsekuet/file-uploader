package crons

import (
	"file-uploader/src/helpers"
	"fmt"
	"github.com/go-co-op/gocron"
	logger "github.com/sirupsen/logrus"
	"os/exec"
	"runtime"
	"time"
)

type Cron struct {
	*gocron.Scheduler
	callBacks []func() (*gocron.Job, error)
}

func NewCron() *Cron {
	return &Cron{
		Scheduler: gocron.NewScheduler(time.UTC),
		callBacks: nil,
	}
}
func (c *Cron) StartScheduler() {
	c.refreshClamAv()
	for _, callback := range c.callBacks {
		j, err := callback()
		logger.Infof("Schedule Job: %v, Error: %v", j, err)
	}
	c.StartAsync()
}

func (c *Cron) refreshClamAv() {
	cronStr := fmt.Sprintf("0 %d * * *", helpers.RandomInt(1, 24))
	logger.Infof("ClamAv cron scheduled at: %s", cronStr)
	j, err := c.Cron(cronStr).Do(func() {
		if runtime.GOOS == "windows" {
			return
		}
		//Refresh the calmav databse
		out, err := exec.Command("freshclam").Output()
		logger.Infof("ClamAv freshclam Command Output: %v, Error: %v", string(out), err)
	})
	logger.Infof("Schedule Job: %v, Error: %v", j, err)
}

func (c *Cron) Add(callback func() (*gocron.Job, error)) {
	c.callBacks = append(c.callBacks, callback)
}
