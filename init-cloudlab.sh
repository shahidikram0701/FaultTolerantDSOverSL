source VMs.sh

for vm in "${vms[@]}"; do
    echo "Connecting to $vm..."
    echo "Setting up go and other dependecies  $vm"
    ssh -o "StrictHostKeyChecking no" "$ssh_user@$vm" "rm -rf FaultTolerantDSOverSL && git clone https://github.com/shahidikram0701/FaultTolerantDSOverSL.git && cd ~/FaultTolerantDSOverSL && git checkout feature-ds-scalog && chmod 777 init.sh && ./init.sh"

    echo "-----------------------------------------------------------------------------------"
done
