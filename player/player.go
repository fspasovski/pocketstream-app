package player

import (
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"syscall"

	"github.com/fspasovski/pocketstream-app/config"
)

type Player struct {
	Cfg                      *config.Config
	Process                  *exec.Cmd
	BroadcasterStreamingUrls map[string]string
}

func (p *Player) Play(broadcasterLogin string) error {
	streamUrl, exists := p.BroadcasterStreamingUrls[broadcasterLogin]
	if !exists {
		streamUrl, _ = p.Cfg.TwitchService.GetStreamingUrl(broadcasterLogin)
		p.BroadcasterStreamingUrls[broadcasterLogin] = streamUrl
	}
	var cmd *exec.Cmd
	cmd = exec.Command("ffplay",
		"-vf", fmt.Sprintf("scale=%d:%d", p.Cfg.Player.StreamWidth, p.Cfg.Player.StreamHeight),
		"-window_title", "fst",
		"-autoexit",
		"-x", strconv.Itoa(p.Cfg.Player.StreamWidth),
		"-y", strconv.Itoa(p.Cfg.Player.StreamHeight),
		streamUrl,
	)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	if err := cmd.Start(); err != nil {
		log.Printf("ffplay exited with error: %v", err)
		return err
	}

	p.Process = cmd
	return nil
}

func (p *Player) Stop() {
	if p.Process != nil {
		pgid, err := syscall.Getpgid(p.Process.Process.Pid)
		if err == nil {
			syscall.Kill(-pgid, syscall.SIGKILL)
		} else {
			p.Process.Process.Kill()
		}
		p.Process.Wait()
		p.Process = nil
	}
}

func (p *Player) IsPlaying() bool {
	return p.Process != nil
}
