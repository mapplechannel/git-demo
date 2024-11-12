package route

import (
	"hsm-scheduling-back-end/internal/api/healthCheck"
	jobschedule "hsm-scheduling-back-end/internal/job_schedule"
	"hsm-scheduling-back-end/middleware"

	"github.com/gin-gonic/gin"
)

func CollectRoute(r *gin.Engine) *gin.Engine {
	r.Use(middleware.CORSMiddleware())

	// job_schedule
	jobRouter := r.Group("api/hsm-ds/task")
	{
		jobRouter.POST("/add", jobschedule.Add)
		jobRouter.POST("/update", jobschedule.Update)
		jobRouter.POST("/delete", jobschedule.Delete)
		jobRouter.POST("/find", jobschedule.Find)
		jobRouter.GET("/findAll", jobschedule.FindAll)
		jobRouter.POST("/run", jobschedule.Run)
		jobRouter.POST("/manualRun", jobschedule.ManualRun)
		jobRouter.POST("/stop", jobschedule.Stop)
	}

	logRouter := r.Group("api/hsm-ds/log")
	{
		logRouter.GET("/findAll", jobschedule.FindLog)
	}

	instanceRouter := r.Group("api/hsm-ds/instance")
	{
		instanceRouter.GET("/findAll", jobschedule.FindRunningLog)
		instanceRouter.POST("/findInstanceLog", jobschedule.FindInstanceLog)
		instanceRouter.POST("/delete", jobschedule.DeleteInstance)
	}

	openRouter := r.Group("api/hsm-ds/openApi")
	{
		openRouter.POST("/autoAddJob", jobschedule.AutoAddJob)
	}

	executorRouter := r.Group("api/hsm-ds/executor")
	{
		executorRouter.POST("/add", jobschedule.ExecutorAdd)
		executorRouter.POST("/update", jobschedule.ExecutorEdite)
		executorRouter.POST("/delete", jobschedule.ExecutorDelete)
		executorRouter.POST("/find", jobschedule.ExecutorFind)
		executorRouter.GET("/findAll", jobschedule.ExecutorFindAll)
		executorRouter.GET("/print", jobschedule.ExecutorPrint)
	}

	healthRouter := r.Group("api/hsm-ds/health")
	{
		healthRouter.GET("/check", jobschedule.CheckHealth)
	}

	healthCheckRouter := r.Group("test")
	{
		healthCheckRouter.GET("/alive", healthCheck.Alive)
	}

	return r
}
