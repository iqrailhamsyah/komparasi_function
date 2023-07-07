package reconapp

import (
	"komparasi/util/db"
	"komparasi/util/ftpserver"
	"time"
)

type Anomalies_Table struct {
	Tablesettlement      int
	Tabletransactionid   string
	Tablepaymentdate     time.Time
	Tablesettlementdate  time.Time
	Tablesettlementtotal string
}

func NewCompare() UseCaseFaspaylogpayment {
	//membuat objek ftps
	return useCaseFaspaylogpayment{}
}

type UseCaseFaspaylogpayment interface {
	paymentTableCheck(alldataexcel []ftpserver.AllDataExcel, alldatadatabase []db.AllDataDatabase) []Anomalies_Table
	GetPaymentAnomalies() []Anomalies_Table
}

type useCaseFaspaylogpayment struct {
}

func (u useCaseFaspaylogpayment) paymentTableCheck(alldataexcel []ftpserver.AllDataExcel, alldatadatabase []db.AllDataDatabase) []Anomalies_Table {
	//membuat returning table
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
				Tablesettlement:      1,
				Tabletransactionid:   dataExcel.Idtransaksi,
				Tablepaymentdate:     dataExcel.Datepayment,
				Tablesettlementdate:  dataExcel.Datesettlement,
				Tablesettlementtotal: dataExcel.Totalsettlement,
			})
		} else if alldatadatabase[rekonindex].RekonStatus == false {
			tabel_anomali = append(tabel_anomali, Anomalies_Table{
				Tablesettlement:      0,
				Tabletransactionid:   dataExcel.Idtransaksi,
				Tablepaymentdate:     time.Now(),
				Tablesettlementdate:  time.Now(),
				Tablesettlementtotal: " ",
			})
		} else {
			tabel_anomali = append(tabel_anomali, Anomalies_Table{
				Tablesettlement:      2,
				Tabletransactionid:   dataExcel.Idtransaksi,
				Tablepaymentdate:     time.Now(),
				Tablesettlementdate:  time.Now(),
				Tablesettlementtotal: " ",
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
				Tablesettlement:      3,
				Tabletransactionid:   dataDatabase.Idtransaksi,
				Tablepaymentdate:     time.Now(),
				Tablesettlementdate:  time.Now(),
				Tablesettlementtotal: " ",
			})
		} else if alldatadatabase[rekonindex].RekonStatus == false {
			tabel_anomali = append(tabel_anomali, Anomalies_Table{
				Tablesettlement:      0,
				Tabletransactionid:   dataDatabase.Idtransaksi,
				Tablepaymentdate:     time.Now(),
				Tablesettlementdate:  time.Now(),
				Tablesettlementtotal: " ",
			})
		} else {
			tabel_anomali = append(tabel_anomali, Anomalies_Table{
				Tablesettlement:      2,
				Tabletransactionid:   dataDatabase.Idtransaksi,
				Tablepaymentdate:     time.Now(),
				Tablesettlementdate:  time.Now(),
				Tablesettlementtotal: " ",
			})
		}

		checking_id_payment_ke_excel = false
	}

	return tabel_anomali
}

func (u useCaseFaspaylogpayment) GetPaymentAnomalies() []Anomalies_Table {
	//tambah depedensi repo ke usecase, panggil method getPaymentByDate
	data_database := []db.AllDataDatabase{{
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
