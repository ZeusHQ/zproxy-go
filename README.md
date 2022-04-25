# zproxy

zproxy is a zero-configuration Golang ReverseProxy for Turborepos that automagically proxies all http and websocket requests from local domains to their respective localhost:port url.

## Usage

Once installed, go to the root of your turborepo and run "sudo zproxy" to start the proxy server. 

`cd ~/dev/my-mono-repo/`
`sudo zproxy`

## Automatic Configuration

zproxy assumes you have an apps folder in the root of your monorepo. Any projects inside the apps folder of your turborepo will get proxied to their respective localhost:port based on the folder name and the port used inside the dev script in package.json.

For example, let's assume you have two NextJS apps in the apps folder:

my-mono-repo/apps/web - set up to run on port 3000
my-mono-repo/apps/docs - set up to run on port 3001

zproxy will add web.z and docs.z to your /etc/hosts file and proxy requests from web.z to http://localhost:3000 and docs.z to http://localhost:3001.


## Wildcards


## Build from source
`git clone git@github.com:ZeusHQ/zproxy-go.git`
`cd zproxy-go`
`go build`
`go install`

## Running source
`git clone git@github.com:ZeusHQ/zproxy-go.git`
`cd zproxy-go`
`sudo go run . --dir PATH_TO_TURBOREPO`