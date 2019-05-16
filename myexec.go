package main

import (
	"fmt"
	"os/exec"
)

func main() {
    c:="ffmpeg -y -i test1.aac -acodec pcm_s16le -f s16le -ac 1 -ar 16000 test2.pcm"
	cmd := exec.Command("sh", "-c",c)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(out))
}

