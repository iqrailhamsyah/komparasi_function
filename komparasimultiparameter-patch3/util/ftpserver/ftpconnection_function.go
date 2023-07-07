package ftpserver

import (
	"github.com/webguerilla/ftps"
)

func NewFtps() FtpsconnectionFunction {
	//membuat objek ftps
	ftps := new(ftps.FTPS)
	return ftpsconnectionfunction{
		libs: ftps,
	}
}

type FtpsconnectionFunction interface {
	DownloadExcel(cred *Credentials, dir *Directorystring)
}

type ftpsconnectionfunction struct {
	libs *ftps.FTPS
}

type Credentials struct {
	Host     string
	Port     int
	Username string
	Password string
}

type Directorystring struct {
	Localfiledirectory string
	Ftpfiledirectory   string
	Filename           string
}

func (f ftpsconnectionfunction) DownloadExcel(cred *Credentials, dir *Directorystring) {

	f.libs.TLSConfig.InsecureSkipVerify = true // often necessary in shared hosting environments
	f.libs.Debug = true

	//setting connect ke host
	err := f.libs.Connect(cred.Host, cred.Port)
	if err != nil {
		panic(err)
	}

	//menutup koneksi FTP
	defer f.libs.Quit()

	//setting credential username & password
	err = f.libs.Login(cred.Username, cred.Password)
	if err != nil {
		panic(err)
	}

	//proses download file dari remote FTP server ke local directory
	f.libs.RetrieveFile(dir.Ftpfiledirectory+dir.Filename, dir.Localfiledirectory+dir.Filename)

}
