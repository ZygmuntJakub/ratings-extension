package database

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"path"
	"sync"
)

type Entity interface {
	GetId() uint
}

type DB[T Entity] struct {
	mu             sync.Mutex
	File           string
	RowSerialize   func(T) ([]string, error)
	RowDeserialize func([]string) (T, error)
}

func (db *DB[T]) SaveAll(entities []T) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	file, err := os.OpenFile(path.Join(wd, db.File), os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, enitity := range entities {
		row, err := db.RowSerialize(enitity)
		fmt.Println(row)
		if err != nil {
			return err
		}
		if err := writer.Write(row); err != nil {
			return err
		}
	}

	return nil
}

func readCsv(filepath string) ([][]string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	file, err := os.Open(path.Join(wd, filepath))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1
	return reader.ReadAll()
}

func (db *DB[T]) ReadAll() ([]T, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	data, err := readCsv(db.File)
	if err != nil {
		return nil, err
	}

	var res []T
	for _, row := range data[1:] {
		value, err := db.RowDeserialize(row)
		if err != nil {
			return nil, err
		}
		res = append(res, value)
	}

	return res, nil
}

func (db *DB[T]) ReadFirst(condition func(entity T) bool) (T, error) {
	db.mu.Lock()
	defer db.mu.Unlock()
	var result T

	data, err := readCsv(db.File)
	if err != nil {
		return result, err
	}

	for i := len(data) - 1; i >= 1; i-- {
		value, err := db.RowDeserialize(data[i])
		if err != nil {
			return result, err
		}

		if condition(value) {
			return value, nil
		}
	}

	return result, errors.New("Cannot meet condition")
}

func (db *DB[T]) Read(condition func(entity T) bool) ([]T, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	data, err := readCsv(db.File)
	if err != nil {
		return nil, err
	}

	res := make([]T, 0)
	for _, row := range data[1:] {
		value, err := db.RowDeserialize(row)
		if err != nil {
			return nil, err
		}

		if condition(value) {
			res = append(res, value)
		}
	}

	return res, nil
}

func (db *DB[T]) Update(entity T) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	file, err := os.OpenFile(path.Join(wd, db.File), os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1

	data, err := reader.ReadAll()
	if err != nil {
		return err
	}

	for idx, row := range data[1:] {
		value, err := db.RowDeserialize(row)
		if err != nil {
			return err
		}

		if value.GetId() == entity.GetId() {
			updatedValue, err := db.RowSerialize(entity)
			if err != nil {
				return err
			}
			fmt.Println(updatedValue)
			data[idx+1] = updatedValue
		}
	}

	file.Seek(0, 0)
	writer := csv.NewWriter(file)
	err = writer.WriteAll(data)
	if err != nil {
		return err
	}
	writer.Flush()

	return nil
}
