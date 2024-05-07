/* !!
 * File: auth_users_do.go
 * File Created: Friday, 22nd November 2019 5:02:53 pm
 * Author: KimEricko™ (phamkim.pr@gmail.com)
 * -----
 * Last Modified: Monday, 2nd December 2019 11:57:02 am
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

// AuthUsersDO type
type AuthUsersDO struct {
	ID            int32     `db:"id" cql:"id"`
	AuthID        int64     `db:"auth_key_id" cql:"auth_key_id"`
	VivaSessionID string    `db:"viva_session_id" cql:"viva_session_id"`
	UserID        int32     `db:"user_id" cql:"user_id"`
	Hash          int64     `db:"hash" cql:"hash"`
	Layer         int32     `db:"layer" cql:"layer"`
	DeviceModel   string    `db:"device_model" cql:"device_model"`
	Platform      string    `db:"platform" cql:"platform"`
	SystemVersion string    `db:"system_version" cql:"system_version"`
	ApiID         int32     `db:"api_id" cql:"api_id"`
	AppName       string    `db:"app_name" cql:"app_name"`
	AppVersion    string    `db:"app_version" cql:"app_version"`
	DateCreated   int32     `db:"date_created" cql:"date_created"`
	DateActived   int32     `db:"date_actived" cql:"date_actived"`
	IP            string    `db:"ip" cql:"ip"`
	Country       string    `db:"country" cql:"country"`
	Region        string    `db:"region" cql:"region"`
	IsDeleted     int8      `db:"deleted" cql:"deleted"`
	CreatedAt     time.Time `db:"created_at" cql:"created_at"`
	UpdatedAt     time.Time `db:"updated_at" cql:"updated_at"`
}
