# Opachat
Video and audio broadcasting WebRTC server is based on

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
##### You can check the "opachat/bin" folder. It should contain the necessary structure of folders and files
```bash
ls -lash --group-directories-first bin
```
##### Start the server
```bash
./r
```
