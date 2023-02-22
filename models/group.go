package models

import (
	"encoding/json"
	"os"
	"time"

	"github.com/Massad/gin-boilerplate/db"
)

//UserModel ...
type GroupModel struct{}

type GroupMember struct {
	Id          string
	IsAdmin     bool
	JoiningDate time.Time
}

type GroupDetails struct {
	Id                      string `db:"group_id, primarykey, autoincrement" json:"Id"`
	Name                    string `db:"name" json:"Name"`
	Icon                    string `db:"icon" json:"Icon"`
	CreatedDate, LastActive time.Time
	Members                 []GroupMember
}

func (m GroupModel) CheckgroupExists(groupId string) (exists bool) {

	exists = false

	_, err := db.GetRedis().Get(groupId).Result()

	if err == nil {

		exists = true

	}

	return exists

}

func (m GroupModel) GetgroupDetails(groupId string) (groupDetailsa GroupDetails, error error) {

	groupDetailsSerialized, err := db.GetRedis().Get(groupId).Result()

	if err != nil {
		return groupDetailsa, err
	}

	json.Unmarshal([]byte(groupDetailsSerialized), &groupDetailsa)

	return groupDetailsa, nil

}

func (m GroupModel) AddGroupDetails(groupId string, groupData GroupDetails) (inserted bool) {

	inserted = false

	groupDataSerialized, _ := json.Marshal(groupData)

	err := db.GetRedis().Set(groupId, groupDataSerialized, 0).Err()

	if err == nil {

		inserted = true

	}

	return inserted
}

func (m GroupModel) DeleteGroup(groupId string) (deleted bool) {

	deleted = false

	err := db.GetRedis().Del(groupId).Err()

	if err == nil {

		deleted = true

	}

	return deleted
}

func (m GroupModel) GroupsList(u Conversations) []GroupDetails {

	var list []GroupDetails
	var columnGroupPrefix = os.Getenv("COLUMN_GROUP_PREFIX")

	for _, user := range u {

		groupId := columnGroupPrefix + user.GroupID

		groupDetailsa, _ := m.GetgroupDetails(groupId)

		list = append(list, groupDetailsa)
	}

	return list
}
