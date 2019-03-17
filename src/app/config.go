package app

import (
	"io/ioutil"
	"math/rand"
	"strconv"
	"strings"
	"time"

	goconfig "github.com/hello-bytes/goconfig"
)

var _iniContent *map[string]interface{}

type WebConfig struct {
	_themeName string
	_httpPort  int
	_cdn       string

	_isDebug bool

	_iniReader *goconfig.IniReader
}

var _config *WebConfig

func InitConfig() {
	_config = &WebConfig{}
	_config.Init()
}

func WebConfigInst() *WebConfig {
	return _config
}

func (self *WebConfig) Init() {
	self._httpPort = 80
	self._themeName = "gitbook"
	self._isDebug = false

	path := currentDirectory() + "/resource/config/env.config"
	iniReader := goconfig.NewIniReader()
	if iniReader.LoadIni(path) {
		self._iniReader = iniReader
		websiteConfigMap := iniReader.GetSectionValues("WebSite")
		self.parseWebSite(websiteConfigMap)
	}
}

func (self *WebConfig) Port() int {
	return self._httpPort
}

func (self *WebConfig) CDN() string {
	return self._cdn
}

func (self *WebConfig) ThemeName() string {
	return self._themeName
}

func (self *WebConfig) IsDebug() bool {
	return self._isDebug
}

func (self *WebConfig) GetSectionValues(section string) map[string]string {
	return self._iniReader.GetSectionValues(section)
}

func (self *WebConfig) parseWebSite(keyVals map[string]string) {
	for k, v := range keyVals {
		if strings.Compare(k, "port") == 0 {
			port, err := strconv.Atoi(v)
			if err == nil {
				self._httpPort = port
			}
		} else if strings.Compare(k, "theme") == 0 {
			self._themeName = v
		} else if strings.Compare(k, "cdn") == 0 {
			self._cdn = v
		} else if strings.Compare(k, "debug") == 0 {
			self._isDebug = strings.Compare(v, "1") == 0 || strings.Compare(v, "true") == 0
		}
	}
}

var _securityKey string
var _randomKey = "QWERTYUIOPASDFGHJKLZXCVBNMqwertyuiopasdfghjklzxcvbnm1234567890"

func randomKey(length int) string {
	result := ""
	r := rand.New(rand.NewSource(time.Now().Unix()))
	rkl := len(_randomKey)
	for i := 0; i < length; i++ {
		randInt := r.Intn(rkl)
		result = result + _randomKey[randInt:randInt+1]
	}
	return result
}

func MakeSecurityKey() {
	if HasSecurityKey() {
		return
	}

	sk := randomKey(48)

	filePath := Path("/resource/config/sk.config")
	ioutil.WriteFile(filePath, []byte(sk), 0644)
}

func SecurityKey() string {
	if len(_securityKey) > 0 {
		return _securityKey
	}

	filePath := Path("/resource/config/sk.config")
	skbin, err := ioutil.ReadFile(filePath)
	if err == nil {
		_securityKey = string(skbin)
	}

	return _securityKey
}

func SkFile() string {
	return Path("/resource/config/sk.config")
}

func HasSecurityKey() bool {
	return len(SecurityKey()) > 0
}
