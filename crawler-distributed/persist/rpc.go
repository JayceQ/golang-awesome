package persist

import (
	"golang-awesome/crawler/engine"
	"golang-awesome/crawler/persist"
	"gopkg.in/olivere/elastic.v5"
	"log"
)

type ItemSaverService struct{
	Client *elastic.Client
	Index string
}

func (s *ItemSaverService) Save(item engine.Item, result *string) error {

	err := persist.Save(s.Client, s.Index, item)
	log.Printf("save profile %s", item.Url)
	if err == nil {
		*result = "ok"
	}else {
		log.Printf("item:%s save error: %s", item, err)
	}

	return err
}
