package main

import (
	"fmt"
	"sql-migrate-demo/pkg/confer"
	"sql-migrate-demo/pkg/mysql"

	migrate "github.com/rubenv/sql-migrate"
)

func main() {
	err := confer.Init("config.yaml")
	if err != nil {
		panic(err)
	}
	initMysql()
	sqlMigrate()
}

func initMysql() {
	mysql.InitMysqlPool(confer.GetGlobalConfig().Mysql, false)
}

func sqlMigrate() {
	migrations := &migrate.FileMigrationSource{
		Dir: "./db",
	}
	Orm := mysql.NewDaoMysql().GetOrm()
	sqlDb, err := Orm.DB.DB()
	if err != nil {
		panic(err)
	}
	fmt.Println("sql migrate start...")
	code, err := migrate.Exec(sqlDb, "mysql", migrations, migrate.Up)
	if err != nil {
		panic(err)
	}
	fmt.Println("sql migrate code is ", code)
}
