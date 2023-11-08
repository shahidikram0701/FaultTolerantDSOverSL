sudo apt install snapd
sudo snap install go --classic

echo 'export GOPATH=$HOME/go' >> ~/.bashrc
echo 'export PATH=$PATH:$GOPATH/bin' >> ~/.bashrc

go install github.com/mattn/goreman@latest

source ~/.bashrc