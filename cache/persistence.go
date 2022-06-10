package cache

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

const FileSUFFIX = "_ffb.cdb"

// Persistence  policy
type Persistence int

const (
	// FFB Full File Backup
	FFB Persistence = iota
	// AOF Append Only File
	//AOF
)

func (persistence *persistenceOption) startPersistence(data interface{}) error {
	switch persistence.persistencePolicy {
	case FFB:
		err := persistence.read(data)
		if err != nil {
			return err
		}
		go persistence.backup(data)
	}
	return nil
}

// load file
func (persistence *persistenceOption) read(data interface{}) error {
	file := filepath.Join(persistence.persistencePath, fmt.Sprintf("%s%s", persistence.persistenceName, FileSUFFIX))
	_, err := os.Stat(file)
	// Skip this step if the file exists
	if err != nil {
		return nil
	}
	fileData, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	if len(fileData) == 0 {
		return nil
	}
	return json.Unmarshal(fileData, data)
}

// If an error occurs, it fails the backup
// If the main process ends, it fails the backup and the file are 0 bytes
func (persistence *persistenceOption) backup(data interface{}) {
	file := filepath.Join(persistence.persistencePath, fmt.Sprintf("%s%s", persistence.persistenceName, FileSUFFIX))
	ticker := time.NewTicker(time.Second * 5)
	for {
		select {
		case <-ticker.C:
			backupData, err := json.Marshal(data)
			if err != nil {
				continue
			}
			err = judgeAndCreate(file)
			if err != nil {
				continue
			}
			err = ioutil.WriteFile(file, backupData, 0755)
			if err != nil {
				continue
			}
		}
	}
}

//Judge whether a file or folder exists. If it does not exist, create it
func judgeAndCreate(path string) error {
	_, err := os.Stat(path)
	if err == nil {
		return nil
	}
	err = os.MkdirAll(filepath.Dir(path), 0755)
	if err != nil {
		return err
	}
	_, err = os.Create(path)
	if err != nil {
		return err
	}
	return nil
}
