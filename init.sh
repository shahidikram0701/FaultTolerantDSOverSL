# sudo apt install snapd
# sudo snap install go --classic

# echo 'export GOPATH=$HOME/go' >> ~/.bashrc
# echo 'export PATH=$PATH:$GOPATH/bin' >> ~/.bashrc

# go install github.com/mattn/goreman@latest


wget https://go.dev/dl/go1.21.3.linux-arm64.tar.gz
sudo tar -C /usr/local/bin/ -xzf go1.21.3.linux-arm64.tar.gz
echo 'export PATH=$PATH:/usr/local/bin/go/bin' >> ~/.bashrc
source ~/.bashrc
rm -f go1.21.3.linux-arm64.tar.gz
alias go="/usr/local/bin/go/bin/go"

echo "Setup go alias and go in PATH"

go mod init github.com/scalog/scalog && go mod tidy && go mod vendor
echo "Setup vendor files"

sed -i 's/1\.21\.1/1\.33\.0/g' go.mod

echo "Updated gRPC version in go.mod"


go mod tidy && go mod vendor && go build .
echo "Syncing libraries and building pkg"
