package main

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"test-data-convert/dbconfig"
	"test-data-convert/model"
	"time"

	// "test-data-convert/model"
	"test-data-convert/repository"
	"test-data-convert/utils"

	"k8s.io/klog/v2"
)

var cacheIdTData = map[string][]int32{}

var cacheIdChannels = []int32{}

func main() {

	workLoad := 10
	numWorkers := 2
	workPool := repository.NewPool(numWorkers, workLoad)

	// Declear Job insert or Update data in db destination
	for i := 0; i < workLoad; i++ {
		workPool.AddJob(*repository.NewJob(UpsertData("ChannelsDO")))
	}

	// Start the work pool
	workPool.Wait.Add(workLoad)
	go workPool.Listener()

	workPool.Start()
	workPool.Wait.Wait()
	klog.Infof(">>> Done!")

	// Listener for data check
	TestDataUpsert()
	// TestDataDelete()
}

func UpsertData(modelName string) func() error {
	db := dbconfig.GetSourceDB()
	query := ""

	return func() error {
		var mu sync.Mutex
		switch modelName {
		case "ChannelsDO":
			newData := utils.GenerateChannelsDO()
			query = dbconfig.QueryBuilder(newData, "channels")
			mu.Lock()
			cacheIdTData[modelName] = append(cacheIdTData[modelName], newData.ID)
			mu.Unlock()
			klog.Infof(">> Executed query id = : %v", newData.ID)
		default:
			return errors.New("Model name not found")
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

func TestDataUpsert() {
	db := dbconfig.GetDestinationDB()
	for _, tableName := range model.REPLICATE_TABLE_LISTS {
		// Check data upsert in db destination
		listIdReplicate := cacheIdTData[tableName]
		query := dbconfig.QueryChecker(listIdReplicate, tableName)
		ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
		defer cancel()

		// klog.Infof(">> Executed query - temp: %v", query)
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

		klog.Info(">> List ID Destination: ", listIdDetination)
		klog.Info(">> List ID Replicate: ", listIdReplicate)
		if len(listIdDetination) != len(listIdReplicate) || len(listIdDetination) == 0 {
			klog.Errorf(">> Has Errors in replicate")
		} else {
			klog.Infof(">> Upserting replication is successful")
		}
	}
}

func TestDataDelete() {
	dbSource := dbconfig.GetSourceDB()
	dbDes := dbconfig.GetDestinationDB()

	for _, tableName := range model.REPLICATE_TABLE_LISTS {
		queryFromSource := "SELECT ID FROM " + tableName + " ;"
		ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
		defer cancel()
		rows, err := dbSource.QueryContext(ctx, queryFromSource)
		if err != nil {
			klog.Errorf("Select ID From Source error: %s", err)
		}

		listID := []string{}
		for rows != nil && rows.Next() {
			var id int32
			err = rows.Scan(&id)
			if err != nil {
				klog.Errorf("Error scanning row in db source: %v", err)
			}
			listID = append(listID, fmt.Sprintf("%v", id))
		}

		queryFromDestination := "SELECT ID FROM " + tableName + "WHERE ID NOT IN (" + strings.Join(listID, ",") + ");"
		ctxV2, cancelV2 := context.WithTimeout(context.Background(), 300*time.Second)
		defer cancelV2()
		rowV2s, err := dbDes.QueryContext(ctxV2, queryFromDestination)
		if err != nil {
			klog.Errorf("Select ID From Destination error: %s", err)
		}

		listIDV2 := []string{}
		for rowV2s != nil && rowV2s.Next() {
			var id int32
			err = rowV2s.Scan(&id)
			if err != nil {
				klog.Errorf("Error scanning row in db source: %v", err)
			}
			listIDV2 = append(listIDV2, fmt.Sprintf("%v", id))
		}

		if len(listIDV2) > 0 {
			klog.Errorf(">> Has Errors in replicate")
		} else {
			klog.Infof(">> Deleting replication is successful")
		}
	}
}

func getTableName(payload interface{}) string {
	switch payload.(type) {
	case model.AuthUsersDO:
		return "auth_users"
	case model.ChatsDO:
		return "chats"
	case model.ChannelsDO:
		return "channels"
	case model.MessageDataDO:
		return "message_data"
	default:
		return ""
	}
}
