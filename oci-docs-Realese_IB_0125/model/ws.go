package model

type WSMessage struct {
	Name string      `json:"name"`
	Data interface{} `json:"data"`
}

type WaitRegFinish struct {
	Phone string
}

type DBEvent struct {
	Table  string `json:"table"`
	Action string `json:"action"`
	Data   struct {
		ID     int    `json:"id"`
		ChatID int    `json:"chat_id"`
		Phone  string `json:"phone"`
		FIO    string `json:"fio"`
	} `json:"data"`
}
