package model

var MapModel = map[string]interface{}{
	"AuthUsersDO":   AuthUsersDO{},
	"ChatsDO":       ChatsDO{},
	"ChannelsDO":    ChannelsDO{},
	"MessageDataDO": MessageDataDO{},
}

var REPLICATE_TABLE_LISTS = []string{
	// "auth_users",
	// "chats",
	"channels",
	// "message_data",
}
