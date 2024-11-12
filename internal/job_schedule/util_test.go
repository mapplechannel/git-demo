package jobschedule

import (
	"fmt"
	"hsm-scheduling-back-end/internal/vo"
	"hsm-scheduling-back-end/pkg/logger"
	"testing"
	"time"
)

func TestAssemblyParams(t *testing.T) {
	job := vo.CreateJobRequest{}
	job.Integrated = "IOIT@hsm_io_it@job_dba3954o@e1"
	data := AssemblyParams(&job)
	logger.Info("data:%v", data)
}

func Test(t *testing.T) {
	endDate, err := time.Parse("2006-01-02 15:04:05", "2024-08-16 19:48:05")
	if err != nil {
		logger.Info("parse effectiveDate error:%v", err)
	}

	tarTime := time.Date(endDate.Year(), endDate.Month(), endDate.Day(), endDate.Hour(), endDate.Minute(), endDate.Second(), 0, time.Local)

	now := time.Now()
	fmt.Println("now.", endDate)
	fmt.Println("now.", now)

	fmt.Println("now.", now.Before(tarTime))

}
