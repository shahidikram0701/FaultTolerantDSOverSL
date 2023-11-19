# sudo apt install snapd
# sudo snap install go --classic

# echo 'export GOPATH=$HOME/go' >> ~/.bashrc
# echo 'export PATH=$PATH:$GOPATH/bin' >> ~/.bashrc

# go install github.com/mattn/goreman@latest

which go
if [ "$?" != "0" ]
then
    wget https://go.dev/dl/go1.21.3.linux-arm64.tar.gz
    sudo tar -C /usr/local/bin/ -xzf go1.21.3.linux-arm64.tar.gz
    echo 'export PATH=$PATH:/usr/local/bin/go/bin' >> ~/.bashrc
    source ~/.bashrc
    rm -f go1.21.3.linux-arm64.tar.gz

    echo "Setup go alias and go in PATH"
fi

/usr/local/bin/go/bin/go mod init github.com/scalog/scalog && /usr/local/bin/go/bin/go mod tidy && /usr/local/bin/go/bin/go mod vendor || exit 1
echo "Setup vendor files"

sed -i 's/1\.21\.1/1\.33\.0/g' go.mod

echo "Updated gRPC version in go.mod"


/usr/local/bin/go/bin/go mod tidy && /usr/local/bin/go/bin/go mod vendor && /usr/local/bin/go/bin/go build .
echo "Syncing libraries and building pkg"

if [ "${DEPLOY}" == "1" ]
then
    if [ -f *_process_id.log ]
    then
        PID=$(cat *_process_id.log)
        kill -9 $PID || true
        rm -rf *_process_id.log
    fi

    PID=$(ps aux | grep "./scalog.*--config" | head -1 | awk '{print $2}')
    kill -9 $PID || true
    
    echo "Starting the component based on the .scalog.yaml config file and the machine IP"
    python3 init-cloudlab.py
fi
