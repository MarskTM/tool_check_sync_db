/* !!
 * File: chats_do.go
 * File Created: Friday, 22nd November 2019 5:02:57 pm
 * Author: KimEricko™ (phamkim.pr@gmail.com)
 * -----
 * Last Modified: Monday, 2nd December 2019 12:09:53 pm
 * Modified By: KimEricko™ (phamkim.pr@gmail.com>)
 * -----
 * Copyright 2018 - 2019 mySoha Platform, VCCorp
 * All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * Developer: NhokCrazy199 (phamkim.pr@gmail.com)
 */

// kingtalk team
// data convert from mysql -> mysql big data tool
// thanhnt
package model

import "time"

// ChatsDO type
type ChatsDO struct {
	ID               int32  `db:"id" cql:"id"`
	CreatorUserID    int32  `db:"creator_user_id" cql:"creator_user_id"`
	AccessHash       int64  `db:"access_hash" cql:"access_hash"`
	RandomID         int64  `db:"random_id" cql:"random_id"`
	ParticipantCount int32  `db:"participant_count" cql:"participant_count"`
	Title            string `db:"title" cql:"title"`
	About            string `db:"about" cql:"about"`
	Link             string `db:"link" cql:"link"` // Thuộc tính này đc lưu trong bảng username
	// ListAvatar       string    `db:"list_avatar" cql:"list_avatar"` // thuộc tính này để lưu avatar quick chat
	PhotoID       int64     `db:"photo_id" cql:"photo_id"`
	AdminsEnabled int8      `db:"admins_enabled" cql:"admins_enabled"`
	MigratedTo    int32     `db:"migrated_to" cql:"migrated_to"`
	Deactivated   int8      `db:"deactivated" cql:"deactivated"`
	Version       int32     `db:"version" cql:"version"`
	Date          int32     `db:"date" cql:"date"`
	CreatedAt     time.Time `db:"created_at" cql:"created_at"`
	UpdatedAt     time.Time `db:"updated_at" cql:"updated_at"`
}
