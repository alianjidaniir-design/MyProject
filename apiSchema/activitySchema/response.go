package activitySchema

import "MyProject/models/activity/dataModel"

type InformationActivity struct {
	Massage     string             `json:"massage"`
	Information dataModel.Activity `json:"information"`
}
