# Opachat
Video and audio broadcasting WebRTC server. With the ability of recording a chat session, both on the server and on the client. An example of a client for the server You can find here [opaweb](https://github.com/opaldone/opaweb)

## It depends on these great packages
![@pion](https://avatars.githubusercontent.com/u/38192892?s=15&v=4) https://github.com/pion/webrtc \
![@gorilla](https://avatars.githubusercontent.com/u/489566?s=15&v=4) https://github.com/gorilla/csrf \
![@gorilla](https://avatars.githubusercontent.com/u/489566?s=15&v=4) https://github.com/gorilla/websocket \
![@julienschmidt](https://avatars.githubusercontent.com/u/944947?s=15&v=4) https://github.com/julienschmidt/httprouter \
![@letsencrypt](https://avatars.githubusercontent.com/u/9289019?s=15&v=4) https://pkg.go.dev/golang.org/x/crypto/acme/autocert

## How to install and compile
##### Clonning
```bash
git clone https://github.com/opaldone/opachat.git
```
##### Go to the root "opachat" directory
```bash
cd opachat
```
##### Set the GOPATH variable to the current directory "opachat" to avoid cluttering the global GOPATH directory
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
## About config
The config file is located here __opachat/bin/config.json__
```JavaScript
{
  // Just a name of application
  "appname": "opachat",
  // IP address of the server, zeros mean current host
  "address": "0.0.0.0",
  // Port, don't forget to open for firewall
  "port": 7778,
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
  ],
  // These are parameters for server session saving via ffmpeg
  // You can investigate it in
  // src/opachat/serv/recorder.go
  // src/opachat/scr/s_s
  "recorder": {
    "urlVirt": "https://admigo.so/virt",
    "soundLib": "pulse",
    "iHw": "0",
    "scrRes": "2560x1080",
    "logLevel": "error",
    "timeout": 10
  }
}
```
