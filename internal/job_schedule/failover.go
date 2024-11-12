package jobschedule

import (
	"hsm-scheduling-back-end/config"
	"hsm-scheduling-back-end/internal/constants"
	"hsm-scheduling-back-end/internal/util"
	"hsm-scheduling-back-end/internal/vo"
	"hsm-scheduling-back-end/pkg/logger"
	"io/ioutil"
	"net/http"
	"time"
)

func CheckAlive() bool /*isAlive*/ {
	// 1.确定当前主从
	logger.Info("util.IsMaster:%v", util.IsMaster)
	defer func() {
		if r := recover(); r != nil {
			logger.Info("Health check is error:%v", r)
		}
	}()
	// 2.监听主机状态
	url := constants.Http + util.BackupOrMasterIp + constants.AddressSplicingSymbols + config.ConfigAll.Port + constants.HealthPath

	logger.Info("url:%v", url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logger.Info("Fail to send request:%v", err)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	// 3.监听不到，执行任务转移
	if err != nil || resp.StatusCode != http.StatusOK {
		logger.Info("Failed to make request:%v", err)
		return false
	} else {
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		logger.Info("response body is:%v", string(body))
		if err != nil {
			logger.Info("Failed to read response body:%v", err)
		}

		if string(body) == "alive" {
			logger.Info("Master is alive")
		}
		return true
	}
}

func FailOver() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			isAlive := CheckAlive()
			// 接受不到另一个机器的心跳，同时运行节点不再当前机器上
			if !isAlive {
				// 4.获取主机上正在运行的任务（status=1)
				runningIds, err := GetRunningTaskId()
				if err != nil {
					logger.Info("Failed to get running ids:%v", err)
				}
				logger.Info("runningIds:%v", runningIds)
				if runningIds == nil {
					logger.Info("Running task id is nil")
				} else {
					// 5.将任务放入列表中
					// 如果任务已经转移，则不再执行
					// 运行任务
					for i := range runningIds {
						idVo := vo.GetID{ID: runningIds[i]}
						if !CheckRunningNode(runningIds[i]) {
							RunTaskForRedency(&idVo)
						}
					}
				}
			} else {
				logger.Info("CheckAlive is alive")
			}
		}
	}
}

func GetRunningTaskId() ([]string, error) {
	db, err := GetPgConn()
	if err != nil {
		logger.Info("Delete Fail to connect pg:%v", err)
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT id FROM hsm_scheduling.tasks WHERE status = $1", "1")
	if err != nil {
		logger.Info("Delete Fail to connect pg:%v", err)
		return nil, err
	}
	defer rows.Close()

	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			logger.Info("Rows scan id error:%v", err)
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}


