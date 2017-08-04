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
	targetFile      string
	platformsConfig string
	wg              sync.WaitGroup
)

func compile(target Target, compileTarget string, ch chan bool) {

	for _, arch := range target.Arch {
		wg.Add(1)
		go func(arch string, target string) {

			env := os.Environ()
			env = append(env, fmt.Sprintf("GOOS=%s", target))
			env = append(env, fmt.Sprintf("GOARCH=%s", arch))

			fmt.Printf("Compiling platform %s with arch: %s\n", target, arch)
			cmd := exec.Command("go", "build", compileTarget)
			cmd.Env = env

			err := cmd.Run()
			if err != nil {
				fmt.Println(err.Error())
			} else {
				fmt.Println("Compilation success")
			}

			defer wg.Done()
		}(arch, target.Id)
	}

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

	wg.Wait()
	fmt.Println("Compilations finished!")

}
