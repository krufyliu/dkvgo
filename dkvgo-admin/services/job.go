package services

import (
	"github.com/astaxie/beego/orm"
	"github.com/krufyliu/dkvgo/dkvgo-admin/models"
)

type jobService struct{}

func (this *jobService) GetTotal() (int64, error) {
	return o.QueryTable(&models.Job{}).Count()
}

func (this *jobService) GetPage(current, pageSize int) (*Page, error) {
	total, err := this.GetTotal()
	if err != nil {
		return nil, err
	}
	return &Page{Total: total, Current: current, PageSize: pageSize}, nil
}

func (this *jobService) CreateJob(name, videoDir, outDir string, startFrame, endFrame, priority int,
	cameraType, enableTop, enableBottom, saveDebugImg string) (*models.Job, error) {
	job := &models.Job{
		Name:         name,
		VideoDir:     videoDir,
		OutputDir:    outDir,
		StartFrame:   startFrame,
		EndFrame:     endFrame,
		Priority:     priority,
		CameraType:   cameraType,
		Algorithm:    "3D_" + cameraType,
		EnableTop:    enableTop,
		EnableBottom: enableBottom,
		SaveDebugImg: saveDebugImg,
	}
	_, err := o.Insert(job)
	return job, err
}

func (this *jobService) AddJob(job *models.Job) error {
	_, err := o.Insert(job)
	return err
}

func (this *jobService) GetJobList(page, pageSize int) orm.QuerySeter {
	offset := (page - 1) * pageSize
	if offset < 0 {
		offset = 0
	}
	qs := o.QueryTable(&models.Job{}).Limit(pageSize, offset)
	return qs
}

func (this *jobService) GetJob(id int, withState bool) (*models.Job, error) {
	job := &models.Job{Id: id}
	err := o.Read(job)
	if err != nil {
		return nil, err
	}
	if withState {
		err = o.Read(job.State)
	}
	return job, err
}

func (this *jobService) Update(job *models.Job, fields ...string) (int64, error) {
	return o.Update(job, fields...)
}
