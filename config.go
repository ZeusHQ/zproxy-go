package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type PackageJson struct {
	ZProxy  ZProxyConfig  `json:"zproxy,omitempty"`
	Scripts ScriptsConfig `json:"scripts,omitempty"`
}

type ScriptsConfig struct {
	Dev string `json:"dev,omitempty"`
}

type ZProxyConfig struct {
	Port       *int      `json:"port,omitempty"`
	Name       *string   `json:"name,omitempty"`
	Subdomains *[]string `json:"subdomains,omitempty"`
}

func LoadPackageJson(path string) PackageJson {
	// Open our jsonFile
	jsonFile, err := os.Open(path)
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}

	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	// read our opened xmlFile as a byte array.
	byteValue, _ := ioutil.ReadAll(jsonFile)

	// we initialize our Users array
	var packageJson PackageJson

	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'users' which we defined above
	json.Unmarshal(byteValue, &packageJson)

	return packageJson
}
