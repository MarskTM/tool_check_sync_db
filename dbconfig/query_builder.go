package dbconfig

import (
	"fmt"
	"reflect"
	"strings"
)

func QueryBuilder(payload interface{}, tableName string) string {
	reflectVal := reflect.ValueOf(payload)
	reflectType := reflect.TypeOf(payload)

	var isInsert bool

	onDuplicateFields := []string{}
	for i := 0; i < reflectType.NumField(); i++ {
		tempField := reflectType.Field(i).Tag.Get("cql")
		tempVal := fmt.Sprintf("%v", reflectVal.Field(i).Interface())

		if tempField == "id" && reflectVal.Field(i).Interface().(int32) == 0 {
			isInsert = true
			continue
		}

		if tempField == "created_at" || tempField == "updated_at" || tempField == "" || tempVal == "" {
			continue
		}

		if reflectVal.Field(i).Kind() == reflect.String {
			onDuplicateFields = append(onDuplicateFields, tempField+" = '"+tempVal+"'")
			continue
		}
		onDuplicateFields = append(onDuplicateFields, tempField+" = "+tempVal)
	}

	query := ""
	if isInsert {
		query += "INSERT INTO " + tableName + " SET " + strings.Join(onDuplicateFields, ",") + " ON DUPLICATE KEY UPDATE " + strings.Join(onDuplicateFields, ",") + ";"
	} else {
		query += "UPDATE " + tableName + " SET " + strings.Join(onDuplicateFields, ",") + " WHERE id = " + fmt.Sprintf("%v", reflectVal.Field(0).Interface()) + ";"
	}

	return query
}

func QueryCheckExistID(payload interface{}, tableName string) string {
	if payload == nil || tableName == "" {
		return ""
	}
	query := "SELECT ID FROM " + tableName + " WHERE id IN "
	reflectVal := reflect.ValueOf(payload).Interface().([]int32)

	listId := []string{}
	for i := 0; i < len(reflectVal); i++ {
		listId = append(listId, fmt.Sprintf("%v", reflectVal[i]))
	}
	query += "(" + strings.Join(listId, ",") + ");"
	return query
}

func QueryDataByID(listId []int32, tableName string) string {
	if listId == nil || tableName == "" {
		return ""
	}
	query := "SELECT * FROM " + tableName + " WHERE id IN "
	listIdStr := []string{}
	for i := 0; i < len(listId); i++ {
		listIdStr = append(listIdStr, fmt.Sprintf("%v", listId[i]))
	}
	query += "(" + strings.Join(listIdStr, ",") + ");"
	return query
}
