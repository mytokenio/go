package mysql

import (
	"database/sql"
	"fmt"
	"math"
	"reflect"
	"strings"

	"github.com/mytokenio/go/log"
)

// 模型：每个业务模型都需要继承该模型
type Model struct {
	TableName string  `db:"-" json:"-"`
	Tx        *sql.Tx `db:"-" json:"-"`
}

// ---------------------------------------------------------------------------------------------------------------------

// 获取DB
func (m *Model) GetDB() *sql.DB {
	return db
}

// 构建并获取查询语句
func (m *Model) Select(fields string) *Query {
	fieldSlice := strings.Split(fields, ",")
	selectFields := make([]string, len(fieldSlice))
	for i, field := range fieldSlice {
		selectFields[i] = fmt.Sprintf("`%s`", strings.TrimSpace(field))
	}
	return &Query{
		Sql: fmt.Sprintf("SELECT %s", strings.Join(selectFields, ", ")),
	}
}

// 基于SQL查询
func (m *Model) SelectBySql(cmd string, value ...interface{}) (*sql.Rows, error) {
	if m.Tx == nil {
		return db.Query(cmd, value)
	} else {
		return m.Tx.Query(cmd, value)
	}
}

// 查询记录
func (m *Model) SelectWhere(query *Query, exp interface{}) (*sql.Rows, error) {
	if query == nil {
		return nil, fmt.Errorf("params error")
	}

	var err error
	if query.Where, err = m.getWhereByInterface(exp); err != nil {
		return nil, err
	}

	cmd := query.Combination()
	log.Infof("[MySQL]: %s", cmd)

	if m.Tx == nil {
		return db.Query(cmd)
	} else {
		return m.Tx.Query(cmd)
	}
}

// 插入数据：支持 对象指针类型 和 Map 类型
func (m *Model) Insert(data interface{}) (int64, error) {
	t := reflect.TypeOf(data)

	switch t.Kind() {
	case reflect.Ptr:
		if mapping, err := struct2Map(data); err != nil {
			return 0, err
		} else {
			return m.insert(mapping)
		}
	case reflect.Map:
		switch data.(type) {
		case map[string]interface{}:
			return m.insert(data.(map[string]interface{}))
		default:
		}
	default:
	}

	return 0, errTypeInvalid
}

// 插入多条记录：支持 对象指针类型 和 Map 类型
func (m *Model) MInsert(data ...interface{}) (int64, error) {
	var dataLen int
	if dataLen = len(data); dataLen == 0 {
		return 0, errParamsBad
	}

	t := reflect.TypeOf(data[0])
	columns, err := getColumns(data[0])
	if err != nil {
		return 0, err
	}

	switch t.Kind() {
	case reflect.Ptr:
		values := make([]interface{}, 0, dataLen)
		for i := 0; i < dataLen; i++ {
			if ptrValues, err := getValues(data[i]); err != nil {
				return 0, err
			} else {
				values = append(values, ptrValues)
			}
		}
		return m.BatchInsert(columns, values)
	case reflect.Map:
		switch data[0].(type) {
		case map[string]interface{}:
			values := make([]interface{}, 0, dataLen)
			columns, err := getColumns(data[0])
			if err != nil {
				return 0, err
			}
			subMapLen := len(data[0].(map[string]interface{}))
			for i := 0; i < dataLen; i++ {
				if len(data[i].(map[string]interface{})) != subMapLen {
					return 0, fmt.Errorf("params map key is not the same")
				}
				subMapValues := make([]interface{}, 0, subMapLen)
				for _, column := range columns {
					subMapValues = append(subMapValues, data[i].(map[string]interface{})[column])
				}
				values = append(values, subMapValues)
			}
			return m.BatchInsert(columns, values)
		}
	}

	return 0, errTypeInvalid
}

// 更新：基于exp表达式更新data数据
func (m *Model) Update(data interface{}, exp interface{}) (int64, error) {
	t := reflect.TypeOf(data)

	switch t.Kind() {
	case reflect.Ptr:
		if mapping, err := struct2Map(data); err != nil {
			return 0, err
		} else {
			return m.update(mapping, exp)
		}
	case reflect.Map:
		switch data.(type) {
		case map[string]interface{}:
			return m.update(data.(map[string]interface{}), exp)
		default:
		}
	default:
	}

	return 0, errTypeInvalid
}

// 删除：基于exp表达式删除数据
func (m *Model) Delete(exp interface{}) (int64, error) {
	var result sql.Result

	retWhere, err := m.getWhereByInterface(exp)
	if err != nil {
		return 0, err
	}

	cmd := fmt.Sprintf("DELETE FROM `%v` %v", m.TableName, retWhere)
	log.Infof("[MySQL]: %s", cmd)

	if m.Tx == nil {
		if result, err = db.Exec(cmd); err != nil {
			return 0, err
		}
	} else {
		if result, err = m.Tx.Exec(cmd); err != nil {
			return 0, err
		}
	}

	return result.RowsAffected()
}

// 批量插入数据
func (m *Model) BatchInsert(columns []string, params []interface{}) (int64, error) {
	var err error
	var lastInsertId int64

	paramsLen := len(params)
	count, fraction := math.Modf(float64(paramsLen) / maxBatchLimit)
	if fraction > 0.000001 {
		count += 1
	}

	for i := 0; i < int(count); i++ {
		var endIndex int
		if (i+1)*maxBatchLimit > paramsLen {
			endIndex = paramsLen
		} else {
			endIndex = (i + 1) * maxBatchLimit
		}

		lastInsertId, err = m.batchInsertByLimit(columns, params[i*maxBatchLimit:endIndex])
		if err != nil {
			return 0, err
		}
	}

	return lastInsertId, err
}

// 加载一个值
func (m *Model) LoadValue(rows *sql.Rows, value interface{}) error {
	if count, err := m.Load(rows, value); err != nil {
		return err
	} else if count == 0 {
		return m.ErrNoRows()
	} else {
		return nil
	}
}

// 加载多个值
func (m *Model) LoadValues(rows *sql.Rows, value interface{}) (int, error) {
	return m.Load(rows, value)
}

// 加载结构体
func (m *Model) LoadStruct(rows *sql.Rows, value interface{}) error {
	if count, err := m.Load(rows, value); err != nil {
		return err
	} else if count == 0 {
		return m.ErrNoRows()
	} else {
		return nil
	}
}

// 批量加载结构体
func (m *Model) LoadStructs(rows *sql.Rows, value interface{}) (int, error) {
	return m.Load(rows, value)
}

// 万能加载
func (m *Model) Load(rows *sql.Rows, value interface{}) (int, error) {
	if rows == nil {
		return 0, errParamsBad
	}

	v := reflect.ValueOf(value)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return 0, errParamsBad
	}

	columns, err := rows.Columns()
	if err != nil {
		return 0, err
	}

	count := 0
	v = v.Elem()
	isSlice := v.Kind() == reflect.Slice && v.Type().Elem().Kind() != reflect.Uint8
	for rows.Next() {
		var elem reflect.Value
		if isSlice {
			elem = reflect.New(v.Type().Elem()).Elem()
		} else {
			elem = v
		}
		if ptr, err := findPtr(columns, elem); err != nil {
			return 0, err
		} else {
			if err = rows.Scan(ptr...); err != nil {
				log.Infof("scan err: %v", err)
				return 0, err
			}
		}
		count++
		if isSlice {
			v.Set(reflect.Append(v, elem))
		} else {
			break
		}
	}

	return count, nil
}

// 基于条件表达式判断数据是否存在
func (m *Model) IsExist(exp interface{}, field string, value string) (bool, error) {
	var key string
	query := m.Select(field).Form(m.TableName)
	if row, err := m.SelectWhere(query, exp); err != nil {
		return false, err
	} else {
		if err = row.Scan(&key); err != nil {
			return false, err
		}
		if key == value {
			return false, nil
		}
	}
	return true, nil
}

// 统计
func (m *Model) Count(exp interface{}) (int, error) {
	var total int
	query := m.Select("COUNT(0)").Form(m.TableName)
	if row, err := m.SelectWhere(query, exp); err != nil {
		return 0, err
	} else {
		if err = row.Scan(&total); err != nil {
			return 0, err
		}
	}
	return total, nil
}

func (m *Model) ErrNoRows() error {
	return sql.ErrNoRows
}

// ---------------------------------------------------------------------------------------------------------------------

// 插入params数据
func (m *Model) insert(params map[string]interface{}) (int64, error) {
	if len(params) == 0 {
		return 0, errParamsBad
	}

	var err error
	var result sql.Result

	length := len(params)
	columns := make([]string, 0, length)
	values := make([]string, 0, length)
	for key, value := range params {
		switch value.(type) {
		case NullString:
			value = value.(NullString).String
		case sql.NullString:
			value = value.(sql.NullString).String
		case NullBool:
			value = value.(NullBool).Bool
		case sql.NullBool:
			value = value.(sql.NullBool).Bool
		case NullInt64:
			value = value.(NullInt64).Int64
		case sql.NullInt64:
			value = value.(sql.NullInt64).Int64
		case NullFloat64:
			value = value.(NullFloat64).Float64
		case sql.NullFloat64:
			value = value.(sql.NullFloat64).Float64
		}
		columns = append(columns, key)
		values = append(values, fmt.Sprintf("%v", value))
	}

	fields := fmt.Sprintf("`%s`", strings.Join(columns, "`,`"))
	fieldValues := fmt.Sprintf("'%s'", strings.Join(values, "','"))
	cmd := fmt.Sprintf("INSERT INTO `%s` (%s) VALUES(%s)", m.TableName, fields, fieldValues)
	log.Infof("[MySQL]: %s", cmd)

	if m.Tx == nil {
		if result, err = db.Exec(cmd); err != nil {
			return 0, err
		}
	} else {
		if result, err = m.Tx.Exec(cmd); err != nil {
			return 0, err
		}
	}

	return result.LastInsertId()
}

// 更新：基于exp表达式更新params数据
func (m *Model) update(params map[string]interface{}, exp interface{}) (int64, error) {
	var result sql.Result

	retWhere, err := m.getWhereByInterface(exp)
	if err != nil {
		return 0, err
	}

	length := len(params)
	setValues := make([]string, 0, length)
	for key, value := range params {
		set := fmt.Sprintf("`%v`='%v'", key, value)
		setValues = append(setValues, set)
	}

	retSet := strings.Join(setValues, ", ")
	cmd := fmt.Sprintf("UPDATE `%s` SET %s %s", m.TableName, retSet, retWhere)
	log.Infof("[MySQL]: %s", cmd)

	if m.Tx == nil {
		if result, err = db.Exec(cmd); err != nil {
			return 0, err
		}
	} else {
		if result, err = m.Tx.Exec(cmd); err != nil {
			return 0, err
		}
	}

	return result.RowsAffected()
}

// 基于表达式获取并构建where语句
func (m *Model) getWhereByInterface(exp interface{}) (string, error) {
	var result string

	if exp == nil {
		return "", nil
	}

	switch exp.(type) {
	case map[string]interface{}:
		if len(exp.(map[string]interface{})) > 0 {
			result = fmt.Sprintf(" WHERE %s", m.getWhereItem("AND", exp.(map[string]interface{})))
		}

	case map[string]map[string]interface{}:
		length := len(exp.(map[string]map[string]interface{}))
		if length > 0 {
			wheres := make([]string, 0, length)
			for key, value := range exp.(map[string]map[string]interface{}) {
				keyToUpper := strings.ToUpper(key)
				if keyToUpper == "AND" || keyToUpper == "OR" {
					wheres = append(wheres, m.getWhereItem(keyToUpper, value))
				} else {
					return "", errParamsBad
				}
			}
			result = fmt.Sprintf(" WHERE %s", strings.Join(wheres, " AND "))
		}

	default:
		return "", errParamsBad
	}

	return result, nil
}

// 获取并构建where中的每个子项
func (m *Model) getWhereItem(join string, exp map[string]interface{}) string {
	var result string

	if length := len(exp); length > 0 {
		where := make([]string, 0, length)
		for key, value := range exp {
			where = append(where, strings.Replace(key, "?", fmt.Sprintf("'%v'", value), -1))
		}
		result = fmt.Sprintf("(%s)", strings.Join(where, fmt.Sprintf(" %s ", join)))
	}

	return result
}

func (m *Model) batchInsertByLimit(columns []string, params []interface{}) (int64, error) {
	paramsLen := len(params)
	if paramsLen > maxBatchLimit {
		return 0, fmt.Errorf("batch insert too large, length: %v", paramsLen)
	}

	// 防止字段是关键字，所以加上转义符号，如：`status`
	for key, value := range columns {
		columns[key] = fmt.Sprintf("`%s`", value)
	}

	data := make([]string, paramsLen)
	for i, v := range params {
		val := reflect.ValueOf(v)
		if val.Kind() != reflect.Slice {
			return 0, fmt.Errorf("params error, insert data must be slice")
		}

		var subVal string
		subData := make([]string, 0, val.Len())
		for j := 0; j < val.Len(); j++ {
			switch t := val.Index(j).Interface().(type) {
			case string:
				subVal = fmt.Sprintf("'%s'", t)
			case NullString:
				subVal = fmt.Sprintf("'%s'", t.String)
			case sql.NullString:
				subVal = fmt.Sprintf("'%s'", t.String)
			default:
				subVal = fmt.Sprintf("%v", t)
			}
			subData = append(subData, subVal)
		}
		data[i] = fmt.Sprintf("(%s)", strings.Join(subData, ","))
	}

	var err error
	var result sql.Result
	cmd := fmt.Sprintf("INSERT INTO `%s` (%s) VALUES %s",
		m.TableName, strings.Join(columns, ","), strings.Join(data, ","))
	log.Infof("[MySQL]: %s", cmd)

	if m.Tx == nil {
		if result, err = db.Exec(cmd); err != nil {
			return 0, err
		}
	} else {
		if result, err = m.Tx.Exec(cmd); err != nil {
			return 0, err
		}
	}

	return result.LastInsertId()
}
