package vo

import (
	"database/sql/driver"
	"encoding/json"
	"hsm-scheduling-back-end/pkg/logger"

	"github.com/robfig/cron/v3"
)

type RuningLog struct {
	ID        string `gorm:"type:varchar(255);column:id" json:"id"`
	JobName   string `gorm:"type:varchar(255);column:jobname" json:"jobname"`
	TaskId    string `gorm:"type:varchar(255);column:taskid;index:task" json:"taskid"`
	Status    string `gorm:"type:varchar(255);column:status" json:"status"`
	Executor  string `gorm:"type:varchar(255);column:executor;index:executor" json:"executor"`
	StartTime string `gorm:"type:varchar(255);column:starttime" json:"starttime"`
	EndTime   string `gorm:"type:varchar(255);column:endtime" json:"endtime"`
	Operator  string `gorm:"type:varchar(255);column:operator" json:"operator"`
}

func (RuningLog) TableName() string {
	return "hsm_scheduling.running_logs"
}

type SysLog struct {
	ID         string `gorm:"type:varchar(255);primaryKey" json:"id"`
	JobName    string `gorm:"type:varchar(255);column:jobname" json:"jobname"`
	LogType    string `gorm:"type:varchar(255);column:logtype" json:"logtype"`
	TaskId     string `gorm:"type:varchar(255);column:taskid;index:task" json:"taskid"`
	InstanceId string `gorm:"type:varchar(255);column:instanceid" json:"instanceid"`
	Content    string `gorm:"type:text;column:content" json:"content"`
	Executor   string `gorm:"type:varchar(255);column:executor;index:executor" json:"executor"`
	LogTime    string `gorm:"type:varchar(255);column:logtime" json:"logtime"`
}

func (SysLog) TableName() string {
	return "hsm_scheduling.sys_logs"
}

type FailOverNode struct {
	TaskId      string `gorm:"type:varchar(19);primaryKey;column:taskid"`
	RunningNode string `gorm:"type:varchar(255);column:runningnode"`
}

func (FailOverNode) TableName() string {
	return "hsm_scheduling.failover_nodes"
}

type Task struct {
	ID             string       `gorm:"type:varchar(19);primaryKey" json:"id"`
	JobName        string       `gorm:"type:varchar(255);column:jobname" json:"jobname"`                  //任务名称
	CreateUser     string       `gorm:"type:varchar(255);column:createuser" json:"createuser"`            //创建人
	Executor       string       `gorm:"type:varchar(255);column:executor;index:executor" json:"executor"` //执行器
	Describe       string       `gorm:"type:text;column:describe" json:"describe"`                        //描述
	JobClass       string       `gorm:"type:varchar(255);column:jobclass" json:"jobclass"`                //任务类别
	URL            string       `gorm:"type:text;column:url" json:"url"`                                  //url
	JobType        string       `gorm:"type:varchar(255);column:jobtype" json:"jobtype"`                  //任务类型
	Status         string       `gorm:"type:varchar(255);column:status" json:"status"`                    //任务状态
	ScheduleType   string       `gorm:"type:varchar(255);column:scheduletype" json:"scheduletype"`        //调度类型
	CronExpression string       `gorm:"type:varchar(255);column:cronexpression" json:"cronexpression"`    //CRON表达式
	JobParams      string       `gorm:"type:varchar(500);column:jobparams" json:"jobparams"`              //任务参数
	RequestType    string       `gorm:"type:varchar(255);column:requesttype" json:"requesttype"`          //请求类型
	ResponseCode   string       `gorm:"type:varchar(255);column:responsecode" json:"responsecode"`        //响应码
	ResponseType   string       `gorm:"type:varchar(255);column:responsetype" json:"responsetype"`        //响应方式
	ReturnKey      string       `gorm:"type:varchar(255);column:returnkey" json:"returnkey"`              //返回值key
	ReturnValue    string       `gorm:"type:varchar(255);column:returnvalue" json:"returnvalue"`          //返回值value
	CreateTime     string       `gorm:"type:timestamptz(6);column:createtime" json:"createtime"`          //创建时间
	EditeTime      string       `gorm:"type:timestamptz(6);column:editetime" json:"editetime"`            //更新时间
	Integrated     string       `gorm:"type:varchar(255);column:integrated" json:"integrated"`
	IsAutoStart    bool         `gorm:"type:bool;column:isautostart" json:"isautostart"`
	EntryId        cron.EntryID `gorm:"-"`
	Every          int          `gorm:"type:int4;column:every" json:"every"`
	Each           CronExp      `gorm:"type:jsonb;column:each" json:"each"`
	Token          string       `gorm:"type:varchar(500);column:token" json:"token"`
	Authorization  string       `gorm:"type:varchar(500);column:authorization" json:"authorization"`
	AuthType       string       `gorm:"type:varchar(255);column:authtype" json:"authtype"`
	LastInstanceId string       `gorm:"type:varchar(255);column:lastinstanceid" json:"lastinstanceid"`
	Timeout        int          `gorm:"type:int;column:timeout" json:"timeout"`
	RetryCount     int          `gorm:"type:int;column:retrycount" json:"retrycount"`
	RetryInterval  int          `gorm:"type:int;column:retryinterval" json:"retryinterval"`
	// EnId           int          `gorm:"type:int;column:enid" json:"enid"`
}

type Executor struct {
	ID         string     `gorm:"type:varchar(19);primaryKey" json:"id"`
	Name       string     `gorm:"type:varchar(255);column:name" json:"name"`
	CreateUser string     `gorm:"type:varchar(255);column:createuser" json:"createuser"`
	EditeTime  string     `gorm:"type:timestamp(6);column:editetime" json:"editetime"` //更新时间
	Desc       string     `gorm:"type:varchar(500);column:desc" json:"desc"`
	MaxTasks   int        `gorm:"type:int2;column:max_tasks" json:"max_tasks"`
	ExistTasks int        `gorm:"type:int2;column:exist_tasks" json:"exist_tasks"`
	IsDefault  bool       `gorm:"type:bool;column:is_default" json:"is_default"`
	Cron       *cron.Cron `gorm:"-"`
}

type CronExp struct {
	Type          string `json:"type"`
	EffectiveDate string `json:"effectiveDate"`
	Time          string `json:"time"`
	DayOfWeek     string `json:"dayOfWeek"`
	DayOfMonth    string `json:"dayOfMonth"`
}

type CronExpIoit struct {
	Type          string `json:"type"`          // day,week,month
	EffectiveDate string `json:"effectiveDate"` // 年月日时分秒
	EndDate       string `json:"endDate"`       // date | endless(无限期)
	EndType       string `json:"endType"`       // date | endless(无限期)
	DayOfWeek     string `json:"dayOfWeek"`     //1，2，3
	DayOfMonth    string `json:"dayOfMonth"`    //1，2，3，20
	ExecuteMonth  string `json:"executeMonth"`  // 1，2，3
}

type PeriodIoit struct {
	Every         int    `json:"every"` //
	Times         int    `json:"times"`
	EffectiveDate string `json:"effectiveDate"`
	Unit          string `json:"unit"`    //s | m | h
	EndDate       string `json:"endDate"` // endless(无限期) | date(结束时间)
	EndType       string `json:"endType"` // endless(无限期) | date(结束时间)
}

type JobInfoIoit struct {
	CronDetail   CronExpIoit `gorm:"type:jsonb;column:cron_detail" json:"cron_detail"`
	PeriodDetail PeriodIoit  `gorm:"type:jsonb;column:period_detail" json:"period_detail"`
	NameSpace    string      `gorm:"type:varchar(255);column:namespace" json:"namespace"`
	Code         string      `gorm:"type:varchar(255);column:code" json:"code"`
}

func (h CronExp) Value() (driver.Value, error) {
	return json.Marshal(h)
}

func (h *CronExp) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		logger.Info("type assertion failed")
	}
	return json.Unmarshal(bytes, h)
}

func (h PeriodIoit) Value() (driver.Value, error) {
	return json.Marshal(h)
}

func (h *PeriodIoit) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		logger.Info("type assertion failed")
	}
	return json.Unmarshal(bytes, h)
}

func (h CronExpIoit) Value() (driver.Value, error) {
	return json.Marshal(h)
}

func (h *CronExpIoit) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		logger.Info("type assertion failed")
	}
	return json.Unmarshal(bytes, h)
}
