package jobschedule

import (
	"context"
	"hsm-scheduling-back-end/internal/constants"
	"hsm-scheduling-back-end/internal/util"
	"hsm-scheduling-back-end/internal/vo"
	"hsm-scheduling-back-end/pkg/logger"
	"strconv"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

// 添加任务
func AddTask(ctx context.Context, req *vo.Task) (bool, string, string) {
	logger.Info("ScheduleType:%v", req.ScheduleType)

	if !strings.HasPrefix(req.Integrated, "IOIT") && req.Integrated != "KangYuan" {
		if req.ScheduleType == "period" {
			req.CronExpression = "@every " + strconv.Itoa(req.Every) + "s"
		} else {
			if req.ScheduleType == "cron" {
				req.CronExpression = generateCronExpression(req.Each)
			}
		}
	}

	if strings.HasPrefix(req.Integrated, "IOIT") {
		if req.ScheduleType == "once" {
			req.ScheduleType = "cron"
			req.Each.Type = "once"
		}
	}

	if req.Integrated == "KangYuan" {
		logger.Info("kangyuan url is:%v", req.URL)
		suf := strings.Split(req.URL, ":")[2]
		logger.Info("suf is :%v", suf)
		req.URL = constants.Http + util.ServerCode + constants.DomainName + ":" + suf
	}

	id := GetId()

	logger.Info("cronexp is:%v", req.CronExpression)

	req.ID = id
	req.CreateTime = time.Now().Format("2006-01-02 15:04:05")
	req.EditeTime = time.Now().Format("2006-01-02 15:04:05")
	req.Status = "0"

	if strings.HasPrefix(req.URL, "https://") {
		return false, "作业添加失败，不支持https类型的url", id
	}

	logger.Info("New job add success")

	data, err := GlobalScheduler.AddTask(*req)
	if err != nil {
		return false, data, id
	}
	if req.IsAutoStart {
		err := GlobalScheduler.StartTask(*req)
		if err == nil {
			req.Status = "1"
		}
		if err != nil {
			return false, "任务添加成功，自动启动失败", id
		}
	}
	executor, err := GlobalScheduler.GetExecutorByName(req.Executor)
	if err != nil {
		logger.Error("作业添加，获取执行器信息失败:%v", err)
		return false, "作业添加失败，获取执行器信息失败", id
	}
	tx := GlobalDb.Begin()
	executor.ExistTasks = executor.ExistTasks + 1
	if err = tx.Table("hsm_scheduling.executors").Save(executor).Error; err != nil {
		tx.Rollback()
		logger.Error("作业添加，执行器已存在作业+1失败:%v", err)
		return false, "作业添加失败，执行器增加失败", id
	}
	//GlobalDb.Table("hsm_scheduling.tasks").Create(req)
	if err = tx.Table("hsm_scheduling.tasks").Create(req).Error; err != nil {
		tx.Rollback()
		logger.Error("作业添加失败:%s", err.Error())
		return false, "作业添加失败", id
	}
	tx.Commit()
	return true, "添加任务成功", id
}

// 更新任务
func UpdateTask(ctx context.Context, req *vo.Task) (bool, string) {

	_, oldTask := FindTask(context.Background(), req.ID)

	req.CreateTime = oldTask.CreateTime
	req.EditeTime = time.Now().Format("2006-01-02 15:04:05")
	req.Status = "0"
	logger.Info("each:%v", req.Each)
	logger.Info("req:%v", req.JobParams)
	logger.Info("ScheduleType:%v", req.ScheduleType)
	if !strings.HasPrefix(req.Integrated, "IOIT") && req.Integrated != "KangYuan" {
		if req.ScheduleType == "period" {
			req.CronExpression = "@every " + strconv.Itoa(req.Every) + "s"
		} else {
			if req.ScheduleType == "cron" {
				req.CronExpression = generateCronExpression(req.Each)
			}
		}
	}

	if req.Integrated == "KangYuan" {
		// 查询状态。停止作业
		if oldTask.Status == "1" {
			// 先停止
			GlobalScheduler.StopTask(*oldTask)
		}
	}

	if req.Integrated == "KangYuan" {
		logger.Info("kangyuan url is:%v", req.URL)
		suf := strings.Split(req.URL, ":")[2]
		logger.Info("suf is :%v", suf)
		req.URL = constants.Http + util.ServerCode + constants.DomainName + ":" + suf
	}

	err := GlobalScheduler.UpdateTask(oldTask.Executor, *req)
	if err != nil {
		return false, "执行器作业数量达到上限"
	}
	if req.IsAutoStart {
		err := GlobalScheduler.StartTask(*req)
		if err == nil {
			req.Status = "1"
		}
		if err != nil {
			return false, "任务添加成功，自动启动失败"
		}
	}

	executor, err := GlobalScheduler.GetExecutorByName(req.Executor)
	if err != nil {
		logger.Error("作业更新，获取执行器信息失败:%v", err)
		return false, "作业更新失败，获取执行器信息失败"
	}

	tx := GlobalDb.Begin()

	if oldTask.Executor != req.Executor {
		executor.ExistTasks = executor.ExistTasks + 1
		if err = tx.Table("hsm_scheduling.executors").Save(executor).Error; err != nil {
			tx.Rollback()
			logger.Error("作业更新，执行器已存在作业+1失败:%v", err)
			return false, "作业更新失败，执行器增加失败"
		}

		oldexecutor, err := GlobalScheduler.GetExecutorByName(oldTask.Executor)
		if err != nil {
			logger.Error("作业更新，获取执行器信息失败:%v", err)
		}

		logger.Info("oldexecutor is :%v", oldexecutor)
		oldexecutor.ExistTasks = oldexecutor.ExistTasks - 1

		err = tx.Table("hsm_scheduling.executors").Where("id = ?", oldexecutor.ID).Update("exist_tasks", oldexecutor.ExistTasks).Error
		if err != nil {
			logger.Error("作业更新，旧的执行器信息失败:%v", err)
		}
	}

	if err = tx.Table("hsm_scheduling.tasks").Model(&vo.Task{}).Where("id = ?", req.ID).Save(req).Error; err != nil {
		tx.Rollback()
		logger.Info("update err is :%v", err)
		return false, "更新失败"
	}
	if req.JobName != oldTask.JobName {
		if err = tx.Table("hsm_scheduling.running_logs").Model(&vo.RuningLog{}).Where("taskid = ?", req.ID).Update("jobname", req.JobName).Error; err != nil {
			tx.Rollback()
			logger.Info("update err is :%v", err)
			return false, "更新失败"
		}
		if err = tx.Table("hsm_scheduling.sys_logs").Model(&vo.SysLog{}).Where("taskid = ?", req.ID).Update("jobname", req.JobName).Error; err != nil {
			tx.Rollback()
			logger.Info("update err is :%v", err)
			return false, "更新失败"
		}
	}
	tx.Commit()

	logger.Info("Jobinfo update for success")
	return true, "任务更新成功"
}

// 删除任务
func DeleteTask(ctx context.Context, req *vo.GetID) (bool, string) {

	var count int64
	err := GlobalDb.Table("hsm_scheduling.running_logs").Model(&vo.RuningLog{}).Where("taskid = ?", req.ID).Count(&count).Error
	if err != nil {
		logger.Info("check name error:%v", err)
	}

	if count > 0 {
		return false, "执行列表中存在作业实例，请删除实例后再删除作业"
	}
	_, data := FindTask(context.Background(), req.ID)
	deletesql := `DELETE FROM hsm_scheduling.tasks WHERE id=$1`
	db, err := GetPgConn()
	if err != nil {
		logger.Info("Delete Fail to connect pg:%v", err)
		return false, "删除失败"
	}
	defer db.Close()

	executor, err := GlobalScheduler.GetExecutorByName(data.Executor)
	if err != nil {
		logger.Error("获取执行器信息失败:%v", err)
		return false, "获取执行器信息失败"
	}
	tx, err := db.Begin()
	if err != nil {
		return false, "新建事物失败"
	}
	executor.ExistTasks = executor.ExistTasks - 1
	executorsql := `update hsm_scheduling.executors set exist_tasks =$1 WHERE id=$2`
	executorStmt, err := tx.PrepareContext(ctx, executorsql)
	if err != nil {
		tx.Rollback()
		logger.Error("更新执行器信息失败%v", err)
		return false, "删除失败"
	}
	defer executorStmt.Close()
	_, err = executorStmt.ExecContext(ctx, executor.ExistTasks, executor.ID)
	if err != nil {
		tx.Rollback()
		logger.Error("更新执行器信息失败%v", err)
		return false, "删除失败"
	}

	stmt, err := tx.PrepareContext(ctx, deletesql)
	if err != nil {
		tx.Rollback()
		logger.Error("Delete Fail to prepare statement:%v", err)
		return false, "删除失败"
	}
	defer stmt.Close()
	logger.Info("删除的ID为:%v", req.ID)
	_, err = stmt.ExecContext(ctx, req.ID)
	if err != nil {
		tx.Rollback()
		logger.Error("Delete Fail to exec statement:%v", err)
		return false, "删除失败"
	}
	tx.Commit()
	logger.Info("Delete Delete job code:%v worker success", req.ID)

	GlobalScheduler.DeleteTask(*data)
	return true, "删除成功"
}

// 查询任务
func FindTask(ctx context.Context, id string) (bool, *vo.Task) {
	var task vo.Task
	if err := GlobalDb.Table("hsm_scheduling.tasks").First(&task, id).Error; err != nil {
		return false, nil
	}
	return true, &task
}

// 查询任务
func FindInstanceLogById(ctx context.Context, id string) (bool, []vo.SysLog) {
	var log []vo.SysLog
	if err := GlobalDb.Table("hsm_scheduling.sys_logs").Where("instanceid = ?", id).Order("logtime DESC").Find(&log).Error; err != nil {
		return false, nil
	}
	return true, log
}

type FindAllResp struct {
	ID             string `json:"id"`
	JobName        string `json:"jobname"`
	CreateUser     string `json:"createuser"`     //创建人
	Executor       string `json:"executor"`       //执行器
	Status         string `json:"status"`         //任务状态
	ScheduleType   string `json:"scheduletype"`   //调度类型
	CronExpression string `json:"cronexpression"` //CRON表达
	JobType        string `json:"jobtype"`        //任务类型
	Integrated     string `json:"integrated"`     //集成标识
	IsAutoStart    bool   `json:"isautostart"`    //集成标识
}

// 获取所有任务
func FindAllTask(ctx context.Context) (bool, *[]FindAllResp) {
	query := `SELECT 
				id,
				jobname,
				createuser,
				executor,
				status,
				scheduletype,
				cronexpression,
				jobtype,
				integrated,
				isautostart
	 		FROM hsm_scheduling.tasks ORDER BY createtime DESC`

	// 获取数据库连接
	db, err := GetPgConn()
	if err != nil {
		logger.Info("FindAll Fail to connect pg:%v", err)
		return false, nil
	}
	defer db.Close()

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		logger.Info("FindAll Fail to query:%v", err)
		return false, nil
	}
	defer rows.Close()

	var jobs []FindAllResp
	for rows.Next() {
		var job FindAllResp
		err := rows.Scan(&job.ID, &job.JobName, &job.CreateUser, &job.Executor,
			&job.Status, &job.ScheduleType, &job.CronExpression, &job.JobType, &job.Integrated, &job.IsAutoStart)
		if err != nil {
			logger.Info("FindAll Fail to scan:%v", err)
			return false, nil
		}
		jobs = append(jobs, job)
	}

	if err = rows.Err(); err != nil {
		logger.Info("FindAll Rows iteration failed:%v", err)
		return false, nil
	}
	logger.Info("FindAll Find all task success")
	return true, &jobs
}

func RunTaskForRedency(idVo *vo.GetID) {
	logger.Info("idVo:%v", idVo)
	_, task := FindTask(context.Background(), idVo.ID)
	err := GlobalScheduler.StartTask(*task)

	logger.Info("err is:%v", err)
}

func UpsertFailOverNode(failover *vo.FailOverNode) {
	GlobalDb.AutoMigrate(&vo.FailOverNode{})
	if err := GlobalDb.Table("hsm_scheduling.failover_nodes").Save(failover).Error; err != nil {
		logger.Info("update failover table failed:%v", err)
	}
}

// 自动注册
func AutoRegisterJob(ctx context.Context, task *vo.Task) bool {
	flag, _, _ := AddTask(context.Background(), task)
	return flag
}

func GetRuningLogById(id string) (vo.RuningLog, error) {
	var runlog vo.RuningLog
	if err := GlobalDb.Table("hsm_scheduling.running_logs").Where("id = ?", id).Find(&runlog).Error; err != nil {
		return runlog, err
	}
	return runlog, nil
}

func DeleteInstanceById(id string) bool {
	err := GlobalDb.Exec(`DELETE FROM hsm_scheduling.running_logs WHERE id=$1`, id).Error
	if err != nil {
		logger.Info("delete failed:%v", err)
		return false
	}
	return true
}

func UpdateInstanceIdByTaskId(id string, instanceId string) {
	err := GlobalDb.Table("hsm_scheduling.tasks").Where("id = ?", id).Update("lastinstanceid", instanceId).Error
	if err != nil {
		logger.Info("update lastinstanceid failed:%v", err)
	}
}

func UpdateInstanceEndTimeByInstanceId(id string) {
	logger.Info(id)
	err := GlobalDb.Table("hsm_scheduling.running_logs").Where("id = ?", id).Update("endtime", time.Now().Format("2006-01-02 15:04:05")).Error
	if err != nil {
		logger.Info("update lastinstanceid failed:%v", err)
	}
}

func UpdateInstanceEndTimeIfEmpty(id string) {
	var log vo.RuningLog
	res := GlobalDb.Where(`id = ?`, id).First(&log)
	if res.Error != nil {
		logger.Info("update lastinstanceid failed:%v", res.Error)
	}
	if log.EndTime == "" {
		err := GlobalDb.Table("hsm_scheduling.running_logs").Where(`id = ?`, id).Update("endtime", time.Now().Format("2006-01-02 15:04:05")).Error
		if err != nil {
			logger.Info("update lastinstanceid failed:%v", err)
		}
	}
}
