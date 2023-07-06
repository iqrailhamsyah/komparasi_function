package ftpserver

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"time"
)

/*
func main() {
	excelgetstruct := NewExcelize()
	arrayexcel := excelgetstruct.AmbilStructAllDataExcel("C:/Users/User/Downloads/JUNIOSMART-2022-10-04-settlement.xlsx")
	for _, excel := range arrayexcel {
		fmt.Print(excel.Idtransaksi)
		fmt.Print(" " + excel.Datepayment.String())
		fmt.Print(" " + excel.Datesettlement.String())
		fmt.Println(" " + excel.Totalsettlement)
	}
}

*/

type AllDataExcel struct {
	Idtransaksi     string
	Datepayment     time.Time
	Datesettlement  time.Time
	Totalsettlement string
}

func NewExcelize() ExcelizeFunction {
	//membuat objek exelcize
	excelize := new(excelize.File)
	return excelizefunction{
		libs:         excelize,
		libsopenfile: excelizewrap{},
	}
}

type ExcelizeFunction interface {
	AmbilStructAllDataExcel(filedirectory string) []AllDataExcel
}

type excelizefunction struct {
	libs         *excelize.File
	libsopenfile excelizewrap
}

type excelizewrap struct {
}

func (e excelizewrap) OpenFileFunction(fileDr string) (*excelize.File, error) {
	return excelize.OpenFile(fileDr)
}

func (e excelizefunction) AmbilStructAllDataExcel(filedirectory string) []AllDataExcel {
	//Melakukan Pembukaan File Excel
	file, err := e.libsopenfile.OpenFileFunction(filedirectory)
	if err != nil {
		fmt.Println("Error opening file:", err)
	}

	//Melakukan Pengambilan Baris-baris pada suatu Sheet
	rows := file.GetRows("Worksheet")

	var alldataexcel []AllDataExcel
	for _, row := range rows {
		//mengecualikan cell kosong dan cell berisi string "Transaction ID"
		if row[4] != "Transaction ID" && row[4] != "" {
			//membuat transactionID untuk ditambahkan ke kumpulan_transactionID
			transactionID := row[4]
			paymentdate := row[6]
			settlementdate := row[13]
			settlementtotal := row[14]

			//parsing tipe data waktu
			parsedpaymentdate, err := time.Parse("2006-01-02 15:04:05", paymentdate)
			if err != nil {
				fmt.Println("Error parsing date:", err)
			}
			parsedsettlementdate, err := time.Parse("2006-01-02 15:04:05", settlementdate)
			if err != nil {
				fmt.Println("Error parsing date:", err)
			}

			alldataexcel = append(alldataexcel, AllDataExcel{
				Idtransaksi:     transactionID,
				Datepayment:     parsedpaymentdate,
				Datesettlement:  parsedsettlementdate,
				Totalsettlement: settlementtotal,
			})

		}
	}

	//melakukan pengembalian array struct id
	return alldataexcel

}
