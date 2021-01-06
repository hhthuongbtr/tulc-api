package model

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


