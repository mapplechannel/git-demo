if job.Schedule == "2" {
	if job.PeriodDetail.EndType == "times" {
		startDate, err := time.Parse("2006-01-02 15:04:05", job.CronDetail.EffectiveDate)
		if err != nil {
			logger.Info("parse effectiveDate error:%v", err)
		}

		tarTime := time.Date(startDate.Year(), startDate.Month(), startDate.Day(), startDate.Hour(), startDate.Minute(), startDate.Second(), 0, time.Local)

		duration := time.Until(tarTime)
		logger.Info("duration:%v", duration)

		var secondInte int
		cycle := job.PeriodDetail.ScheduleCycle
		switch job.PeriodDetail.ScheduleUnit {
		case "second":
			secondInte = cycle
		case "minute":
			secondInte = cycle * 60
		case "hour":
			secondInte = cycle * 3600
		default:
			logger.Info("default")
		}

		if duration <= 0 {
			logger.Info("time is in the past")
			duration = time.Duration(job.PeriodDetail.Times*secondInte+3) * time.Second
		} else {
			duration = duration + time.Duration(job.PeriodDetail.Times*secondInte+3)
		}

		logger.Info("距离停止时间还需:%v", duration)

		go func() {
			<-time.After(duration)
			run.Stop(job.Code)
			job.Status = 0
			util.UpdateJob(job)
		}()
		return
	}
}


type Period struct {
	ScheduleCycle int    `json:"schedule_cycle"`
	ScheduleUnit  string `json:"schedule_unit"`
	Times         int    `json:"times"`
	EffectiveDate string `json:"effectiveDate"`
	EndDate       string `json:"endDate"` // endless(无限期) | date(结束时间)
	EndType       string `json:"endType"` // endless(无限期) | date(结束时间)
}

type CronExp struct {
	Type          string `json:"type"`          // day,week,month
	EffectiveDate string `json:"effectiveDate"` // 年月日时分秒
	EndDate       string `json:"endDate"`       // date | endless(无限期)
	EndType       string `json:"endType"`       // date | endless(无限期)
	DayOfWeek     string `json:"dayOfWeek"`     //1，2，3
	DayOfMonth    string `json:"dayOfMonth"`    //1，2，3，20
	ExecuteMonth  string `json:"executeMonth"`  // 1，2，3
}