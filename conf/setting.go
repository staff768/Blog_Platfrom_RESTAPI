package conf

import (
	"encoding/json"
	"log"
	"os"
)

type setting struct {
	ServerHost string
	ServerPort string
	PgHost     string
	PgPort     string
	PgUser     string
	PgPass     string
	PgBase     string
}

var Cfg setting

func init() {
	file, err := os.Open("conf/setting.cfg")
	if err != nil {
		log.Fatalf("не удалось открыть файл %s", err)
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		log.Fatalf("не удалось прочитать конфигурацию %s", err)
	}

	readByte := make([]byte, stat.Size())
	file.Read(readByte)

	err = json.Unmarshal(readByte, &Cfg)
	if err != nil {
		log.Fatalf("не удалось преобразовать в json %s", err)
	}

}
