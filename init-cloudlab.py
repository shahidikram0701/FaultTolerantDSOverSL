'''
    order-0: ./scalog order --config=.scalog.yaml --oid=0
    order-1: ./scalog order --config=.scalog.yaml --oid=1
    order-2: ./scalog order --config=.scalog.yaml --oid=2

    data-0-0: ./scalog data --config=.scalog.yaml --rid=0 --sid=0
    data-0-1: ./scalog data --config=.scalog.yaml --rid=1 --sid=0
    data-1-0: ./scalog data --config=.scalog.yaml --rid=0 --sid=1
    data-1-1: ./scalog data --config=.scalog.yaml --rid=1 --sid=1

    zk-0: ./scalog zookeeper --config=.scalog.yaml --zid=0
    zk-1: ./scalog zookeeper --config=.scalog.yaml --zid=1
    zk-2: ./scalog zookeeper --config=.scalog.yaml --zid=2

    discovery: ./scalog discovery --config=.scalog.yaml
'''

import yaml, socket, subprocess

CONFIG_FILE = '.scalog.yaml'


def run_cmd(command, component):
    command = command.split()
    print(f'Running: {command}')
    process = subprocess.Popen(command, stdout=subprocess.PIPE, stderr=subprocess.PIPE, stdin=subprocess.PIPE)
    with open(f'{component}_process_id.log', 'w') as f:
        f.write(str(process.pid))
    return


try:
    with open(CONFIG_FILE, 'r') as f:
        scalog_config_data = yaml.load(f.read(), Loader=yaml.FullLoader)
        ip_addresses = {
            'order': scalog_config_data['order-ip-address'],
            'data-0': scalog_config_data['data-ip-address-0'],
            'data-1': scalog_config_data['data-ip-address-1'],
            'discovery': scalog_config_data['disc-ip-address'],
            'zookeeper': scalog_config_data['zk-servers']
        }

        current_ip_address = socket.gethostbyname(socket.gethostname())
        for component in ip_addresses:
            if current_ip_address in ip_addresses[component]:
                print(f'Running Component: {component} on {current_ip_address}')
                
                if component == 'order':
                    oid = ip_addresses[component].index(current_ip_address)
                    run_cmd(f'./scalog order --config=.scalog.yaml --oid={oid}', component)

                elif component == 'data-0':
                    sid = 0
                    rid = ip_addresses[component].index(current_ip_address)
                    run_cmd(f'./scalog data --config=.scalog.yaml --rid={rid} --sid={sid}', component)
                
                elif component == 'data-1':
                    sid = 1
                    rid = ip_addresses[component].index(current_ip_address)
                    run_cmd(f'./scalog data --config=.scalog.yaml --rid={rid} --sid={sid}', component)

                elif component == 'discovery':
                    run_cmd(f'./scalog discovery --config=.scalog.yaml', component)
                
                elif component == 'zookeeper':
                    zid = ip_addresses[component].index(current_ip_address)
                    run_cmd(f'./scalog zookeeper --config=.scalog.yaml --zid={zid}', component)
        
                
except Exception as e:
    print(e)
    exit(1)
