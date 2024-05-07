/*
 *  Copyright (c) 2018, http://103.69.195.249/kimpv/kingTalk
 *  All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

// kingtalk team
// data convert from mysql -> mysql big data tool
// thanhnt
package model

import "time"

// ChannelsDO type
type ChannelsDO struct {
	ID               int32     `db:"id" cql:"id"`
	CreatorUserID    int32     `db:"creator_user_id" cql:"creator_user_id"`
	AccessHash       int64     `db:"access_hash" cql:"access_hash"`
	RandomID         int64     `db:"random_id" cql:"random_id"`
	Type             int32     `db:"type" cql:"type"`
	TopMessage       int32     `db:"top_message" cql:"top_message"`
	ParticipantCount int32     `db:"participant_count" cql:"participant_count"`
	Title            string    `db:"title" cql:"title"`
	About            string    `db:"about" cql:"about"`
	PhotoID          int64     `db:"photo_id" cql:"photo_id"`
	Public           int8      `db:"public" cql:"public"`
	Link             string    `db:"link" cql:"link"`
	Broadcast        int8      `db:"broadcast" cql:"broadcast"`
	Verified         int8      `db:"verified" cql:"verified"`
	MegaGroup        int8      `db:"megagroup" cql:"megagroup"`
	Democracy        int8      `db:"democracy" cql:"democracy"`
	Signatures       int8      `db:"signatures" cql:"signatures"`
	AdminsEnabled    int8      `db:"admins_enabled" cql:"admins_enabled"`
	Deactivated      int8      `db:"deactivated" cql:"deactivated"`
	Version          int32     `db:"version" cql:"version"`
	Date             int32     `db:"date" cql:"date"`
	CreatedAt        time.Time `db:"created_at" cql:"created_at"`
	UpdatedAt        time.Time `db:"updated_at" cql:"updated_at"`

	//fix
	SlowMode            int32 `db:"slow_mode" cql:"slow_mode"`
	DefaultBannedRights int32 `db:"default_banned_rights" cql:"default_banned_rights"`
	MigratedFrom        int32 `db:"migrated_from" cql:"migrated_from"`
	GlobalSearch        int8  `db:"global_search" cql:"global_search"`
	PreHistory          int8  `db:"pre_history" cql:"pre_history"`
	HiddenParticipants  int8  `db:"hidden_participants" cql:"hidden_participants"`
	JoinRequest         int8  `db:"join_request" cql:"join_request"`
	RequestPendings     int32 `db:"request_pendings" cql:"request_pendings"`
	Date2               int32 `db:"date2" cql:"date2"`
}
