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
	"3D_AURA":  "test_3d_aura",
	"3D_GOPRO": "test_3d_gopro",
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
	// if strings.ToUpper(cg.job.Algorithm) == "TOP_BOTTOM" {
	// 	return outputDir
	// }
	return path.Join(outputDir, cg.job.Algorithm)
}

func (cg *CmdGenerator) getCameraSettingFileName() string {
	cameraType := strings.ToLower(cg.job.CameraType)
	// algo := strings.ToLower(cg.job.Algorithm)
	// if strings.ToUpper(algo) == "PREVIEW" {
	// 	return cameraType + "_camera_setting_facebook.xml"
	// }
	// algoBase := strings.Split(algo, "_")[0]
	return fmt.Sprintf("%s/%s/camera_setting_default.xml", cg.settingDirectory, cameraType)
}

func (cg *CmdGenerator) getCameraSettingDir() string {
	cameraType := strings.ToLower(cg.job.CameraType)
	return fmt.Sprintf("%s/%s", cg.settingDirectory, cameraType)
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
		"time_alignment_file": path.Join(videoDir, "time.txt"),
		"ring_rectify_file":   path.Join(cg.settingDirectory, "ring_rectify.xml"),
		"top_rectify_file":    path.Join(videoDir, "top_rectify.xml"),
		"bottom_rectify_file": path.Join(videoDir, "bottom_rectify.xml"),
		"mix_rectify_file":    path.Join(videoDir, "mix_rectify.xml"),
		//"camera_setting_file": cg.getCameraSettingFileName(),
		"camera_setting_file_dir": cg.getCameraSettingDir(),
		"enable_top":              cg.job.EnableTop,
		"enable_bottom":           cg.job.EnableBottom,
		"save_debug_img":          cg.job.SaveDebugImg,
		"start_frame":             strconv.Itoa(startFrame),
		"end_frame":               strconv.Itoa(cg.segOptions.EndFrame),
		"save_type":               cg.getSaveType(),
		"thread_num":              strconv.Itoa(cg.threadNum),
	}
	return opts
}

// GetCmd get a exec.Cmd
func (cg *CmdGenerator) GetCmd() *exec.Cmd {
	fields := []string{
		"video_dir", "output_dir", "time_alignment_file",
		"camera_setting_file_dir", "ring_rectify_file",
		"top_rectify_file", "bottom_rectify_file", "mix_rectify_file",
		"start_frame", "end_frame", "enable_top", "enable_bottom",
		"save_type", "save_debug_img", "thread_num",
	}
	args := util.MapToCmdArgs(cg.getCmdOpts(), "-", fields...)
	return exec.Command(path.Join(cg.binDirecotry, cg.getBinName()), args...)
}

func (cg *CmdGenerator) GetCmdLine() string {
	args := util.MapToCmdArgs(cg.getCmdOpts(), "-")
	executable := path.Join(cg.binDirecotry, cg.getBinName())
	return executable + " " + strings.Join(args, " ")
}
