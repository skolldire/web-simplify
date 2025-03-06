package data_converter

import (
	"encoding/json"
	"fmt"
	"github.com/mitchellh/mapstructure"
	errorHandler "github.com/skolldire/web-simplify/pkg/utilities/error_handler"
	"net/http"
)

func BytesToModel[O any](c []byte) (O, error) {
	h := *new(O)
	e := map[string]interface{}{}
	err := json.Unmarshal(c, &e)
	cfg := &mapstructure.DecoderConfig{
		Metadata: nil,
		Result:   &h,
		TagName:  "json",
	}
	decoder, _ := mapstructure.NewDecoder(cfg)
	err = decoder.Decode(e)
	if err != nil {
		return h, errorHandler.NewCommonApiError("TRF-001",
			"[Convert Data To Response]Failed to convert byte array to struct", err, http.StatusInternalServerError)
	}
	return h, nil
}

func ModelToBytes[O any](c O) ([]byte, error) {
	b, err := json.Marshal(c)
	if err != nil {
		return nil, errorHandler.NewCommonApiError("TRF-002",
			"[Convert Data To Response]Failed to convert struct to bytes", err, http.StatusInternalServerError)
	}
	return b, nil
}

func MapToStructure[O any](c map[string]interface{}) (O, error) {
	h := *new(O)
	cfg := &mapstructure.DecoderConfig{
		Metadata: nil,
		Result:   &h,
		TagName:  "json",
	}
	decoder, _ := mapstructure.NewDecoder(cfg)
	err := decoder.Decode(c)
	if err != nil {
		return h, errorHandler.NewCommonApiError("TRF-003",
			"[Convert Data To Response]Failed to convert map to struct", err, http.StatusInternalServerError)
	}
	return h, nil
}

func StructToMap(data interface{}) (map[string]interface{}, error) {

	var mapa map[string]interface{}

	if data == nil {
		return mapa, fmt.Errorf("[StructToMap] input value is nil")
	}

	dadsEmJson, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(dadsEmJson, &mapa)
	if err != nil {
		return nil, err
	}

	return mapa, nil
}

func DBToDriverMap(motor string) string {
	switch motor {
	case SQLite:
		return SQLiteDriver
	case Oracle:
		return OracleDriver
	case Mysql:
		return MysqlDriver
	case Postgres:
		return PostgresDriver
	default:
		return ""
	}
}
