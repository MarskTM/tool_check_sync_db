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

	workLoad := 10000
	numWorkers := 2
	workPool := repository.NewPool(numWorkers, workLoad)

	// Declear Job insert or Update data in db destination
	for i := 0; i < workLoad; i++ {
		workPool.AddJob(*repository.NewJob(InsertData("channels")))
		// workPool.AddJob(*repository.NewJob(UpdateData("channels")))
		// workPool.AddJob(*repository.NewJob(DeleteData("channels")))
	}

	// Start the work pool
	workPool.Wait.Add(len(workPool.WorkList))
	go workPool.Listener()

	workPool.Start()
	workPool.Wait.Wait()

	klog.Infof("================================================================")
	klog.Infof("+ WorkerPool done!")
	for _, job := range workPool.WorkList {
		if job.Err != nil {
			klog.Errorf("- Error: %v", job.Err)
		}
	}

	klog.Infof("================================================================")
	klog.Infof("Start check data in db destination")

	// Start check data in db destination
	CheckDataInsert()
	// CheckDataUpdate()
	CheckDataDelete()

}

func InsertData(tableName string) func() error {
	klog.Infof(">> Insert data to table: %s", tableName)
	db := dbconfig.GetSourceDB()
	query := ""

	return func() error {
		var mu sync.Mutex

		switch tableName {
		case "channels":
			newData := utils.GenerateChannelsDO()
			newData.ID = 0 // make sure ID is auto increment
			query = dbconfig.QueryBuilder(newData, tableName)
		}

		ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
		defer cancel()

		newRow, err := db.ExecContext(ctx, query)
		if err != nil {
			klog.Errorf("Error creating new row: %v", err)
			return err
		}

		newRowId, err := newRow.LastInsertId()
		if err != nil {
			klog.Errorf("Error creating new row: %v", err)
			return err
		}
		// klog.Infof(">> Executed query - temp: %v", newRowId)

		mu.Lock()
		cacheIdDataInsert[tableName] = append(cacheIdDataInsert[tableName], int32(newRowId))
		mu.Unlock()

		return nil
	}
}

func UpdateData(tableName string) func() error {
	db := dbconfig.GetSourceDB()
	query := ""
	var mu sync.Mutex

	return func() error {
		switch tableName {
		case "channels":
			newData := utils.GenerateChannelsDO()
			query = dbconfig.QueryBuilder(newData, tableName)

			mu.Lock()
			id := newData.ID
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

		query := dbconfig.QueryCheckExistID(listIdReplicate, tableName)
		ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
		defer cancel()

		rows, err := db.QueryContext(ctx, query)
		if err != nil {
			klog.Errorf("QueryChecker error: %s", err)
		}
		defer rows.Close()

		listIdDetination := []int32{}
		for rows != nil && rows.Next() {
			var id int32
			err = rows.Scan(&id)
			if err != nil {
				klog.Errorf("Error scanning row in db source: %v", err)
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
		// klog.Info(">> List ID Destination: ", listIdDetination)
		klog.Info(">> List ID Insert: ", listIdReplicate)
		if len(listIdDetination) != len(listIdReplicate) || len(listIdDetination) == 0 {
			klog.Errorf(">> Has Errors in replicate Insert")
		} else {
			klog.Infof(">> Insert replication is successful")
		}
	}
}

func CheckDataUpdate() {
	db := dbconfig.GetDestinationDB()

	// Wait for data to be replicated in 8s is time large enough for context timeout
	retry := 4
	delay := 2 * time.Second

	for _, tableName := range model.REPLICATE_TABLE_LISTS {
		// Check data upsert in db destination
		tableDatas := cacheIdDataUpdate[tableName]

		listIdReplicate := []int32{}
		for id := range tableDatas {
			listIdReplicate = append(listIdReplicate, id)
		}
		klog.Info(">> List ID Update: ", listIdReplicate)

		query := dbconfig.QueryDataByID(listIdReplicate, tableName)
		ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
		defer cancel()

		// klog.Infof(">> Executed query - temp: %v", query)
		rows, err := db.QueryxContext(ctx, query)
		if err != nil {
			klog.Errorf("QueryChecker error: %s", err)
		}
		defer rows.Close()

		// Retry to get new data from db source
		if rows == nil {
			for i := 0; i < retry; i++ {
				klog.Infof(">> Retry to get new data from db source: %v", i)
				time.Sleep(delay)
				klog.Infof(">> Executed query - temp: %v", query)
				ctxT, cancelT := context.WithTimeout(context.Background(), 300*time.Second)
				defer cancelT()
				rowsTry, err := db.QueryxContext(ctxT, query)
				if err != nil {
					klog.Errorf("QueryChecker error: %s", err)
				}
				defer rows.Close()

				if rowsTry != nil {
					rows = rowsTry
					break
				}
			}
		}

		//
		var isValid bool
		for rows != nil && rows.Next() {
			dataConvert := model.MapModel[tableName]
			if err := rows.StructScan(dataConvert); err != nil {
				klog.Infof("err convert: %s - %s", tableName, err)
			}

			klog.Infof(">> Data ID: %d", dataConvert.(model.ChannelsDO).ID)
			isValid = compareData(dataConvert, tableName)

			if !isValid {
				klog.Errorf(">> Has Errors in replicate Update")
				break
			}
		}
		klog.Infof(">> Update replication is successful")
	}
}

func CheckDataDelete() {
	db := dbconfig.GetDestinationDB()

	for _, tableName := range model.REPLICATE_TABLE_LISTS {
		listId := cacheIdDataDelete[tableName]
		klog.Info(">> List ID Delete: ", listId)

		query := dbconfig.QueryCheckExistID(listId, tableName)
		ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
		defer cancel()

		rows, err := db.QueryContext(ctx, query)
		if err != nil {
			klog.Errorf("Get list ID error: %s", err)
		}
		defer rows.Close()

		// klog.Info(">> Executed query - temp: %v", query)
		if rows == nil {
			klog.Errorf(">> Has Errors in replicate")
		} else {
			klog.Infof(">> Deleting replication is successful")
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
	reflectSourceVal := reflect.ValueOf(dataReplicate)

	// i = 1 to skip ID field
	isValid := false
	for i := 1; i < reflectDesType.NumField(); i++ {
		if reflectDesVal.Field(i) == reflectSourceVal.Field(i) {
			isValid = true
		} else {
			isValid = false
		}
		klog.Infof(">> Field: %s - %v", reflectDesType.Field(i).Name, isValid)
	}

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
