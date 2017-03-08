package store

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"

	"encoding/json"

	"github.com/krufyliu/dkvgo/job"
)

type DatabaseStore struct {
	dbType string
	dbAddr string
	db     *sql.DB
}

func NewDatabaseStore(dbType string, addr string) *DatabaseStore {
	var ds = &DatabaseStore{dbType: dbType, dbAddr: addr}
	ds.init()
	return ds
}

func (ds *DatabaseStore) init() {
	db, err := sql.Open(ds.dbType, ds.dbAddr)
	if err != nil {
		panic(err)
	}
	ds.db = db
}

func (ds *DatabaseStore) GetJob() *job.Job {
	var query = `
	select 
		id, name, priority, progress, status, start_frame, end_frame,
		camera_type, algorithm, video_dir, output_dir, enable_top, 
		enable_bottom, quality, enable_color_adjust 
	from jobs 
	where status = 0 
	order by priority desc 
	limit 1
	`
	var _job = job.Job{}
	var row = ds.db.QueryRow(query)
	err := row.Scan(&_job.ID, &_job.Name, &_job.Priority, &_job.Progress, &_job.Status,
		&_job.StartFrame, &_job.EndFrame, &_job.CameraType, &_job.Algorithm,
		&_job.VideoDir, &_job.OutputDir, &_job.EnableTop, &_job.EnableBottom,
		&_job.Quality, &_job.EanbleColorAdjust)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		panic(err)
	}
	return &_job
}

func (ds *DatabaseStore) UpdateJob(_job *job.Job) bool {
	_, err := ds.db.Exec("update jobs set status=?, progress=? where id=?",
		_job.Status, _job.Progress, _job.ID)
	if err != nil {
		panic(err)
	}
	return true
}

func (ds *DatabaseStore) SaveJobState(_job *job.Job) bool {
	var taskOpts = _job.TaskOpts
	if len(taskOpts) != 0 {
		return true
	}
	content, err := json.Marshal(taskOpts)
	if err != nil {
		panic(err)
	}
	result, err := ds.db.Exec("update job_states set content=? where job_id=?", content, _job.ID)
	if err != nil {
		panic(err)
	}
	count, err := result.RowsAffected()
	if err != nil {
		panic(err)
	}
	if count == 0 {
		_, err = ds.db.Exec("insert into job_states(job_id, content) values(?, ?)", _job.ID, string(content))
		if err != nil {
			panic(err)
		}
	}
	return true
}

func (ds *DatabaseStore) LoadJobState(_job *job.Job) bool {
	var content []byte
	var row = ds.db.QueryRow("select content from job_states where job_id=?", _job.ID)
	err := row.Scan(&content)
	if err != nil {
		if err == sql.ErrNoRows {
			return false
		}
		panic(err)
	}
	var taskOpts []*job.TaskOptions
	if err := json.Unmarshal(content, taskOpts); err != nil {
		panic(err)
	}
	if len(taskOpts) == 0 {
		return false
	}
	_job.TaskOpts = taskOpts
	return true
}
