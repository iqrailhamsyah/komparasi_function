package main

import (
	"fmt"
	"komparasi/util/ftpserver"
	"time"
)

func main() {

	compare := NewCompare()
	tabel_anomali := compare.GetPaymentAnomalies()

	for _, anomali := range tabel_anomali {
		fmt.Print(anomali.tablesettlement)
		fmt.Print(" " + anomali.tabletransactionid)
		fmt.Print(" " + anomali.tablepaymentdate.String())
		fmt.Print(" " + anomali.tablesettlementdate.String())
		fmt.Println(" " + anomali.tablesettlementtotal)

	}

}

type Anomalies_Table struct {
	tablesettlement      int
	tabletransactionid   string
	tablepaymentdate     time.Time
	tablesettlementdate  time.Time
	tablesettlementtotal string
}

func NewCompare() UseCaseFaspaylogpayment {
	//membuat objek ftps
	return useCaseFaspaylogpayment{}
}

type UseCaseFaspaylogpayment interface {
	paymentTableCheck(alldataexcel []ftpserver.AllDataExcel, alldatadatabase []ftpserver.AllDataDatabase) []Anomalies_Table
	GetPaymentAnomalies() []Anomalies_Table
}

type useCaseFaspaylogpayment struct {
}

func (u useCaseFaspaylogpayment) paymentTableCheck(alldataexcel []ftpserver.AllDataExcel, alldatadatabase []ftpserver.AllDataDatabase) []Anomalies_Table {
	//membuat returning table
	/*
		settlement_status := []int{}
		trasactionid := []string{}
		paymentdate := []time.Time{}
		settlementdate := []time.Time{}
		settlementtotal := []string{}

	*/
	var tabel_anomali []Anomalies_Table

	//rekonindex
	var rekonindex int

	//membuat tabel untuk checking
	checking_id_excel_ke_payment := false
	checking_id_payment_ke_excel := false
	row := len(alldataexcel)
	column := len(alldatadatabase)
	checking_table := make([][]bool, row)
	for i := 0; i < row; i++ {
		checking_table[i] = make([]bool, column)
	}

	//logic untuk menambahkan id excel ke tabel payment jika tidak ada di fastpay
	for e, dataExcel := range alldataexcel {
		for f, dataDatabase := range alldatadatabase {

			if dataDatabase.Idtransaksi == dataExcel.Idtransaksi {
				checking_table[e][f] = true
				rekonindex = f
			} else {
				checking_table[e][f] = false
			}

			checking_id_excel_ke_payment = checking_id_excel_ke_payment || checking_table[e][f]
		}

		if checking_id_excel_ke_payment == false {
			//tambahkan id excel ke table payment, kondisi di tabel excel ada tapi di tabel payment tidak ada (Flagging Paid)
			tabel_anomali = append(tabel_anomali, Anomalies_Table{
				tablesettlement:      1,
				tabletransactionid:   dataExcel.Idtransaksi,
				tablepaymentdate:     dataExcel.Datepayment,
				tablesettlementdate:  dataExcel.Datesettlement,
				tablesettlementtotal: dataExcel.Totalsettlement,
			})
		} else if alldatadatabase[rekonindex].RekonStatus == false {
			tabel_anomali = append(tabel_anomali, Anomalies_Table{
				tablesettlement:      0,
				tabletransactionid:   dataExcel.Idtransaksi,
				tablepaymentdate:     time.Now(),
				tablesettlementdate:  time.Now(),
				tablesettlementtotal: " ",
			})
		} else {
			tabel_anomali = append(tabel_anomali, Anomalies_Table{
				tablesettlement:      2,
				tabletransactionid:   dataExcel.Idtransaksi,
				tablepaymentdate:     time.Now(),
				tablesettlementdate:  time.Now(),
				tablesettlementtotal: " ",
			})
		}
		checking_id_excel_ke_payment = false
	}

	//membuat tabel(Transpose) untuk checking
	row = len(alldatadatabase)
	column = len(alldataexcel)
	checking_table = make([][]bool, row)
	for i := 0; i < row; i++ {
		checking_table[i] = make([]bool, column)
	}

	//logic untuk menghapus id dari tabel payment jika tidak ada di excel
	for f, dataDatabase := range alldatadatabase {
		for e, dataExcel := range alldataexcel {

			if dataExcel.Idtransaksi == dataDatabase.Idtransaksi {
				checking_table[f][e] = true
				rekonindex = f
			} else {
				checking_table[f][e] = false
			}

			checking_id_payment_ke_excel = checking_id_payment_ke_excel || checking_table[f][e]
		}

		if checking_id_payment_ke_excel == false {
			//hapus id dari table payment, kondisi di table payment ada tapi di excel tidak ada (Flagging Unpaid)
			tabel_anomali = append(tabel_anomali, Anomalies_Table{
				tablesettlement:      3,
				tabletransactionid:   dataDatabase.Idtransaksi,
				tablepaymentdate:     time.Now(),
				tablesettlementdate:  time.Now(),
				tablesettlementtotal: " ",
			})
		} else if alldatadatabase[rekonindex].RekonStatus == false {
			tabel_anomali = append(tabel_anomali, Anomalies_Table{
				tablesettlement:      0,
				tabletransactionid:   dataDatabase.Idtransaksi,
				tablepaymentdate:     time.Now(),
				tablesettlementdate:  time.Now(),
				tablesettlementtotal: " ",
			})
		} else {
			tabel_anomali = append(tabel_anomali, Anomalies_Table{
				tablesettlement:      2,
				tabletransactionid:   dataDatabase.Idtransaksi,
				tablepaymentdate:     time.Now(),
				tablesettlementdate:  time.Now(),
				tablesettlementtotal: " ",
			})
		}

		checking_id_payment_ke_excel = false
	}

	/*
		for i := 0; i < len(settlement_status); i++ {
			tabel_anomali = append(tabel_anomali, Anomalies_Table{
				tablesettlement:      settlement_status[i],
				tabletransactionid:   trasactionid[i],
				tablepaymentdate:     paymentdate[i],
				tablesettlementdate:  settlementdate[i],
				tablesettlementtotal: settlementtotal[i],
			})
		}

	*/

	return tabel_anomali
}

func (u useCaseFaspaylogpayment) GetPaymentAnomalies() []Anomalies_Table {
	//tambah depedensi repo ke usecase, panggil method getPaymentByDate
	data_database := []ftpserver.AllDataDatabase{{
		Idtransaksi: "E83B8A27-10AC-42E1-91A5-23F5E0B6A4A7",
		RekonStatus: false,
	}, {
		Idtransaksi: "FF235235-B14F-49C3-9789-DC6F7D86F19F",
		RekonStatus: true,
	}, {
		Idtransaksi: "FDDE7B91-76CE-417A-A2A9-AEACE350C5A0",
		RekonStatus: false,
	}, {
		Idtransaksi: "Mengarang Sendiri",
		RekonStatus: false,
	},
	}
	ftpconnection := ftpserver.NewFtps()
	excel := ftpserver.NewExcelize()

	kredensial := &ftpserver.Credentials{
		Host:     "fsrv.bri.co.id",
		Port:     50021,
		Username: "H2H_JSRFaspay_USR2023",
		Password: "w@4rTCaBmpZL",
	}
	direktori := &ftpserver.Directorystring{
		Localfiledirectory: "C:/Users/User/Downloads/",
		Ftpfiledirectory:   "/faspay/settlement_xlsx/Processed/",
		Filename:           "JUNIOSMART-2022-10-04-settlement.xlsx",
	}
	ftpconnection.DownloadExcel(kredensial, direktori)
	return u.paymentTableCheck(excel.AmbilStructAllDataExcel("C:/Users/User/Downloads/JUNIOSMART-2022-10-04-settlement.xlsx"), data_database)
}
