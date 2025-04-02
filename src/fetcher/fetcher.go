package fetcher

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// --------------------------------------------------------------------------------------------------------------------

type System struct {
	Distribution string
	Hostname     string
	Address      string
	Uptime       string
	Kernel       string
}

// --------------------------------------------------------------------------------------------------------------------

type CronJob struct {
	Description string
	Schedule    string
}

// --------------------------------------------------------------------------------------------------------------------

type Tasks struct {
	Lines []string
}

// --------------------------------------------------------------------------------------------------------------------

type Information struct {
	Sys  *System
	Cron []CronJob
	Task *Tasks
}

// --------------------------------------------------------------------------------------------------------------------

func getDistribution() string {
	file, err := os.Open("/etc/os-release")
	if err != nil {
		return "Unknown"
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "PRETTY_NAME=") {
			return strings.Trim(line[13:], "\"")
		}
	}
	return "Unknown"
}

// --------------------------------------------------------------------------------------------------------------------

func getHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		return "Unknown"
	}
	return hostname
}

// --------------------------------------------------------------------------------------------------------------------

// TODO: This function will cause issues when there are multiple NICs on a system
func getAddress() string {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "Unknown"
	}

	for _, iface := range interfaces {

		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		addresses, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, address := range addresses {
			switch ip := address.(type) {
			case *net.IPNet:
				if ip.IP.To4() != nil {
					return ip.String()
				}
			}
		}
	}
	return "127.0.0.1"
}

// --------------------------------------------------------------------------------------------------------------------

func getUptime() string {
	uptimeBytes, err := os.ReadFile("/proc/uptime")
	if err != nil {
		return "Unknown"
	}

	uptimeStr := strings.Split(string(uptimeBytes), " ")[0]
	uptimeFloat, err := strconv.ParseFloat(uptimeStr, 64)
	if err != nil {
		return "Unknown"
	}

	uptime := int(uptimeFloat)
	hours := uptime / 3600
	minutes := (uptime % 3600) / 60
	days := hours / 24
	hours = hours % 24
	return fmt.Sprintf("%dd %dh %dm", days, hours, minutes)
}

// --------------------------------------------------------------------------------------------------------------------

func getKernel() string {
	cmd := exec.Command("uname", "-r")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return "Unknown"
	}
	return strings.TrimSpace(out.String())
}

// --------------------------------------------------------------------------------------------------------------------

func getCronjobs() []CronJob {
	cmd := exec.Command("crontab", "-l")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return nil
	}

	var jobs []CronJob
	lines := strings.Split(out.String(), "\n")
	var currentDescription string = "Sample Task"

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "#") {
			currentDescription = line[1:]
		} else if line != "" {
			parts := strings.Fields(line)
			if len(parts) >= 5 {
				schedule := strings.Join(parts[:5], " ")
				jobs = append(jobs, CronJob{
					Description: strings.TrimSpace(currentDescription),
					Schedule:    schedule,
				})
			}
			currentDescription = "Task"
		}
	}
	return jobs
}

// --------------------------------------------------------------------------------------------------------------------

func getTasks(filename string) *Tasks {
	file, err := os.Open(filename)
	if err != nil {
		return nil
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return &Tasks{Lines: lines}
}

// --------------------------------------------------------------------------------------------------------------------

func Fetch() *Information {
	sys := System{
		Distribution: getDistribution(),
		Hostname:     getHostname(),
		Address:      getAddress(),
		Uptime:       getUptime(),
		Kernel:       getKernel(),
	}

	jobs := getCronjobs()
	tasks := getTasks(os.Getenv("HOME") + "/.info")

	return &Information{&sys, jobs, tasks}
}
