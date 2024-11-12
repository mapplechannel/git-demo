package jobschedule

import (
	"hsm-scheduling-back-end/internal/vo"
	"hsm-scheduling-back-end/pkg/logger"

	"gorm.io/gorm"
)

func UpsertSysLog(syslog *vo.SysLog) error {
	var existingLog vo.SysLog
	result := GlobalDb.Where("id = ?", syslog.ID).First(&existingLog)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			logger.Info("insert data")
			syslog.ID = GetId()
			result = GlobalDb.Create(syslog)
		} else {
			return result.Error
		}
	} else {
		logger.Info("update data")
		existingLog.ID = syslog.ID
		existingLog.JobName = syslog.JobName
		existingLog.LogType = syslog.LogType
		existingLog.TaskId = syslog.TaskId
		existingLog.Content = syslog.Content
		existingLog.LogTime = syslog.LogTime
		existingLog.Executor = syslog.Executor
		result = GlobalDb.Save(&existingLog)
	}
	return result.Error
}

func GetAllSysLog() ([]vo.SysLog, error) {
	var sysLogs []vo.SysLog
	result := GlobalDb.Order("logtime desc").Limit(1000).Find(&sysLogs)
	logger.Info("result.Error is :%v", result.Error)
	return sysLogs, result.Error
}
