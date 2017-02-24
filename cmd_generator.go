package dkvgo

import (
	"fmt"
	"os/exec"
	"path"
	"runtime"
	"strconv"
	"strings"
)

// BinMap define composition algorithm to executable
const BinMap = map[string]string{
	"VISIONDK_3D": "test_3d_visiondk",
	"VISIONDK_2D": "test_2d_visiondk",
	"FACEBOOK_3D": "test_3d_facebook",
	"FACEBOOK_2D": "test_2d_facebook",
	"PREVIEW":     "test_preview",
	"TOP_BOTTOM":  "test_top_and_bottom",
}

// CmdGenerator hold task options to generate exec.Cmd or shell command
type CmdGenerator struct {
	taskOptions      *TaskOptions
	binDirecotry     string
	settingDirectory string
}

func (cg *CmdGenerator) getBinName() string {
	return BinMap[cg.taskOptions.Algorithm]
}

func (cg *CmdGenerator) getSaveType() string {
	if strings.ToUpper(cg.taskOptions.CameraType) == "BMPCC" {
		return ".tiff"
	}
	return ".jpg"
}

func (cg *CmdGenerator) getFinalOuptutDir() string {
	outputDir := cg.taskOptions.OutputDir
	if strings.ToUpper(cg.taskOptions.Algorithm) == "TOP_BOTTOM" {
		return outputDir
	}
	return path.Join(outputDir, cg.taskOptions.Algorithm)
}

func (cg *CmdGenerator) getCameraSettingFileName() string {
	cameraType := strings.ToLower(cg.taskOptions.CameraType)
	algo := strings.ToLower(cg.taskOptions.Algorithm)
	if strings.ToUpper(algo) == "PREVIEW" {
		return cameraType + "_camera_setting_facebook.xml"
	}
	algoBase := strings.Split(algo, "_")[0]
	return fmt.Sprintf("%s_camera_setting_%s.xml", cameraType, algo)
}

func (cg *CmdGenerator) getCmdOpts() map[string]string {
	var videoDir = cg.taskOptions.VideoDir
	var opts = map[string]string{
		"video_dir":           videoDir,
		"output_dir":          cg.getFinalOuptutDir(),
		"ring_rectify_file":   path.Join(videoDir, "ring_rectify.xml"),
		"top_rectify_file":    path.Join(videoDir, "top_rectify.xml"),
		"bottom_rectify_file": path.Join(videoDir, "bottom_rectify.xml"),
		"mix_rectify_file":    path.Join(videoDir, "mix_rectify.xml"),
		"mix_rectify_file":    path.Join(cg.settingDirectory, cg.getCameraSettingFileName()),
		"enable_top":          cg.taskOptions.EnableTop,
		"enable_bottom":       cg.taskOptions.EnableBottom,
		"start_frame":         strconv.Itoa(cg.taskOptions.StartFrame),
		"end_frame":           strconv.Itoa(cg.taskOptions.EndFrame),
		"time_alignment_file": path.Join(videoDir, "time.txt"),
		"save_type":           cg.getSaveType(),
		"thread_num":          runtime.NumCPU(),
	}
}

func (cg *CmdGenerator) getCmd() *exec.Cmd {
	args = mapToCmdArgs(cg.getCmdOpts(), "-")
	return &exec.Command(path.Join(cg.binDirecotry, cg.getBinName), args...)
}

func (cg *CmdGenerator) getCmdLine() string {
	args = mapToCmdArgs(cg.getCmdOpts(), "-")
	executable = path.Join(cg.binDirecotry, cg.getBinName())
	return executable + " " + strings.Join(args, " ")
}
