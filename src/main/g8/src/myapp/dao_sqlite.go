package myapp

import (
	"fmt"
	"github.com/btnguyen2k/consu/reddo"
	"github.com/btnguyen2k/godal"
	"github.com/btnguyen2k/godal/sql"
	"github.com/btnguyen2k/prom"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"strings"
)

func newSqliteConnection(dir, dbName string) *prom.SqlConnect {
	err := os.MkdirAll(dir, 0711)
	if err != nil {
		panic(err)
	}
	sqlc, err := prom.NewSqlConnect("sqlite3", dir+"/"+dbName+".db", 10000, nil)
	if err != nil {
		panic(err)
	}
	return sqlc
}

func sqliteInitTableGroup(sqlc *prom.SqlConnect, tableName string) {
	sql := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s VARCHAR(64), %s VARCHAR(255), PRIMARY KEY (%s))",
		tableName, colGroupId, colGroupName, colGroupId)
	_, err := sqlc.GetDB().Exec(sql)
	if err != nil {
		panic(err)
	}
}

func sqliteInitTableUser(sqlc *prom.SqlConnect, tableName string) {
	sql := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s VARCHAR(64), %s VARCHAR(64), %s VARCHAR(64), %s VARCHAR(64), PRIMARY KEY (%s))",
		tableName, colUserUsername, colUserPassword, colUserName, colUserGroupId, colUserUsername)
	_, err := sqlc.GetDB().Exec(sql)
	if err != nil {
		panic(err)
	}
}

/*----------------------------------------------------------------------*/

func newUserDaoSqlite(sqlc *prom.SqlConnect, tableName string) UserDao {
	dao := &UserDaoSqlite{tableName: tableName}
	dao.GenericDaoSql = sql.NewGenericDaoSql(sqlc, godal.NewAbstractGenericDao(dao))
	dao.SetRowMapper(&sql.GenericRowMapperSql{
		NameTransformation:          sql.NameTransfLowerCase,
		GboFieldToColNameTranslator: map[string]map[string]interface{}{tableName: mapFieldToColNameUser},
		ColNameToGboFieldTranslator: map[string]map[string]interface{}{tableName: mapColNameToFieldUser},
		ColumnsListMap:              map[string][]string{tableName: colsUser},
	})
	return dao
}

const (
	tableUser       = namespace + "_user"
	colUserUsername = "uname"
	colUserPassword = "upwd"
	colUserName     = "display_name"
	colUserGroupId  = "gid"
)

var (
	colsUser              = []string{colUserUsername, colUserPassword, colUserName, colUserGroupId}
	mapFieldToColNameUser = map[string]interface{}{fieldUserUsername: colUserUsername, fieldUserPassword: colUserPassword, fieldUserName: colUserName, fieldUserGroupId: colUserGroupId}
	mapColNameToFieldUser = map[string]interface{}{colUserUsername: fieldUserUsername, colUserPassword: fieldUserPassword, colUserName: fieldUserName, colUserGroupId: fieldUserGroupId}
)

type UserDaoSqlite struct {
	*sql.GenericDaoSql
	tableName string
}

/*----------------------------------------------------------------------*/

func newGroupDaoSqlite(sqlc *prom.SqlConnect, tableName string) GroupDao {
	dao := &GroupDaoSqlite{tableName: tableName}
	dao.GenericDaoSql = sql.NewGenericDaoSql(sqlc, godal.NewAbstractGenericDao(dao))
	dao.SetRowMapper(&sql.GenericRowMapperSql{
		NameTransformation:          sql.NameTransfLowerCase,
		GboFieldToColNameTranslator: map[string]map[string]interface{}{tableName: mapFieldToColNameGroup},
		ColNameToGboFieldTranslator: map[string]map[string]interface{}{tableName: mapColNameToFieldGroup},
		ColumnsListMap:              map[string][]string{tableName: colsGroup},
	})
	return dao
}

const (
	tableGroup   = namespace + "_group"
	colGroupId   = "gid"
	colGroupName = "gname"
)

var (
	colsGroup              = []string{colGroupId, colGroupName}
	mapFieldToColNameGroup = map[string]interface{}{fieldGroupId: colGroupId, fieldGroupName: colGroupName}
	mapColNameToFieldGroup = map[string]interface{}{colGroupId: fieldGroupId, colGroupName: fieldGroupName}
)

type GroupDaoSqlite struct {
	*sql.GenericDaoSql
	tableName string
}

// it is recommended to have a function that transforms godal.IGenericBo to business object and vice versa.
func (dao *GroupDaoSqlite) toBo(gbo godal.IGenericBo) *Group {
	if gbo == nil {
		return nil
	}
	bo := &Group{
		Id:   gbo.GboGetAttrUnsafe(fieldGroupId, reddo.TypeString).(string),
		Name: gbo.GboGetAttrUnsafe(fieldGroupName, reddo.TypeString).(string),
	}
	return bo
}

// it is recommended to have a function that transforms godal.IGenericBo to business object and vice versa.
func (dao *GroupDaoSqlite) toGbo(bo *Group) godal.IGenericBo {
	if bo == nil {
		return nil
	}
	gbo := godal.NewGenericBo()
	gbo.GboSetAttr(fieldGroupId, bo.Id)
	gbo.GboSetAttr(fieldGroupName, bo.Name)
	return gbo
}

// Get implements GroupDao.Create
func (dao *GroupDaoSqlite) Create(id, name string) (bool, error) {
	bo := &Group{
		Id:   strings.ToLower(strings.TrimSpace(id)),
		Name: strings.TrimSpace(name),
	}
	numRows, err := dao.GdaoCreate(dao.tableName, dao.toGbo(bo))
	return numRows > 0, err
}

// Get implements GroupDao.Get
func (dao *GroupDaoSqlite) Get(id string) (*Group, error) {
	gbo, err := dao.GdaoFetchOne(dao.tableName, map[string]interface{}{colGroupId: id})
	if err != nil {
		return nil, err
	}
	return dao.toBo(gbo), nil
}
