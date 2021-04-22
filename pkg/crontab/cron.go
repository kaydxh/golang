package cron

import (
	"fmt"

	"github.com/robfig/cron"
)

type CronProcessor struct {
	cronRunner *cron.Cron
}

type Job interface {
	Run()
}

func NewCronProcessor(specTime string, job Job) (*CronProcessor, error) {
	cp := &CronProcessor{
		cronRunner: cron.New(),
	}

	if cp.cronRunner == nil {
		return nil, fmt.Errorf("failed to init cronRunner")
	}

	if specTime == "" {
		return nil, fmt.Errorf("specTime is empty")
	}

	err := cp.cronRunner.AddJob(specTime, job)
	if err != nil {
		return nil, fmt.Errorf("failed to odd job in specTime: [%v]", specTime)
	}

	cp.cronRunner.Start()

	return cp, nil
}

func (c *CronProcessor) Stop() {
	c.cronRunner.Stop()
}
