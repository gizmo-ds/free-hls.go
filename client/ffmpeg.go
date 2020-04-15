package client

import (
	"bytes"
	"os"
	"os/exec"
	"strconv"

	"free-hls.go/utils"
)

var (
	FFmpegBin  = ""
	FFprobeBin = ""
)

func init() {
	var err error
	FFprobeBin, err = exec.LookPath("ffprobe")
	if err != nil {
		panic(err)
	}
	FFmpegBin, err = exec.LookPath("ffmpeg")
	if err != nil {
		panic(err)
	}
}

// 获取视频总秒数
func getVideoDuration(filename string) int64 {
	cmd := exec.Command(FFprobeBin,
		"-i", filename,
		"-show_entries", "format=duration",
		"-v", "quiet",
		"-of", "default=noprint_wrappers=1:nokey=1",
	)
	var b bytes.Buffer
	cmd.Stdout = &b
	if err := cmd.Run(); err != nil {
		panic(err)
	}
	data, err := utils.ToUnix(b.Bytes())
	if err != nil {
		panic(err)
	}
	data = bytes.ReplaceAll(data, []byte{'\n'}, nil)
	lengthF, err := strconv.ParseFloat(string(data), 64)
	if err != nil {
		panic(err)
	}
	return int64(lengthF)
}

// 获取视频编码
func getVideoCodec(filename string) string {
	cmd := exec.Command(FFprobeBin,
		"-v", "error",
		"-select_streams", "v:0",
		"-show_entries", "stream=codec_name",
		"-of", "default=noprint_wrappers=1:nokey=1",
		filename)
	var b bytes.Buffer
	cmd.Stdout = &b
	if err := cmd.Run(); err != nil {
		panic(err)
	}
	data, err := utils.ToUnix(b.Bytes())
	if err != nil {
		panic(err)
	}
	data = bytes.ReplaceAll(data, []byte{'\n'}, nil)
	return string(data)
}

// 分片视频
func FF(filename string, s int) error {
	_ = os.RemoveAll("./tmp")
	_ = os.MkdirAll("./tmp", 0666)
	cmd := exec.Command(FFmpegBin,
		"-i", filename,
		"-vcodec", getVideoCodec(filename),
		"-acodec", "aac", "-map", "0", "-f", "segment",
		"-segment_list", "tmp/out.m3u8",
		"-segment_time", "5", "tmp/out%05d.ts",
	)
	return cmd.Run()
}
