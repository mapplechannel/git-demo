package jobschedule

import (
	"context"
	"hsm-scheduling-back-end/internal/constants"
	"hsm-scheduling-back-end/internal/util"
	"hsm-scheduling-back-end/internal/vo"
	"hsm-scheduling-back-end/pkg/logger"
	"hsm-scheduling-back-end/response"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// 服务注册接口******************
func RegistryWorkerHost(ctx *gin.Context) {
	var registryHost RegistryHost
	if err := ctx.ShouldBind(&registryHost); err != nil {
		response.Api_Code_Fail(ctx, response.REQUEST_PARAM_INVALID, time.Now())
		return
	}
	flag := RegistryNewHost(context.Background(), &registryHost)
	if flag {
		response.Api_Success(ctx, flag, "registry success", time.Now())
		return
	}
	// response.Api_Code_Fail(ctx, code, time.Now())
}

func UpdateWorkerHeartbeat(ctx *gin.Context) {
	var registryHost RegistryHost
	if err := ctx.ShouldBind(&registryHost); err != nil {
		response.Api_Code_Fail(ctx, response.REQUEST_PARAM_INVALID, time.Now())
		return
	}
	flag := UpdateHostHeartbeat(context.Background(), registryHost.Addr)
	if flag {
		response.Api_Success(ctx, flag, "heartbeat update success", time.Now())
		return
	}
	// // response.Api_Fail(ctx, msg, time.Now())
	// response.Api_Code_Fail(ctx, code, time.Now())
}

// 作业操作接口******************
func Add(ctx *gin.Context) {
	GlobalScheduler.LoadState()

	var jobRequest vo.Task
	if err := ctx.ShouldBind(&jobRequest); err != nil {
		logger.Info("err is :%v", err)
		response.Api_Code_Fail(ctx, response.REQUEST_PARAM_INVALID, time.Now())
		return
	}
	//jobRequest.JobType = strings.ToUpper(jobRequest.JobType)
	flag, msg, id := AddTask(context.Background(), &jobRequest)
	if flag {
		response.Api_Success(ctx, id, msg, time.Now())
		return
	} else {
		response.Api_Fail(ctx, msg, time.Now())
	}
}

func Update(ctx *gin.Context) {
	var updateJob vo.Task
	if err := ctx.ShouldBind(&updateJob); err != nil {
		response.Api_Code_Fail(ctx, response.REQUEST_PARAM_INVALID, time.Now())
		return
	}
	flag, msg := UpdateTask(context.Background(), &updateJob)
	if flag {
		response.Api_Success(ctx, flag, "update task success", time.Now())
		return
	} else {
		response.Api_Fail(ctx, msg, time.Now())
	}
}

func Delete(ctx *gin.Context) {
	GlobalScheduler.LoadState()
	defer func() {
		if msg := recover(); msg != nil {
			logger.Error("删除作业发生panic:%v", msg)
			logger.Error(string(debug.Stack()))
			response.Api_Code_Fail(ctx, response.DELETE_FAILED, time.Now())
		}
	}()
	var id vo.GetID
	if err := ctx.ShouldBind(&id); err != nil {
		response.Api_Code_Fail(ctx, response.REQUEST_PARAM_INVALID, time.Now())
		return
	}
	flag, msg := DeleteTask(context.Background(), &id)
	if flag {
		response.Api_Success(ctx, flag, "Delete task success", time.Now())
		return
	} else {
		response.Api_Fail(ctx, msg, time.Now())
	}
}

func Find(ctx *gin.Context) {
	var id vo.GetID
	if err := ctx.ShouldBind(&id); err != nil {
		response.Api_Code_Fail(ctx, response.REQUEST_PARAM_INVALID, time.Now())
		return
	}
	flag, data := FindTask(context.Background(), id.ID)
	if flag {
		response.Api_Success(ctx, data, "Find task success", time.Now())
		return
	} else {
		response.Api_Code_Fail(ctx, response.DATABASE_READ_FAILED, time.Now())
	}
}

func FindAll(ctx *gin.Context) {
	// GlobalScheduler.LoadState()
	flag, data := FindAllTask(context.Background())
	if flag {
		response.Api_Success(ctx, data, "Find task success", time.Now())
		return
	} else {
		response.Api_Code_Fail(ctx, response.DATABASE_READ_FAILED, time.Now())
	}
}

func FindLog(ctx *gin.Context) {
	data, err := GetAllSysLog()
	if err == nil {
		response.Api_Success(ctx, data, "Find log success", time.Now())
		return
	} else {
		response.Api_Code_Fail(ctx, response.DATABASE_READ_FAILED, time.Now())
	}
}

func FindRunningLog(ctx *gin.Context) {
	data, err := GetAllRuningLog()
	if err == nil {
		response.Api_Success(ctx, data, "Find running log success", time.Now())
		return
	} else {
		response.Api_Code_Fail(ctx, response.DATABASE_READ_FAILED, time.Now())
	}
}

func FindInstanceLog(ctx *gin.Context) {
	var instance vo.RuningLog
	if err := ctx.ShouldBind(&instance); err != nil {
		response.Api_Code_Fail(ctx, response.REQUEST_PARAM_INVALID, time.Now())
		return
	}
	logger.Info("instance is :%v", instance)
	flag, data := FindInstanceLogById(context.Background(), instance.ID)
	if flag {
		response.Api_Success(ctx, data, "Find task success", time.Now())
		return
	} else {
		response.Api_Code_Fail(ctx, response.DATABASE_READ_FAILED, time.Now())
	}
}

func DeleteInstance(ctx *gin.Context) {
	var instance vo.RuningLog
	if err := ctx.ShouldBind(&instance); err != nil {
		response.Api_Code_Fail(ctx, response.REQUEST_PARAM_INVALID, time.Now())
		return
	}
	logger.Info("instance is :%v", instance)
	runingLog, err := GetRuningLogById(instance.ID)
	if err != nil {
		logger.Error("查询 runingLog 失败：%s", err.Error())
		response.Api_Code_Fail(ctx, response.DATABASE_READ_FAILED, time.Now())
		return
	}
	exist, taskf := FindTask(context.Background(), runingLog.TaskId)
	if exist && taskf.Status == "1" && taskf.LastInstanceId == instance.ID {
		response.Api_Code_Fail(ctx, response.RUNNING_JOB_NOT_OP, time.Now())
		return
	}

	//删除实例就认为停止作业
	flag := DeleteInstanceById(instance.ID)
	if flag {
		response.Api_Success(ctx, flag, "delete instance success", time.Now())
		return
	} else {
		response.Api_Code_Fail(ctx, response.DATABASE_READ_FAILED, time.Now())
	}
}

// 自动注册
func AutoAddJob(ctx *gin.Context) {
	GlobalScheduler.LoadState()
	var task vo.Task
	if err := ctx.ShouldBind(&task); err != nil {
		response.Api_Code_Fail(ctx, response.REQUEST_PARAM_INVALID, time.Now())
		return
	}
	flag := AutoRegisterJob(context.Background(), &task)
	if flag {
		response.Api_Success(ctx, flag, "Auto add success", time.Now())
		return
	} else {
		response.Api_Code_Fail(ctx, response.OPERATION_FAILED, time.Now())
	}
}

// 执行任务
func Run(ctx *gin.Context) {

	var idVo vo.GetID
	if err := ctx.ShouldBind(&idVo); err != nil {
		response.Api_Code_Fail(ctx, response.REQUEST_PARAM_INVALID, time.Now())
		return
	}
	_, task := FindTask(context.Background(), idVo.ID)
	_ = UpdateStatus(task.ID, "1")

	if strings.HasPrefix(task.Integrated, "IOIT") {
		if task.ScheduleType == "cron" && task.Each.Type == "once" {
			uniqueSlice := strings.Split(task.Integrated, "@")
			namespace := uniqueSlice[1]
			code := uniqueSlice[2]
			var ioitJob vo.JobInfoIoit
			if err := GlobalDb.Table("ioit.job_info").Where("code = ? AND namespace = ?", code, namespace).First(&ioitJob).Error; err != nil {
				logger.Info("select ioitinfo err :%v", err)
			}
			logger.Info("ioitJob.CronDetail is :%v", ioitJob.CronDetail)
			logger.Info("ioitJob.PeriodDetail is :%v", ioitJob.PeriodDetail)
			failover := vo.FailOverNode{
				TaskId:      task.ID,
				RunningNode: util.ServerCode,
			}
			UpsertFailOverNode(&failover)

			startDate, err := time.Parse("2006-01-02 15:04:05", ioitJob.CronDetail.EffectiveDate)
			if err != nil {
				logger.Info("parse effectiveDate error:%v", err)
			}

			tarTime := time.Date(startDate.Year(), startDate.Month(), startDate.Day(), startDate.Hour(), startDate.Minute(), startDate.Second(), 0, time.Local)

			duration := time.Until(tarTime)
			logger.Info("duration:%v", duration)
			if duration <= 0 {
				logger.Info("time is in the past")
			}
			time.AfterFunc(duration, func() { ManualRunTask(task, nil) })
			response.Api_Success(ctx, task.ID, "", time.Now())
			return
		}
	}

	if task.ScheduleType == "none" {
		code := GenerateCode()
		runningLog := vo.RuningLog{
			ID:       code,
			Operator: task.CreateUser,
			JobName:  task.JobName,
			Executor: task.Executor,
			TaskId:   task.ID,
		}
		if !strings.HasPrefix(task.Integrated, "IOIT") {
			UpdateInstanceEndTimeIfEmpty(task.LastInstanceId)
			UpdateInstanceIdByTaskId(task.ID, code)
		}
		failover := vo.FailOverNode{
			TaskId:      task.ID,
			RunningNode: util.ServerCode,
		}
		UpsertFailOverNode(&failover)
		go ManualRunTask(task, &runningLog)
	} else {
		if task.Each.Type == "once" {
			date, err := time.Parse("2006-01-02", task.Each.EffectiveDate)
			if err != nil {
				logger.Info("failed to parse date:%v", err)
			}

			parseTime, err := time.Parse("15:04:05", task.Each.Time)
			if err != nil {
				logger.Info("failed to parse time:%v", err)
			}

			failover := vo.FailOverNode{
				TaskId:      task.ID,
				RunningNode: util.ServerCode,
			}
			UpsertFailOverNode(&failover)

			tarTime := time.Date(date.Year(), date.Month(), date.Day(), parseTime.Hour(), parseTime.Minute(), parseTime.Second(), 0, time.Local)

			duration := time.Until(tarTime)
			if duration <= 0 {
				logger.Info("time is in the past")
			}
			code := GenerateCode()
			runningLog := vo.RuningLog{
				ID:       code,
				Operator: task.CreateUser,
				JobName:  task.JobName,
				Executor: task.Executor,
				TaskId:   task.ID,
			}
			UpdateInstanceEndTimeIfEmpty(task.LastInstanceId)
			UpdateInstanceIdByTaskId(task.ID, code)
			time.AfterFunc(duration, func() { ManualRunTask(task, &runningLog) })
		} else {
			go GlobalScheduler.StartTask(*task)
		}
	}
	// if flag {
	response.Api_Success(ctx, task.ID, "", time.Now())
	// return
	// } else {
	// response.Api_Fail(ctx, msg, time.Now())
	// }
}

func ManualRun(ctx *gin.Context) {
	var idVo vo.GetID
	if err := ctx.ShouldBind(&idVo); err != nil {
		response.Api_Code_Fail(ctx, response.REQUEST_PARAM_INVALID, time.Now())
		return
	}
	_, task := FindTask(context.Background(), idVo.ID)
	code := GenerateCode()
	runningLog := vo.RuningLog{
		ID:       code,
		Operator: task.CreateUser,
		JobName:  task.JobName,
		Executor: task.Executor,
		TaskId:   task.ID,
	}
	_ = UpdateStatus(task.ID,"1")
	ManualRunTask(task, &runningLog)
	// if flag {
	response.Api_Success(ctx, task.ID, "", time.Now())
	// return
	// } else {
	// response.Api_Fail(ctx, msg, time.Now())
	// }
}

// 停止任务
func Stop(ctx *gin.Context) {
	var task vo.Task
	if err := ctx.ShouldBind(&task); err != nil {
		response.Api_Code_Fail(ctx, response.REQUEST_PARAM_INVALID, time.Now())
		return
	}
	_, taskf := FindTask(context.Background(), task.ID)
	err := GlobalScheduler.StopTask(*taskf)
	if err != nil {
		response.Api_Fail(ctx, "作业停止失败", time.Now())
		return
	}
	// if flag {
	response.Api_Success(ctx, task.ID, "作业停止成功", time.Now())
}

// 执行器增删改查开始**************************************
func ExecutorAdd(ctx *gin.Context) {
	var executor vo.Executor
	if err := ctx.ShouldBind(&executor); err != nil {
		response.Api_Code_Fail(ctx, response.REQUEST_PARAM_INVALID, time.Now())
		return
	}
	id := GetId()
	err := GlobalScheduler.AddExecutor(id, executor.Name, executor.Desc, executor.CreateUser, executor.MaxTasks, false)
	if err == nil {
		response.Api_Success(ctx, id, "add success", time.Now())
		return
	}
	// response.Api_Code_Fail(ctx, response.DATABASE_OPERATION_FAILED, time.Now())
	response.Api_Fail(ctx, "存在同名的执行器", time.Now())
}

func ExecutorEdite(ctx *gin.Context) {
	GlobalScheduler.LoadState()

	var executor vo.Executor
	if err := ctx.ShouldBind(&executor); err != nil {
		response.Api_Code_Fail(ctx, response.REQUEST_PARAM_INVALID, time.Now())
		return
	}

	//判断库中的执行器最大作业数是否小于已存在作业数
	if executor.ID == constants.EmptyContent {
		response.Api_Code_Fail(ctx, response.REQUEST_PARAM_INVALID, time.Now())
	}
	executorById, err := GlobalScheduler.GetExecutorById(executor.ID)
	if err != nil {
		response.Api_Code_Fail(ctx, response.DATABASE_OPERATION_FAILED, time.Now())
		return
	}
	if executorById.ExistTasks > executor.MaxTasks {
		response.Api_Code_Fail(ctx, response.EXECUTOR_MAXTASKS, time.Now())
		return
	}

	//同名执行器
	if executor.Name != executorById.Name && CheckExecutorName(executor.Name) {
		response.Api_Code_Fail(ctx, response.EXECUTOR_SAME_NAME, time.Now())
		return
	}

	err = GlobalScheduler.UpdateExecutor(executor.ID, executor.Name, executor.Desc, executor.CreateUser, executor.MaxTasks, false)
	if err == nil {
		response.Api_Success(ctx, executor.ID, "edite success", time.Now())
		return
	}
	response.Api_Code_Fail(ctx, response.DATABASE_OPERATION_FAILED, time.Now())
}

func ExecutorDelete(ctx *gin.Context) {
	GlobalScheduler.LoadState()

	var executor vo.Executor
	if err := ctx.ShouldBind(&executor); err != nil {
		response.Api_Code_Fail(ctx, response.REQUEST_PARAM_INVALID, time.Now())
		return
	}

	err, code := GlobalScheduler.RemoveExecutor(executor.ID)
	if err == nil {
		response.Api_Success(ctx, executor.ID, "delete success", time.Now())
		return
	}
	response.Api_Code_Fail(ctx, code, time.Now())
}

func ExecutorFind(ctx *gin.Context) {
	var executor vo.Executor
	if err := ctx.ShouldBind(&executor); err != nil {
		response.Api_Code_Fail(ctx, response.REQUEST_PARAM_INVALID, time.Now())
		return
	}

	data, err := GlobalScheduler.GetExecutorById(executor.ID)
	if err == nil {
		response.Api_Success(ctx, data, "find success", time.Now())
		return
	}
	response.Api_Code_Fail(ctx, response.DATABASE_OPERATION_FAILED, time.Now())
}

func ExecutorFindAll(ctx *gin.Context) {
	//GlobalScheduler.PrintExecutorsAndTasks()
	GlobalScheduler.LoadState()
	data, err := GlobalScheduler.GetAllExecutors()

	if err == nil {
		response.Api_Success(ctx, data, "find all success", time.Now())
		return
	}
	response.Api_Code_Fail(ctx, response.DATABASE_OPERATION_FAILED, time.Now())
}

func ExecutorPrint(ctx *gin.Context) {
	var executor vo.Executor
	if err := ctx.ShouldBind(&executor); err != nil {
		response.Api_Code_Fail(ctx, response.REQUEST_PARAM_INVALID, time.Now())
		return
	}

	GlobalScheduler.PrintExecutorsAndTasks()
	response.Api_Success(ctx, "res", "find all success", time.Now())
}

// 执行器增删改查结束**************************************

func CheckHealth(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "alive")
}
