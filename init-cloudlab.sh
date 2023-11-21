# vms=("ms0308.utah.cloudlab.us" "ms0313.utah.cloudlab.us" "ms0341.utah.cloudlab.us" "ms0610.utah.cloudlab.us" "ms0315.utah.cloudlab.us" "ms0301.utah.cloudlab.us" "ms0322.utah.cloudlab.us" "ms0607.utah.cloudlab.us" "ms0620.utah.cloudlab.us" "ms0635.utah.cloudlab.us" "ms0330.utah.cloudlab.us" "ms0321.utah.cloudlab.us")

vms=("ms0330.utah.cloudlab.us")

ssh_user="ritwikd2"

for vm in "${vms[@]}"; do
    echo "Connecting to $vm..."
    echo "Setting up go and other dependecies  $vm"
    ssh -o "StrictHostKeyChecking no" "$ssh_user@$vm" "set DEPLOY=${DEPLOY} rm -rf FaultTolerantDSOverSL && git clone https://github.com/shahidikram0701/FaultTolerantDSOverSL.git && cd ~/FaultTolerantDSOverSL && git checkout feature-ds-scalog && chmod 777 init.sh && ./init.sh"

    echo "-----------------------------------------------------------------------------------"
done