package helper

import (
	"database/sql"
	"fmt"
	"log"
	"sync"

	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
)

var mu sync.Mutex

type Mysql struct {
	DB *sql.DB
}

var instance *Mysql

func GetMysqlInstance() *Mysql {
	mu.Lock()
	defer mu.Unlock()

	if instance == nil {
		instance = &Mysql{}
		instance.connect()
	}
	return instance
}
func (m *Mysql) connect() {
	dbhost := viper.Get("dbhost")
	dbport := viper.Get("dbport")
	dbuser := viper.Get("dbuser")
	dbpwd := viper.Get("dbpwd")
	dbname := viper.Get("dbname")
	db, err := sql.Open("mysql", dbuser.(string)+":"+dbpwd.(string)+"@tcp("+dbhost.(string)+":"+dbport.(string)+")/"+dbname.(string)+"?charset=utf8")
	if err != nil {
		log.Fatal(err)
	}
	m.DB = db
}
func (m *Mysql) MysqlWrite(Name string, content map[string]string) bool {
	sql := "insert " + Name + "("
	var first = true
	for k, _ := range content {
		if first {
			first = false
		} else {
			sql += ","
		}
		sql += k
	}
	sql += ") values ("
	first = true
	for _, v := range content {
		if first {
			first = false
		} else {
			sql += ","
		}
		sql += fmt.Sprintf(`"%s"`, v)
	}
	sql += ")"
	_, err := m.DB.Exec(sql)
	if err != nil {
		log.Fatal(err)
		return false
	} else {
		return true
	}

}
func (m *Mysql) CreateTable(name string, fields map[string]string) bool {
	sql := "create table IF NOT EXISTS " + name + " ("
	sql += "id bigint unsigned not null auto_increment,"
	for k, v := range fields {
		sql += k + " " + v + " ,"
	}
	sql += "primary key (id)) engine = InnoDB default charset = utf8mb4;"
	_, err := m.DB.Exec(sql)
	if err != nil {
		log.Fatal(err)
		return false
	} else {
		return true
	}
	//rowAffected, _ := result.RowsAffected()
	//lastInsertId, _ := result.LastInsertId()
}
