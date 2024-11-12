package jobschedule

import (
	"context"
	"database/sql"
	"fmt"
	"hsm-scheduling-back-end/config"
	"hsm-scheduling-back-end/pkg/logger"
	"net"
	"time"
)

type RegistryHost struct {
	Addr         string `json:"addr" binding:"required"`
	HostName     string `json:"hostname" binding:"required"`
	RunningTasks string `json:"runningtasks" binding:"omitempty"`
}

// 获取pg连接
func GetPgConn() (*sql.DB, error) {
	dsn := config.ConfigAll.Postgres.Url + " dbname=" + config.ConfigAll.Postgres.DbName

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// 生成唯一主键
func GetId() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

// 注册新的worker
func RegistryNewHost(ctx context.Context, req *RegistryHost) bool {
	hostsql := `INSERT INTO hsm_scheduling.hosts (id,addr,hostname,last_time_unix) VALUES ($1,$2,$3,$4)`
	addr := req.Addr
	hostname := req.HostName
	// 检查是否重复的执行器
	hosts := getHosts(ctx, addr)
	if len(hosts) == 1 {
		logger.Info("Addr already registered:%v", addr)
		return false
	}
	// 获取数据库连接
	db, err := GetPgConn()
	if err != nil {
		logger.Info("Fail to connect pg:%v", err)
		return false
	}
	defer db.Close()
	// 开启事务
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		logger.Info("Fail to start transaction:%v", err)
		return false
	}
	// 预编译sql
	stmt, err := tx.PrepareContext(ctx, hostsql)
	if err != nil {
		tx.Rollback()
		logger.Info("Fail to prepare statement:%v", err)
		return false
	}
	defer stmt.Close()
	// 获取主键ID
	id := GetId()
	logger.Info("主键id为:%v", id)
	_, err = stmt.ExecContext(ctx, id, addr, hostname, time.Now().Unix())
	if err != nil {
		tx.Rollback()
		logger.Info("Fail to exec statement:%v", err)
		return false
	}
	// 事务提交
	if err := tx.Commit(); err != nil {
		logger.Info("Fail to transaction commit:%v", err)
		return false
	}

	logger.Info("New worker registered:%v", addr)
	return true
}

// 更新已经存在的worker
func UpdateHostWeight(ctx context.Context, weight int, addr string) {
	updatesql := `UPDATE hsm_scheduling.hosts SET weight=$1 WHERE addr=$2`
	db, err := GetPgConn()
	if err != nil {
		logger.Info("Fail to connect pg:%v", err)
		return
	}
	defer db.Close()
	stmt, err := db.PrepareContext(ctx, updatesql)
	if err != nil {
		logger.Info("Fail to prepare statement:%v", err)
		return
	}
	defer stmt.Close()
	_, err = stmt.ExecContext(ctx, weight, addr)
	if err != nil {
		logger.Info("Fail to exec statement:%v", err)
		return
	}
	logger.Info("Weight update for addr:%v success", addr)
}

// 更新心跳
func UpdateHostHeartbeat(ctx context.Context, addr string) bool {
	// addr := req.Addr
	// runningTasks := req.RunningTasks
	updatesql := `UPDATE hsm_scheduling.hosts SET last_time_unix=$1 WHERE addr=$2`
	db, err := GetPgConn()
	if err != nil {
		logger.Info("Fail to connect pg:%v", err)
		return false
	}
	defer db.Close()
	stmt, err := db.PrepareContext(ctx, updatesql)
	if err != nil {
		logger.Info("Fail to prepare statement:%v", err)
		return false
	}
	defer stmt.Close()
	result, err := stmt.ExecContext(ctx, time.Now().Unix(), addr)
	if err != nil {
		logger.Info("Fail to exec statement:%v", err)
		return false
	}
	line, err := result.RowsAffected()
	if err != nil {
		logger.Info("Fail to exec rows affected:%v", err)
		return false
	}
	if line <= 0 {
		logger.Info("Host %v has not registered", addr)
		return false
	}
	logger.Info("heartbeat update for addr:%v success", addr)
	return true
}

// 删除worker
func DeleteHost(ctx context.Context, id string) {
	deletesql := `DELETE FROM hsm_scheduling.hosts WHERE id=$1`
	db, err := GetPgConn()
	if err != nil {
		logger.Info("Fail to connect pg:%v", err)
		return
	}
	defer db.Close()
	stmt, err := db.PrepareContext(ctx, deletesql)
	if err != nil {
		logger.Info("Fail to prepare statement:%v", err)
		return
	}
	defer stmt.Close()
	_, err = stmt.ExecContext(ctx, time.Now().Unix(), id)
	if err != nil {
		logger.Info("Fail to exec statement:%v", err)
		return
	}
	logger.Info("Delete id:%v worker success", id)
}

// 判断是否存在重复的addr
func getHosts(ctx context.Context, addr string) []string {
	var hosts []string
	query := `SELECT addr FROM hsm_scheduling.hosts WHERE addr=$1`
	db, err := GetPgConn()
	if err != nil {
		return nil
	}
	defer db.Close()
	rows, err := db.QueryContext(ctx, query, addr)
	if err != nil {
		return nil
	}
	defer rows.Close()
	for rows.Next() {
		var host string
		if err := rows.Scan(&host); err != nil {
			return nil
		}
		hosts = append(hosts, host)
	}
	if err := rows.Err(); err != nil {
		return nil
	}
	return hosts
}

type RegisterResponse struct {
	Message string `json:"message"`
}

const (
	timeout       = 90 * time.Second
	checkInterval = 30 * time.Second
)

func CleanupExpiredServices(ctx context.Context) {
	logger.Info("start cleanup expired services")
	db, err := GetPgConn()
	if err != nil {
		logger.Info("Fail to connect pg:%v", err)
		return
	}
	defer db.Close()

	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(checkInterval):
			now := time.Now().Unix()
			cutoff := now - int64(timeout.Seconds())
			_, err := db.Exec("DELETE FROM hsm_scheduling.hosts WHERE last_time_unix < $1 AND hostname != 'http-executor'", cutoff)
			if err != nil {
				logger.Info("Failed to delete:%v", err)
			} else {
				logger.Info("success delete no update addr")
			}
		}
	}
}

func UpdateHB(ctx context.Context) bool {
	ip := findIp()
	addr := ip + ":" + "6120"
	updatesql := `UPDATE hsm_scheduling.hosts SET last_time_unix=$1 WHERE addr=$2`
	db, err := GetPgConn()
	if err != nil {
		logger.Info("Fail to connect pg:%v", err)
		return false
	}
	defer db.Close()
	stmt, err := db.PrepareContext(ctx, updatesql)
	if err != nil {
		logger.Info("Fail to prepare statement:%v", err)
		return false
	}
	defer stmt.Close()
	result, err := stmt.ExecContext(ctx, time.Now().Unix(), addr)
	if err != nil {
		logger.Info("Fail to exec statement:%v", err)
		return false
	}
	line, err := result.RowsAffected()
	if err != nil {
		logger.Info("Fail to exec rows affected:%v", err)
		return false
	}
	if line <= 0 {
		logger.Info("Host %v has not registered", addr)
		return false
	}
	logger.Info("heartbeat update for addr:%v success", addr)
	return true
}

func findIp() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}

	var ip string

	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				ip = ipNet.IP.String()
			}
		}
	}
	return ip
}

func ExecutorUpdateHB(ctx context.Context) {
	ticker := time.NewTicker(30 * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				UpdateHB(ctx)
			case <-ctx.Done():
				ticker.Stop()
				return
			}
		}
	}()
}

// func GetAllExcutor() (*[]vo.IoitJsHost, bool) {
// 	db, err := GetPgConnByGORM()
// 	if err != nil {
// 		logger.Info("Fail to connect pg:%v", err)
// 		return nil, false
// 	}

// 	var ioitJsHost []vo.IoitJsHost
// 	result := db.Table("hsm_scheduling.hosts").Select("hostname", "addr").Find(&ioitJsHost)
// 	if result.Error != nil {
// 		logger.Info("Fail to fetch records:%v", result.Error)
// 		return nil, false
// 	}

// 	logger.Info("ioitJsHost:%v", ioitJsHost)
// 	return &ioitJsHost, true
// }
