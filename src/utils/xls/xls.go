package xls

import (
	"github.com/spf13/cast"
)

type Xls struct{}

func (xls *Xls) XlsGetCellRow(k int, row int) string {
	darr := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}
	dmarr := []string{"AA", "AB", "AC", "AD", "AE", "AF", "AG", "AH", "AI", "AJ", "AK", "AL", "AM", "AN", "AO", "AP", "AQ", "AR", "AS", "AT", "AU", "AV", "AW", "AX", "AY", "AZ"}
	if k < 26 {
		return darr[k] + cast.ToString(row)
	} else {
		return dmarr[k] + cast.ToString(row)
	}
}

func (xls *Xls) XlsGetCell(k int) string {
	darr := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}
	dmarr := []string{"AA", "AB", "AC", "AD", "AE", "AF", "AG", "AH", "AI", "AJ", "AK", "AL", "AM", "AN", "AO", "AP", "AQ", "AR", "AS", "AT", "AU", "AV", "AW", "AX", "AY", "AZ"}
	if k < 26 {
		return darr[k]
	} else {
		return dmarr[k]
	}
}
