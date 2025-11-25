echo "Welcome to agent-go installer"
mkdir ~/.tmp
cd ~/.tmp
git clone https://github.com/finettt/agent-go agent-go
cd ./agent-go
make
chown $USER:$USER ./agent-go
sudo mv ./agent-go /usr/local/bin/
rm -rf ~/.tmp/*
echo "Agent-go installed!"
