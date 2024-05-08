package model

var MapModel = map[string]interface{}{
	"channels": ChannelsDO{},
}

var MapListModel = map[string]interface{}{
	// "auth_users":      AuthUsersDO{},
	// "chats":           ChatsDO{},
	"channels": []ChannelsDO{},
	// "message_data":    MessageDataDO{},
}

var REPLICATE_TABLE_LISTS = []string{
	// "auth_users",
	// "chats",
	"channels",
	// "message_data",
}
