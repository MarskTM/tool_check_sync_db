package utils

import (
	"math/rand"
	"test-data-convert/model"
	"time"
)

// Set seed for random number
func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandomInt(min, max int) int {
	return min + rand.Intn(max-min)
}

func RandomString(length int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	b := make([]rune, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func GenerateAuthUsersDO() model.AuthUsersDO {
	return model.AuthUsersDO{
		ID:            int32(RandomInt(1, 1000)),
		AuthID:        int64(RandomInt(1, 100)),
		VivaSessionID: RandomString(10),
		UserID:        int32(RandomInt(1, 100)),
		Hash:          int64(RandomInt(1, 100)),
		Layer:         int32(RandomInt(1, 100)),
		DeviceModel:   RandomString(10),
		Platform:      RandomString(10),
		SystemVersion: RandomString(10),
	}
}

func GenerateChatsDO() model.ChatsDO {
	return model.ChatsDO{
		ID:               int32(RandomInt(1, 1000)),
		CreatorUserID:    int32(RandomInt(1, 100)),
		AccessHash:       int64(RandomInt(1, 100)),
		RandomID:         int64(RandomInt(1, 100)),
		ParticipantCount: int32(RandomInt(1, 100)),
		Title:            RandomString(10),
		About:            RandomString(10),
		Link:             RandomString(10),
		PhotoID:          int64(RandomInt(1, 1000)),
		AdminsEnabled:    int8(RandomInt(1, 100)),
		MigratedTo:       int32(RandomInt(1, 100)),
		Deactivated:      int8(RandomInt(1, 100)),
		Version:          int32(RandomInt(1, 100)),
		Date:             int32(RandomInt(1, 100)),
	}
}

func GenerateChannelsDO() model.ChannelsDO {
	return model.ChannelsDO{
		ID:               int32(RandomInt(1, 100)),
		CreatorUserID:    int32(RandomInt(1, 100)),
		AccessHash:       int64(RandomInt(1, 100)),
		RandomID:         int64(RandomInt(1, 100)),
		Type:             int32(RandomInt(1, 100)),
		TopMessage:       int32(RandomInt(1, 100)),
		ParticipantCount: int32(RandomInt(1, 100)),
		Title:            RandomString(10),
		About:            RandomString(10),
		PhotoID:          int64(RandomInt(1, 100)),
		Public:           int8(RandomInt(1, 100)),
		Link:             RandomString(10),
		Broadcast:        int8(RandomInt(1, 100)),
		Verified:         int8(RandomInt(1, 100)),
		MegaGroup:        int8(RandomInt(1, 100)),
		Democracy:        int8(RandomInt(1, 100)),
		Signatures:       int8(RandomInt(1, 100)),
		AdminsEnabled:    int8(RandomInt(1, 100)),
		Deactivated:      int8(RandomInt(1, 100)),
		Version:          int32(RandomInt(1, 100)),
		Date:             int32(RandomInt(1, 100)),
		MigratedFrom:    int32(RandomInt(1, 100)),
	}
}

func GenerateMessageDataDO() model.MessageDataDO {
	return model.MessageDataDO{
		ID:              int32(RandomInt(200, 210)),
		MessageDataID:   int64(RandomInt(100, 2000)),
		DialogID:        int64(RandomInt(100, 2000)),
		DialogMessageID: int32(RandomInt(100, 2000)),
		SenderUserID:    int32(RandomInt(100, 2000)),
		PeerType:        int8(RandomInt(100, 200)),
		PeerID:          int32(RandomInt(100, 2000)),
		RandomID:        int64(RandomInt(100, 2000)),
		MessageType:     int8(RandomInt(1, 100)),
		MessageData:     RandomString(10),
		MediaUnread:     int8(RandomInt(1, 100)),
		Views:           int32(RandomInt(1, 100)),
		HasMediaUnread:  int8(RandomInt(1, 100)),
		Date:            int32(RandomInt(1, 100)),
		EditMessage:     RandomString(10),
		EditDate:        int32(RandomInt(1, 100)),
		IsReaction:      int8(RandomInt(1, 100)),
		IsPinned:        int8(RandomInt(1, 100)),
	}
}