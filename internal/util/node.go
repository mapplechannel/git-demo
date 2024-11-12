package util

import (
	"fmt"
	"hsm-scheduling-back-end/internal/constants"
	"hsm-scheduling-back-end/pkg/logger"
	"hsm-scheduling-back-end/public/nodeapp"
	"path/filepath"
)

var ServerCode string
var IsMaster bool
var MasterNode string
var BackupNode string
var BackupOrMasterIp string

func GetNodeServer() []string {

	currPath := constants.HSM_OS_ROOT
	path := filepath.Join(currPath, constants.NodeManifest)
	var entity nodeapp.NodeConfig
	err := nodeapp.ParseNodeConfig(path, &entity)

	if err != nil {
		fmt.Println(err)
		return nil
	}
	var res []string

	var index int
	nodeLenthlen := len(entity.Nodes)
	logger.Info("node节点长度，%v", nodeLenthlen)

	currentNode := entity.CurrentNode
	logger.Info("currentNode:%v", currentNode)

	for i, node := range entity.Nodes {
		if node.Code == currentNode {
			index = i
		}
	}

	code := entity.Nodes[index].Code
	name := "数据服务器"
	// nodeName = name
	ip := entity.Nodes[index].Ip1
	masterNode := entity.Nodes[index].MasterNodes
	if masterNode == nil {
		// 表明是主机
		IsMaster = true
		for i, node := range entity.Nodes {
			if node.MasterNodes != nil {
				if node.MasterNodes[0] == currentNode {
					BackupNode = entity.Nodes[i].Code
					BackupOrMasterIp = entity.Nodes[i].Ip1
					logger.Info("BackupNode:%v,BackupOrMasterIp is:%v", BackupNode, BackupOrMasterIp)
				}
			}
		}
	} else {
		// 表明是从机
		for i, node := range entity.Nodes {
			if node.Code == masterNode[0] {
				MasterNode = masterNode[0]
				BackupOrMasterIp = entity.Nodes[i].Ip1
			}
		}
	}

	res = append(res, ip, code, name)
	logger.Info("res:%v", res)
	ServerCode = currentNode
	return res
}