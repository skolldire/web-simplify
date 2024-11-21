package xml_file

import (
	"encoding/xml"
	"fmt"
	"github.com/skolldire/cash-manager-toolkit/pkg/kit/error_wrapper"
	"github.com/skolldire/cash-manager-toolkit/pkg/kit/file_util"
	"os"
)

func Write[I any](data I, path string) error {
	out, _ := xml.MarshalIndent(data, " ", "  ")
	fmt.Println(string(out))
	f, err := os.Create(path)
	if err != nil {
		return error_wrapper.NewCommonApiError(file_util.CFE.Code,
			fmt.Sprintf(file_util.CFE.Msg, "xml"), err, file_util.CFE.HttpCode)
	}
	defer f.Close()
	_, err = f.WriteString(string(out))
	if err != nil {
		return error_wrapper.NewCommonApiError(file_util.WRE.Code,
			fmt.Sprintf(file_util.WRE.Msg, "xml"), err, file_util.WRE.HttpCode)
	}
	return nil
}
