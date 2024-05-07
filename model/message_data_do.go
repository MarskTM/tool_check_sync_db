// kingtalk team
// data convert from mysql -> mysql big data tool
// thanhnt
package model

import "time"

// MessageDataDO type
type MessageDataDO struct {
	ID              int32     `db:"id" cql:"id"`
	MessageDataID   int64     `db:"message_data_id" cql:"message_data_id"`
	DialogID        int64     `db:"dialog_id" cql:"dialog_id"`
	DialogMessageID int32     `db:"dialog_message_id" cql:"dialog_message_id"`
	SenderUserID    int32     `db:"sender_user_id" cql:"sender_user_id"`
	PeerType        int8      `db:"peer_type" cql:"peer_type"`
	PeerID          int32     `db:"peer_id" cql:"peer_id"`
	RandomID        int64     `db:"random_id" cql:"random_id"`
	MessageType     int8      `db:"message_type" cql:"message_type"`
	MessageData     string    `db:"message_data" cql:"message_data"`
	MediaUnread     int8      `db:"media_unread" cql:"media_unread"`
	Views           int32     `db:"views" cql:"views"`
	HasMediaUnread  int8      `db:"has_media_unread" cql:"has_media_unread"`
	Date            int32     `db:"date" cql:"date"`
	EditMessage     string    `db:"edit_message" cql:"edit_message"`
	EditDate        int32     `db:"edit_date" cql:"edit_date"`
	IsReaction 		int8 	  `db:"is_reaction" cql:"is_reaction"`
	IsPinned  	 	int8 	  `db:"is_pinned" cql:"is_pinned"`

	Deleted         int8      `db:"deleted" cql:"deleted"`
	CreatedAt       time.Time `db:"created_at" cql:"created_at"`
	UpdatedAt       time.Time `db:"updated_at" cql:"updated_at"`
}
