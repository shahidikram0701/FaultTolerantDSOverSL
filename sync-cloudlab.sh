vms=("ms0308.utah.cloudlab.us" "ms0313.utah.cloudlab.us" "ms0341.utah.cloudlab.us" "ms0610.utah.cloudlab.us" "ms0315.utah.cloudlab.us" "ms0301.utah.cloudlab.us" "ms0322.utah.cloudlab.us" "ms0607.utah.cloudlab.us" "ms0620.utah.cloudlab.us" "ms0635.utah.cloudlab.us" "ms0330.utah.cloudlab.us" "ms0321.utah.cloudlab.us")

ssh_user="ritwikd2"

for vm in "${vms[@]}"; do
    echo "Connecting to $vm..."
    echo "Setting up go and other dependecies  $vm"
    ssh -o "StrictHostKeyChecking no" "$ssh_user@$vm" "cd ~/FaultTolerantDSOverSL && git checkout feature-ds-scalog && git pull origin feature-ds-scalog"

    echo "-----------------------------------------------------------------------------------"
done