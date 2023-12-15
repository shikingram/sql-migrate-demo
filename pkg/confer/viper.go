package confer

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/spf13/viper"
)

var globalConfig Server
var mutex sync.RWMutex

func Init(configURL string) (err error) {
	v := viper.New()
	v.SetConfigFile(configURL)
	err = v.ReadInConfig()
	if err != nil {
		err = fmt.Errorf("Fatal error config file: %w", err)
		return
	}
	if err := v.Unmarshal(&globalConfig); err != nil {
		return err
	}
	return changeDataByEnv()
}

func changeDataByEnv() (err error) {
	if mysqlDbname := os.Getenv(globalConfig.Mysql.DBName); len(mysqlDbname) > 0 {
		globalConfig.Mysql.DBName = mysqlDbname
	}
	if mysqlWriteAddr := os.Getenv(globalConfig.Mysql.Write.Host); len(mysqlWriteAddr) > 0 {
		globalConfig.Mysql.Write.Host = mysqlWriteAddr
	}
	// 处理mysql地址
	host, port, err := net.SplitHostPort(globalConfig.Mysql.Write.Host)
	if err != nil {
		if !strings.Contains(err.Error(), "missing port in address") {
			err = fmt.Errorf("mysql host port is wrong :%w,%s", err, globalConfig.Mysql.Write.Host)
			return
		}
		err = nil
		if mysqlPort := os.Getenv(fmt.Sprint(globalConfig.Mysql.Write.Port)); len(mysqlPort) > 0 {
			globalConfig.Mysql.Write.Port = mysqlPort
		}
	} else {
		globalConfig.Mysql.Write.Host = host
		portInt, _ := strconv.Atoi(port)
		globalConfig.Mysql.Write.Port = portInt
	}

	if mysqlWriteUser := os.Getenv(globalConfig.Mysql.Write.User); len(mysqlWriteUser) > 0 {
		globalConfig.Mysql.Write.User = mysqlWriteUser
	}
	if mysqlWritePwd := os.Getenv(globalConfig.Mysql.Write.Password); len(mysqlWritePwd) > 0 {
		globalConfig.Mysql.Write.Password = mysqlWritePwd
	}
	globalConfig.Mysql.Write.DBName = globalConfig.Mysql.DBName
	globalConfig.Mysql.Write.Prefix = globalConfig.Mysql.Prefix
	return
}

func GetGlobalConfig() *Server {
	mutex.RLock()
	defer mutex.RUnlock()
	return &globalConfig
}
