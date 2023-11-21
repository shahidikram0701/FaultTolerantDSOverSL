source VMs.sh

for vm in "${vms[@]}"; do

    # Do no start anything on client VM
    if [ "${vm}" != "ms0321.utah.cloudlab.us" ]
    then
        echo "Connecting to $vm..."
        # Start the respective components on the VMs
        ssh -o "StrictHostKeyChecking no" "$ssh_user@$vm" "sh -c 'cd ~/FaultTolerantDSOverSL; ./kill_cmd.sh; rm -f cmd.sh; cp ../cmd.sh .; ./cmd.sh > /dev/null 2>&1 &'"

        echo "-----------------------------------------------------------------------------------"
    fi
done