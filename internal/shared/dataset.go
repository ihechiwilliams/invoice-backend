package shared

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"sort"
	"time"

	"github.com/samber/lo"
	"gorm.io/gorm"
)

const (
	defaultPaginationLimit = 100
)

type AnyTime struct{}

func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

func BuildDataset(ctx context.Context, db *gorm.DB, tableName string, filters any) (*gorm.DB, error) {
	dataset := db.WithContext(ctx).Table(tableName)

	return FilterDataset(dataset, filters)
}

func FilterDataset(dataset *gorm.DB, filters any) (*gorm.DB, error) {
	newDataset := dataset
	bytesFilters := lo.Must(json.Marshal(filters))
	rawFilters := map[string]interface{}{}

	err := json.Unmarshal(bytesFilters, &rawFilters)
	if err != nil {
		return nil, err
	}

	keys := lo.Keys(rawFilters)
	sort.Strings(keys)

	for _, key := range keys {
		value := rawFilters[key]
		valueType := reflect.TypeOf(value).Kind()

		if key == "created_after" {
			newDataset = newDataset.Where("DATE(created_at) > ?", value)

			continue
		}

		if key == "created_before" {
			newDataset = newDataset.Where("DATE(created_at) < ?", value)

			continue
		}

		if key == "created_after_datetime" {
			newDataset = newDataset.Where("created_at > ?", value)

			continue
		}

		if key == "created_before_datetime" {
			newDataset = newDataset.Where("created_at < ?", value)

			continue
		}

		if key == "transactions_metadata_operation_types" {
			newDataset = newDataset.Where("metadata ->> 'operation_type' IN ?", value)

			continue
		}

		if !lo.Contains([]reflect.Kind{reflect.Interface, reflect.Map, reflect.Slice}, valueType) {
			newDataset = newDataset.Where(map[string]interface{}{key: value})

			continue
		}

		result, jsonErr := handleJSONFilters(newDataset, key, value)
		if jsonErr != nil {
			return nil, jsonErr
		}

		newDataset = result
	}

	return newDataset, nil
}

func PaginateDataset(dataset *gorm.DB, pagination Pagination) *gorm.DB {
	if pagination.Limit != nil {
		return dataset.Limit(*pagination.Limit)
	}

	return dataset.Limit(defaultPaginationLimit)
}

func ExcludeFilterDataset(dataset *gorm.DB, filters any) (*gorm.DB, error) {
	newDataset := dataset
	bytesFilters := lo.Must(json.Marshal(filters))
	rawFilters := map[string]interface{}{}

	err := json.Unmarshal(bytesFilters, &rawFilters)
	if err != nil {
		return nil, err
	}

	keys := lo.Keys(rawFilters)
	sort.Strings(keys)

	for _, key := range keys {
		value := rawFilters[key]
		valueType := reflect.TypeOf(value).Kind()

		if !lo.Contains([]reflect.Kind{reflect.Interface, reflect.Map, reflect.Slice}, valueType) {
			newDataset = newDataset.Not(map[string]interface{}{key: value})

			continue
		}

		result, jsonErr := handleExcludeJSONFilters(newDataset, key, value)
		if jsonErr != nil {
			return nil, jsonErr
		}

		newDataset = result
	}

	return newDataset, nil
}

func handleJSONFilters(dataset *gorm.DB, columnName string, filters any) (*gorm.DB, error) {
	newDataset := dataset

	switch filter := filters.(type) {
	case map[string]interface{}:
		for k, v := range filter {
			valueType := reflect.TypeOf(v).Kind()
			switch valueType {
			case reflect.Slice:
				newDataset = newDataset.Where(fmt.Sprintf("%s ->> '%s' IN (?)", columnName, k), v)
			default:
				newDataset = newDataset.Where(fmt.Sprintf("%s ->> '%s' = ?", columnName, k), v)
			}
		}
	case []any:
		newDataset = newDataset.Where(fmt.Sprintf("%s IN ?", columnName), filter)
	default:
		return nil, errors.New("filters are invalid")
	}

	return newDataset, nil
}

func handleExcludeJSONFilters(dataset *gorm.DB, columnName string, filters any) (*gorm.DB, error) {
	newDataset := dataset

	switch filter := filters.(type) {
	case map[string]interface{}:
		for k, v := range filter {
			newDataset = newDataset.Where(fmt.Sprintf("%s ->> '%s' = ?", columnName, k), v)
		}
	case []any:
		newDataset = newDataset.Where(fmt.Sprintf("%s NOT IN ?", columnName), filter)
	default:
		return nil, errors.New("filters are invalid")
	}

	return newDataset, nil
}
