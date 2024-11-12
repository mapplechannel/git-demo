package jobschedule

import (
	"errors"
	"hsm-scheduling-back-end/config"
	"hsm-scheduling-back-end/internal/util"
	"hsm-scheduling-back-end/internal/vo"
	"hsm-scheduling-back-end/pkg/logger"
	"hsm-scheduling-back-end/response"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Scheduler struct {
	executors map[string]*vo.Executor
	tasks     map[string]map[string]vo.Task
	db        *gorm.DB
	mu        sync.Mutex
}

var GlobalScheduler *Scheduler

func InitExecutor() {
	GetPgConnByGORM()
	scheduler, err := NewScheduler()
	if err != nil {
		logger.Info("init executor failed")
	}
	GlobalScheduler = scheduler
	scheduler.LoadState()
}

func NewScheduler() (*Scheduler, error) {
	db, err := gorm.Open(postgres.Open(config.ConfigAll.Postgres.Url+" dbname="+config.ConfigAll.Postgres.DbName), &gorm.Config{})
	if err != nil {
		logger.Info("connect pg failed:%v", err)
		return nil, err
	}

	db.AutoMigrate(&vo.Executor{})

	scheduler := &Scheduler{
		executors: make(map[string]*vo.Executor),
		tasks:     make(map[string]map[string]vo.Task),
		db:        db,
	}

	scheduler.ensureDefaultExecutor()
	return scheduler, nil
}

func (s *Scheduler) ensureDefaultExecutor() {
	var defaultExecutor vo.Executor
	err := s.db.Table("hsm_scheduling.executors").Where("name = ?", "默认执行器").First(&defaultExecutor).Error
	logger.Info("defaultExecutor:%v", defaultExecutor)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			s.AddExecutor(GetId(), "默认执行器", "系统默认执行器", "admin", 199, true)
		} else {
			logger.Info("failed to query default executor:%v", err)
		}
	} else {
		defaultExecutor.Cron = cron.New(cron.WithSeconds())
		s.executors[defaultExecutor.ID] = &defaultExecutor
		s.tasks[defaultExecutor.ID] = make(map[string]vo.Task)
		go defaultExecutor.Cron.Start()
	}
}

func (s *Scheduler) AddExecutor(id, name, desc, createUser string, maxTasks int, isDefault bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	flag := CheckExecutorName(name)
	if flag {
		return errors.New("存在同名的执行器")
	}
	executor := &vo.Executor{
		ID:         id,
		Name:       name,
		Desc:       desc,
		MaxTasks:   maxTasks,
		CreateUser: createUser,
		IsDefault:  isDefault,
		EditeTime:  time.Now().Format("2006-01-02 15:04:05"),
		Cron:       cron.New(cron.WithSeconds()),
	}
	s.executors[id] = executor
	s.tasks[id] = make(map[string]vo.Task)
	s.db.Table("hsm_scheduling.executors").Create(executor)
	go executor.Cron.Start()
	return nil
}

func CheckExecutorName(name string) bool {
	var count int64
	err := GlobalDb.Table("hsm_scheduling.executors").Model(&vo.Executor{}).Where("name = ?", name).Count(&count).Error
	if err != nil {
		logger.Info("check name error:%v", err)
	}

	return count > 0
}

func (s *Scheduler) UpdateExecutor(id, name, desc, createUser string, maxTasks int, isDefault bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 查询旧的执行器
	oldExecutorName := GetNameById(id)

	executor, ok := s.executors[id]
	if !ok {
		logger.Info("executor with ID does not exist node memory")
		var err error
		executor, err = s.GetExecutorById(id)
		if err != nil {
			return err
		}
	}
	executor.Name = name
	executor.Desc = desc
	executor.MaxTasks = maxTasks
	executor.CreateUser = createUser
	executor.IsDefault = isDefault
	executor.EditeTime = time.Now().Format("2006-01-02 15:04:05")
	tx := s.db.Begin()
	if err := tx.Table("hsm_scheduling.executors").Save(executor).Error; err != nil {
		tx.Rollback()
		logger.Error("update executor failed:%s", err.Error())
		return err
	}
	//执行器没有做id关联，临时处理使用名称更新
	if oldExecutorName != name {
		if err := tx.Table("hsm_scheduling.tasks").Model(&vo.Task{}).Where("executor = ?", oldExecutorName).Update("executor", name).Error; err != nil {
			tx.Rollback()
			logger.Error("update task executor name failed:%s", err.Error())
			return err
		}
		if err := tx.Table("hsm_scheduling.running_logs").Model(&vo.RuningLog{}).Where("executor = ?", oldExecutorName).Update("executor", name).Error; err != nil {
			tx.Rollback()
			logger.Error("update runingLog executor name failed:%s", err.Error())
			return err
		}
		if err := tx.Table("hsm_scheduling.sys_logs").Model(&vo.SysLog{}).Where("executor = ?", oldExecutorName).Update("executor", name).Error; err != nil {
			tx.Rollback()
			logger.Error("update syslog executor name failed:%s", err.Error())
			return err
		}
	}
	tx.Commit()

	// 修改对应作业的执行器名称
	//err := s.db.Table("hsm_scheduling.tasks").Model(&vo.Task{}).Where("executor = ?", oldExecutorName).Update("executor", name).Error
	//if err != nil {
	//	logger.Info("task with executor does not update:%v", err)
	//}
	return nil
}

func (s *Scheduler) GetExecutorById(id string) (*vo.Executor, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var executor vo.Executor
	if err := s.db.Table("hsm_scheduling.executors").First(&executor, id).Error; err != nil {
		logger.Info("find executor failed:%v", err)
		return nil, err
	}
	parseTIme, _ := time.Parse(time.RFC3339, executor.EditeTime)
	executor.EditeTime = parseTIme.Format("2006-01-02 15:04:05")
	return &executor, nil
}

func (s *Scheduler) GetExecutorByName(name string) (*vo.Executor, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var executor vo.Executor
	if err := s.db.Table("hsm_scheduling.executors").Where("name = ?", name).Find(&executor).Error; err != nil {
		logger.Info("find executor failed:%v", err)
		return nil, err
	}
	return &executor, nil
}

func (s *Scheduler) GetAllExecutors() ([]vo.Executor, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var executors []vo.Executor
	if err := s.db.Table("hsm_scheduling.executors").Order("editetime desc").Find(&executors).Error; err != nil {
		logger.Info("find executors failed:%v", err)
		return nil, err
	}

	for i := range executors {
		parseTIme, _ := time.Parse(time.RFC3339, executors[i].EditeTime)
		executors[i].EditeTime = parseTIme.Format("2006-01-02 15:04:05")
	}

	var res []vo.Executor
	for i := range executors {
		if executors[i].Name != "信息系统集成" {
			res = append(res, executors[i])
		}
	}

	return res, nil
}

func (s *Scheduler) RemoveExecutor(id string) (error, int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	executor, ok := s.executors[id]
	if !ok {
		logger.Info("executor with ID does not exist")
	}

	if len(s.tasks[id]) > 0 {
		logger.Info("executor has task")
		return errors.New("执行器下有作业，无法删除"), response.EXIST_JOB
	}

	executor.Cron.Stop()
	delete(s.tasks, id)
	delete(s.executors, id)
	s.db.Table("hsm_scheduling.executors").Delete(&vo.Executor{}, id)
	return nil, 0
}

func GetIdByName(name string) string {
	var executor vo.Executor

	if err := GlobalDb.Table("hsm_scheduling.executors").Where("name = ?", name).First(&executor).Error; err != nil {
		logger.Info("failed to select:%v", err)
	}
	return executor.ID
}

func GetNameById(id string) string {
	var executor vo.Executor
	logger.Info("id:%v", id)

	if err := GlobalDb.Table("hsm_scheduling.executors").Where("id = ?", id).First(&executor).Error; err != nil {
		logger.Info("failed to select:%v", err)
	}
	return executor.Name
}

func (s *Scheduler) AddTask(task vo.Task) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	executorId := GetIdByName(task.Executor)
	logger.Info("executorId:%v", executorId)
	executor, ok := s.executors[executorId]
	if !ok {
		logger.Info("executor with ID does not exist")
	}

	if len(s.tasks[executorId]) >= executor.MaxTasks {
		logger.Info("executor'task has reached full")
		return "执行器作业数量达到上限", errors.New("执行器作业数量达到上限")
	}

	s.tasks[executorId][task.ID] = task
	return task.ID, nil
}

func (s *Scheduler) UpdateTask(oldExecutor string, newTask vo.Task) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	// 新的执行器ID
	logger.Info("newTask.Executor is:%v", newTask.Executor)

	newExecutorId := GetIdByName(newTask.Executor)
	newExecutor, ok := s.executors[newExecutorId]
	logger.Info("new executor is :%v", newExecutor.Name)
	if !ok {
		logger.Info("executor with ID does not exist")
	}
	logger.Info("newExecutorId:%v", newExecutorId)

	// 获取旧任务的执行器
	oldExecutorId := GetIdByName(oldExecutor)
	// oldExecutor, ok := s.executors[executorId]
	if !ok {
		logger.Info("executor with ID does not exist")
	}
	// 执行器相同
	if newExecutorId == oldExecutorId {
		tasks := s.tasks[newExecutorId]
		tasks[newTask.ID] = newTask
		s.tasks[newExecutorId] = tasks
	} else {
		// 执行器不同
		if len(s.tasks[newExecutorId]) >= newExecutor.MaxTasks {
			logger.Info("executor'task has reached full")
			return errors.New("执行器作业数量达到上限")
		}

		tasks := s.tasks[oldExecutorId]
		delete(tasks, newTask.ID)

		s.tasks[newExecutorId][newTask.ID] = newTask
	}
	// tasks[newTask.ID] = newTask
	// s.tasks[executorId] = tasks
	return nil
}

func (s *Scheduler) DeleteTask(req vo.Task) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	executorId := GetIdByName(req.Executor)
	logger.Info("executorId:%v", executorId)
	tasks := s.tasks[executorId]
	task, _ := tasks[req.ID]
	delete(tasks, task.ID)
	s.tasks[executorId] = tasks
	return nil
}

func (s *Scheduler) StartTask(task vo.Task) error {
	// GlobalScheduler.LoadState()
	failover := vo.FailOverNode{
		TaskId:      task.ID,
		RunningNode: util.ServerCode,
	}
	UpsertFailOverNode(&failover)
	s.mu.Lock()
	defer s.mu.Unlock()
	executorId := GetIdByName(task.Executor)
	executor, ok := s.executors[executorId]
	if !ok {
		logger.Info("executor with ID does not exist")
	}

	taskOfExe, ok := s.tasks[executorId][task.ID]
	if !ok {
		logger.Info("task with ID does not exist")
	}

	if strings.HasPrefix(task.Integrated, "IOIT") {
		uniqueSlice := strings.Split(task.Integrated, "@")
		namespace := uniqueSlice[1]
		code := uniqueSlice[2]
		var ioitJob vo.JobInfoIoit
		if err := GlobalDb.Table("ioit.job_info").Where("code = ? AND namespace = ?", code, namespace).First(&ioitJob).Error; err != nil {
			logger.Info("select ioitinfo err :%v", err)
		}
		logger.Info("ioitJob.CronDetail is :%v", ioitJob.CronDetail)
		logger.Info("ioitJob.PeriodDetail is :%v", ioitJob.PeriodDetail)

		var executedTimes int
		entryId, err := executor.Cron.AddFunc(task.CronExpression, func() {
			if task.ScheduleType == "cron" {
				effectiveDateArr := strings.Split(ioitJob.CronDetail.EffectiveDate, " ")
				startDate, err := time.Parse("2006-01-02", effectiveDateArr[0])
				if err != nil {
					logger.Info("parse effectiveDate error:%v", err)
				}

				now := time.Now()
				if now.After(startDate) {
					if ioitJob.CronDetail.EndType == "endless" {
						RunTask(&task, nil)
					} else {
						endDate, err := time.Parse("2006-01-02 15:04:05", ioitJob.CronDetail.EndDate)
						if err != nil {
							logger.Info("parse effectiveDate error:%v", err)
						}
						tarTime := time.Date(endDate.Year(), endDate.Month(), endDate.Day(), endDate.Hour(), endDate.Minute(), endDate.Second(), 0, time.Local)

						if now.Before(tarTime) {
							RunTask(&task, nil)
						} else {
							logger.Info("stop schedule")
							// 停止作业
							SendStopToIOit(task)
							return
						}
					}
				}
			} else {
				if task.ScheduleType == "period" {

					effectiveDateArr, err := time.Parse("2006-01-02 15:04:05", ioitJob.CronDetail.EffectiveDate)
					if err != nil {
						logger.Info("parse effectiveDate error:%v", err)
					}
					tarTime := time.Date(effectiveDateArr.Year(), effectiveDateArr.Month(), effectiveDateArr.Day(), effectiveDateArr.Hour(), effectiveDateArr.Minute(), effectiveDateArr.Second(), 0, time.Local)

					if err != nil {
						logger.Info("parse effectiveDate error:%v", err)
					}
					now := time.Now()
					if now.After(tarTime) {
						if ioitJob.PeriodDetail.EndType == "times" {
							if executedTimes < ioitJob.PeriodDetail.Times {
								RunTask(&task, nil)
								executedTimes++
							} else {
								logger.Info("停止调度，超过次数")
								// 停止作业
								SendStopToIOit(task)
								return
							}
						} else {
							if ioitJob.PeriodDetail.EndType == "endless" {
								RunTask(&task, nil)
							} else {
								endDate, err := time.Parse("2006-01-02 15:04:05", ioitJob.PeriodDetail.EndDate)
								if err != nil {
									logger.Info("parse effectiveDate error:%v", err)
								}
								tarTime := time.Date(endDate.Year(), endDate.Month(), endDate.Day(), endDate.Hour(), endDate.Minute(), endDate.Second(), 0, time.Local)

								if now.Before(tarTime) {
									RunTask(&task, nil)
								} else {
									logger.Info("stop schedule")
									// 停止作业
									SendStopToIOit(task)
									return
								}
							}
						}
					}
				}
			}
		})
		if err != nil {
			logger.Info("failed to start")
		}
		taskOfExe.EntryId = entryId
		s.tasks[executorId][task.ID] = taskOfExe
	} else {
		if task.Integrated == "" {
			// 生成ID
			code := GenerateCode()
			runningLog := vo.RuningLog{
				ID:        code,
				Operator:  task.CreateUser,
				JobName:   task.JobName,
				Executor:  task.Executor,
				TaskId:    task.ID,
				EndTime:   "",
				StartTime: time.Now().Format("2006-01-02 15:04:05"),
			}
			UpdateInstanceEndTimeIfEmpty(task.LastInstanceId)
			UpdateInstanceIdByTaskId(task.ID, code)
			logger.Info("run id is :%v", task.ID)
			entryId, err := executor.Cron.AddFunc(task.CronExpression, func() {
				if task.ScheduleType == "cron" {
					if task.Each.DayOfMonth == "0" {
						now := time.Now()
						if !isLastDayOfMonth(now) {
							logger.Info("not last day of month")
							return
						}
					}
					parseDate, err := time.Parse("2006-01-02", task.Each.EffectiveDate)
					if err != nil {
						logger.Info("parse effectiveDate error:%v", err)
					}
					now := time.Now()
					if now.Before(parseDate) {
						logger.Info("skipping execution, not yet effective")
						return
					}
				}
				RunTask(&task, &runningLog)
			})

			if err != nil {
				logger.Info("failed to start:%v", err)
			}

			// err1 := GlobalDb.Table("hsm_scheduling.tasks").Where("id = ?", task.ID).Update("enid", int(entryId)).Error
			// if err1 != nil {
			// 	logger.Info("err is :%v", err1)
			// }
			taskOfExe.EntryId = entryId
			logger.Info("taskOfExe.EntryId:%v", taskOfExe)
			s.tasks[executorId][task.ID] = taskOfExe
		} else {
			if task.Integrated == "KangYuan" {
				logger.Info("run id is :%v", task.ID)
				entryId, err := executor.Cron.AddFunc(task.CronExpression, func() {
					RunTask(&task, nil)
				})

				if err != nil {
					logger.Info("failed to start: %v", err)
				}

				taskOfExe.EntryId = entryId
				s.tasks[executorId][task.ID] = taskOfExe
			}
		}
	}
	return nil
}

func isLastDayOfMonth(t time.Time) bool {
	nextDay := t.AddDate(0, 0, 1)
	return nextDay.Day() == 1
}

func GenerateCode() string {
	rand.Seed(time.Now().UnixNano())

	sec := time.Now().Unix()
	prefix := strconv.FormatInt(sec, 10)
	const letters = "abcdefghijklmnopqrstuvwxyz"
	code := prefix + "-"
	code += string(letters[rand.Intn(26)]) // 第二位为小写字母
	for i := 0; i < 7; i++ {
		code += string(letters[rand.Intn(len(letters))])
	}
	return code
}

func (s *Scheduler) StopTask(task vo.Task) error {
	// GlobalScheduler.LoadState()
	s.mu.Lock()
	defer s.mu.Unlock()
	if !strings.HasPrefix(task.Integrated, "IOIT") {
		db, err := GetPgConn()

		if err != nil {
			logger.Info("UpdateStatus Update Fail to connect pg:%v", err)
		}
		defer db.Close()

		queryRunningNode := `SELECT runningnode FROM hsm_scheduling.failover_nodes WHERE taskid = $1`
		var runningNOde string

		err = db.QueryRow(queryRunningNode, task.ID).Scan(&runningNOde)

		if err != nil {
			logger.Info("查询作业runningNode失败%v", err)
		}
		logger.Info("runningNOde:%v", runningNOde)

		if runningNOde != util.ServerCode {
			return SendToAnotherSysStop(task, runningNOde)
		}
	}

	executorId := GetIdByName(task.Executor)
	executor, ok := s.executors[executorId]
	if !ok {
		logger.Info("executor with ID does not exist")
	}

	taskOfExe, ok := s.tasks[executorId][task.ID]
	if !ok {
		logger.Info("task with ID does not exist")
	}

	logger.Info("taskOfExe.EntryId:%v", taskOfExe)
	executor.Cron.Remove(taskOfExe.EntryId)
	_ = UpdateStatus(task.ID, "0")
	if !strings.HasPrefix(task.Integrated, "IOIT") {
		logger.Info("endtime***")
		UpdateInstanceEndTimeByInstanceId(task.LastInstanceId)
	}
	// GlobalDb.Table("hsm_scheduling.running_logs").Model(&vo.RuningLog{}).Where("id = ?", req.ID).Updates(req)
	return nil
}

func (s *Scheduler) PrintExecutorsAndTasks() {
	s.mu.Lock()
	defer s.mu.Unlock()

	resMap := make(map[string]int)

	for _, executor := range s.executors {
		logger.Info("executor name is:%v", executor.Name)
		tasks, ok := s.tasks[executor.ID]
		if !ok || len(tasks) == 0 {
			logger.Info("no ok or no tasks")
		} else {
			for _, task := range tasks {
				logger.Info("task name is:%v", task.JobName)
			}
		}
		resMap[executor.ID] = len(tasks)
	}

	for executorId, existTask := range resMap {
		err := GlobalDb.Table("hsm_scheduling.executors").Model(&vo.Executor{}).Where("id = ?", executorId).Update("exist_tasks", existTask).Error
		if err != nil {
			logger.Info("failed to update existask")
		}
	}

}

func (s *Scheduler) LoadState() error {
	var executors []vo.Executor
	var tasks []vo.Task

	logger.Info("loadstat")

	s.db.Table("hsm_scheduling.executors").Find(&executors)
	s.db.Table("hsm_scheduling.tasks").Find(&tasks)

	for _, exec := range executors {
		executor := exec
		executor.Cron = cron.New(cron.WithSeconds())
		s.executors[executor.ID] = &executor
		s.tasks[executor.ID] = make(map[string]vo.Task)
		go executor.Cron.Start()
	}
	for _, task := range tasks {
		executorId := GetIdByName(task.Executor)
		s.tasks[executorId][task.ID] = task
		// s.StartTask(task)
	}

	return nil
}
