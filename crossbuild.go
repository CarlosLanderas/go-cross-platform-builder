package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"sync"
)

type TargetConfig struct {
	Targets []Target
}

type CompileConfig struct {
	CompileTarget  string
	PlatformConfig string
}

type Target struct {
	Id   string
	Arch []string
}

var (
	targetFile       string
	platformsConfig  string
	waitSync         sync.WaitGroup
	binaryExtensions map[string]string
)

func compile(target Target, compileTarget string, ch chan bool) {

	for _, arch := range target.Arch {
		waitSync.Add(1)

		go func(arch string, target string) {

			fmt.Printf("Compiling platform %s with arch: %s\n", target, arch)

			env := os.Environ()
			env = append(env, fmt.Sprintf("GOOS=%s", target))
			env = append(env, fmt.Sprintf("GOARCH=%s", arch))

			targetAbsPath, _ := filepath.Abs(compileTarget)
			ext := binaryExtensions[target]
			fmt.Printf(ext)
			targetPath := fmt.Sprintf("%s%s", getCleanFileName(filepath.Base(targetAbsPath)), binaryExtensions[target])

			cmd := exec.Command("go", "build", "-o", fmt.Sprintf("%s\\%s\\%s",
				target, arch, targetPath), compileTarget)
			cmd.Env = env

			err := cmd.Run()
			if err != nil {
				fmt.Println(err.Error())
			}

			defer waitSync.Done()
		}(arch, target.Id)
	}

}

func getCleanFileName(fileName string) string {
	extension := filepath.Ext(fileName)
	return fileName[0 : len(fileName)-len(extension)]
}

func configureCLI() {
	flag.StringVar(&targetFile, "compile", "", "Target go file to compile")
	flag.StringVar(&platformsConfig, "config", "targets.json", "Json file with target platform configuration")
	flag.Parse()
}

func unmarshallConfig(file string, targetConfig *TargetConfig) {
	bytesConfig, _ := ioutil.ReadFile(file)
	json.Unmarshal(bytesConfig, targetConfig)
}

func ensureConfigsExist(config *CompileConfig) {

	config.PlatformConfig, _ = filepath.Abs(config.PlatformConfig)
	config.CompileTarget, _ = filepath.Abs(config.CompileTarget)

	refl := reflect.ValueOf(*config)
	values := make([]interface{}, refl.NumField())
	for i := 0; i < len(values); i++ {
		propValue := refl.Field(i).Interface()
		_, err := os.Stat(fmt.Sprintf("%v", propValue))
		if err != nil {
			fmt.Printf(fmt.Sprintf("%s %s", propValue, err.Error()))
			os.Exit(2)
		}
	}

}

func main() {
	configureCLI()

	binaryExtensions = map[string]string{
		"windows": ".exe",
		"darwin":  "",
		"linux":   ""}

	config := CompileConfig{
		CompileTarget:  targetFile,
		PlatformConfig: platformsConfig,
	}

	ensureConfigsExist(&config)

	var configuredTargets TargetConfig
	unmarshallConfig(config.PlatformConfig, &configuredTargets)

	ch := make(chan bool)
	fmt.Println("Compiling path: ", config.CompileTarget)
	for _, target := range configuredTargets.Targets {
		compile(target, config.CompileTarget, ch)
	}

	waitSync.Wait()
	fmt.Println("Compilations finished!")

}
