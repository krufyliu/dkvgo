package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"io/ioutil"

	"encoding/json"

	"github.com/astaxie/beego"
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
	qs := services.JobService.GetJobList(page, pageSize)
	field := this.GetString("field")
	keyword := this.GetString("keyword")
	if field != "" && keyword != "" {
		if field == "Name" || field == "VideoDir" {
			qs = qs.Filter(field+"__contains", keyword)
		}
	}
	_, err = qs.OrderBy("-UpdateAt").RelatedSel("Creator", "Operator").All(&jobs)
	this.CheckError(err)
	pager, err := services.JobService.GetPage(page, pageSize)
	this.CheckError(err)
	this.DataJsonResponseWithPage(jobs, pager)
}

func (this *JobsController) Post() {
	startFrame, err := this.GetInt("StartFrame")
	this.CheckError(err)
	endFrame, err := this.GetInt("EndFrame")
	this.CheckError(err)
	job := models.Job{
		Name:         this.GetString("Name"),
		VideoDir:     this.GetString("VideoDir"),
		OutputDir:    this.GetString("OutputDir"),
		StartFrame:   startFrame,
		EndFrame:     endFrame,
		CameraType:   this.GetString("CameraType"),
		Quality:      this.GetString("Quality"),
		EnableTop:    this.GetString("EnableTop"),
		EnableBottom: this.GetString("EnableBottom"),
		SaveDebugImg: this.GetString("SaveDebugImg"),
	}
	valid := validation.Validation{}
	valid.Required(job.Name, "Name")
	valid.Required(job.VideoDir, "VideoDir")
	valid.Required(job.OutputDir, "OutputDir")
	valid.Required(job.CameraType, "Cameratype")
	valid.Required(job.EnableTop, "EnableTop")
	valid.Required(job.EnableBottom, "EnableTop")
	valid.Required(job.SaveDebugImg, "SaveDebugImg")
	if startFrame >= endFrame {
		valid.SetError("StartFrame", "StartFrame必须小于EndFrame")
	}
	if valid.HasErrors() {
		this.ErrorJsonResponse("参数不合符要求", valid.Errors)
	}
	job.EnableBottom = job.EnableTop
	job.Algorithm = "3D_" + job.CameraType
	job.EnableColorAdjust = "1"
	job.Creator = this.LoginUser()
	job.Operator = this.LoginUser()
	services.JobService.AddJob(&job)
	this.DataJsonResponse(job)
}

func (this *JobsController) Resume() {
	jobId, _ := strconv.Atoi(this.Ctx.Input.Param(":id"))
	job, err := services.JobService.GetJob(jobId, false)
	this.CheckError(err)
	if job.Status != 0x04 && job.Status != 0x06 {
		this.ShowErrorMsg("当前操作不允许")
	}
	job.Status = 0x00
	job.Operator = this.LoginUser()
	services.JobService.Update(job, "Status", "Operator")
	this.DataJsonResponse(job)
}

func (this *JobsController) Stop() {
	jobId, _ := strconv.Atoi(this.Ctx.Input.Param(":id"))
	job, err := services.JobService.GetJob(jobId, false)
	this.CheckError(err)
	if job.Status == 0x00 {
		job.Status = 0x04
		job.Operator = this.LoginUser()
		services.JobService.Update(job, "Status", "UpdateAt")
	} else if job.Status == 0x01 || job.Status == 0x02 {
		var scheAPIAddr = beego.AppConfig.String("scheduler.apiaddr")
		if scheAPIAddr == "" {
			scheAPIAddr = "http://localhost:9999"
		} else if !strings.HasPrefix(scheAPIAddr, "http://") {
			scheAPIAddr = "http://" + scheAPIAddr
		}
		stopURL := fmt.Sprintf("%s/api/jobs/%d/action/stop", scheAPIAddr, jobId)
		beego.Info("stopURL", stopURL)
		res, err := http.Post(stopURL, "application/json", nil)
		defer res.Body.Close()
		if err != nil {
			this.ShowErrorMsg("操作失败, 请稍后重试")
		}
		if res.StatusCode != 200 {
			this.ShowErrorMsg("call api: " + res.Status)
		}
		body, err := ioutil.ReadAll(res.Body)
		this.CheckError(err)
		result := make(map[string]interface{})
		err = json.Unmarshal(body, &result)
		this.CheckError(err)
		if value, ok := result["success"].(bool); ok && value {
			job.Status = 0x03
			job.Operator = this.LoginUser()
			services.JobService.Update(job, "Status", "Operator", "UpdateAt")
		} else {
			this.ShowErrorMsg("当前操作不允许")
		}
	} else {
		this.ShowErrorMsg("当前操作不允许")
	}
	this.DataJsonResponse(job)
}
