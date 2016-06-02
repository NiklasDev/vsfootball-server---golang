package models

import (
//		"fmt"
	"strconv"
	"strings"
)

type DeviceRegister struct {
	Register IDbDeviceRegister
}
type IDbDeviceRegister interface {
	getTableIdentifyName() string
	getTableName() string
	AddDeviceIdentity(identify string) (bool, int64)
	getUserTableDevItemName() string
}
type AndroidDeviceRegister struct {
}
type IosDeviceRegister struct {
}

func (d *AndroidDeviceRegister) getTableIdentifyName() string {
	return "Deviceid"
}
func (d *AndroidDeviceRegister) getTableName() string {
	return "androiddevice"
}
func (d *AndroidDeviceRegister) getUserTableDevItemName() string {
	return "Androiddevices"
}
func (d *AndroidDeviceRegister) AddDeviceIdentity(identify string) (bool, int64) {
	dbmap := dbPool.GetConnection()
	defer dbPool.ReleaseConnection(dbmap)
	newDevice := &Androiddevice{
		Deviceid: identify}
	err := dbmap.Insert(newDevice)
	if nil == err {
		return true, newDevice.Id
	}
	return false, newDevice.Id
}
func (d *IosDeviceRegister) getTableIdentifyName() string {
	return "Devicetoken"
}
func (d *IosDeviceRegister) getTableName() string {
	return "iosdevice"
}
func (d *IosDeviceRegister) getUserTableDevItemName() string {
	return "Iosdevices"
}
func (d *IosDeviceRegister) AddDeviceIdentity(identify string) (bool, int64) {
	dbmap := dbPool.GetConnection()
	defer dbPool.ReleaseConnection(dbmap)
	newDevice := &Iosdevice{
		Devicetoken: identify}
	err := dbmap.Insert(newDevice)
	if nil == err {
		return true, newDevice.Id
	}
	return false, newDevice.Id
}

func (r *DeviceRegister) RegisteDevice(guid, identity string) (string, string, int) {
	dbmap := dbPool.GetConnection()
	defer dbPool.ReleaseConnection(dbmap)
	selectUser := "select Guid," + r.Register.getUserTableDevItemName() + " as Devices from user where user.Guid=? limit 1"

	row := db.QueryRow(selectUser, guid)
	var Guid, Devices string
	if err := row.Scan(&Guid, &Devices); err != nil {
		//fmt.Println(err)
		return "false", "User not found.", 401
	}
	//remove the old device
	var oldDeviceId = r.GetDeviceIdByIdentity(identity)
	if oldDeviceId > -1 {
		if RemoveDeviceFromStr(Devices,oldDeviceId)==Devices{
			r.RemoveDeviceFromUserTable(oldDeviceId)
			r.RemoveDeviceIdentity(identity)
		}else{
			return "true", "Already registered.", 200
		}
	}
	//add device to new user
	addADTRes, newDeviceId := r.Register.AddDeviceIdentity(identity)
	if addADTRes {
		newDevices := RemoveDeviceFromStr(Devices, oldDeviceId)
		newDevices += "," + strconv.FormatInt(newDeviceId, 10)
		//fmt.Printf("%s:new devices:%s,%d,%d",Devices, newDevices,oldDeviceId,newDeviceId)
		if r.UpdateDevices(Guid, newDevices) {
			return "true", "Successfully registered.", 200
		}
	}
	return "false", "Database failure.", 501
}

func (r *DeviceRegister) RemoveDeviceIdentity(identify string) bool {
	dbmap := dbPool.GetConnection()
	defer dbPool.ReleaseConnection(dbmap)
	sqlRemoveDev := "delete FROM " + r.Register.getTableName() + " where " + r.Register.getTableIdentifyName() + " = ?"
	_, removeDevError := dbmap.Exec(sqlRemoveDev, identify)
	if nil == removeDevError {
		return true
	}
	return false
}

func (r *DeviceRegister) RemoveDeviceFromUserTable(deviceId int64) bool {
	dbmap := dbPool.GetConnection()
	defer dbPool.ReleaseConnection(dbmap)
	if deviceId > 0 {
		var deviceIdstr = strconv.FormatInt(deviceId, 10)
		sqlQueryDev := "SELECT Guid," + r.Register.getUserTableDevItemName() + " as Devices FROM user where user." + r.Register.getUserTableDevItemName() + " like '%," + deviceIdstr + ",%' or user." + r.Register.getUserTableDevItemName() +" like '%," + deviceIdstr + "' limit 1"
		rows, err := db.Query(sqlQueryDev)
		if nil == err {
			defer rows.Close()
			for rows.Next() {
				var Guid, Devices string
				if err := rows.Scan(&Guid, &Devices); err != nil {
					return false
				}
				removedDevs := RemoveDeviceFromStr(Devices, deviceId)
				//update to db
				r.UpdateDevices(Guid, removedDevs)
			}
			return true
		}
	}
	return false
}

/*
Clear not existed device from user
*/
func (r *DeviceRegister) Clear(){
	dbmap := dbPool.GetConnection()
	defer dbPool.ReleaseConnection(dbmap)
	
	sqlQuery := "select guid,"+r.Register.getUserTableDevItemName()+" from user"
	rows, err := db.Query(sqlQuery)
	if nil == err {
		defer rows.Close()
		for rows.Next() {
			var Guid, Devices string
			err := rows.Scan(&Guid, &Devices);
			if  err == nil {
				devicesbuff := Devices
				devids := strings.Split(Devices,",")
				for i:=0;i<len(devids);i++{
					deviceId,converr := strconv.ParseInt(devids[i],10,64)
					if converr == nil && r.GetDeviceIdentityById(deviceId)==""{
						Devices = RemoveDeviceFromStr(Devices, deviceId)
					}
				}
				if devicesbuff != Devices{
					//update to db
					r.UpdateDevices(Guid, Devices)
				}
			}
			
		}
	}
}

func (r *DeviceRegister) GetDeviceIdentityById(id int64) string {
	dbmap := dbPool.GetConnection()
	defer dbPool.ReleaseConnection(dbmap)
	sqlQueryId := "SELECT "+r.Register.getTableIdentifyName()+" FROM " + r.Register.getTableName() + " where id=? limit 1"
	var devids []*string
	_, selectDevError := dbmap.Select(&devids, sqlQueryId, id)
	if nil == selectDevError && len(devids) > 0 {
		return *devids[0]
	}
	return ""
}

func (r *DeviceRegister) UpdateDevices(guid string, devices string) bool {
	dbmap := dbPool.GetConnection()
	defer dbPool.ReleaseConnection(dbmap)
	updateQuery := "update user set user." + r.Register.getUserTableDevItemName() + "=? where user.Guid=?"
	_, updateFailure := dbmap.Exec(updateQuery, devices, guid)
	if updateFailure != nil {
		return false
	} else {
		return true
	}
}

func (r *DeviceRegister) GetDeviceIdByIdentity(identity string) int64 {
	dbmap := dbPool.GetConnection()
	defer dbPool.ReleaseConnection(dbmap)
	sqlQueryId := "SELECT id FROM " + r.Register.getTableName() + " where " + r.Register.getTableIdentifyName() + "=? limit 1"
	var devids []*int64
	_, selectDevError := dbmap.Select(&devids, sqlQueryId, identity)
	if nil == selectDevError && len(devids) > 0 {
		return *devids[0]
	}
	return -1
}

func RemoveDeviceFromStr(devstr string, devid int64) string {
	rmvDeviceIdStr := "," + strconv.FormatInt(devid, 10)
	newDevices := strings.Replace(devstr, rmvDeviceIdStr+",", ",", -1)
	newDevices = strings.Replace(newDevices, ",,", ",", -1)
	buffOffset := len(newDevices) - len(rmvDeviceIdStr)
	if buffOffset > -1 {
		buffcomp := string([]byte(newDevices)[buffOffset:])
		if buffcomp == rmvDeviceIdStr {
			newDevices = string([]byte(newDevices)[0:buffOffset])
		}
	}
	return newDevices
}

func GetUserByGuid(guid string) *User {
	dbmap := dbPool.GetConnection()
	defer dbPool.ReleaseConnection(dbmap)
	selectUser := "select * from user where user.Guid=? limit 1"
	var userResult []*User
	_, selectUserError := dbmap.Select(&userResult, selectUser, guid)
	if selectUserError == nil && len(userResult) > 0 {
		return userResult[0]
	}
	return nil
}
