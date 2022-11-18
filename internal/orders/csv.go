package orders

import (
	"encoding/csv"
	"os"
)

func FormatData(header []string, total []ServiceTotal) [][]string {
	data := [][]string{header}
	for _, value := range total {
		var tmp []string
		data = append(data, append(tmp, value.ServiceName, value.TotalAmount))
	}
	return data
}

func CreateData(pathFile string, header []string, value []ServiceTotal) error {
	file, err := os.Create(pathFile)
	if err != nil {
		return err
	}
	data := FormatData(header, value)
	writer := csv.NewWriter(file)
	err = writer.WriteAll(data)
	if err != nil {
		return err
	}
	return nil
}
