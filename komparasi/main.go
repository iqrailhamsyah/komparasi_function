package main

import "fmt"

func main() {
	transactionid_excel := []string{"apa", "mengapa", "bagaimana"}
	transactionid_faspaylog := []string{"mengapa", "bagaimana", "siapa"}

	compare := NewCompare()
	compare.Fastpay_table_check(transactionid_excel, transactionid_faspaylog)
}

func NewCompare() UseCaseFaspaylogpayment {
	//membuat objek ftps
	return useCaseFaspaylogpayment{}
}

type UseCaseFaspaylogpayment interface {
	Fastpay_table_check(transactionid_excel []string, transactionid_faspaylog []string)
}

type useCaseFaspaylogpayment struct {
}

func (u useCaseFaspaylogpayment) Fastpay_table_check(transactionid_excel []string, transactionid_faspaylog []string) {
	//membuat tabel untuk checking
	checking_id_excel_ke_fastpay := false
	checking_id_fastpay_ke_excel := false
	row := len(transactionid_excel)
	column := len(transactionid_faspaylog)
	checking_table := make([][]bool, row)
	for i := 0; i < row; i++ {
		checking_table[i] = make([]bool, column)
	}

	//logic untuk menambahkan id excel ke tabel fastpay jika tidak ada di fastpay
	for e, idexcel := range transactionid_excel {
		for f, idfastpay := range transactionid_faspaylog {

			if idfastpay == idexcel {
				checking_table[e][f] = true
			} else {
				checking_table[e][f] = false
			}
			fmt.Println(checking_table[e][f])

			checking_id_excel_ke_fastpay = checking_id_excel_ke_fastpay || checking_table[e][f]
		}

		if checking_id_excel_ke_fastpay == false {
			fmt.Println("tambahkan id excel ke table fastpay")
		}
		checking_id_excel_ke_fastpay = false
	}

	//logic untuk menghapus id dari tabel fastpay jika tidak ada di excel
	for f, idfastpay := range transactionid_faspaylog {
		for e, idexcel := range transactionid_excel {

			if idexcel == idfastpay {
				checking_table[f][e] = true
			} else {
				checking_table[f][e] = false
			}
			fmt.Println(checking_table[f][e])

			checking_id_fastpay_ke_excel = checking_id_fastpay_ke_excel || checking_table[f][e]

		}

		if checking_id_fastpay_ke_excel == false {
			fmt.Println("hapus id dari table fastpay")
		}
		checking_id_fastpay_ke_excel = false
	}
}
