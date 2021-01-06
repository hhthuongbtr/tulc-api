package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"time"
	"./model"
	"./utils"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"./configuration"
	"./worker"
)

type WebProxy struct {
	Host	string	`json:"host"`
	Port	int16	`json:"port"`
	UriApiFromPartner	string	`json:"uri_api_from_partner"`
	SecretKey	string	`json:"secret_key"`
	ServerlistFilePath	string	`json:"serverlist_file_path"`
	ConcurrencyThread int	`json:"concurrency_thread"`
}

func main()  {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	var conf  configuration.Conf
	conf.LoadConf()
	serverlistFilePtr := flag.String("i", conf.Server.ServerListFilePath, "Serverlist file")
	uriPtr := flag.String("u", conf.PartnerApi.Uri, "Api get ccu url from partner")
	modePtr := flag.String("m", "http", "Run mode: http(default)/urgent")
	serverIdPtr := flag.String("s", "", "ServerID")
	flag.Parse()

	//log.Print(getNowAsUnixTimestamp())
	fmt.Printf(*modePtr)
	var concurrencyThread int
	if conf.PartnerApi.Concurrent == 0 {
		concurrencyThread = 10
	} else {
		concurrencyThread = conf.PartnerApi.Concurrent
	}
	switch *modePtr {
	case "http":
		log.Println("http mode, please wait")
		webContext := WebProxy{
			Host: conf.Server.Host,
			Port: conf.Server.Port,
			UriApiFromPartner: *uriPtr,
			SecretKey: conf.PartnerApi.SecretKey,
			ServerlistFilePath: conf.Server.ServerListFilePath,
			ConcurrencyThread: concurrencyThread,
		}
		server := initializeServer()
		setupRoute(server, &webContext)
		log.Print("begin run http server...")
		listenAdd := fmt.Sprintf("%s:%d", webContext.Host, webContext.Port)
		_ = server.Run(listenAdd)

	case "urgent":
		log.Println("Urgent mode")
		switch *serverIdPtr {
		case "":
			log.Printf("Get all server from server list file %s\n", *serverlistFilePtr)
			ccuResponse, err := GetCCUFromServerListFile(*uriPtr, conf.PartnerApi.SecretKey, conf.Server.ServerListFilePath, concurrencyThread)
			if err != nil {
				log.Print(err)
				return
			}
			responseMessage, err := GetResponeApiJson(ccuResponse)
			if err != nil {
				log.Printf("Error: %s", err)
				return
			}
			log.Printf("Response: %#v", responseMessage)
		default:
			log.Printf("Get ccu from server %s\n", *serverIdPtr)
			ccuInfo, err := worker.GetCcuInfoByServerID(*serverIdPtr, *uriPtr, conf.PartnerApi.SecretKey)
			if err != nil {
				log.Print(err)
				return
			} else {
				ccu := model.CCU{
					Serverid:   *serverIdPtr,
					ServerName: "",
					Ccu:       ccuInfo.Count,
				}
				var ServerListRspCCU []model.CCU
				ServerListRspCCU = append(ServerListRspCCU, ccu)
				ccuResponse := model.CcuResponse{
					Unixtime:   utils.GetNowAsUnixTimestamp(),
					Groupid:    "",
					GroupName:  "",
					ServerList: ServerListRspCCU,
				}
				responseMessage, err := GetResponeApiJson(ccuResponse)
				if err != nil {
					log.Printf("Error: %s", err)
					return
				}
				log.Printf("Response: %#v", responseMessage)
			}
		}
	default:
		log.Fatal("Please choice http mode or urgent mode")
	}

}


func setupRoute(server *gin.Engine, webContext *WebProxy) {
	v1 := server.Group("/api/v1")
	{
		//----------------CCU-------------------
		users := v1.Group("/getCCU")
		{
			users.GET("", webContext.getCCU)
		}

	}
}

func initializeServer() *gin.Engine {
	server := gin.New()
	gin.SetMode(gin.ReleaseMode)
	server.Use(gin.Logger())
	server.Use(gin.Recovery())

	// CORS for https://foo.com and https://github.com origins, allowing:
	// - PUT and PATCH methods
	// - Origin header
	// - Credentials share
	// - Preflight requests cached for 30 seconds
	server.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT"},
		AllowHeaders:     []string{"Access-Control-Allow-Headers", "Access-Control-Allow-Origin", "Origin", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           30 * time.Second,
	}))
	return server
}

func (w *WebProxy) getCCU(ctx *gin.Context) {
	var responseMessage string
	ccuRespone, err := GetCCUFromServerListFile(w.UriApiFromPartner, w.SecretKey, w.ServerlistFilePath, w.ConcurrencyThread)
	if err != nil {
		ctx.String(500, err.Error())
		return
	}
	responseMessage, err2 := GetResponeApiJson(ccuRespone)
	if err2 != nil {
		ctx.String(500, err.Error())
		return
	}
	ctx.String(200, responseMessage)
	return
}


func GetCCUFromServerListFile(uri string, secretKey string, serverlistFile string, concurrencyThread int) (CCUResponse model.CcuResponse, err error) {
	// read serverlist file
	var ServerListRspCCU []model.CCU
	var ccuRespone  model.CcuResponse
	file, err := os.Open(serverlistFile)
	if err != nil {
		log.Fatal(err)
		return ccuRespone, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	if err := scanner.Err(); err != nil {
		return ccuRespone, err
	}
	var tasks []*worker.Task
	for scanner.Scan() {
		serverId := scanner.Text()
		tasks = append(tasks, worker.NewTask(serverId))
	}
	p := worker.NewPool(tasks, concurrencyThread, uri, secretKey)
	p.Run()
	var numErrors int
	var totalTasks int
	var is_too_many_err bool
	for _, task := range p.Tasks {
		totalTasks++
		if task.Err != nil {
			if is_too_many_err != true {
				log.Print(task.Err)
			}
			numErrors++
		} else {
			if numErrors >= 10 {
				if is_too_many_err {
					break
				}
				log.Print("Too many errors.")
				is_too_many_err = true
			}
			ccu := model.CCU{
				Serverid:   task.ServerID,
				ServerName: "",
				Ccu:        task.CCUFromPartner.Count,
			}
			ServerListRspCCU = append(ServerListRspCCU, ccu)
		}
	}
	ccuRespone.Unixtime = utils.GetNowAsUnixTimestamp()
	ccuRespone.ServerList = ServerListRspCCU
	return ccuRespone, nil
}


func GetResponeApiJson(serverRespone model.CcuResponse) (ResponeApiJson string, err error) {
	b, err := json.Marshal(serverRespone)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
