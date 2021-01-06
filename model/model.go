package model

import "encoding/json"

type ApiServerResponse struct {
	Count	int `json:"count"`
	ReturnMessage	string `json:"returnMessage"`
	ReturnCode	int	`json:"returnCode"`
}

type CCU struct {
	Serverid	string	`json:"serverid"`
	ServerName	string	`json:"server_name"`
	status		bool	`json:"status"`
	Ccu	int	`json:"ccu"`
}

type CcuResponse struct {
	Unixtime	int64	`json:"unixtime"`
	Groupid	string	`json:"groupid"`
	GroupName	string	`json:"group_name"`
	ServerList	[]CCU	`json:"server_list"`
}

type ServerListResponse struct {
	ReturnMessage	string	`json:"returnMessage" yaml:"returnMessage"`
	ReturnCode	int	`json:"returnCode" yaml:"returnCode"`
	Data map[string]ServerListElement	`json:"data" yaml:"data"`
}

type ServerListElement struct {
	Status	int	`json:"status" yaml:"status"`
	Info struct{
		MergeTargetServerID int `json:"merge_target_server_id" yaml:"merge_target_server_id"`
		OpenTime	int	`json:"open_time" yaml:"open_time"`
	}	`json:"info" yaml:"info"`
	ServerID	string	`json:"serverID" yaml:"serverID"`
	ServerName	string	`json:"serverName" yaml:"serverName"`
}

func (svlrsp *ServerListResponse) GetJsonString() (JsonString string, err error) {
	b, err := json.Marshal(svlrsp)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (svlrsp *ServerListResponse) LoadFromJsonString(JsonString string) (err error) {
	err = json.Unmarshal([]byte(JsonString), svlrsp)
	if err != nil {
		return err
	}
	return
}