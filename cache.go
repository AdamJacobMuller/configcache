package configcache

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
)

type ConfigCache struct {
	CacheFile string
}

func (a *ConfigCache) Store(config interface{}) (bool, error) {
	mData, err := json.Marshal(config)
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Error("unable to marshal data")
		return false, err
	}
	wrote, err := writeBytes(a.CacheFile, mData)
	if err != nil {
		log.WithFields(log.Fields{"err": err, "cacheFile": a.CacheFile}).Error("WriteBytes failed")
		return false, err
	} else {
		log.WithFields(log.Fields{"cacheFile": a.CacheFile, "wrote": wrote}).Info("WriteBytes completed")
	}
	return true, nil
}

func (a *ConfigCache) Load(config interface{}) error {
	var err error
	_, err = os.Stat(a.CacheFile)
	if err != nil {
		log.WithFields(log.Fields{"error": err, "cacheFile": a.CacheFile}).Error("unable to read existing file")
		return err
	}
	fileBytes, err := ioutil.ReadFile(a.CacheFile)
	if err != nil {
		log.WithFields(log.Fields{"error": err, "cacheFile": a.CacheFile}).Error("unable to read existing file")
		return err
	}
	err = json.Unmarshal(fileBytes, &config)
	if err != nil {
		log.WithFields(log.Fields{"error": err, "cacheFile": a.CacheFile}).Error("json.Unmarshal")
		return err
	}

	log.WithFields(log.Fields{"cacheFile": a.CacheFile, "size": len(fileBytes)}).Info("loaded sites data from cache")
	return nil
}
