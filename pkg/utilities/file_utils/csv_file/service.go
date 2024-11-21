package csv_file

import (
	"github.com/gocarina/gocsv"
	"os"
)

func Read[O any](path string) ([]O, error) {
	csvFile, err := os.OpenFile(path, os.O_RDWR, os.ModePerm)
	if err != nil {
		return nil, err
	}
	defer csvFile.Close()
	var registries []O
	if err := gocsv.UnmarshalFile(csvFile, &registries); err != nil {
		return nil, err
	}

	return registries, nil
}

func Write[I any](path string, data []I) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	if err := gocsv.MarshalFile(&data, file); err != nil {
		return err
	}
	return nil
}
