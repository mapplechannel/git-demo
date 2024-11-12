package jobschedule

import (
	"context"
	"encoding/json"
	"errors"
	"hsm-scheduling-back-end/internal/constants"
	"hsm-scheduling-back-end/internal/util"
	"hsm-scheduling-back-end/internal/vo"
	"hsm-scheduling-back-end/pkg/logger"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

var c *cron.Cron

func UpsertRunningLog(log *vo.RuningLog) error {
	var existingLog vo.RuningLog
	result := GlobalDb.Where("id = ?", log.ID).First(&existingLog)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			result = GlobalDb.Create(log)
		} else {
			return result.Error
		}
	} else {
		existingLog.JobName = log.JobName
		existingLog.Status = log.Status
		existingLog.TaskId = log.TaskId
		existingLog.StartTime = log.StartTime
		existingLog.EndTime = log.EndTime
		existingLog.Operator = log.Operator
		existingLog.Executor = log.Executor

		result = GlobalDb.Save(&existingLog)
	}
	return result.Error
}

func GetAllRuningLog() ([]vo.RuningLog, error) {
	var runningLog []vo.RuningLog
	result := GlobalDb.Order("starttime desc").Find(&runningLog)
	return runningLog, result.Error
}

var scheduleMap sync.Map
var cronMap sync.Map

func SendRetry(retryCount int, retryInterval, timeout time.Duration, url, reqbody, httpType, token, authtype string) ([]byte, error) {
	errS := ""
	for index := 0; index <= retryCount; index++ {
		res, err := Send(httpType, url, reqbody, token, authtype, timeout)
		if err == nil {
			return res, nil
		}

		if index != 0 {
			errS += "|" + err.Error()
		} else {
			errS += err.Error()
		}
		if index != 0 {
			logger.Info("正在尝试第%d次连接，共%d次", index, retryCount)
		}
		time.Sleep(retryInterval)
	}
	return nil, errors.New("send err:" + errS)
}

func Send(httpType, url, reqBody, token, authtype string, timeout time.Duration) ([]byte, error) {
	logger.Info("timeout:%v", timeout)
	client := &http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				conn, err := net.DialTimeout(netw, addr, timeout)
				if err != nil {
					return nil, err
				}

				_ = conn.SetDeadline(time.Now().Add(timeout))
				return conn, nil
			},
			ResponseHeaderTimeout: timeout,
		},
	}

	logger.Info("reqBody:%v", reqBody)

	requestDO, err := http.NewRequest(httpType, url, strings.NewReader(reqBody))
	if err != nil {
		logger.Info("err is:%v", err)
		return nil, errors.New("http err" + err.Error())
	}

	if token != "" {
		logger.Info("add token")
		requestDO.Header.Add(authtype, token)
	}

	requestDO.Header.Add("Content-Type", "application/json")

	res, err := client.Do(requestDO)
	if err != nil {
		logger.Info("err is:%v", err)
		return nil, errors.New("http err" + err.Error())
	}
	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logger.Info("err is:%v", err)
		return nil, errors.New("readall err" + err.Error())
	}

	return data, nil

}

func RunTask(task *vo.Task, runningLog *vo.RuningLog) {
	logger.Info("running id:%v task", task.ID)

	defer func() {
		if r := recover(); r != nil {
			logger.Info("捕捉到panic:%v", r)
		}
	}()

	_, job := FindTask(context.Background(), task.ID)
	//组装IOIT任务运行参数
	uniqueSlice := strings.Split(job.Integrated, "@")
	logger.Info("uniqueSlice is :%v", uniqueSlice)
	logger.Info("job.Integrated is :%v", job.Integrated)

	var IOITParams vo.IOIT
	var requestBody []byte
	var err1 error
	if uniqueSlice[0] == "IOIT" {
		IOITParams.NameSpace = uniqueSlice[1]
		newUrl := constants.Http + util.ServerCode + constants.DomainName + ":6120" + job.URL
		job.URL = newUrl
		// job.URL = "http://172.21.220.111:6120" + job.URL
		IOITParams.Code = uniqueSlice[2]
		IOITParams.Name = job.JobName
		requestBody, err1 = json.Marshal(IOITParams)
	} else {
		if job.Integrated == "" || job.Integrated == "KangYuan" {
			var mapBody map[string]interface{}

			if job.JobParams != "" {
				err := json.Unmarshal([]byte(job.JobParams), &mapBody)
				if err != nil {
					logger.Info("Fail to mapping:%v", err)
				}
			}

			logger.Info("job.JobParams is :%v", job.JobParams)
			logger.Info("mapBody is :%v", mapBody)

			requestBody, err1 = json.Marshal(mapBody)
		}
	}
	if err1 != nil {
		logger.Info("Fail to mapping:%v", err1)
	}

	logger.Info("Send request to :%v", job.URL)
	logger.Info("requestBody string is :%v", string(requestBody))
	var log vo.SysLog
	log.JobName = task.JobName
	log.TaskId = task.ID
	log.LogTime = time.Now().Format("2006-01-02 15:04:05")
	log.Executor = task.Executor

	var token string
	if job.AuthType != "none" {
		authtype := job.AuthType
		if authtype == "token" {
			token = job.Token
		}
		if authtype == "authorization" {
			token = job.Authorization
		}
	}

	data, err := SendRetry(job.RetryCount, time.Duration(job.RetryInterval)*time.Millisecond, time.Duration(job.Timeout)*time.Millisecond, job.URL, string(requestBody), job.RequestType, token, job.AuthType)
	logger.Info("response data is:%v", string(data))
	if err != nil {
		log.LogType = "error"
		log.Content = "执行失败: " + err.Error()
		logger.Info("Failed to read respBody")
		scheduleMap.Store(task.ID, false)
		if !strings.HasPrefix(task.Integrated, "IOIT") && task.Integrated != "KangYuan" {
			runningLog.Status = "0"
			log.InstanceId = runningLog.ID
			UpsertRunningLog(runningLog)
		}
		UpsertSysLog(&log)
	} else {
		log.LogType = "info"
		log.Content = "执行成功: " + string(data)
		log.ID = GetId()
		if !strings.HasPrefix(task.Integrated, "IOIT") && task.Integrated != "KangYuan" {
			log.InstanceId = runningLog.ID
			runningLog.Status = "1"
			UpsertRunningLog(runningLog)
		}
		UpsertSysLog(&log)
	}
}

func ManualRunTask(task *vo.Task, runningLog *vo.RuningLog) {
	logger.Info("running id:%v task", task.ID)

	uniqueSlice := strings.Split(task.Integrated, "@")
	logger.Info("uniqueSlice is :%v", uniqueSlice)
	logger.Info("job.Integrated is :%v", task.Integrated)

	var IOITParams vo.IOIT
	var requestBody []byte
	var err1 error
	if uniqueSlice[0] == "IOIT" {
		IOITParams.NameSpace = uniqueSlice[1]
		newUrl := constants.Http + util.ServerCode + constants.DomainName + ":6120" + task.URL
		task.URL = newUrl
		IOITParams.Code = uniqueSlice[2]
		IOITParams.Name = task.JobName

		requestBody, err1 = json.Marshal(IOITParams)
	} else {
		if task.Integrated == "" || task.Integrated == "KangYuan" {
			var mapBody map[string]interface{}

			if task.JobParams != "" {
				err := json.Unmarshal([]byte(task.JobParams), &mapBody)
				if err != nil {
					logger.Info("Fail to mapping:%v", err)
				}
			}

			logger.Info("job.JobParams is :%v", task.JobParams)
			logger.Info("mapBody is :%v", mapBody)

			requestBody, err1 = json.Marshal(mapBody)
		}
	}
	if err1 != nil {
		logger.Info("Fail to mapping:%v", err1)
	}

	logger.Info("Send request to :%v", task.URL)
	logger.Info("requestBody string is :%v", string(requestBody))

	var log vo.SysLog
	log.JobName = task.JobName
	log.TaskId = task.ID
	log.LogTime = time.Now().Format("2006-01-02 15:04:05")
	log.Executor = task.Executor

	var token string
	if task.AuthType != "none" {
		authtype := task.AuthType
		if authtype == "token" {
			token = task.Token
		}
		if authtype == "authorization" {
			token = task.Authorization
		}
	}

	if !strings.HasPrefix(task.Integrated, "IOIT") && task.Integrated != "KangYuan" {
		runningLog.StartTime = time.Now().Format("2006-01-02 15:04:05")
	}
	data, err := SendRetry(task.RetryCount, time.Duration(task.RetryInterval)*time.Millisecond, time.Duration(task.Timeout)*time.Millisecond, task.URL, string(requestBody), task.RequestType, token, task.AuthType)
	logger.Info("response data is:%v", string(data))
	if err != nil {
		log.LogType = "error"
		log.Content = "执行失败: " + err.Error()
		logger.Info("Failed to read respBody")
		scheduleMap.Store(task.ID, false)
		if !strings.HasPrefix(task.Integrated, "IOIT") && task.Integrated != "KangYuan" {
			runningLog.Status = "0"
			log.InstanceId = runningLog.ID
			if task.ScheduleType == constants.SCHEDULE_TYPE_NONE || (task.ScheduleType == constants.SCHEDULE_TYPE_CRON &&
				task.CronExpression == constants.EmptyContent) {
				runningLog.EndTime = time.Now().Format("2006-01-02 15:04:05")
			}
			UpsertRunningLog(runningLog)
		}
		UpsertSysLog(&log)
	} else {
		log.LogType = "info"
		log.Content = "执行成功: " + string(data)
		log.ID = GetId()
		//UpsertSysLog(&log)
		if !strings.HasPrefix(task.Integrated, "IOIT") && task.Integrated != "KangYuan" {
			log.InstanceId = runningLog.ID
			runningLog.Status = "1"
			if task.ScheduleType == constants.SCHEDULE_TYPE_NONE || (task.ScheduleType == constants.SCHEDULE_TYPE_CRON &&
				task.CronExpression == constants.EmptyContent) {
				runningLog.EndTime = time.Now().Format("2006-01-02 15:04:05")
			}
			UpsertRunningLog(runningLog)
		}
		UpsertSysLog(&log)
	}
	if strings.HasPrefix(task.Integrated, "IOIT") {
		SendStopToIOit(*task)
	}
	_ = UpdateStatus(task.ID, "0")
}

func DeleteCron(examId cron.EntryID) {
	// 移除定时任务
	c.Remove(examId)
}

// func StopTask(id string) bool {
// 	if cronId, ok := cronMap.Load(id); ok {
// 		DeleteCron(cronId.(cron.EntryID))
// 		logger.Info("delete cron:%v, id:%v", cronId, id)
// 		return true
// 	}
// 	return false
// }
