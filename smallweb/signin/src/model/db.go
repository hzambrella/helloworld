package model

import (
	"database/sql"
	"engine/datastore"
	"fmt"
	"os"
	"sync"

	_ "github.com/go-sql-driver/mysql"
)

//前面的杠表示，运行这个包的init函数。
//参考http://www.tuicool.com/articles/jyqq63
//	_ "github.com/go-sql-driver/mysql"

type DBLinkPool interface {
	GetDB(dbName string) *sql.DB
}

//连接数据库对系统而言是很大的开销。缓存连接的目的是减少开销。这里模仿 柳丁的git.ot24.net/go/engine/datastore包,可能不周全
var mu sync.Mutex

type DBstore struct {
	Ms map[string]*sql.DB
	Mu sync.Mutex
}

// 实例化一个缓存器
var LinkStore DBstore

var driver string = "mysql"
var dbname string
var dsnCfg string

func init() {
	m := make(map[string]*sql.DB, 0)
	LinkStore = DBstore{m, mu}

	// 路径花式出错:no such file or director。这里go test对的， 但是main的init就出错。按照main的改，可以运行，但XXNewDB又出错。不得已规定了一个环境变量，用绝对路径
	//data, err := datastore.ParseDataFromFile("../../db.cfg")
	//data, err := datastore.ParseDataFromFile( "../db.cfg")
	data, err := datastore.ParseDataFromFile(os.Getenv("ETCDIR") + "/db.cfg")
	if err != nil {
		panic(err)
	}

	dbname = data["dbname"]
	dsnCfg = dbname + data["dsn"]
	fmt.Println("数据库配置是:", dsnCfg)

}

//获取数据库连接
func (d DBstore) GetDB(linkName string) (*sql.DB, error) {
	d.Mu.Lock()
	defer d.Mu.Unlock()

	db, ok := d.Ms[linkName]
	if !ok || db == nil {
		newlink, err := sql.Open(driver, dsnCfg)
		if err != nil {
			return nil, err
		}
		d.Ms[dbname] = newlink
		return newlink, nil
	}
	return db, nil
}

//compare
func (d DBstore) GetDBNoPool(linkName string) (*sql.DB, error) {
	d.Mu.Lock()
	defer d.Mu.Unlock()

	newlink, err := sql.Open(driver, dsnCfg)
	if err != nil {
		return nil, err
	}
	return newlink, nil
}
