package apiHandler

import (
	"fmt"
	"horus/config"
	"horus/internal"
	"io/ioutil"
	"os/exec"
	"strconv"
	"strings"

	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

const CUR_RASP = true

func HandleLogin(c *gin.Context) {
	session := sessions.Default(c)

	username := c.PostForm("Username")
	password := c.PostForm("Password")

	hasError := false

	if config.UserConfiguration.Security.UserInput {
		if username != config.UserConfiguration.UserInfo.Username {
			hasError = true
		}
	}

	err := bcrypt.CompareHashAndPassword([]byte(config.UserConfiguration.UserInfo.Password), []byte(password))
	if err != nil {
		hasError = true
	}

	if hasError {
		c.JSON(http.StatusUnauthorized, gin.H{
			"Status": "error",
		})
	} else {
		session.Clear()
		session.Set("LoggedIn", true)
		session.Save()
		c.JSON(http.StatusOK, gin.H{
			"Status": "OK",
		})
	}
}
func HandleLogout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
	c.Redirect(http.StatusTemporaryRedirect, "/")
}
func HandleGetStats(c *gin.Context) {
	var temperature float64
	var cpuUsage int
	var ramUsage int
	var diskUsage [2]float64
	var sysUptime float64

	if CUR_RASP {
		temperature, _ = getTemp()
		cpuUsage, _ = getCpuUsage()
		ramUsage, _ = getRamUsage()
		diskUsage, _ = getDiskUsage()
		sysUptime, _ = getSysUptime()
	} else {
		temperature = float64(internal.RandomValue(0, 85))
		cpuUsage = internal.RandomValue(0, 100)
		ramUsage = internal.RandomValue(0, 100)
		diskUsage = [2]float64{float64(internal.RandomValue(0, 120000)), 120000}
		sysUptime = float64(internal.RandomValue(0, 100000))
	}

	fmt.Println(gin.H{ // TODO : Change to receive real data. Placeholder for now.
		"Temperature": temperature,
		"CPU":         cpuUsage,
		"RAM":         ramUsage,
		"Disk":        diskUsage[0],
		"DiskMax":     diskUsage[1],
		"Uptime":      sysUptime,
	})
	c.JSON(http.StatusOK, gin.H{ // TODO : Change to receive real data. Placeholder for now.
		"Temperature": temperature,
		"CPU":         cpuUsage,
		"RAM":         ramUsage,
		"Disk":        diskUsage[0],
		"DiskMax":     diskUsage[1],
		"Uptime":      sysUptime,
	})
}

func getTemp() (float64, error) {
	cmd := exec.Command("vcgencmd", "measure_temp")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return 0, err
	}

	value := strings.TrimPrefix(string(output), "temp=")
	value = strings.Replace(value, "'C", "", 1)
	value = strings.TrimSpace(value)

	valueF64, _ := strconv.ParseFloat(value, 64)
	return valueF64, nil
}
func getCpuUsage() (int, error) {
	return 0, nil
}
func getRamUsage() (int, error) {
	return 0, nil
}
func getDiskUsage() ([2]float64, error) {
	cmd := exec.Command("df", "-h", "/")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return [2]float64{-1, -1}, err
	}

	lines := strings.Split(string(output), "\n")
	dataLine := lines[1]

	fields := strings.Fields(dataLine)

	totalSize := strings.Replace(fields[1], "G", "", 1)
	usedSize := strings.Replace(fields[2], "G", "", 1)

	totalSizeF64, _ := strconv.ParseFloat(totalSize, 64)
	usedSizeF64, _ := strconv.ParseFloat(usedSize, 64)

	return [2]float64{usedSizeF64, totalSizeF64}, nil
}
func getSysUptime() (float64, error) {
	data, err := ioutil.ReadFile("/proc/uptime")
	if err != nil {
		return 0, err
	}

	uptime := strings.Fields(string(data))[0]

	uptimeF64, _ := strconv.ParseFloat(uptime, 64)

	return uptimeF64, nil
}
