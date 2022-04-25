package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http/httputil"
	"os/exec"
	"regexp"

	"github.com/txn2/txeh"
)

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

	// Nextjs dev port regex
	r, _ := regexp.Compile("next dev --port ([\\d]+)")

	for _, f := range files {
		packagePath := fmt.Sprintf("%s/%s/package.json", appsDir, f.Name())

		// Open our package.json
		packageJson, err := ioutil.ReadFile(packagePath)
		// if we os.Open returns an error then handle it
		if err != nil {
			fmt.Println(err)
		}

		// Find the port
		match := r.FindStringSubmatch(string(packageJson))
		port := "3000"

		if len(match) > 0 {
			port = match[len(match)-1]
		}

		host := fmt.Sprintf("%s.z", f.Name())
		proxiedHost := fmt.Sprintf("http://localhost:%s", port)
		reverseProxy, err := NewProxy(proxiedHost)
		if err != nil {
			panic(err)
		}
		handler.proxies[host] = reverseProxy

		// Add the custom domain to /etc/hosts
		hosts.AddHost("127.0.0.1", host)

		fmt.Println(fmt.Sprintf("-> Proxying http://%s", host))
	}

	// Save /etc/hosts
	hosts.Save()

	return handler
}

func OpenHosts(handler *proxy) {
	for host := range handler.proxies {
		exec.Command("open", fmt.Sprintf("http://%s", host)).Start()
	}
}
