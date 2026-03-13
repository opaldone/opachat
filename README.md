<h1 align="center">
  <img src="./lo.svg" alt="Opachat">
  <br />
  Opachat
  <br />
</h1>
<h4 align="center">
  WebRTC server for video and audio broadcasting. A client example can be found here: <a href="https://github.com/opaldone/opaweb">opaweb</a>
</h4>
<p align="center">
<img src="https://img.shields.io/badge/opaldone-opachat-gray.svg?longCache=true&colorB=brightgreen" alt="Opachat" />
<a href="https://sourcegraph.com/github.com/opaldone/opachat?badge">
  <img src="https://sourcegraph.com/github.com/opaldone/opachat/-/badge.svg" alt="Sourcegraph Widget" />
</a>
</p>
<br />

### Built with these excellent Go libraries
<p>
  <a href="https://github.com/pion/webrtc" style="text-decoration:none;">
    <img src="https://avatars.githubusercontent.com/u/38192892?s=30&v=4" alt="pion" style="vertical-align:middle;margin-right:10px;height:30px" />
    <span style="vertical-align:middle">pion/webrtc</span>
  </a>
  <a href="https://github.com/gorilla/csrf" style="text-decoration:none;">
    <img src="https://avatars.githubusercontent.com/u/489566?s=30&v=4" alt="gorilla CSRF" style="vertical-align:middle;margin-right:10px;margin-left:10px;height:30px" />
    <span style="vertical-align:middle">gorilla/csrf</span>
  </a>
  <a href="https://github.com/gorilla/websocket" style="text-decoration:none;">
    <img src="https://avatars.githubusercontent.com/u/489566?s=30&v=4" alt="gorilla websocket" style="vertical-align:middle;margin-right:10px;margin-left:10px;height:30px" />
    <span style="vertical-align:middle">gorilla/websocket</span>
  </a>
  <a href="https://github.com/julienschmidt/httprouter" style="text-decoration:none;">
    <img src="https://avatars.githubusercontent.com/u/944947?s=30&v=4" alt="julienschmidt httprouter" style="vertical-align:middle;margin-right:10px;margin-left:10px;height:30px" />
    <span style="vertical-align:middle">julienschmidt/httprouter</span>
  </a>
  <a href="https://pkg.go.dev/golang.org/x/crypto/acme/autocert" style="text-decoration:none;">
    <img src="https://avatars.githubusercontent.com/u/9289019?s=30&v=4" alt="acme autocert" style="vertical-align:middle;margin-right:10px;margin-left:10px;height:30px" />
    <span style="vertical-align:middle">acme/autocert</span>
  </a>
</p>
<h1></h1>

### How to install and compile
##### Clonning
```bash
git clone https://github.com/opaldone/opachat.git
```
##### Go to the root "opachat" directory
```bash
cd opachat
```
##### Set your GOPATH to the "opachat" directory to keep your global GOPATH clean
```bash
export GOPATH=$(pwd)
```
##### Go to the source folder
```bash
cd src/opachat
```
##### Installing the required Golang packages
```bash
go mod init
```
```bash
go mod tidy
```
##### Return to the "opachat" root directory, You can see the "opachat/pkg" folder that contains the required Golang packages
```bash
cd ../..
```
##### Compiling by the "r" bash script
> r - means "run", b - means "build"
```bash
./r b
```
##### Creating the required folders structure and copying the frontend part by the "u" bash script
> The "u" script is a watching script then for stopping press Ctrl+C \
> u - means "update"
```bash
./u
```
> The "u" script reads sub file "watch_files" \
> E_FOLDERS - the array of creating empty folders \
> C_FOLDERS - the array of folders to simple copy \
> W_FILES - the array of files whose changes are tracked
```bash
./watch_files
```
##### You can check the "opachat/bin" folder. It should contain the necessary structure of folders and files
```bash
ls -lash --group-directories-first bin
```
##### Start the server
```bash
./r
```
### About config
The config file is located here __opachat/bin/config.json__
```JavaScript
{
  // Just a name of application
  "appname": "opachat",

  // IP address of the server, zeros mean current host
  "address": "0.0.0.0",

  // Port, don't forget to open for firewall
  "port": 8080,

  // The folder that stores the frontend part of the site
  "static": "static",

  // Set "acme": true if You need to use acme/autocert
  // false - if You use self-signed certificates
  "acme": false,

  // The array of domain names, set "acme": true
  "acmehost": [
    "opaldone.click",
    "206.189.101.23",
    "www.opaldone.click"
  ],

  // The folder where acme/autocert will store the keys, set "acme": true
  "dirCache": "./certs",

  // The paths to your self-signed HTTPS keys, set "acme": false
  "crt": "/server.crt",
  "key": "/server.key",

  // array of STUN or TURN servers for web rtc connection
  "iceList": [
    {
      // Example turn:192.177.0.555:3478
      "urls": "turn:[some ip]:[some port]",
      // The login for turn server if exists
      "username": "login",
      // The password for turn server
      "credential": "password"
    }
  ]
}
```

### License
MIT License - see [LICENSE](LICENSE) for full text
