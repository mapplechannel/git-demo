package constants

import "os"

func getHsmHome() string {
	hsmTemplate := os.Getenv("HSM_HOME")
	if hsmTemplate == "" {
		return "/usr/local/hsm-os"
	}
	return hsmTemplate
}

const (
	Suffix    = ".json"
	Hierarchy = "/" //文件夹层
	// 地址拼接符号
	AddressSplicingSymbols = ":"
	// 空字符串
	EmptyContent = ""
	// PORTAL项目名称
	PORTALNAME = "hsm_io_it"
	// SK配置文件
	NodeManifest = "nodemanifest.ini"
	// /usr/local/hsm-os
	Http = "http://"

	HealthPath = "/api/hsm-ds/health/check"

	IOITPort = "6120"

	IOITExcutor = "hsm-io-it"
	HttpExcutor = "http-executor"

	DomainName = ".local"

	SCHEDULE_TYPE_NONE = "none"
	SCHEDULE_TYPE_CRON = "cron"
)

var HSM_OS_ROOT = getHsmHome()
