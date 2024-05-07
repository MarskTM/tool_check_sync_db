package dbconfig

import (
	"fmt"

	"github.com/BurntSushi/toml"
	"github.com/jmoiron/sqlx"
	// "gorm.io/gorm"
)

var (
	conf DataBaseConfig
	sourceDB      *sqlx.DB
	destinationDB *sqlx.DB

	// sourceDB      *gorm.DB
	// destinationDB *gorm.DB
)

func init() {
	// load config from file .toml
	_, err := toml.DecodeFile("conf.dev.toml", &conf)
	if err != nil {
		err = fmt.Errorf("decode file %s error: %v", "conf.dev.toml", err)
		panic(err)
	}

	sourceDB = NewSqlxDB(&conf.MySql)
	destinationDB = NewSqlxDB(&conf.MySqlDes)

	// sourceDB = NewDBGorm(&conf.MySql)
	// destinationDB = NewDBGorm(&conf.MySqlDes)
}

// func GetSourceDB() *gorm.DB{
// 	return sourceDB
// }

// func GetDestinationDB() *gorm.DB {
// 	return destinationDB
// }

func GetSourceDB() *sqlx.DB {
	return sourceDB
}

func GetDestinationDB() *sqlx.DB {
    return destinationDB
}