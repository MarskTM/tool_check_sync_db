package dbconfig

import (
	"fmt"
	"time"

	"k8s.io/klog/v2"

	"github.com/jmoiron/sqlx"
	_ "github.com/go-sql-driver/mysql"

	// "gorm.io/gorm"
    // "gorm.io/driver/mysql"
)

var (
	// DBConfig is the configuration for the database
	DBConfig DataBaseConfig
)

type DataBaseConfig struct {
	MySql    MySqlConfig
	MySqlDes MySqlConfig
}

type MySqlConfig struct {
	Name        string
	Environment string
	DSN         string // Data Source Name
	Host        string
	Port        int
	Username    string
	Password    string
	Database    string
	Active      int
	Idle        int
	Lifetime    time.Duration
}

func (mysqlc *MySqlConfig) buildDSN() string {
	// dsn = "testuser:Mysohapass@tcp(10.5.45.76:3306)/devtalk?parseTime=true&charset=utf8mb4"
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&charset=utf8mb4&multiStatements=true", mysqlc.Username, mysqlc.Password, mysqlc.Host, mysqlc.Port, mysqlc.Database)
}

// This function is only use for lib sqlx
func NewSqlxDB(c *MySqlConfig) (db *sqlx.DB) {
	fmt.Println("DNS", c.buildDSN())
	db, err := sqlx.Connect("mysql", c.buildDSN())
	if err != nil {
		klog.Errorf("Connect db error: %s", err)
	}

	klog.Infof("NewSqlxDB: %s - %+v", c.Name, db.Stats())
	db.SetConnMaxLifetime(c.Lifetime * time.Second)
	db.SetMaxOpenConns(c.Active)
	db.SetMaxIdleConns(c.Idle)
	return
}

// func NewDBGorm(c *MySqlConfig) (db *gorm.DB) {
// 	dsn := c.buildDSN()
// 	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
// 	if err != nil {
// 		klog.Errorf("Connect db error: %s", err)
// 	}
//     klog.Infof("NewSqlxDB: %s - %s", c.Name, c.Database)
// 	return
// }
