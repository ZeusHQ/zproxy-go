package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http/httputil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"

	"github.com/txn2/txeh"
)

func GetMonorepoProjectName(packageJson PackageJson, dirName string) string {
	if packageJson.ZProxy.Name != nil {
		return *packageJson.ZProxy.Name
	}

	if dirName == "web" {
		return "www"
	}

	return dirName
}

func GetProjectName(packageJson PackageJson, dirName string) string {
	if packageJson.ZProxy.Name != nil {
		return *packageJson.ZProxy.Name
	}

	if dirName == "web" {
		return "www"
	}

	return dirName
}

func GetProjectPort(packageJson PackageJson) string {
	if packageJson.ZProxy.Port != nil {
		return strconv.Itoa(*packageJson.ZProxy.Port)
	}

	// `yarn dev` port regex
	r, _ := regexp.Compile("-{1,2}po?r?t?\\s+(\\d+)")

	// Find the port
	match := r.FindStringSubmatch(packageJson.Scripts.Dev)
	port := "3000"

	if len(match) > 0 {
		port = match[len(match)-1]
	}

	return port
}

func GetMonorepoProjectHosts(packageJson PackageJson, name string, monorepoName string) []string {
	var hosts []string

	hosts = append(hosts, fmt.Sprintf("%s.%s.z", name, monorepoName))

	if packageJson.ZProxy.Subdomains != nil {
		for _, subdomain := range *packageJson.ZProxy.Subdomains {
			hosts = append(hosts, fmt.Sprintf("%s.%s.z", subdomain, monorepoName))
		}
	}

	return hosts
}

func GetProjectHosts(packageJson PackageJson, monorepoName string) []string {
	var hosts []string

	hosts = append(hosts, fmt.Sprintf("%s.z", monorepoName))

	if packageJson.ZProxy.Subdomains != nil {
		for _, subdomain := range *packageJson.ZProxy.Subdomains {
			hosts = append(hosts, fmt.Sprintf("%s.%s.z", subdomain, monorepoName))
		}
	}

	return hosts
}

func (handler *proxy) AddMonorepoProjectProxy(monorepoName string, appsDir string, dirName string) []string {
	packageJsonPath := fmt.Sprintf("%s/%s/package.json", appsDir, dirName)
	packageJson := LoadPackageJson(packageJsonPath)

	name := GetMonorepoProjectName(packageJson, dirName)
	port := GetProjectPort(packageJson)
	projectHosts := GetMonorepoProjectHosts(packageJson, name, monorepoName)

	proxiedHost := fmt.Sprintf("http://localhost:%s", port)
	reverseProxy, err := NewProxy(proxiedHost)
	if err != nil {
		panic(err)
	}

	for _, host := range projectHosts {
		handler.proxies[host] = reverseProxy
		fmt.Println(fmt.Sprintf("-> Proxying http://%s to %s", host, proxiedHost))
	}

	return projectHosts
}

func (handler *proxy) AddProjectProxy(monorepoName string, appsDir string) []string {
	packageJsonPath := fmt.Sprintf("%s/package.json", appsDir)
	packageJson := LoadPackageJson(packageJsonPath)

	name := GetProjectName(packageJson, monorepoName)
	port := GetProjectPort(packageJson)
	projectHosts := GetProjectHosts(packageJson, name)

	proxiedHost := fmt.Sprintf("http://localhost:%s", port)
	reverseProxy, err := NewProxy(proxiedHost)
	if err != nil {
		panic(err)
	}

	for _, host := range projectHosts {
		handler.proxies[host] = reverseProxy
		fmt.Println(fmt.Sprintf("-> Proxying http://%s to %s", host, proxiedHost))
	}

	return projectHosts
}

func CreateHandler() *proxy {
	handler := &proxy{
		proxies: map[string]*httputil.ReverseProxy{},
	}

	return handler
}

func (handler *proxy) AddMonorepoHosts(dir string) {
	hosts, err := txeh.NewHostsDefault()
	if err != nil {
		panic(err)
	}

	monorepoName := filepath.Base(dir)

	// Add Zeus Host
	// hosts.AddHost("127.0.0.1", "dev.z")

	appsDir := fmt.Sprintf("%s/apps", dir)

	// Does the appsDir exist?
	if _, err := os.Stat(appsDir); os.IsNotExist(err) {
		// If not, is there a package.json in the current directory?
		packageJsonPath := fmt.Sprintf("%s/package.json", dir)
		if _, err := os.Stat(packageJsonPath); !os.IsNotExist(err) {
			projectHosts := handler.AddProjectProxy(monorepoName, dir)

			for _, host := range projectHosts {
				// Add the custom domain to /etc/hosts
				hosts.AddHost("127.0.0.1", host)
			}
		}
		hosts.Save()
		return
	}

	files, err := ioutil.ReadDir(appsDir)

	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		if f.IsDir() {
			projectHosts := handler.AddMonorepoProjectProxy(monorepoName, appsDir, f.Name())

			for _, host := range projectHosts {
				// Add the custom domain to /etc/hosts
				hosts.AddHost("127.0.0.1", host)
			}
		}
	}

	// Save /etc/hosts
	hosts.Save()
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
