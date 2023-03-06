package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"text/template"

	"gopkg.in/yaml.v3"
)

func main() {
	// Load Rules
	rul, _ := os.ReadFile("rules.yml")
	err := yaml.Unmarshal(rul, &buildsConfig)
	if err != nil {
		log.Println(err)
	}

	// template
	result, err := template.ParseFiles("templates/main.j2")
	if err != nil {
		panic(err)
	}

	// Update modules
	updateModules(buildsConfig)
	of, err := os.OpenFile("mains/all_main.go", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm)
	if err != nil {
		panic(err)
	}
	result.Execute(of, buildsConfig)
	build(buildsConfig, "mains/all_main.go", "all")

	for key, value := range buildsConfig.CommunityModules {
		result, err := template.ParseFiles("templates/main.j2")
		if err != nil {
			panic(err)
		}
		buildConfigs := Builds{}
		buildConfigs.CommunityModules = make(map[string]Modules)
		buildConfigs.CommunityModules[key] = value

		mainName := fmt.Sprintf("mains/%s_main.go", key)
		of, _ := os.OpenFile(mainName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm)
		result.Execute(of, buildConfigs)
		build(buildConfigs, mainName, key)
	}

	// Compress
	cmd := exec.Command("bash", []string{"upx_compress.sh"}...)
	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("out:", outb.String(), "err:", errb.String())
}

func updateModules(m Builds) {
	for _, value := range m.CommunityModules {
		fmt.Printf("%s@%s", value.Import, value.Version)
		cmd := exec.Command("go", []string{"get", "-u", fmt.Sprintf("%s@%s", value.Import, value.Version)}...)
		var outb, errb bytes.Buffer
		cmd.Stdout = &outb
		cmd.Stderr = &errb
		err := cmd.Run()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("out:", outb.String(), "err:", errb.String())
	}
}

func build(m Builds, mainName string, module string) {
	newMain, err := os.ReadFile(mainName)
	if err != nil {
		panic(err)
	}

	os.WriteFile("main.go", newMain, 0755)
	if module == "all" {
		cmd := exec.Command("bash", []string{"buildmulti.sh", "all"}...)
		var outb, errb bytes.Buffer
		cmd.Stdout = &outb
		cmd.Stderr = &errb
		err := cmd.Run()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("out:", outb.String(), "err:", errb.String())
	} else {
		for key := range m.CommunityModules {
			cmd := exec.Command("bash", []string{"buildmulti.sh", key}...)
			var outb, errb bytes.Buffer
			cmd.Stdout = &outb
			cmd.Stderr = &errb
			err := cmd.Run()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("out:", outb.String(), "err:", errb.String())
		}
	}
}
