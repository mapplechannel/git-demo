package jobschedule

import (
	"bytes"
	"encoding/json"
	"fmt"
	"hsm-scheduling-back-end/config"
	"hsm-scheduling-back-end/internal/constants"
	"hsm-scheduling-back-end/internal/util"
	"hsm-scheduling-back-end/internal/vo"
	"hsm-scheduling-back-end/pkg/logger"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var GlobalDb *gorm.DB

func GetPgConnByGORM() (*gorm.DB, error) {
	dsn := config.ConfigAll.Postgres.Url + " dbname=" + config.ConfigAll.Postgres.DbName
	// logger.Info("Pg connect dsn:%v", dsn)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Info("Pg connect failed:%v", err)
		return nil, err
	}
	GlobalDb = db
	return db, nil
}

func AssemblyParams(job *vo.CreateJobRequest) *vo.IOIT {
	uniqueSlice := strings.Split(job.Integrated, "@")
	var IOITParams vo.IOIT
	if uniqueSlice[0] == "IOIT" {
		IOITParams.NameSpace = uniqueSlice[1]
		IOITParams.Code = uniqueSlice[2]
		IOITParams.SKNode = uniqueSlice[3]
		IOITParams.Name = job.JobName
	}
	return &IOITParams
}

func UpdateStatus(ID string, status string) bool {

	err := GlobalDb.Table("hsm_scheduling.tasks").Where("id = ?", ID).Update("status", status).Error

	if err != nil {
		logger.Info("Update Fail to exec statement:%v", err)
		return false
	}
	logger.Info("Jobinfo update for success")
	return true
}

// func UpdateStatusAndEndTime(instanceId, status string) {
// 	db, err := GetPgConnByGORM()

// 	if err != nil {
// 		logger.Info("UpdateStatus Update Fail to connect pg:%v", err)
// 	}

// 	db.Table("hsm_scheduling.instances").Updates(TaskInstance{Status: status, EndTime: time.Now().Format("2006-01-02 15:04:05")})
// 	logger.Info("UpdateStatusAndEndTime update for success")
// }

// 同名校验
func CheckName(addName string) bool {
	db, err := GetPgConn()

	if err != nil {
		logger.Info("Update Fail to connect pg:%v", err)
	}

	query := `SELECT jobname FROM hsm_scheduling.tasks`

	rows, err := db.Query(query)
	if err != nil {
		logger.Info("Update Fail to connect pg:%v", err)
	}
	defer rows.Close()

	var names []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			logger.Info("Update Fail to connect pg:%v", err)
		}
		names = append(names, name)
	}
	for _, name := range names {
		if name == addName {
			return true
		}
	}
	return false
}

func CheckExecutor(hostname string) bool {
	db, err := GetPgConn()
	if err != nil {
		logger.Info("数据库连接失败")
	}
	defer db.Close()

	var exists bool
	query := `SELECT EXISTS ( SELECT 1 FROM hsm_scheduling.hosts WHERE hostname = $1)`

	err = db.QueryRow(query, hostname).Scan(&exists)
	if err != nil {
		logger.Info("数据库查询失败%v", err)
	}
	return exists
}

// 判断运行节点
func CheckRunningNode(id string) bool {
	db, err := GetPgConnByGORM()
	if err != nil {
		logger.Info("Pg connect failed:%v", err)
	}
	var failoverNode vo.FailOverNode
	result := db.Table("hsm_scheduling.failover_nodes").First(&failoverNode, id)
	if result.Error != nil {
		logger.Info("Pg connect failed:%v", err)
	}
	if failoverNode.RunningNode != util.ServerCode {
		logger.Info("check running node is finish false")
		return false
	} else {
		logger.Info("check running node is finish true")
		return true
	}
}

func GetInstanceId() string {
	sec := time.Now().Unix()
	str := strconv.FormatInt(sec, 10)

	rand.Seed(time.Now().UnixNano())

	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	code := ""
	code += string(letters[rand.Intn(26)]) // 第二位为小写字母
	for i := 0; i < 7; i++ {
		code += string(letters[rand.Intn(len(letters))])
	}

	res := str + "-" + code
	return res
}

// 调用调度系统接口
func SendStopToIOit(task vo.Task) {
	uniqueSlice := strings.Split(task.Integrated, "@")
	var IOITParams vo.IOIT
	var requestBody []byte
	IOITParams.NameSpace = uniqueSlice[1]
	IOITParams.Code = uniqueSlice[2]
	requestBody, err := json.Marshal(IOITParams)
	if err != nil {
		logger.Info("Fail to mapping:%v", err)
	}

	ip := util.ServerCode + constants.DomainName
	logger.Info("ip is :%v", ip)

	req, err := http.NewRequest("POST", constants.Http+ip+":6120/api/hsm-io-it/job/stop", bytes.NewBuffer(requestBody))
	if err != nil {
		logger.Info("Fail to send request:%v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Info("Failed to excute job:%v", err)
	} else {
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			logger.Info("Failed to excute job:%v", err)
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			logger.Info("Failed to excute job:%v", err)
		}
		logger.Info("body:%v", string(body))
	}
}

func SendToAnotherSysStop(task vo.Task, runningNode string) error {
	idMap := make(map[string]string)
	idMap["id"] = task.ID
	requestBody, err := json.Marshal(idMap)
	if err != nil {
		logger.Info("Fail to mapping:%v", err)
		return fmt.Errorf("fail to mapping:%s", err.Error())
	}

	ip := runningNode + constants.DomainName

	req, err := http.NewRequest("POST", constants.Http+ip+":6222/api/hsm-ds/task/stop", bytes.NewBuffer(requestBody))
	if err != nil {
		logger.Info("Fail to send request:%v", err)
		return fmt.Errorf("fail to send request:%s", err.Error())
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Info("Failed to excute job:%v", err)
		return fmt.Errorf("failed to excute job:%s", err.Error())
	} else {
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			logger.Info("Failed to excute job:%v", err)
			return fmt.Errorf("通知执行节点停止作业返回错误:%s", err.Error())
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			logger.Info("Failed to excute job:%v", err)
			return fmt.Errorf("解析返回信息失败:%s", err.Error())
		}
		logger.Info("body:%v", string(body))
	}
	return nil
}
