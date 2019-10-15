package config

import (
	"../common"
	"fmt"
	"github.com/sirupsen/logrus"
	"gopkg.in/ini.v1"
	"os"
	"path/filepath"
	"strings"
)

var Log = *logrus.New()

func LoadConfig(configLoaded bool, version string) {
	Log.Out = os.Stdout
	Log.Formatter = &logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "02-01-2006 15:04:05",
	}
	var files []string
	dir, err := os.Getwd()
	if err != nil {
		Log.Fatal("Could not find current working directory. Cannot continue.")
	}
	configDir := dir + "/config/"
	err = filepath.Walk(configDir, func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		if strings.Contains(file, ".ini") {
			cfg, err := ini.Load(file)
			if err != nil {
				fmt.Printf("Fail to read file: %v", err)
				os.Exit(1)
			}
			iniFile := cfg.Section("")
			ModuleName, err := iniFile.GetKey("ModuleName")
			if err != nil {

			} else {
				Log.Info("Loading module " + ModuleName.String() + " config from file " + file)
				configKeys := iniFile.KeyStrings()
				for _, v := range configKeys {
					shouldSave := false
					if strings.HasSuffix(v, ".required") {

						saveField := strings.Split(v, ".required")[0]
						safeKey, err := iniFile.GetKey(saveField)
						if err != nil {
							//Log.Warn(err)
						} else {
							shouldSave = true
							currentVal := safeKey.Value()
							if currentVal != "" {
								continue
							}
						}

						question, err := iniFile.GetKey(v + ".question")
						if err != nil {
							continue
						} else {
							var defaultString string
							stringDefault, err := iniFile.GetKey(v + ".default")

							if err != nil {
								defaultString = ""
							} else {
								defaultString = stringDefault.Value()
							}
							var questionText string
							if strings.Contains(question.Value(), "ZMVAR_AdminUserName") {
								SavedUsername, err := iniFile.GetKey("AdminUserName")
								if err != nil {
									logrus.Fatal(err)
								}
								questionText = strings.Replace(question.Value(), "ZMVAR_AdminUserName", SavedUsername.Value(), -1)
							} else {
								questionText = question.Value()
							}
							returnValue := common.PromptForInput("[Module: "+ModuleName.String()+"] "+questionText, false, defaultString)
							if shouldSave == true {
								safeKey.SetValue(returnValue)
								err := cfg.SaveTo(file)
								if err != nil {
									Log.Fatal(err)
								}
							}
						}
					}
				}
			}

		}
	}
}



