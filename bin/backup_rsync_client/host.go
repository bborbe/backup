package main

import "fmt"

type host struct {
	Active      bool   `json:"active"`
	User        string `json:"user"`
	Host        string `json:"host"`
	Port        int    `json:"port"`
	Directory   string `json:"dir"`
	ExcludeFrom string `json:"exclude_from"`
}

func (h *host) Backup(targetDirectory targetDirectory) error {
	return runRsync(
		"-azP",
		fmt.Sprintf("-e \"ssh -p %d\"", h.Port),
		"--delete",
		"--delete-excluded",
		fmt.Sprintf("--port=%d", h.Port),
		fmt.Sprintf("--exclude-from=%s", h.ExcludeFrom),
		fmt.Sprintf("--link-dest=%s", h.linkDest(targetDirectory)),
		h.from(),
		h.to(targetDirectory),
	)
}

func (h *host) from() string {
	return fmt.Sprintf("%s@%s:%s", h.User, h.Host, h.Directory)
}

func (h *host) to(targetDirectory targetDirectory) string {
	return fmt.Sprintf("%s/%s/incomplete%s", targetDirectory, h.Host, h.Directory)
}

func (h *host) linkDest(targetDirectory targetDirectory) string {
	return fmt.Sprintf("%s/%s/current/%s", targetDirectory, h.Host, h.Directory)
}

func (h *host) Validate() error {
	if len(h.User) == 0 {
		return fmt.Errorf("user invalid")
	}
	if len(h.Host) == 0 {
		return fmt.Errorf("host invalid")
	}
	if h.Port <= 0 {
		return fmt.Errorf("port invalid")
	}
	if len(h.Directory) == 0 {
		return fmt.Errorf("directory invalid")
	}
	return nil
}
