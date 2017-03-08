package job

import (
	"fmt"
	"os/exec"
	"path"
	"runtime"
	"strconv"
	"strings"

	"github.com/krufyliu/dkvgo/util"
)

// BinMap define composition algorithm to executable
var BinMap = map[string]string{
	"VISIONDK_3D": "test_3d_visiondk",
	"VISIONDK_2D": "test_2d_visiondk",
	"FACEBOOK_3D": "test_3d_facebook",
	"FACEBOOK_2D": "test_2d_facebook",
	"PREVIEW":     "test_preview",
	"TOP_BOTTOM":  "test_top_and_bottom",
}

// CmdGenerator hold task options to generate exec.Cmd or shell command
type CmdGenerator struct {
	job              *Job
	segOptions       *TaskOptions
	threadNum        int
	binDirecotry     string
	settingDirectory string
}

// NewCmdGeneratorFromTaskSegment create CmdGenerator from Task instance
func NewCmdGeneratorFromTaskSegment(task *Task, threadNum int, binDirectory string, settingDirectory string) *CmdGenerator {
	if threadNum <= 0 {
		threadNum = runtime.NumCPU()
	}
	return &CmdGenerator{
		job:              task.Job,
		segOptions:       task.Options,
		threadNum:        threadNum,
		binDirecotry:     binDirectory,
		settingDirectory: settingDirectory,
	}
}

func (cg *CmdGenerator) getBinName() string {
	return BinMap[cg.job.Algorithm]
}

func (cg *CmdGenerator) getSaveType() string {
	if strings.ToUpper(cg.job.CameraType) == "BMPCC" {
		return ".tiff"
	}
	return ".jpg"
}

func (cg *CmdGenerator) getFinalOuptutDir() string {
	outputDir := cg.job.OutputDir
	if strings.ToUpper(cg.job.Algorithm) == "TOP_BOTTOM" {
		return outputDir
	}
	return path.Join(outputDir, cg.job.Algorithm)
}

func (cg *CmdGenerator) getCameraSettingFileName() string {
	cameraType := strings.ToLower(cg.job.CameraType)
	algo := strings.ToLower(cg.job.Algorithm)
	if strings.ToUpper(algo) == "PREVIEW" {
		return cameraType + "_camera_setting_facebook.xml"
	}
	algoBase := strings.Split(algo, "_")[0]
	return fmt.Sprintf("%s_camera_setting_%s.xml", cameraType, algoBase)
}

func (cg *CmdGenerator) getCmdOpts() map[string]string {
	var videoDir = cg.job.VideoDir
	var startFrame = cg.segOptions.StartFrame
	if cg.segOptions.FrameAt > cg.segOptions.StartFrame {
		startFrame = cg.segOptions.FrameAt
	}
	var opts = map[string]string{
		"video_dir":           videoDir,
		"output_dir":          cg.getFinalOuptutDir(),
		"ring_rectify_file":   path.Join(videoDir, "ring_rectify.xml"),
		"top_rectify_file":    path.Join(videoDir, "top_rectify.xml"),
		"bottom_rectify_file": path.Join(videoDir, "bottom_rectify.xml"),
		"mix_rectify_file":    path.Join(videoDir, "mix_rectify.xml"),
		"camera_setting_file": path.Join(cg.settingDirectory, cg.getCameraSettingFileName()),
		"enable_top":          cg.job.EnableTop,
		"enable_bottom":       cg.job.EnableBottom,
		"start_frame":         strconv.Itoa(startFrame),
		"end_frame":           strconv.Itoa(cg.segOptions.EndFrame),
		"time_alignment_file": path.Join(videoDir, "time.txt"),
		"save_type":           cg.getSaveType(),
		"thread_num":          strconv.Itoa(cg.threadNum),
	}
	return opts
}

// GetCmd get a exec.Cmd
func (cg *CmdGenerator) GetCmd() *exec.Cmd {
	args := util.MapToCmdArgs(cg.getCmdOpts(), "-")
	return exec.Command(path.Join(cg.binDirecotry, cg.getBinName()), args...)
}

func (cg *CmdGenerator) GetCmdLine() string {
	args := util.MapToCmdArgs(cg.getCmdOpts(), "-")
	executable := path.Join(cg.binDirecotry, cg.getBinName())
	return executable + " " + strings.Join(args, " ")
}
