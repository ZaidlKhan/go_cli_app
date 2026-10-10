package parser

import (
	"bufio"
	"fmt"
	"main/internal/types"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

func GetProcessArgs(pid int) string {
	cmd := exec.Command("ps", "-p", strconv.Itoa(pid), "-o", "args=")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return ""
	}
	return string(out)
}

func ParseOutput(input string) map[int]types.Process {
	processes := make(map[int]types.Process)
	var currPID int
	var currName string
	var currAddress string
	var currPort uint16

	scanner := bufio.NewScanner(strings.NewReader(input))

	for scanner.Scan() {
		line := scanner.Text()
		if len(line) < 2 {
			continue
		}

		tag := line[0]
		val := line[1:]

		switch tag {
		case 'p':
			if currPID != 0 {
				processes[currPID] = types.Process{
					PID:         currPID,
					ProcessName: currName,
					Address:     currAddress,
					Port:        currPort,
				}
			}

			pid, err := strconv.Atoi(val)
			if err != nil {
				currPID = 0
				continue
			}
			currPID = pid
			currName = ""
			currAddress = ""
			currPort = 0

		case 'c':
			currName = val

		case 'n':
			if strings.Contains(val, ":") || strings.HasPrefix(val, "*") {
				currAddress = val
				parts := strings.Split(val, ":")
				if len(parts) > 1 {
					if p, err := strconv.ParseUint(parts[len(parts)-1], 10, 16); err == nil {
						currPort = uint16(p)
					}
				}
			}
		}
	}

	if currPID != 0 {
		processes[currPID] = types.Process{
			PID:         currPID,
			ProcessName: currName,
			Address:     currAddress,
			Port:        currPort,
		}
	}

	for pid, proc := range processes {
		cmdArgs := GetProcessArgs(pid)

		if !evalProcess(&proc, cmdArgs) {
			delete(processes, pid)
		} else {
			if proc.GitBranch == "" {
				proc.GitBranch = "N/A"
			}
			processes[pid] = proc
		}
	}

	return processes
}

func evalProcess(proc *types.Process, cmdArgs string) bool {
	if !inspectPath(cmdArgs) {
		return false
	}

	resolveGitInfo(proc, cmdArgs)
	resolveURL(proc)

	hasRuntime := false
	lowerArgs := strings.ToLower(cmdArgs)
	for _, r := range types.DevRuntimes {
		if strings.Contains(lowerArgs, r) {
			hasRuntime = true
			break
		}
	}

	hasDevPort := types.DevPorts[proc.Port]
	hasGit := proc.GitBranch != ""

	if hasGit && hasRuntime {
		return true
	}

	score := 0
	if hasGit {
		score += 40
	}
	if hasRuntime {
		score += 30
	}
	if hasDevPort {
		score += 10
	}
	if proc.Port > 1024 {
		score += 10
	}

	return score >= 40
}

func inspectPath(cmdArgs string) bool {
	path := strings.TrimSpace(cmdArgs)

	for _, prefix := range types.InvalidPrefixes {
		if strings.HasPrefix(path, prefix) {
			return false
		}
	}

	return true
}

func resolveGitInfo(proc *types.Process, cmdArgs string) {
	var repoRoot string

	if proc.PID != 0 {
		cmd := exec.Command("lsof", "-a", "-p", strconv.Itoa(proc.PID), "-d", "cwd", "-Fn")
		if out, err := cmd.Output(); err == nil {
			scanner := bufio.NewScanner(strings.NewReader(string(out)))
			for scanner.Scan() {
				line := scanner.Text()
				if strings.HasPrefix(line, "n/") {
					cwd := line[1:]
					curr := filepath.Clean(cwd)
					for curr != "/" && curr != "." {
						if _, err := os.Stat(filepath.Join(curr, ".git")); err == nil {
							repoRoot = curr
							break
						}
						curr = filepath.Dir(curr)
					}
					break
				}
			}
		}
	}

	if repoRoot == "" {
		fields := strings.Fields(cmdArgs)
		for _, field := range fields {
			if strings.HasPrefix(field, "/") {
				cleanPath := filepath.Clean(field)

				if info, err := os.Stat(cleanPath); err == nil && !info.IsDir() {
					cleanPath = filepath.Dir(cleanPath)
				}

				curr := cleanPath
				for curr != "/" && curr != "." {
					if _, err := os.Stat(filepath.Join(curr, ".git")); err == nil {
						repoRoot = curr
						break
					}
					curr = filepath.Dir(curr)
				}

				if repoRoot != "" {
					break
				}
			}
		}
	}

	if repoRoot != "" {
		cmd := exec.Command("git", "-C", repoRoot, "branch", "--show-current")
		if out, err := cmd.Output(); err == nil {
			proc.GitBranch = strings.TrimSpace(string(out))
		}
	}
}

func resolveURL(proc *types.Process) {
	if proc.Port == 0 {
		return
	}

	host := "localhost"
	if strings.HasPrefix(proc.Address, "127.0.0.1") || strings.HasPrefix(proc.Address, "[::1]") {
		parts := strings.Split(proc.Address, ":")
		if len(parts) > 0 && parts[0] != "" {
			host = parts[0]
		}
	}

	if strings.Contains(host, ":") && !strings.HasPrefix(host, "[") {
		proc.URL = fmt.Sprintf("http://[%s]:%d", host, proc.Port)
	} else {
		proc.URL = fmt.Sprintf("http://%s:%d", host, proc.Port)
	}
}
