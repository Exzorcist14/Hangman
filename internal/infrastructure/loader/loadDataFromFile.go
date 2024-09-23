package loader

import (
	"encoding/json"
	"fmt"
	"os"
)

// LoadDataFromFile считывает файл по указанному path и десереализует данные в target.
func LoadDataFromFile(path string, target any) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("can`t read file: %w", err)
	}

	err = json.Unmarshal(data, target)
	if err != nil {
		return fmt.Errorf("can`t unmarshal data: %w", err)
	}

	return nil
}
