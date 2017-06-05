package configcache

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
)

func writeBytes(fileName string, data []byte) (bool, error) {
	var err error

	_, err = os.Stat(fileName)
	if err != nil {
		if os.IsNotExist(err) == false {
			log.WithFields(log.Fields{"error": err, "fileName": fileName}).Error("unable to read existing file")
			return false, err
		}
	} else {
		fileBytes, err := ioutil.ReadFile(fileName)
		if err != nil {
			log.WithFields(log.Fields{"error": err, "fileName": fileName}).Error("unable to read file")
			return false, err
		}
		if len(fileBytes) == len(data) {
			fullEqual := true
			for i := 0; i < len(data); i++ {
				if fileBytes[i] != data[i] {
					fullEqual = false
				}
			}
			if fullEqual {
				return false, nil
			}
		}
	}

	tempFileName := fmt.Sprintf("%s.tmp", fileName)
	f, err := os.Create(tempFileName)
	if err != nil {
		log.WithFields(log.Fields{"error": err, "tempFileName": tempFileName}).Error("unable to create temporary file")
		return false, err
	}

	_, err = f.Write(data)
	if err != nil {
		log.WithFields(log.Fields{"error": err, "tempFileName": tempFileName}).Error("unable to write to temporary file")
		return false, err
	}

	err = f.Sync()
	if err != nil {
		log.WithFields(log.Fields{"error": err, "tempFileName": tempFileName}).Error("unable to sync temporary file to disk")
		return false, err
	}

	err = f.Close()
	if err != nil {
		log.WithFields(log.Fields{"error": err, "tempFileName": tempFileName}).Error("unable to close temporary file")
		return false, err
	}

	err = os.Rename(tempFileName, fileName)
	if err != nil {
		log.WithFields(log.Fields{"error": err, "tempFileName": tempFileName, "fileName": fileName}).Error("unable to rename temporary file")
		return false, err
	}

	return true, nil
}
