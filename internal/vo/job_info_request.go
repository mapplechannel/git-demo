package vo

import "github.com/robfig/cron/v3"

// 添加任务请求
type UpdateJob struct {
	GetID
	CreateJobRequest
}

type CreateJobRequest struct {
	JobName        string `json:"jobname"`        //任务名称
	CreateUser     string `json:"createuser"`     //创建人
	Executor       string `json:"executor"`       //执行器
	Describe       string `json:"describe"`       //描述
	JobClass       string `json:"jobclass"`       //任务类别
	URL            string `json:"url"`            //url
	JobType        string `json:"jobtype"`        //任务类型
	Status         string `json:"status"`         //任务状态
	ScheduleType   string `json:"scheduletype"`   //调度类型
	CronExpression string `json:"cronexpression"` //CRON表达式
	JobParams      string `json:"jobparams"`      //任务参数
	RequestType    string `json:"requesttype"`    //请求类型
	ResponseCode   string `json:"responsecode"`   //响应码
	ResponseType   string `json:"responsetype"`   //响应方式
	ReturnKey      string `json:"returnkey"`      //返回值key
	ReturnValue    string `json:"returnvalue"`    //返回值value
	CreateTime     string `json:"createtime"`     //创建时间
	EditeTime      string `json:"editetime"`      //更新时间
	Integrated     string `json:"integrated"`     //集成标识
	IsAutoStart    bool   `json:"isautostart"`    //集成标识
	EntryId        cron.EntryID
}

type GetID struct {
	ID string `json:"id"`
}

// 信息系统集成
type IOIT struct {
	NameSpace       string `json:"namespace"`
	SKNode          string `json:"s_node"`
	Code            string `json:"code"`
	Name            string `json:"name"`
	Describe        string `json:"describe"`
	ScheduleDetails string `json:"schedule_details"`
	Status          int    `json:"status"`
	CreateTime      string `json:"create_time"`
	EditeTime       string `json:"edite_time"`
	CreateUse       string `json:"create_user"`
}

type PageReq struct {
	PageSize int `json:"pagesize"`
}
