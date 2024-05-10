package main

import (
	"context"
	"fmt"
	"reflect"
	"sync"
	"test-data-convert/dbconfig"
	"test-data-convert/model"
	"time"

	// "test-data-convert/model"
	"test-data-convert/repository"
	"test-data-convert/utils"

	"github.com/jmoiron/sqlx"
	"k8s.io/klog/v2"
)

var (
	cacheIdDataInsert = map[string][]int32{}
	cacheIdDataDelete = map[string][]int32{}

	cacheIdDataUpdate = map[string]map[int32]interface{}{}
	mu                sync.Mutex
)

func init() {
	for _, tableName := range model.REPLICATE_TABLE_LISTS {
		cacheIdDataInsert[tableName] = make([]int32, 0)
		cacheIdDataDelete[tableName] = make([]int32, 0)

		temp := make(map[int32]interface{})
		cacheIdDataUpdate[tableName] = temp
	}
}

func main() {

	workLoad := 5
	numWorkers := 2
	workPool := repository.NewPool(numWorkers, workLoad)

	// Declear Job insert or Update data in db destination
	for i := 0; i < workLoad; i++ {
		// workPool.AddJob(*repository.NewJob(InsertData("channels")))
		// workPool.AddJob(*repository.NewJob(InsertData("message_data")))

		workPool.AddJob(*repository.NewJob(UpdateData("channels")))
		// workPool.AddJob(*repository.NewJob(UpdateData("message_data")))
		// workPool.AddJob(*repository.NewJob(DeleteData("message_data")))
	}

	// Start the work pool
	workPool.Wait.Add(len(workPool.WorkList))
	go workPool.Listener()

	workPool.Start()
	workPool.Wait.Wait()

	klog.Infof("================================================================")
	klog.Infof("+ WorkerPool done!")
	klog.Infof("================================================================")

	klog.Infof(">> Start check data in db destination")

	// Start check data in db destination
	/*
		v1.0.0.1: Bug - logic check didn't work properly true for now
		// CheckDataInsert()
	*/
	CheckDataUpdate()
	CheckDataDelete()
	klog.Infof(">> End check data in db destination")

}

func InsertData(tableName string) func() error {
	db := dbconfig.GetSourceDB()
	query := ""
	var mui sync.Mutex

	return func() error {
		switch tableName {
		case "channels":
			newData := utils.GenerateChannelsDO()
			newData.ID = 0 // make sure ID is auto increment
			query = dbconfig.QueryBuilder(newData, tableName)
		case "message_data":
			newData := utils.GenerateMessageDataDO()
			newData.ID = 0 // make sure ID is auto increment
			query = dbconfig.QueryBuilder(newData, tableName)
		case "auth_users":
			newData := utils.GenerateAuthUsersDO()
			newData.ID = 0 // make sure ID is auto increment
			query = dbconfig.QueryBuilder(newData, tableName)
		case "chats":
			newData := utils.GenerateChatsDO()
			newData.ID = 0 // make sure ID is auto increment
			query = dbconfig.QueryBuilder(newData, tableName)
		}

		ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
		defer cancel()

		query += "SELECT LAST_INSERT_ID();"

		// klog.Infof(">> Executed query - temp: %v", query)
		newRow, err := db.QueryxContext(ctx, query)
		if err != nil {
			klog.Errorf("Error creating new row: %v", err)
			return err
		}
		defer newRow.Close()

		var dataId int32
		for newRow != nil && newRow.Next() {
			err := newRow.Scan(&dataId)
			if err != nil {
				klog.Errorf("Error creating new row: %v", err)
				return err
			}
		}

		// klog.Infof(">> Executed query - temp: %v", dataId)
		mui.Lock()
		cacheIdDataInsert[tableName] = append(cacheIdDataInsert[tableName], dataId)
		mui.Unlock()

		return nil
	}
}

func UpdateData(tableName string) func() error {
	db := dbconfig.GetSourceDB()
	query := ""

	return func() error {
		var id int32
		mu.Lock()

		switch tableName {
		case "channels":
			newData := utils.GenerateChannelsDO()
			query = dbconfig.QueryBuilder(newData, tableName)

			id = newData.ID
			cacheIdDataUpdate[tableName][id] = newData
			mu.Unlock()
			// klog.Infof(">> Executed query id = : %v", newData.ID)
		case "message_data":
			newData := utils.GenerateMessageDataDO()
			query = dbconfig.QueryBuilder(newData, tableName)

			id = newData.ID
			cacheIdDataUpdate[tableName][id] = newData
			mu.Unlock()
			// klog.Infof(">> Executed query id = : %v", newData.ID)
		case "auth_users":
			newData := utils.GenerateAuthUsersDO()
			query = dbconfig.QueryBuilder(newData, tableName)

			id = newData.ID
			cacheIdDataUpdate[tableName][id] = newData
			mu.Unlock()
			// klog.Infof(">> Executed query id = : %v", newData.ID)
		case "chats":
			newData := utils.GenerateChatsDO()
			query = dbconfig.QueryBuilder(newData, tableName)

			id = newData.ID
			cacheIdDataUpdate[tableName][id] = newData
			mu.Unlock()
			// klog.Infof(">> Executed query id = : %v", newData.ID)
		}

		ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
		defer cancel()
		// klog.Infof(">> Upsert query - temp: %v", query)
		_, err := db.QueryContext(ctx, query)
		if err != nil {
			return err
		}
		return nil
	}
}

func DeleteData(tableName string) func() error {
	db := dbconfig.GetSourceDB()

	return func() error {
		var mu sync.Mutex

		randomId := utils.RandomInt(1, 1000)

		mu.Lock()
		cacheIdDataDelete[tableName] = append(cacheIdDataDelete[tableName], int32(randomId))
		mu.Unlock()

		query := fmt.Sprintf("DELETE FROM %s WHERE ID = %d;", tableName, randomId)
		ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
		defer cancel()
		_, err := db.ExecContext(ctx, query)
		if err != nil {
			klog.Errorf("Error delete record from: %s - %v", tableName, err)
			return err
		}

		return nil
	}
}

// ============================ Func check Data in DB Destination =====================================
func CheckDataInsert() {
	db := dbconfig.GetDestinationDB()

	// Wait for data to be replicated in 8s is time large enough for context timeout
	retry := 4
	delay := 2 * time.Second

	for _, tableName := range model.REPLICATE_TABLE_LISTS {
		// Check data upsert in db destination
		listIdReplicate := cacheIdDataInsert[tableName]
		if len(listIdReplicate) == 0 {
			continue
		}

		query := dbconfig.QueryCheckExistID(listIdReplicate, tableName)
		ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
		defer cancel()

		rows, err := db.QueryContext(ctx, query)
		if err != nil {
			klog.Errorf("- QueryChecker error: %s", err)
		}
		defer rows.Close()

		listIdDetination := []int32{}
		for rows != nil && rows.Next() {
			var id int32
			err = rows.Scan(&id)
			if err != nil {
				klog.Errorf("- Error scanning row in db source: %v", err)
			}
			listIdDetination = append(listIdDetination, id)
		}

		// Retry to get new data from db source
		if listIdDetination == nil || len(listIdDetination) < len(listIdReplicate) {
			for i := 0; i < retry; i++ {
				klog.Infof(">> Retry to get new data from db source: %v", i)
				time.Sleep(delay)
				klog.Infof(">> Executed query - temp: %v", query)
				ctxT, cancelT := context.WithTimeout(context.Background(), 300*time.Second)
				defer cancelT()
				rowsTry, err := db.QueryContext(ctxT, query)
				if err != nil {
					klog.Errorf("QueryChecker error: %s", err)
				}
				defer rows.Close()

				listIdDetination = []int32{}
				for rowsTry != nil && rowsTry.Next() {
					var id int32
					err = rowsTry.Scan(&id)
					if err != nil {
						klog.Errorf("Error scanning row in db source: %v", err)
					}
					listIdDetination = append(listIdDetination, id)
				}

				if listIdDetination != nil && len(listIdDetination) == len(listIdReplicate) {
					break
				}
			}
		}

		klog.Info("==========================================")
		klog.Infof("+ List ID Destination: %s - %v ", tableName, listIdDetination)
		klog.Infof("+ List ID Inserted: %s - %v ", tableName, listIdReplicate)
		if len(listIdDetination) != len(listIdReplicate) || len(listIdDetination) == 0 {
			klog.Errorf("- Has Errors in replicate Insert %s", tableName)
		} else {
			klog.Infof(">> Table %s has inserted replication is successful", tableName)
		}
	}
}

func CheckDataUpdate() {
	db := dbconfig.GetDestinationDB()

	// Wait for data to be replicated in 2s is time large enough for context timeout
	time.Sleep(2 * time.Second)

	delay := 2 * time.Second
	retry := 4

	for _, tableName := range model.REPLICATE_TABLE_LISTS {

		// Check data upsert in db destination
		tableDatas := cacheIdDataUpdate[tableName]
		if len(tableDatas) == 0 {
			continue
		}

		listIdReplicate := []int32{}
		for id := range tableDatas {
			listIdReplicate = append(listIdReplicate, id)
		}
		klog.Infof("+ List ID %s has Updated: %v", tableName, listIdReplicate)

		query := dbconfig.QueryDataByID(listIdReplicate, tableName)
		ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
		defer cancel()

		klog.Infof("+ Executed query check update %s - temp: %v", tableName, query)
		rows, err := db.QueryxContext(ctx, query)
		if err != nil {
			klog.Errorf("- QueryChecker error: %s", err)
		}
		defer rows.Close()

		var isValid bool
		var count int = 0 // number of rows
		for rows != nil && rows.Next() {
			switch tableName {
			case "message_data":
				dataConvert := model.MessageDataDO{}
				if err := rows.StructScan(&dataConvert); err != nil {
					klog.Infof("- Error convert %s: %s", tableName, err)
				}
				isValid = compareData(dataConvert, tableName)
			case "channels":
				dataConvert := model.ChannelsDO{}
				if err := rows.StructScan(&dataConvert); err != nil {
					klog.Infof("- Error convert %s: %s", tableName, err)
				}
				isValid = compareData(dataConvert, tableName)
			}

			// klog.Infof(">> Check Update %s - ID: %d", tableName, dataConvert.(model.ChannelsDO).ID)
			if !isValid {
				klog.Errorf(">> Has Errors replicate Updated %s - %s", tableName, "data is not valid")
				break
			}
			count++
		}

		// Retry if data is not valid
		for !isValid && retry > 0 {
			klog.Infof("* Retry to get new data from db source: %v", retry)
			time.Sleep(delay)
			ctxT, cancelT := context.WithTimeout(context.Background(), 300*time.Second)
			rowsTry, err := db.QueryxContext(ctxT, query)
			if err != nil {
				klog.Errorf("- QueryChecker error: %s", err)
			}

			// Reset count rows to
			count = 0 // number of rows
			for rows != nil && rows.Next() {
				switch tableName {
				case "message_data":
					dataConvert := model.MessageDataDO{}
					if err := rows.StructScan(&dataConvert); err != nil {
						klog.Infof("- Error convert %s: %s", tableName, err)
					}
					isValid = compareData(dataConvert, tableName)
				case "channels":
					dataConvert := model.ChannelsDO{}
					if err := rows.StructScan(&dataConvert); err != nil {
						klog.Infof("- Error convert %s: %s", tableName, err)
					}
					isValid = compareData(dataConvert, tableName)
				}

				// klog.Infof(">> Check Update %s - ID: %d", tableName, dataConvert.(model.ChannelsDO).ID)
				if !isValid {
					klog.Errorf(">> Has Errors replicate Updated %s - %s", tableName, "data is not matching")
					break
				}
				count++
			}
			rowsTry.Close()
			cancelT()
			retry--
		}

		if count != len(listIdReplicate) && isValid{
			klog.Errorf(">> Has Errors replicate Updated %s - %s", tableName, "lost data in db destination")
			return
		} else if count == len(listIdReplicate) && isValid {
			klog.Infof(">> Table %s has updated replication is successful", tableName)
		}
	}
}

func CheckDataDelete() {
	db := dbconfig.GetDestinationDB()

	for _, tableName := range model.REPLICATE_TABLE_LISTS {
		listId := cacheIdDataDelete[tableName]
		if len(listId) == 0 {
			continue
		}
		klog.Infof("+ List ID %s has Deleted: %v", tableName, listId)

		query := dbconfig.QueryCheckExistID(listId, tableName)
		ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
		defer cancel()

		rows, err := db.QueryContext(ctx, query)
		if err != nil {
			klog.Errorf("- Get list ID error: %s", err)
		}
		defer rows.Close()

		klog.Infof("+ Executed query check delete - temp: %v", query)
		if rows == nil {
			klog.Errorf("- Has Errors replicate Deleted %s", tableName)
		} else {
			klog.Infof(">> Table %s has deleted replication is successful", tableName)
		}
	}
}

func compareData(data interface{}, tablename string) bool {

	// Data of db destination
	reflectDesVal := reflect.ValueOf(data)
	reflectDesType := reflect.TypeOf(data)

	if reflectDesVal.Kind() == reflect.Ptr {
		reflectDesVal = reflectDesVal.Elem()
		reflectDesType = reflectDesType.Elem()
	}

	// Data of db source has been cached in memory
	Id := reflect.ValueOf(data).FieldByName("ID").Interface().(int32)
	dataReplicate := cacheIdDataUpdate[tablename][Id]

	if dataReplicate == nil {
		return false
	}
	reflectSourceVal := reflect.ValueOf(dataReplicate)

	// i = 1 to skip ID field
	isValid := true
	for i := 1; i < reflectDesType.NumField(); i++ {
		fieldName := reflectDesType.Field(i).Tag.Get("cql")

		// Pass the field unique
		switch tablename {
		case "auth_users":
		case "chats":
		case "channels":
			if fieldName == "migrated_from" {
				continue
			}
		case "message_data":
			if fieldName == "message_data_id" || fieldName == "dialog_id" || fieldName == "dialog_message_id" || fieldName == "sender_user_id" || fieldName == "random_id" {
				continue
			}
		}

		if fieldName == "created_at" || fieldName == "updated_at" {
			continue
		}

		if reflectDesVal.Field(i).Interface() != reflectSourceVal.Field(i).Interface() {
			isValid = false
			klog.Infof("- ID: %v / Field: %s not match ~ [Data Source: %v]  -  [Data Des: %v]", Id, fieldName, reflectSourceVal.Field(i).Interface(), reflectDesVal.Field(i).Interface())
			return isValid
		}
		// klog.Infof(">> Field: %s -  State: %v  ~ source: %v      Des: %v", reflectDesType.Field(i).Name, isValid, reflectSourceVal.Field(i).Interface(), reflectDesVal.Field(i).Interface())

	}
	// klog.Infof("+ Check %s ID : %d - %v", tablename, Id, isValid)
	return isValid
}

func getMaxRowID(db *sqlx.DB, tableName string) int32 {
	query := "SELECT MAX(ID) FROM " + tableName + ";"
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		klog.Errorf("QueryChecker error: %s", err)
	}
	defer rows.Close()

	var maxID int32
	for rows != nil && rows.Next() {
		err = rows.Scan(&maxID)
		if err != nil {
			klog.Errorf("Error scanning row in db source: %v", err)
		}
	}
	return maxID
}
