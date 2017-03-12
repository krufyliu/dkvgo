package controllers

import (
	"github.com/astaxie/beego/validation"
	"github.com/krufyliu/dkvgo/dkvgo-admin/models"
	"github.com/krufyliu/dkvgo/dkvgo-admin/services"
)

type JobsController struct {
	BaseController
}

func (this *JobsController) Get() {
	var jobs []*models.Job
	page, err := this.GetInt("page", 1)
	this.CheckError(err)
	pageSize, err := this.GetInt("size", 10)
	this.CheckError(err)
	_, err = services.JobService.GetJobList(page, pageSize).RelatedSel("Creator", "Operator").All(&jobs)
	this.CheckError(err)
	this.DataJsonResponse(jobs)
}

func (this *JobsController) Post() {
	job := models.Job{}
	valid := validation.Validation{}
	valid.Required(job.Name, "name")
	valid.Required(job.VideoDir, "video_dir")
	valid.Required(job.OutputDir, "output_dir")
	valid.Required(job.StartFrame, "start_frame")
	valid.Required(job.EndFrame, "end_frame")
	valid.Min(job.StartFrame, 0, "start_frame")
	valid.Min(job.EndFrame, 0, "end_frame")
	valid.Required(job.CameraType, "camera_type")
	valid.Required(job.EnableTop, "enable_top")
	valid.Required(job.EnableBottom, "enable_bottom")
	valid.Required(job.SaveDebugImg, "save_debug_img")
	if valid.HasErrors() {
		//this.ErrorJsonResponse("参数不合符要求", valid.Errors)
	}
	services.JobService.AddJob(&job)
	this.DataJsonResponse(job)
}