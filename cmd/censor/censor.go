package censor

import (
	"comments/pkg/db/postgres"
	"encoding/json"
	"regexp"
	"strings"
	"time"
	"io"
	"log"
	"os"
)

var banlist = "banlist.json"

func Censore(db *postgres.DB){
	
	//читаем бан-лист
	list:=readBanList()
	
	//делаем регулярку на все бан-слова
	join := strings.Join(list, "|")

	banRegex := regexp.MustCompile(join)
	
	//получаем все комменты

	var doneid int

	for {
		comments, err:= db.GetAllComments()

		if err!=nil{
			log.Println("Censor couldn't get list of comments from DB")
		}

		for _,c:=range comments{
			
			if (doneid <c.ID){
				if banRegex.MatchString(c.Text) {
					err := db.SetCensored(c.ID)
					if err!=nil{
						log.Println("Censor couldn't set \"censored\" for comment id=", c.ID)
					}
				}
				doneid = c.ID
			}
		}

		time.Sleep(time.Second*5)
	}

}

func readBanList() []string {
	f, err := os.OpenFile(banlist, os.O_RDONLY, 0777)
	if err != nil {
		log.Fatal("Cannot read ban list file: ", err)
	}
	defer f.Close()

	buf, err := io.ReadAll(f)
	if err != nil {
		log.Fatal("Cannot read ban list file: ", err)
	}

	conf := struct {
		List []string
	}{}

	err = json.Unmarshal(buf, &conf)
	if err != nil {
		log.Fatal("Cannot parse ban list : ", err)
	}

	return conf.List
}
