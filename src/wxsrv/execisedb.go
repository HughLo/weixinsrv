package wxsrv

import (
	"log"
	"fmt"
)

var ConnString string

type ExeciseRecord struct {
	UserName string
	ExeciseTime int //the merit is minute
	ExeciseEnergy int // the merit is kcal
}

type ReportData struct {
	TotalTime int
	TotalEnergy int
}

type ExeciseDB struct {
	dbMgr *DBMgr
}

//return nil when failed to create execise db instance.
//has to set ConnString first. Or it will return nil
func CreateExeciseDB() *ExeciseDB {
	if len(ConnString) <= 0 {
		log.Println("empyt connection string")
		return nil
	}

	dbmgr := CreateDBMgr(ConnString)
	ed := &ExeciseDB{dbMgr:dbmgr}
	ed.dbMgr.UseDB("weixin_hugh")
	return ed
}

func (ed *ExeciseDB) Close() {
	ed.dbMgr.Close()
}

func (ed *ExeciseDB) Insert(rec *ExeciseRecord) (ExecResult, error) {
	cols := []string{"user","execisetime","execiseenergy"}
	vals := []string{
		fmt.Sprintf(`"%s"`,rec.UserName),
		fmt.Sprintf("%d",rec.ExeciseTime),
		fmt.Sprintf("%d",rec.ExeciseEnergy),
	}
	return ed.dbMgr.Cols(cols).Table("execise_records").Values(vals).Insert()
}

func (ed *ExeciseDB) ReportAll() (*ReportData, error) {
	r, err := ed.dbMgr.Call("report_all")
	if err != nil {
		return nil, err
	}

	defer r.Rows.Close()
	for r.Rows.Next() {
		var tt, te int
		err = r.Rows.Scan(tt, te)
		if err != nil {
			return nil, err
		}

		return &ReportData{TotalEnergy:te, TotalTime:tt}, nil
	}

	return nil, nil
}