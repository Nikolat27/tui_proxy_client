package tui

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
)

// disconnect stops the current connection by killing processes on port 1080 (PID method only)
func (tui *TUI) disconnect() {
	go func() {
		var logs []string
		var finalStatus string
		var finalColor tcell.Color

		logs = append(logs, "Checking port 1080 status...")

		// If port is free — nothing to do
		if !tui.isPort1080InUse() {
			tui.handleNoPortUse(&logs, &finalStatus, &finalColor)
			return
		}

		logs = append(logs, "Attempting to disconnect using PID method...")
		pids := tui.getPIDsOnPort1080()

		if len(pids) == 0 {
			tui.handleNoProcessesFound(&logs, &finalStatus, &finalColor)
			return
		} else {
			logs = append(logs, fmt.Sprintf("Found %d process(es) on port 1080. Attempting to kill...", len(pids)))
			success, pidLogs := tui.killProcessesByPID(context.Background(), pids)
			logs = append(logs, pidLogs...)

			if success {
				tui.isConnected = false
				clientType := tui.clientType
				tui.clientType = ""
				tui.connectedConfig = ""
				finalStatus = fmt.Sprintf("Disconnected from %s (port 1080 freed)", clientType)
				finalColor = tcell.ColorGreen
			} else {
				finalStatus = "Failed to free port 1080 - check logs"
				finalColor = tcell.ColorRed
			}
		}

		// Single UI update
		tui.app.QueueUpdateDraw(func() {
			tui.configText.SetText(strings.Join(logs, "\n"))
			tui.updateStatus(finalStatus, finalColor)
			if !tui.isConnected {
				tui.connectionStatus.SetText("Status: Not Connected (Port 1080 Free)").SetTextColor(tcell.ColorRed)
			}
		})
	}()
}

// killProcessesByPID safely kills processes by their PIDs with context (no UI calls here)
func (tui *TUI) killProcessesByPID(ctx context.Context, pids []string) (bool, []string) {
	var pidLogs []string
	successCount := 0

	for _, pid := range pids {
		pid = strings.TrimSpace(pid)
		if pid == "" {
			continue
		}

		killCmd := exec.CommandContext(ctx, "kill", pid)
		err := killCmd.Run()

		if err == nil {
			successCount++
			pidLogs = append(pidLogs, fmt.Sprintf("Successfully killed process %s", pid))
		} else {
			pidLogs = append(pidLogs, fmt.Sprintf("Failed to kill process %s: %v", pid, err))
		}
	}

	time.Sleep(1 * time.Second)
	return !tui.isPort1080InUse(), pidLogs
}

// getPIDsOnPort1080 returns the PIDs of processes using port 1080
func (tui *TUI) getPIDsOnPort1080() []string {
	var pids []string

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Fallback to lsof if netstat gave nothing
	cmd := exec.CommandContext(ctx, "lsof", "-ti:1080")
	output, err := cmd.Output()
	if err != nil {
		return pids
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	for _, line := range lines {
		if line = strings.TrimSpace(line); line != "" {
			pids = append(pids, line)
		}
	}

	return pids
}

func (tui *TUI) handleNoPortUse(logs *[]string, finalStatus *string, finalColor *tcell.Color) {
	*logs = append(*logs, "Port 1080 is not in use.\nNo active connections found.")
	tui.resetConnectionState()
	*finalStatus = "Port 1080 is not in use - nothing to disconnect"
	*finalColor = tcell.ColorYellow
	tui.updateDisconnectUI(*logs, *finalStatus, *finalColor)
}

// handleNoProcessesFound handles the "no PIDs found" case
func (tui *TUI) handleNoProcessesFound(logs *[]string, finalStatus *string, finalColor *tcell.Color) {
	*logs = append(*logs, "No processes found using port 1080.")
	*finalStatus = "No processes to kill"
	*finalColor = tcell.ColorYellow
	tui.updateDisconnectUI(*logs, *finalStatus, *finalColor)
}

func (tui *TUI) resetConnectionState() {
	tui.isConnected = false
	tui.clientType = ""
	tui.connectedConfig = ""
}

func (tui *TUI) updateDisconnectUI(logs []string, status string, color tcell.Color) {
	tui.app.QueueUpdateDraw(func() {
		tui.configText.SetText(strings.Join(logs, "\n"))
		tui.updateStatus(status, color)
		if !tui.isConnected {
			tui.connectionStatus.SetText("Status: Not Connected (Port 1080 Free)").
				SetTextColor(tcell.ColorRed)
		}
	})
}
