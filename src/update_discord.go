package main

import (
	"fmt"
	"os"
	"os/exec"
	"net/http"
	"io"
)


// it will produce the progress bar
type Progress struct {
	Total int64
	Current int64
}

func (p *Progress) Write(b []byte) (int, error) {
	var length = len(b)
	p.Current += int64(length)

	fmt.Printf("\rDownloading... [%.2f%%]", float64(p.Current) * 100 / float64(p.Total))

	return length, nil
}

// function that will perform all the steps from download of file to installation
func updateDiscord(path string) error {
	var url string = "https://discord.com/api/download?platform=linux&format=deb"
	var response, err = http.Get(url) // url is default to the deb file on discord website (obviously, it updates discord not everything)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	output, err := os.Create(path)
	if err != nil {
		return err
	}
	defer output.Close()

	var size = response.ContentLength
	var progress *Progress = &Progress{Total: size}

	_, err = io.Copy(output, io.TeeReader(response.Body, progress))
	if err != nil {
		return err
	}

	var cmd *exec.Cmd = exec.Command("sudo", "dpkg", "i", path)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		return err
	}

	fmt.Println("Update was successful :D")
	return nil
}

func homePath() string {
	return os.Getenv("HOME")
}

func main() {
	if os.Geteuid() != 0 {
		fmt.Println("You must run this program as root")
		return
	}

	var path string = fmt.Sprintf("%s/%s", homePath(), "discord.deb")
	var err error = updateDiscord(path)
	if err != nil {
		fmt.Println("Discord update was unsuccessful")
		return
	}

	os.Remove(path)

	return
}
