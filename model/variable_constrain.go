package model

var MapModel = map[string]interface{}{
	"channels": ChannelsDO{},
	"auth_users": AuthUsersDO{},
	"chats": ChatsDO{},
    "message_data": MessageDataDO{},
}

var REPLICATE_TABLE_LISTS = []string{
	"auth_users",
	"chats",
	"channels",
	"message_data",
}
