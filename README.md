# Opachat
Video and audio broadcasting WebRTC server is based on
* ![@pion](https://avatars.githubusercontent.com/u/38192892?s=15&v=4) https://github.com/pion/webrtc
* ![@gorilla](https://avatars.githubusercontent.com/u/489566?s=15&v=4) https://github.com/gorilla/csrf
* ![@gorilla](https://avatars.githubusercontent.com/u/489566?s=15&v=4) https://github.com/gorilla/websocket
* ![@julienschmidt](https://avatars.githubusercontent.com/u/944947?s=15&v=4) https://github.com/julienschmidt/httprouter
* ![@letsencrypt](https://avatars.githubusercontent.com/u/9289019?s=15&v=4) https://pkg.go.dev/golang.org/x/crypto/acme/autocert

## How to install and compile
```bash
# Clonning
git clone https://github.com/opaldone/opachat.git

# Go to root opachat directory
cd opachat

# Set the GOPATH variable to the current directory opachat
# to avoid cluttering the global GOPATH directory
export GOPATH=$(pwd)

# Go to source folder
cd src/opachat

# Installing the required Golang packages
go mod init
go mod tidy

# Return to the opachat's root directory
cd ../..

# There is a "pkg" folder that contains the required Golang packages
# Compiling by the "r" bash script. r - means "run", b - means "build"
./r b

# Creating the required folders structure
# and copying the frontend part by the "u" bash script. u - means "update"
./u

# The "u" script is a watching script then for stopping press Ctrl+C
Ctrl+C

# You can check the "bin" folder. It should contain the necessary structure of folders and files.

# Start the server
./r
```
