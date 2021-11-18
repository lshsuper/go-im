package models


type  SendToGroupRequest struct {

	GroupID string  `json:"group_id"`
	Msg string `json:"msg"`

}


type GroupMsgDTO struct {

	UserID string  `json:"user_id"`
	Msg string `json:"msg"`

}

type GroupMemberCountRequest struct {
	GroupID string  `json:"group_id"`

}
