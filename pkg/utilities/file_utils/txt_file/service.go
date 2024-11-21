package txt_file

import (
	"encoding/json"
	"fmt"
	"github.com/skolldire/cash-manager-toolkit/pkg/kit/error_wrapper"
	"github.com/skolldire/cash-manager-toolkit/pkg/kit/file_util"
	"io"
	"os"
)

func Write[I any](path string, data []I) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, item := range data {
		jsonData, err := json.Marshal(item)
		if err != nil {
			return err
		}
		str := string(jsonData)
		str = str[1 : len(str)-1]
		str = str + "\n"
		_, err = file.WriteString(str)
		if err != nil {
			return err
		}
	}
	return nil
}

func Read[O any](path string) ([]O, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, error_wrapper.NewCommonApiError(file_util.CFE.Code,
			fmt.Sprintf(file_util.CFE.Msg, "xml"), err, file_util.CFE.HttpCode)
	}
	defer file.Close()
	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	var output []O
	err = json.Unmarshal(bytes, &output)
	if err != nil {
		return nil, err
	}
	return output, nil
}
