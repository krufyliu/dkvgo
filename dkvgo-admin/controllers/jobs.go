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
	valid.Required(job.Name, "Name")
	valid.Required(job.VideoDir, "VideoDir")
	valid.Required(job.OutputDir, "OutputDir")
	valid.Required(job.StartFrame, "StartFrame")
	valid.Required(job.EndFrame, "EndFrame")
	valid.Min(job.StartFrame, 0, "StartFrame")
	valid.Min(job.EndFrame, 0, "Endframe")
	valid.Required(job.CameraType, "Cameratype")
	valid.Required(job.EnableTop, "EnableTop")
	valid.Required(job.SaveDebugImg, "SaveDebugImg")
	if valid.HasErrors() {
		this.ErrorJsonResponse("参数不合符要求", valid.Errors)
	}
	job.EnableBottom = job.EnableTop
	job.Creator = this.LoginUser()
	job.Operator = this.LoginUser()
	services.JobService.AddJob(&job)
	this.DataJsonResponse(job)
}
