package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type connectionInfo struct {
	proto       string
	pid         int
	localAddr   string
	processName string
}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "--kill" {
		if len(os.Args) > 2 {
			killPort(os.Args[2])
		} else {
			showPortsAndKill()
		}
	} else {
		listAllPorts()
	}
}

func listAllPorts() {
	fmt.Println("Active network connections:")
	fmt.Printf("%-6s %-8s %-20s %-10s\n", "Proto", "PID", "Local Address", "Process")

	connections := getAllConnections()
	if len(connections) == 0 {
		fmt.Println("No active network connections found.")
		return
	}

	for _, conn := range connections {
		fmt.Printf("%-6s %-8d %-20s %-10s\n",
			conn.proto,
			conn.pid,
			conn.localAddr,
			conn.processName)
	}
}

func showPortsAndKill() {
	fmt.Println("Scanning active ports...")
	connections := getAllConnections()

	if len(connections) == 0 {
		fmt.Println("No active network connections found.")
		return
	}

	fmt.Println("\nActive network connections:")
	fmt.Printf("%-6s %-8s %-20s %-10s\n", "Proto", "PID", "Local Address", "Process")
	for i, conn := range connections {
		fmt.Printf("[%2d] %-6s %-8d %-20s %-10s\n",
			i+1,
			conn.proto,
			conn.pid,
			conn.localAddr,
			conn.processName)
	}

	fmt.Print("\nEnter number to kill (or 'q' to quit): ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input == "q" {
		return
	}

	selection, err := strconv.Atoi(input)
	if err != nil || selection < 1 || selection > len(connections) {
		fmt.Println("Invalid selection")
		return
	}

	target := connections[selection-1]
	fmt.Printf("Killing %s (PID %d) on %s...\n",
		target.processName,
		target.pid,
		target.localAddr)

	killProcess(target.pid)
}

func getAllConnections() []connectionInfo {
	var connections []connectionInfo
	if isWindows() {
		connections = append(connections, getWindowsConnections("tcp")...)
		connections = append(connections, getWindowsConnections("udp")...)
	} else {
		connections = append(connections, getUnixConnections("tcp")...)
		connections = append(connections, getUnixConnections("udp")...)
	}
	return connections
}

func getUnixConnections(netType string) []connectionInfo {
	var connections []connectionInfo

	cmd := exec.Command("lsof", "-i", netType, "-P", "-n", "-l")
	output, err := cmd.Output()
	if err != nil {
		return connections
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines[1:] {
		if line == "" {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 9 {
			continue
		}

		pid, err := strconv.Atoi(fields[1])
		if err != nil {
			continue
		}

		connections = append(connections, connectionInfo{
			proto:       strings.ToUpper(netType),
			pid:         pid,
			localAddr:   fields[8],
			processName: fields[0],
		})
	}

	return connections
}

func getWindowsConnections(netType string) []connectionInfo {
	var connections []connectionInfo

	cmd := exec.Command("netstat", "-ano", "-p", netType)
	output, err := cmd.Output()
	if err != nil {
		return connections
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if !strings.Contains(line, "LISTENING") && !strings.Contains(line, "ESTABLISHED") {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 5 {
			continue
		}

		// Extract PID (last field)
		pid, err := strconv.Atoi(fields[len(fields)-1])
		if err != nil {
			continue
		}

		// Get process name
		processName := getWindowsProcessName(pid)

		connections = append(connections, connectionInfo{
			proto:       strings.ToUpper(netType),
			pid:         pid,
			localAddr:   fields[1],
			processName: processName,
		})
	}

	return connections
}

func getWindowsProcessName(pid int) string {
	cmd := exec.Command("tasklist", "/FI", fmt.Sprintf("PID eq %d", pid), "/FO", "CSV", "/NH")
	output, err := cmd.Output()
	if err != nil {
		return "unknown"
	}

	lines := strings.Split(string(output), "\n")
	if len(lines) == 0 {
		return "unknown"
	}

	fields := strings.Split(lines[0], ",")
	if len(fields) < 1 {
		return "unknown"
	}

	return strings.Trim(fields[0], "\"")
}

func killPort(port string) {
	port = normalizePort(port)
	connections := getAllConnections()
	var pids []int

	for _, conn := range connections {
		if strings.HasSuffix(conn.localAddr, port) {
			pids = append(pids, conn.pid)
		}
	}

	if len(pids) == 0 {
		fmt.Printf("No processes found using port %s\n", port)
		return
	}

	for _, pid := range pids {
		killProcess(pid)
	}
}

func killProcess(pid int) {
	var cmd *exec.Cmd
	if isWindows() {
		cmd = exec.Command("taskkill", "/PID", strconv.Itoa(pid), "/F")
	} else {
		cmd = exec.Command("kill", "-9", strconv.Itoa(pid))
	}

	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error killing process %d: %v\n", pid, err)
	} else {
		fmt.Printf("Successfully killed process %d\n", pid)
	}
}

func normalizePort(port string) string {
	if !strings.HasPrefix(port, ":") {
		port = ":" + port
	}
	return port
}

func isWindows() bool {
	return os.Getenv("OS") == "Windows_NT"
}