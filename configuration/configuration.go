package configuration

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"github.com/hhthuongbtr/tulc-api/utils"
)


func (c *Conf) LoadConf() *Conf {
	fileName := "/opt/tulc/application.yml"
	if c.ConfigureFile != "" {
		fileName = c.ConfigureFile
	}

	confFile := utils.MyFile{Path:fileName}
	yamlFile, err := confFile.Read()
	if err != nil {
		log.Printf("yamlFile.Get err #%v ", err)
		panic(err)
	}
	err = yaml.Unmarshal([]byte(yamlFile), c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	return c
}

//LoadConfFrom load conf from filename(para)
func LoadConfFromFile(filename string) error {
	c := &Conf{}
	if yamlFile, err := ioutil.ReadFile(filename); err == nil {
		if err := yaml.Unmarshal(yamlFile, c); err != nil {
			return err
		}
		return nil
	} else {
		return err
	}
}

