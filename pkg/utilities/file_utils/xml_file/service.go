package xml_file

import (
	"encoding/xml"
	"fmt"
	"github.com/skolldire/web-simplify/pkg/utilities/error_handler"
	"github.com/skolldire/web-simplify/pkg/utilities/file_utils"
	"os"
)

func Write[I any](data I, path string) error {
	out, _ := xml.MarshalIndent(data, " ", "  ")
	fmt.Println(string(out))
	f, err := os.Create(path)
	if err != nil {
		return error_handler.NewCommonApiError(file_utils.CFE.Code,
			fmt.Sprintf(file_utils.CFE.Msg, "xml"), err, file_utils.CFE.HttpCode)
	}
	defer f.Close()
	_, err = f.WriteString(string(out))
	if err != nil {
		return error_handler.NewCommonApiError(file_utils.WRE.Code,
			fmt.Sprintf(file_utils.WRE.Msg, "xml"), err, file_utils.WRE.HttpCode)
	}
	return nil
}
