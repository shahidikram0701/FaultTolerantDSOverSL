# sudo apt install snapd
# sudo snap install go --classic

# echo 'export GOPATH=$HOME/go' >> ~/.bashrc
# echo 'export PATH=$PATH:$GOPATH/bin' >> ~/.bashrc

# go install github.com/mattn/goreman@latest


wget https://go.dev/dl/go1.21.3.linux-arm64.tar.gz
sudo tar -C /usr/local/bin/ -xzf go1.21.3.linux-arm64.tar.gz
echo 'export PATH=$PATH:/usr/local/bin/go/bin' >> ~/.bashrc

source ~/.bashrc

go mod init github.com/scalog/scalog && go mod tidy && go mod vendor

sed -i 's/1\.21\.1/1\.33\.0/g' go.mod

go mod tidy && go mod vendor && go build .
