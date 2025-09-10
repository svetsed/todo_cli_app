package storage

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/gofrs/flock"
)

func Save(filename string, data any) error {
	lock := flock.New(filename + ".lock")

	locked, err := lock.TryLock()
	if err != nil {
		return err
	}

	if !locked {
		return fmt.Errorf("file is locked by another process")
	}
	defer lock.Unlock()

	fileData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	tmpFile := filename + ".tmp"
	if err := os.WriteFile(tmpFile, fileData, 0644); err != nil {
		return err
	}

	return os.Rename(tmpFile, filename)
}

func Load(filename string, data any) error {
	lock := flock.New(filename + ".lock")

	locked, err := lock.TryRLock()
	if err != nil {
		return err
	}

	if !locked {
		return fmt.Errorf("file is locked by another process")
	}

	defer lock.Unlock()

	fileData, err := os.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	if len(fileData) == 0 {
		return nil
	}
	return json.Unmarshal(fileData, data)
}
