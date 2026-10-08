package filesystems

import (
	"bufio"
	"fmt"
	"main/internal/types"
	"strings"
)

func ParseOutput(input string) []types.Process {
	var processes []types.Process
	var curr *types.Process

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
			if curr != nil {
				processes = append(processes, *curr)
			}
			curr = &types.Process{
				PID:   val,
				Ports: []string{},
			}

		case 'c':
			if curr != nil {
				curr.App = val
			}

		case 'n':
			if curr != nil {
				curr.Ports = append(curr.Ports, val)
			}

		case 'f':
			continue
		}
	}

	if curr != nil {
		processes = append(processes, *curr)
	}

	fmt.Println(processes)

	return processes
}
