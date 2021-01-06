package main

import (
	"github.com/hhthuongbtr/tulc-api/configuration"
	"gopkg.in/yaml.v2"
	"log"
	"github.com/hhthuongbtr/tulc-api/utils"
	"github.com/hhthuongbtr/tulc-api/model"
)

func main()  {
	var conf  configuration.Conf
	conf.LoadConf()

	var severList model.ServerListResponse
	log.Printf("%#v",conf)
	log.Println(conf.Server.ServerListFilePathForStaging)
	svlistFile := utils.MyFile{Path:conf.Server.ServerListFilePathForStaging}
	yamlFile, err := svlistFile.Read()
	if err != nil {
		log.Printf("yamlFile.Get err #%v ", err)
		panic(err)
	}
	err = yaml.Unmarshal([]byte(yamlFile), &severList)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	log.Printf("Serverlist %#v", severList)

	jsonString, err := severList.GetJsonString()
	if err != nil{
		println(err)
	} else {
		println(jsonString)
	}
}

