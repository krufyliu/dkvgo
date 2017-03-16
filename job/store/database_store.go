package store

import (
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"encoding/json"

	"github.com/krufyliu/dkvgo/job"
)

// TimeLayout define database datetime format
const TimeLayout = "2006-01-02 15:04:05"

type DatabaseStore struct {
	dbType string
	dbAddr string
	db     *sql.DB
}

func _checkError(err error) error {
	if err == nil || err == sql.ErrNoRows {
		return err
	}
	panic(err)
}

func NewDatabaseStore(dbType string, addr string) *DatabaseStore {
	var ds = &DatabaseStore{dbType: dbType, dbAddr: addr}
	ds.init()
	return ds
}

func (ds *DatabaseStore) init() {
	db, err := sql.Open(ds.dbType, ds.dbAddr)
	_checkError(err)
	ds.db = db
	//
	_, err = ds.db.Exec("update job set status='0' where status='1' or status='2'")
	_checkError(err)
	_, err = ds.db.Exec("update job set status='4' where status='3'")
	_checkError(err)
}

func (ds *DatabaseStore) GetJob() *job.Job {
	var query = `
	select 
		id, name, priority, progress, status, start_frame, end_frame,
		camera_type, algorithm, video_dir, output_dir, enable_top, 
		enable_bottom, quality, save_debug_img, enable_color_adjust 
	from job 
	where status = 0 
	order by priority desc 
	limit 1
	`
	var _job = job.Job{}
	var row = ds.db.QueryRow(query)
	err := row.Scan(&_job.ID, &_job.Name, &_job.Priority, &_job.Progress, &_job.Status,
		&_job.StartFrame, &_job.EndFrame, &_job.CameraType, &_job.Algorithm,
		&_job.VideoDir, &_job.OutputDir, &_job.EnableTop, &_job.EnableBottom,
		&_job.Quality, &_job.SaveDebugImg, &_job.EanbleColorAdjust)
	if _checkError(err) == sql.ErrNoRows {
		return nil
	}
	return &_job
}

func (ds *DatabaseStore) UpdateJob(_job *job.Job) bool {
	_, err := ds.db.Exec("update job set status=?, progress=?, update_at=? where id=?",
		_job.Status, _job.CalcProgress(), time.Now().Format(TimeLayout), _job.ID)
	_checkError(err)
	return true
}

func (ds *DatabaseStore) SaveJobState(_job *job.Job) bool {
	var taskOpts = _job.TaskOpts
	if len(taskOpts) == 0 {
		return true
	}
	content, err := json.Marshal(taskOpts)
	_checkError(err)
	var _id int
	var querySql = "select id from job_state where job_id=?"
	err = ds.db.QueryRow(querySql, _job.ID).Scan(&_id)
	if _checkError(err) == nil {
		var updateSql = "update job_state set content=?, update_at=? where job_id=?"
		_, err = ds.db.Exec(updateSql, content, time.Now().Format(TimeLayout), _job.ID)
		_checkError(err)
	} else {
		var timeStr = time.Now().Format(TimeLayout)
		var insertSql = `insert into job_state(job_id, content, create_at, update_at) values(?, ?, ?, ?)`
		_, err = ds.db.Exec(insertSql, _job.ID, content, timeStr, timeStr)
		_checkError(err)
	}
	return true
}

func (ds *DatabaseStore) LoadJobState(_job *job.Job) bool {
	var content []byte
	var row = ds.db.QueryRow("select content from job_state where job_id=?", _job.ID)
	err := row.Scan(&content)
	if _checkError(err) == sql.ErrNoRows {
		return false
	}
	var taskOpts []*job.TaskOptions
	if err := json.Unmarshal(content, &taskOpts); err != nil {
		panic(err)
	}
	if len(taskOpts) == 0 {
		return false
	}
	_job.TaskOpts = taskOpts
	return true
}
