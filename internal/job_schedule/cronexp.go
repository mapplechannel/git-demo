package jobschedule

import (
	"fmt"
	"hsm-scheduling-back-end/internal/vo"
	"hsm-scheduling-back-end/pkg/logger"
	"strings"
	"time"
)

func generateCronExpression(cronExp vo.CronExp) string {
	switch cronExp.Type {
	case "hour":
		return generateHourCronExpression(cronExp)
	case "day":
		return generateDayCronExpression(cronExp)
	case "week":
		return generateWeekCronExpression(cronExp)
	case "month":
		return generateMonthCronExpression(cronExp)
	default:
		return ""
	}
}

func generateHourCronExpression(cronExp vo.CronExp) string {
	timeArr := strings.Split(cronExp.Time, ":")
	return fmt.Sprintf("%s %s * * * *", timeArr[1], timeArr[0])
}

func generateDayCronExpression(cronExp vo.CronExp) string {
	timeArr := strings.Split(cronExp.Time, ":")
	return fmt.Sprintf("%s %s %s * * *", timeArr[2], timeArr[1], timeArr[0])
}

func generateWeekCronExpression(cronExp vo.CronExp) string {
	timeArr := strings.Split(cronExp.Time, ":")
	return fmt.Sprintf("%s %s %s * * %s", timeArr[2], timeArr[1], timeArr[0], cronExp.DayOfWeek)
}

func generateMonthCronExpression(cronExp vo.CronExp) string {
	if cronExp.DayOfMonth == "0" {
		parseDate, err := time.Parse("15:04:05", cronExp.Time)
		if err != nil {
			logger.Info("parse effectiveDate error:%v", err)
		}
		h := parseDate.Hour()
		m := parseDate.Minute()
		s := parseDate.Second()
		return fmt.Sprintf("%d %d %d * * *", s, m, h)
	}
	logger.Info("cronexp is :%v", cronExp)
	timeArr := strings.Split(cronExp.Time, ":")
	return fmt.Sprintf("%s %s %s %s * ?", timeArr[2], timeArr[1], timeArr[0], cronExp.DayOfMonth)
}
