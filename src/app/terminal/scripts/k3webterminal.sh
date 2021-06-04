#!/bin/bash
echo "Environment variables:"
env

# Set terminal type 
echo "export TERM=xterm" >> /root/.bashrc

# Set completion
echo "source /usr/share/bash-completion/bash_completion" >> /root/.bashrc
echo 'source <(kubectl completion bash)' >> /root/.bashrc
echo 'complete -F __start_kubectl k' >> /root/.bashrc

# Config namespace
/opt/k3webterminal/config-namespace.sh
