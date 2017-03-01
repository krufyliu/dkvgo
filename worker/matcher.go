package worker

import (
	"bufio"
	"regexp"

	"io"

	"strconv"

	"github.com/krufyliu/dkvgo/task"
)

var (
	frameMatcher         = regexp.MustCompile(`Process frame: ([\d]+)`)
	prepareTimeMatcher   = regexp.MustCompile(`Prepare images time: ([\d]+\.[\d]+)`)
	flowTimeMatcher      = regexp.MustCompile(`Compute flow time: ([\d]+\.[\d]+)`)
	novelViewTimeMatcher = regexp.MustCompile(`Compute novel view time: ([\d]+\.[\d]+)`)
	totalTimeMatcher     = regexp.MustCompile(`Total time: ([\d]+\.[\d]+)`)
)

func matchState(reader *bufio.Reader) (*task.RunState, error) {
	var state = new(task.RunState)
	var (
		line                 string
		err                  error
		frameAt              int
		prepareImagesTime    float64
		computeFlowTime      float64
		computeNovelViewTime float64
		totalTime            float64
	)
	for {
		line, err = reader.ReadString('\n')
		if err != nil && err != io.EOF {
			return nil, err
		}
		if err == io.EOF {
			break
		}
		if matches := frameMatcher.FindStringSubmatch(line); matches != nil {
			frameAt, _ = strconv.Atoi(matches[1])
		} else if matches := prepareTimeMatcher.FindStringSubmatch(line); matches != nil {
			prepareImagesTime, _ = strconv.ParseFloat(matches[1], 32)
		} else if matches := flowTimeMatcher.FindStringSubmatch(line); matches != nil {
			computeFlowTime, _ = strconv.ParseFloat(matches[1], 32)
		} else if matches := novelViewTimeMatcher.FindStringSubmatch(line); matches != nil {
			computeNovelViewTime, _ = strconv.ParseFloat(matches[1], 32)
		} else if matches := totalTimeMatcher.FindStringSubmatch(line); matches != nil {
			totalTime, _ = strconv.ParseFloat(matches[1], 32)
			break
		}
	}
	state.FrameAt = frameAt
	state.PrepareImagesTime = float32(prepareImagesTime)
	state.ComputeFlowTime = float32(computeFlowTime)
	state.ComputeNovelViewTime = float32(computeNovelViewTime)
	state.TotalTime = float32(totalTime)
	return state, err
}
