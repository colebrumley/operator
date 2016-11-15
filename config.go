package operator

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/olebedev/config"
)

// GetConfig loads the configuration from files, flags, or env vars
func GetConfig() (cfg *config.Config, err error) {
	dirname, _ := os.Getwd()
	files, err := ioutil.ReadDir(dirname)
	if err != nil {
		log.Error(err)
	}

	allConfigs := []*config.Config{}
	for _, f := range files {
		if !isCfgFile(f.Name()) {
			continue
		}
		var c *config.Config
		if strings.HasSuffix(f.Name(), ".yml") || strings.HasSuffix(f.Name(), ".yaml") {
			c, err = config.ParseYamlFile(f.Name())
			if err != nil {
				return
			}
			allConfigs = append(allConfigs, c)
		}

		if strings.HasSuffix(f.Name(), ".json") {
			c, err = config.ParseJsonFile(f.Name())
			if err != nil {
				return
			}
			allConfigs = append(allConfigs, c)
		}
	}

	allConfigs = append(allConfigs, &config.Config{
		Root: DefaultConfig,
	})

	cfg = combineConfigs(allConfigs...)

	return
}

// PrettyPrintFlagMap prints the list of flags for use in the --help message
func PrettyPrintFlagMap(m map[string]interface{}, prefix []string) {
	for k, v := range m {
		flagName := "-" + k
		if len(prefix) > 0 {
			flagName = "-" + strings.Join(prefix, "-") + flagName
		}
		switch v.(type) {
		case string, int, bool:
			fmt.Printf("  %s=%+v\n", flagName, v)
		case map[string]interface{}:
			PrettyPrintFlagMap(v.(map[string]interface{}), append(prefix, k))
		}
	}
}

func combineConfigs(cfgs ...*config.Config) (r *config.Config) {
	r = nil
	for _, conf := range cfgs {
		for k, v := range conf.Root.(map[string]interface{}) {
			if r == nil {
				r = &config.Config{
					Root: map[string]interface{}{},
				}
			}
			r.Root.(map[string]interface{})[k] = v
		}
	}
	return
}

func isCfgFile(path string) bool {
	file, err := os.Open(path)
	if err != nil {
		return false
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if scanner.Scan() && scanner.Text() == "#operator" {
		return true
	}
	return false
}
