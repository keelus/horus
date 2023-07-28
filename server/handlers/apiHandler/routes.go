package apiHandler

import (
	"encoding/base64"
	"fmt"
	"horus/config"
	"horus/internal"
	"horus/logger"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"net/http"

	"github.com/fatih/color"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

var previousCpuIdle, previousCpuTotal int

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
		if config.UserConfiguration.Security.UserInput {
			logger.Log(c, logger.LOGIN, fmt.Sprintf("Unsuccessful login attempt [username='%s'].", username))
		} else {
			logger.Log(c, logger.LOGIN, "Unsuccessful login attempt.")
		}
		c.JSON(http.StatusUnauthorized, gin.H{
			"Status": "error",
		})
	} else {
		session.Clear()
		session.Set("LoggedIn", true)
		session.Save()
		logger.Log(c, logger.LOGIN, "Successfull login.")
		c.JSON(http.StatusOK, gin.H{
			"Status": "OK",
		})
	}
}
func HandleLogout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
	logger.Log(c, logger.LOGIN, "User logout.")
	c.Redirect(http.StatusTemporaryRedirect, "/")
}

func HandleGetStats(c *gin.Context) {
	if !internal.IsLogged(c) {
		c.JSON(http.StatusForbidden, gin.H{"details": "User not logged in."})
		return
	}

	var temperature float64
	var cpuUsage float64
	var ramUsage float64
	var diskUsage [2]float64
	var sysUptime float64

	if CUR_RASP { // TODO: Move functions to a single file dedicated for that only.
		temperature, _ = getTemp()
		cpuUsage, _ = getCpuUsage()
		ramUsage, _ = getRamUsage()
		diskUsage, _ = getDiskUsage()
		sysUptime, _ = getSysUptime()
	} else {
		temperature = float64(internal.RandomValue(0, 85))
		cpuUsage = float64(internal.RandomValue(0, 100))
		ramUsage = float64(internal.RandomValue(0, 100))
		diskUsage = [2]float64{float64(internal.RandomValue(0, 120000)), 120000}
		sysUptime = float64(internal.RandomValue(0, 100000))
	}

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
func getCpuUsage() (float64, error) {
	data, err := ioutil.ReadFile("/proc/stat")
	if err != nil {
		return 0, err
	}

	lines := strings.Split(string(data), "\n")
	dataLineB := lines[0]

	dataLine := strings.Fields(string(dataLineB))

	idle, _ := strconv.Atoi(dataLine[4])

	total := 0
	for i := 1; i < len(dataLine); i++ {
		valueInt, _ := strconv.Atoi(dataLine[i])
		total += valueInt
	}

	percent := (1 - float64(idle-previousCpuIdle)/float64(total-previousCpuTotal)) * 100

	previousCpuIdle = idle
	previousCpuTotal = total

	return percent, nil
}

func getRamUsage() (float64, error) {
	data, err := ioutil.ReadFile("/proc/meminfo")
	if err != nil {
		return 0, err
	}

	total := 0
	free := 0

	dataLines := strings.Split(string(data), "\n")

	for _, dataLine := range dataLines {
		fields := strings.Fields(dataLine)

		if len(fields) < 2 {
			continue
		}

		valueInt, err := strconv.Atoi(fields[1])
		if err != nil { // Doesnt contain an number
			continue
		}

		if fields[0] == "MemTotal:" {
			total = valueInt
		} else if fields[0] == "MemFree:" {
			free = valueInt
		}
	}

	used := total - free

	percent := float64(used*100) / float64(total)

	return percent, nil
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

func HandleAvatar(c *gin.Context) {
	allowedFormats := map[string]string{
		"jpg":  "image/jpeg",
		"png":  "image/png",
		"webp": "image/webp",
		"gif":  "image/gif",
	}

	for extension, format := range allowedFormats {
		if fileExists(fmt.Sprintf("./avatar.%s", extension)) {
			avatarData, err := ioutil.ReadFile("avatar.jpg")
			if err == nil {
				c.Data(http.StatusOK, format, avatarData)
				return
			}
		}
	}

	// If no existing avatar, use default one:
	avatarData, err := defaultAvatar()
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Data(http.StatusOK, "image/jpeg", avatarData)
}

func HandleLatestVersion(c *gin.Context) {
	internal.CheckLatestVersion()
	c.JSON(http.StatusOK, gin.H{
		"CurrentVersion":     internal.VERSION_CURRENT,
		"LatestVersion":      internal.VERSION_LAST,
		"UsingLatestVersion": internal.VERSION_UPDATED,
	})
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func defaultAvatar() ([]byte, error) {
	cError := color.New(color.FgRed, color.Bold)

	defaultAvatarB64 := "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"
	defaultAvatarDecoded, err := base64.StdEncoding.DecodeString(defaultAvatarB64)
	if err != nil {
		cError.Println("Error decoding default avatar.")
		return []byte{}, err
	}
	return defaultAvatarDecoded, nil
}
