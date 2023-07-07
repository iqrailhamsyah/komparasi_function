package main

import (
	"fmt"
	"komparasi/modules/reconapp"
)

func main() {

	compare := reconapp.NewCompare()
	tabel_anomali := compare.GetPaymentAnomalies()

	for _, anomali := range tabel_anomali {
		fmt.Print(anomali.Tablesettlement)
		fmt.Print(" " + anomali.Tabletransactionid)
		fmt.Print(" " + anomali.Tablepaymentdate.String())
		fmt.Print(" " + anomali.Tablesettlementdate.String())
		fmt.Println(" " + anomali.Tablesettlementtotal)

	}

}
