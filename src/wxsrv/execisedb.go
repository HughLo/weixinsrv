package wxsrv

import (
	"fmt"
	"log"
	"time"
)

//Must set the ConnString before you start to use ExeciseDB
//CreateExeciseDB will use the ConnString to connect to the
//database
var ConnString string

type ExeciseRecord struct {
	UserName      string
	ExeciseTime   int //the merit is minute
	ExeciseEnergy int // the merit is kcal
}

type ReportData struct {
	TotalTime   int
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
	ed := &ExeciseDB{dbMgr: dbmgr}
	ed.dbMgr.UseDB("weixin_hugh")
	return ed
}

func (ed *ExeciseDB) Close() {
	ed.dbMgr.Close()
}

//insert a ExciseRecord instance into database
func (ed *ExeciseDB) Insert(rec *ExeciseRecord) (ExecResult, error) {
	cols := []string{"user", "execisetime", "execiseenergy"}
	vals := []string{
		fmt.Sprintf(`"%s"`, rec.UserName),
		fmt.Sprintf("%d", rec.ExeciseTime),
		fmt.Sprintf("%d", rec.ExeciseEnergy),
	}
	return ed.dbMgr.Cols(cols).Table("execise_records").Values(vals).Insert()
}

func (ed *ExeciseDB) ReportAll() (*ReportData, error) {
	qs := `select sum(er.execisetime) as all_execise_time,
		sum(er.execiseenergy) as all_execise_energy from execise_records as er;`

	return ed.reportInternal(qs)
}

func (ed *ExeciseDB) ReportSinceThisWeek() (*ReportData, error) {
	year, wk := time.Now().ISOWeek()
	return ed.ReportSinceWeek(year, wk)
}

func (ed *ExeciseDB) ReportSinceLastWeek() (*ReportData, error) {
	t := time.Now().AddDate(0, 0, -7)
	return ed.ReportSinceWeek(t.ISOWeek())
}

func (ed *ExeciseDB) ReportSinceWeek(year, wk int) (*ReportData, error) {
	//mysql supports 7 modes of week representation. please refer to
	//http://dev.mysql.com/doc/refman/5.5/en/date-and-time-functions.html#function_week
	//to get the details about the 7 modes. The mode 3 represenatation
	//is what used by GOLANG.
	qs := `select sum(er.execisetime) as all_execise_time, sum(er.execiseenergy) as
		all_execise_energy from execise_records as er where yearweek(er.record_time, 3)>=%d;`
	qs = fmt.Sprintf(qs, year*100+wk)

	return ed.reportInternal(qs)
}

func (ed *ExeciseDB) ReportSinceThisMonth() (*ReportData, error) {
	now := time.Now()
	return ed.ReportSinceMonth(now.Year(), int(now.Month()))
}

func (ed *ExeciseDB) ReportSinceLastMonth() (*ReportData, error) {
	prev := time.Now().AddDate(0, -1, 0)
	return ed.ReportSinceMonth(prev.Year(), int(prev.Month()))
}

func (ed *ExeciseDB) ReportSinceMonth(year, mon int) (*ReportData, error) {
	qs := `select sum(er.execisetime) as all_execise_time, sum(er.execiseenergy) as
		all_execise_energy from execise_records as er where year(er.record_time)*100+
		month(er.record_time)>=%d;`
	qs = fmt.Sprintf(qs, year*100+mon)

	return ed.reportInternal(qs)
}

func (ed *ExeciseDB) ReportSinceThisYear() (*ReportData, error) {
	return ed.ReportSinceYear(time.Now().Year())
}

func (ed *ExeciseDB) ReportSinceLastYear() (*ReportData, error) {
	return ed.ReportSinceYear(time.Now().Year()-1)
}

func (ed *ExeciseDB) ReportSinceYear(year int) (*ReportData, error) {

	qs := `select sum(er.execisetime) as all_execise_time, sum(er.execiseenergy) as
		all_execise_energy from execise_records as er where year(er.record_time)>=%d;`
	qs = fmt.Sprintf(qs, year)

	return ed.reportInternal(qs)
}

func (ed *ExeciseDB) reportInternal(qs string) (*ReportData, error) {
	r, err := ed.dbMgr.RawQuery(qs)
	if err != nil {
		return nil, err
	}

	defer r.Rows.Close()

	if !r.Rows.Next() {
		return nil, nil
	}

	var rd ReportData
	err = r.Rows.Scan(&(rd.TotalTime), &(rd.TotalEnergy))
	if err != nil {
		return nil, err
	}

	return &rd, nil
}
