echo "Welcome to agent-go installation script"
cd ~/.tmp
git clone https://github.com/finettt/agent-go ./agent-go
cd ./agent-go
make
chown $USER:$USER ./agent-go
mv ./agent-go /usr/local/bin/
rm -rf ~/.tmp/agent-go
echo "Agent-go installed!"