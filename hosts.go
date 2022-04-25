package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http/httputil"
	"os/exec"
	"regexp"
	"strconv"

	"github.com/txn2/txeh"
)

func GetProjectName(packageJson PackageJson, dirName string) string {
	if packageJson.ZProxy.Name != nil {
		return *packageJson.ZProxy.Name
	}
	return dirName
}

func GetProjectPort(packageJson PackageJson) string {
	if packageJson.ZProxy.Port != nil {
		return strconv.Itoa(*packageJson.ZProxy.Port)
	}

	// `yarn dev` port regex
	r, _ := regexp.Compile("--port ([\\d]+)")

	// // Open our package.json
	// packageJson, err := ioutil.ReadFile(packagePath)

	// // if we os.Open returns an error then handle it
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// Find the port
	match := r.FindStringSubmatch(packageJson.Scripts.Dev)
	port := "3000"

	if len(match) > 0 {
		port = match[len(match)-1]
	}

	return port
}

func (handler *proxy) AddProjectProxy(appsDir string, dirName string) string {
	packageJsonPath := fmt.Sprintf("%s/%s/package.json", appsDir, dirName)
	packageJson := LoadPackageJson(packageJsonPath)

	name := GetProjectName(packageJson, dirName)
	port := GetProjectPort(packageJson)
	host := fmt.Sprintf("%s.z", name)

	proxiedHost := fmt.Sprintf("http://localhost:%s", port)
	reverseProxy, err := NewProxy(proxiedHost)
	if err != nil {
		panic(err)
	}
	handler.proxies[host] = reverseProxy

	fmt.Println(fmt.Sprintf("-> Proxying http://%s to %s", host, proxiedHost))
	return host
}

func AddHosts(dir string) *proxy {
	hosts, err := txeh.NewHostsDefault()
	if err != nil {
		panic(err)
	}

	// Add Zeus Host
	hosts.AddHost("127.0.0.1", "dev.z")

	appsDir := fmt.Sprintf("%s/apps", dir)

	files, err := ioutil.ReadDir(appsDir)

	if err != nil {
		log.Fatal(err)
	}

	handler := &proxy{
		proxies: map[string]*httputil.ReverseProxy{},
	}

	for _, f := range files {
		if f.IsDir() {
			host := handler.AddProjectProxy(appsDir, f.Name())

			// Add the custom domain to /etc/hosts
			hosts.AddHost("127.0.0.1", host)
		}
	}

	// Save /etc/hosts
	hosts.Save()

	return handler
}

func RemoveHosts(handler *proxy) {
	hosts, err := txeh.NewHostsDefault()
	if err != nil {
		panic(err)
	}

	for host := range handler.proxies {
		hosts.RemoveHost(host)
		// exec.Command("open", fmt.Sprintf("http://%s", host)).Start()
	}

	hosts.Save()
}

func OpenHosts(handler *proxy) {
	for host := range handler.proxies {
		exec.Command("open", fmt.Sprintf("http://%s", host)).Start()
	}
}
