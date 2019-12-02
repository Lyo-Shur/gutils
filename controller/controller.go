package controller

import (
	"github.com/Lyo-Shur/gorm"
	"github.com/Lyo-Shur/gorm/generate/mvc"
	"github.com/Lyo-Shur/gorm/info"
	"github.com/Lyo-Shur/gorm/tool"
	"github.com/Lyo-Shur/gutils/api"
	"github.com/kataras/iris"
	"log"
	"strconv"
	"strings"
	"time"
)

// 返回状态码以及数据
const Success = 0
const Fail = -1

const DataAnalysisFailMessage = "数据解析失败"
const ExecuteFailMessage = "执行出错"
const SuccessMessage = "成功"

const NullData = ""

// 直接初始化数据库相关所有信息
func DBHelpers() *dbHelpers {
	dbHelpers := dbHelpers{}
	dbHelpers.Map = make(map[string]dbHelper)
	for k := range gorm.Client.DBS {
		dbHelper := dbHelper{}
		dbHelper.database = info.GetMultiDataBase(k)
		dbHelper.serviceMap = make(map[string]mvc.Service)
		for i := 0; i < len(dbHelper.database.Tables); i++ {
			table := dbHelper.database.Tables[i]
			serviceName := tool.ToSmallHump(table.Name)
			dbHelper.serviceMap[serviceName] = mvc.GetService(dbHelper.database, table.Name)
		}
		dbHelpers.Map[k] = dbHelper
	}
	return &dbHelpers
}

// 多数据库dbHelper集合
type dbHelpers struct {
	Map map[string]dbHelper
}

// 数据库信息 数据查询service
type dbHelper struct {
	database   info.DataBase
	serviceMap map[string]mvc.Service
}

// ========================================= 控制器层 ========================================= //

// 简单控制器层
type SimpleApi struct{}

// 主请求体
// clientAlias 链接别名
// tableName 数据库名
// operation 操作命令
// ctx IRIS上下文
// dbHelpers 多数据库查询集合(使用IRIS注入)
func (simpleApi *SimpleApi) PostBy(clientAlias string, tableName string, operation string, ctx iris.Context, dbHelpers *dbHelpers) string {
	// 分发请求
	unknown := "unknown"
	dbHelper, ok := dbHelpers.Map[clientAlias]
	if !ok {
		simpleApi.operationFactory(unknown)(dbHelper, tableName, ctx)
	}
	if !dbHelper.database.TableExist(tableName) {
		simpleApi.operationFactory(unknown)(dbHelper, tableName, ctx)
	}
	return simpleApi.operationFactory(operation)(dbHelper, tableName, ctx)
}

// 操作分发工厂
func (simpleApi *SimpleApi) operationFactory(operation string) func(dbHelper dbHelper, tableName string, ctx iris.Context) string {
	switch operation {
	case "list":
		return simpleApi.list
	case "model":
		return simpleApi.model
	case "insert":
		return simpleApi.insert
	case "update":
		return simpleApi.update
	case "delete":
		return simpleApi.delete
	default:
		return simpleApi.unknown
	}
}

func (simpleApi *SimpleApi) list(dbHelper dbHelper, tableName string, ctx iris.Context) string {
	// 解析参数
	values, err := simpleApi.Values(dbHelper, tableName, ctx)
	if err != nil {
		log.Println(err)
		return api.JsonCodeModeDTO(Fail, DataAnalysisFailMessage, NullData)
	}
	// 执行查询
	table, err := dbHelper.serviceMap[tableName].GetList(values)
	if err != nil {
		log.Println(err)
		return api.JsonCodeModeDTO(Fail, ExecuteFailMessage, NullData)
	}
	count, err := dbHelper.serviceMap[tableName].GetCount(values)
	if err != nil {
		log.Println(err)
		return api.JsonCodeModeDTO(Fail, ExecuteFailMessage, NullData)
	}
	// 组合外键
	tableMap := table.ToMap()
	keys := dbHelper.database.GetTable(tableName).Keys
	for i := 0; i < len(tableMap); i++ {
		for j := 0; j < len(keys); j++ {
			query := tableMap[i][tool.ToBigHump(keys[j].ColumnName)]
			table, err := dbHelper.serviceMap[tool.ToSmallHump(keys[j].RelyTable)].GetModel(query)
			if err != nil {
				log.Println(err)
				return api.JsonCodeModeDTO(Fail, ExecuteFailMessage, NullData)
			}
			tempMap := table.ToMap()
			if len(tempMap) >= 1 {
				tableMap[i][tool.ToBigHump(keys[j].RelyTable)] = tempMap[0]
			}
		}
	}
	return api.JsonCodeModeDTO(Success, SuccessMessage, map[string]interface{}{
		"List":  tableMap,
		"Count": count,
	})
}

func (simpleApi *SimpleApi) model(dbHelper dbHelper, tableName string, ctx iris.Context) string {
	// 解析参数
	values, err := simpleApi.Values(dbHelper, tableName, ctx)
	if err != nil {
		log.Println(err)
		return api.JsonCodeModeDTO(Fail, DataAnalysisFailMessage, NullData)
	}
	// 根据索引确定主键
	key := ""
	indexs := dbHelper.database.GetTable(tableName).Indexs
	for i := 0; i < len(indexs); i++ {
		index := indexs[i]
		if index.Name == "PRIMARY" {
			key = tool.ToBigHump(index.ColumnName)
		}
	}
	// 判断主键是非存在
	v, ok := values[key]
	if !ok {
		return api.JsonCodeModeDTO(Fail, "缺少参数"+key, NullData)
	}
	// 执行查询
	table, err := dbHelper.serviceMap[tableName].GetModel(v)
	if err != nil {
		log.Println(err)
		return api.JsonCodeModeDTO(Fail, ExecuteFailMessage, NullData)
	}
	// 判断数据是否存在
	m := table.ToMap()
	if len(m) < 1 {
		return api.JsonCodeModeDTO(Success, SuccessMessage, NullData)
	}
	data := m[0]
	// 组合外键
	keys := dbHelper.database.GetTable(tableName).Keys
	for i := 0; i < len(keys); i++ {
		query := data[tool.ToBigHump(keys[i].ColumnName)]
		table, err := dbHelper.serviceMap[tool.ToSmallHump(keys[i].RelyTable)].GetModel(query)
		if err != nil {
			log.Println(err)
			return api.JsonCodeModeDTO(Fail, ExecuteFailMessage, NullData)
		}
		tempMap := table.ToMap()
		if len(tempMap) >= 1 {
			data[tool.ToBigHump(keys[i].RelyTable)] = tempMap[0]
		}
	}
	return api.JsonCodeModeDTO(Success, SuccessMessage, data)
}

func (simpleApi *SimpleApi) insert(dbHelper dbHelper, tableName string, ctx iris.Context) string {
	// 解析参数
	values, err := simpleApi.Values(dbHelper, tableName, ctx)
	if err != nil {
		log.Println(err)
		return api.JsonCodeModeDTO(Fail, DataAnalysisFailMessage, NullData)
	}
	// 类型推断格式校验
	table := dbHelper.database.GetTable(tableName)
	for i := 0; i < len(table.Columns); i++ {
		// 栏位信息
		column := table.Columns[i]
		key := tool.ToBigHump(column.Name)
		// 分析字段类型 长度
		sp := strings.Split(column.Type, "(")
		columnType := sp[0]
		switch columnType {
		case "varchar":
			{
				l := int64(len(values[key].(string)))
				strColumnLength := strings.Split(sp[1], ")")[0]
				columnLength, err := strconv.ParseInt(strColumnLength, 10, 64)
				if err != nil {
					log.Println(err)
					return api.JsonCodeModeDTO(Fail, "数据长度错误", NullData)
				}
				if l == 0 || l > columnLength {
					return api.JsonCodeModeDTO(Fail, key+"字段长度错误", NullData)
				}
				break
			}
		default:
			{
				break
			}
		}
	}
	// 数据插入
	key, err := dbHelper.serviceMap[tableName].Insert(values)
	if err != nil {
		log.Println(err)
		return api.JsonCodeModeDTO(Fail, ExecuteFailMessage, NullData)
	}
	return api.JsonCodeModeDTO(Success, SuccessMessage, key)
}

func (simpleApi *SimpleApi) update(dbHelper dbHelper, tableName string, ctx iris.Context) string {
	// 解析参数
	values, err := simpleApi.Values(dbHelper, tableName, ctx)
	if err != nil {
		log.Println(err)
		return api.JsonCodeModeDTO(Fail, DataAnalysisFailMessage, NullData)
	}
	// 数据更新
	key, err := dbHelper.serviceMap[tableName].Update(values)
	if err != nil {
		log.Println(err)
		return api.JsonCodeModeDTO(Fail, ExecuteFailMessage, NullData)
	}
	return api.JsonCodeModeDTO(Success, SuccessMessage, key)
}

func (simpleApi *SimpleApi) delete(dbHelper dbHelper, tableName string, ctx iris.Context) string {
	// 解析参数
	values, err := simpleApi.Values(dbHelper, tableName, ctx)
	if err != nil {
		log.Println(err)
		return api.JsonCodeModeDTO(Fail, DataAnalysisFailMessage, NullData)
	}
	// 根据索引确定主键
	key := ""
	indexs := dbHelper.database.GetTable(tableName).Indexs
	for i := 0; i < len(indexs); i++ {
		index := indexs[i]
		if index.Name == "PRIMARY" {
			key = tool.ToBigHump(index.ColumnName)
		}
	}
	// 判断主键是非存在
	v, ok := values[key]
	if !ok {
		return api.JsonCodeModeDTO(Fail, "缺少参数"+key, NullData)
	}
	// 执行查询
	count, err := dbHelper.serviceMap[tableName].Delete(v)
	if err != nil {
		log.Println(err)
		return api.JsonCodeModeDTO(Fail, ExecuteFailMessage, NullData)
	}
	return api.JsonCodeModeDTO(Success, SuccessMessage, count)
}

// unknown操作
func (simpleApi *SimpleApi) unknown(dbHelper dbHelper, tableName string, ctx iris.Context) string {
	return "unknown request path"
}

// ========================================= 参数处理层 ========================================= //

// 读取数据
func (simpleApi *SimpleApi) Values(dbHelper dbHelper, tableName string, ctx iris.Context) (map[string]interface{}, error) {
	m := make(map[string]string)
	// 读取表单中的数据
	formMap, err := simpleApi.formValues(dbHelper, tableName, ctx)
	if err != nil {
		return nil, err
	}
	for k, v := range formMap {
		m[tool.ToBigHump(k)] = v
	}
	// 读取JSON中的数据
	jsonMap, err := simpleApi.jsonValues(dbHelper, tableName, ctx)
	if err != nil {
		return nil, err
	}
	for k, v := range jsonMap {
		m[tool.ToBigHump(k)] = v
	}
	// 清洗数据
	attr, err := simpleApi.clear(dbHelper, tableName, m)
	if err != nil {
		return nil, err
	}
	return attr, nil
}

// 读取表单数据
func (simpleApi *SimpleApi) formValues(dbHelper dbHelper, tableName string, ctx iris.Context) (map[string]string, error) {
	values := make(map[string]string)
	for k, v := range ctx.FormValues() {
		values[k] = v[0]
	}
	return values, nil
}

// 读取JSON数据
func (simpleApi *SimpleApi) jsonValues(dbHelper dbHelper, tableName string, ctx iris.Context) (map[string]string, error) {
	values := make(map[string]string)
	if ctx.GetHeader("Content-Type") != "application/json" {
		return values, nil
	}
	err := ctx.ReadJSON(&values)
	return values, err
}

// 清洗数据
func (simpleApi *SimpleApi) clear(dbHelper dbHelper, tableName string, m map[string]string) (map[string]interface{}, error) {
	// 清洗数据
	// 规则是遍历表结构信息
	// 并按表列类型初始化attrMap或者转化m中的值
	attrMap := make(map[string]interface{})
	// 获取并遍历表结构信息
	table := dbHelper.database.GetTable(tableName)
	for i := 0; i < len(table.Columns); i++ {
		// 栏位信息
		column := table.Columns[i]
		key := tool.ToBigHump(column.Name)
		v, ok := m[key]
		switch strings.Split(column.Type, "(")[0] {
		case "int":
			{
				if ok {
					value, err := strconv.ParseInt(v, 10, 64)
					if err != nil {
						return nil, err
					}
					attrMap[key] = value
					break
				}
				attrMap[key] = 0
				break
			}
		case "datetime":
			{
				if ok {
					value, err := time.Parse("2006-01-02 15:04:05", v)
					if err != nil {
						return nil, err
					}
					attrMap[key] = value
					break
				}
				attrMap[key] = time.Time{}
				break
			}
		default:
			{
				if ok {
					attrMap[key] = v
					break
				}
				attrMap[key] = nil
				break
			}
		}
	}
	// 处理分页
	v, ok := m["Start"]
	if ok {
		value, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return nil, err
		}
		attrMap["Start"] = value
	} else {
		attrMap["Start"] = 0
	}
	v, ok = m["Length"]
	if ok {
		value, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return nil, err
		}
		attrMap["Length"] = value
	} else {
		attrMap["Length"] = 10
	}
	return attrMap, nil
}
