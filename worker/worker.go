package worker

import (
	"github.com/hhthuongbtr/tulc-api/model"
	"github.com/hhthuongbtr/tulc-api/utils"
	"encoding/json"
	"fmt"
	"sync"
)

type Pool struct {
	Tasks []*Task

	Concurrency int
	TasksChan   chan *Task
	Wg          sync.WaitGroup
	Uri			string
	SecretKey	string
}


func NewPool(tasks []*Task, concurrency int, uri string, secretKey string) *Pool {
	return &Pool{
		Tasks:       tasks,
		Concurrency: concurrency,
		TasksChan:   make(chan *Task),
		Uri:		uri,
		SecretKey: secretKey,
	}
}


func (p *Pool) Run() {
	var worker_num int
	if len(p.Tasks) < p.Concurrency {
		worker_num = len(p.Tasks)
	} else {
		worker_num = p.Concurrency
	}
	//log.Println("Spawn ", worker_num, " worker(s)")
	for i := 0; i <= worker_num; i++ {
		go p.work(i, p.Uri, p.SecretKey)
	}
	p.Wg.Add(len(p.Tasks))
	for _, task := range p.Tasks {
		p.TasksChan <- task
	}

	close(p.TasksChan)

	p.Wg.Wait()
}


func (p *Pool) work(worker_id int, uri string, secretKey string) {
	for task := range p.TasksChan {
		task.Run(&p.Wg, worker_id, uri, secretKey)
	}
}

//--------------------------------------------------------------

type Task struct {
	Err error
	ServerID string
	CCUFromPartner model.ApiServerResponse
}


func NewTask(serverID string) *Task {
	return &Task{ServerID: serverID}
}


func (t *Task) Run(wg *sync.WaitGroup, worker_id int, uri string, secretKey string) {
	//fmt.Println("worker", worker_id, "processing serverid ", t.ServerID)
	t.CCUFromPartner, t.Err = GetCcuInfoByServerID(t.ServerID, uri, secretKey)
	wg.Done()
}



func GetCcuInfoByServerID(serverId string, uri string, secretKey string) (serverCcuInfo model.ApiServerResponse, err error) {
	var ccuInfo model.ApiServerResponse
	ts := utils.GetNowAsUnixTimestamp()
	strToSign := fmt.Sprintf("%s%s%d", secretKey, serverId, ts)
	sig := utils.GetMd5FromString(strToSign)
	url := fmt.Sprintf("%s?server_id=%s&ts=%d&sig=%s", uri, serverId, ts, sig)
	body, err := utils.HttpGet(url)
	if err != nil {
		return ccuInfo, err
	}
	err = json.Unmarshal([]byte(body), &ccuInfo)
	if err != nil {
		return ccuInfo, err
	}
	return ccuInfo, nil
}


