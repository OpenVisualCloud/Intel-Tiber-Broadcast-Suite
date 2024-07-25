# Building commands

In order to properly build commands for Intel® Tiber™ Broadcast Suite's image, it is needed to gather information about the platform used for execution, and properly put those into the command.

## NIC-related settings

This chapter explains how to gather the platform information regarding network card/s and Virtual Functions and use them in commands.

### Container's IP address
As per [`first_run.sh`](../../first_run.sh) script execution, IP address should be taken from the subnet `192.168.2.0/24`.

.e.g. `--ip=192.168.2.2` but also matching `-p_sip 192.168.2.2`

Otherwise, use the network defined for docker containers to determine a pool.
```shell
docker network ls
```
Examplary response:
```text
NETWORK ID     NAME           DRIVER    SCOPE
d25f79b8d31d   bridge         bridge    local
b334bae7126d   host           host      local
b49202de451f   my_net_801f0   bridge    local
c06348a2888f   none           null      local
```
`my_net_801f0` network is created by [`first_run.sh`](../../first_run.sh).

Then determine the Subnet based on the docker's network ID or name:
```shell
docker network inspect my_net_801f0 | grep Subnet
```
or
```shell
docker network inspect b49202de451f | grep Subnet
```
Examplary response:
```text
                    "Subnet": "192.168.2.0/24",
```


### Find a proper E810 Series NIC
Run `lshw` command for network interfaces and find all E810 Series cards.
```shell
sudo lshw -C network
```

Examplary response:
```text
(...)
  *-network:0
       description: Ethernet interface
       product: Ethernet Controller E810-C for QSFP
       vendor: Intel Corporation
       physical id: 0
       bus info: pci@0000:c0:00.0
       logical name: eth0
       version: 02
       serial: 6c:fe:54:5a:18:70
       capacity: 25Gbit/s
       width: 64 bits
       clock: 33MHz
       capabilities: pm msi msix pciexpress vpd bus_master cap_list rom ethernet physical 25000bt-fd autonegotiation
       configuration: autonegotiation=off broadcast=yes driver=ice driverversion=Kahawai_1.13.7_20240220 firmware=4.40 0x8001c967 1.3534.0 latency=0 link=no multicast=yes
       resources: iomemory:2eff0-2efef iomemory:2eff0-2efef irq:16 memory:2efffa000000-2efffbffffff memory:2efffe010000-2efffe01ffff memory:f1600000-f16fffff memory:2efffd000000-2efffdffffff memory:2efffe220000-2efffe41ffff
  *-network:1
       description: Ethernet interface
       product: Ethernet Controller E810-C for QSFP
       vendor: Intel Corporation
       physical id: 0.1
       bus info: pci@0000:c0:00.1
       logical name: eth1
       version: 02
       serial: 6c:fe:54:5a:18:71
       capacity: 25Gbit/s
       width: 64 bits
       clock: 33MHz
       capabilities: pm msi msix pciexpress vpd bus_master cap_list rom ethernet physical fibre 25000bt-fd autonegotiation
       configuration: autonegotiation=on broadcast=yes driver=ice driverversion=Kahawai_1.13.7_20240220 duplex=full firmware=4.40 0x8001c967 1.3534.0 latency=0 link=yes multicast=yes
       resources: iomemory:2eff0-2efef iomemory:2eff0-2efef irq:16 memory:2efff8000000-2efff9ffffff memory:2efffe000000-2efffe00ffff memory:f1500000-f15fffff memory:2efffc000000-2efffcffffff memory:2efffe020000-2efffe21ffff
(...)
```

Note the following information for all of the matching cards:
- bus info
- logical name
- serial
- configuration

Based on noted logical names, check which of the cards' interfaces has state UP:
```shell
ip a
```
Examplary response:
```text
2: eth0: <NO-CARRIER,BROADCAST,MULTICAST,UP> mtu 1500 qdisc mq state DOWN group default qlen 1000
    link/ether 6c:fe:54:5a:18:70 brd ff:ff:ff:ff:ff:ff
    altname enp192s0f0
    altname ens33f0
4: eth1: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc mq state UP group default qlen 1000
    link/ether 6c:fe:54:5a:18:71 brd ff:ff:ff:ff:ff:ff
    altname enp192s0f1
    altname ens33f1
    inet6 fe80::6efe:54ff:fe5a:1871/64 scope link 
       valid_lft forever preferred_lft forever
```
Out of those cards, only one has the UP state - `eth1`.

### Find proper NUMA node and CPUs for the NIC

As found earlier, `eth1` has a `bus info` value set as `pci@0000:c0:00.1`. The PCI address is the same, but without the `pci@` prefix.

Use the PCI address in `sudo lspci -Dvmm | grep [pci_address] -A 10`, so `sudo lspci -Dvmm | grep 0000:c0:00.1 -A 10`.

Examplary response:
```text
Slot:   0000:c0:00.1
Class:  Ethernet controller
Vendor: Intel Corporation
Device: Ethernet Controller E810-C for QSFP
SVendor:        Intel Corporation
SDevice:        Ethernet Network Adapter E810-C-Q2
PhySlot:        33
Rev:    02
NUMANode:       1
IOMMUGroup:     202
```
Please note the NUMANode number.

Finally, using `lscpu`, from the `NUMA` section, check which CPU(s) are assigned to the given NUMA node.

```text
(...)
  L3:                    135 MiB (2 instances)
NUMA:                    
  NUMA node(s):          2
  NUMA node0 CPU(s):     0-35,72-107
  NUMA node1 CPU(s):     36-71,108-143
Vulnerabilities:         
  Itlb multihit:         Not affected
(...)
```

Here, for the NUMA node numbered 1, the values are: `36-71,108-143`.

### Virtual Function's port address
Take the output from `virtual_functions.txt` (if saved) or rerun [`first_run.sh`](../../first_run.sh) script to find which interface was bound to whihc Virtual Function.
```shell
sudo -E ./first_run.sh 
```

```text
(...)
0000:c0:00.0 'Ethernet Controller E810-C for QSFP 1592' if=eth0 drv=ice unused=vfio-pci 
Bind 0000:c0:01.0(eth4) to vfio-pci success
Bind 0000:c0:01.1(eth5) to vfio-pci success
Bind 0000:c0:01.2(eth6) to vfio-pci success
Bind 0000:c0:01.3(eth7) to vfio-pci success
Bind 0000:c0:01.4(eth8) to vfio-pci success
Bind 0000:c0:01.5(eth9) to vfio-pci success
Create 6 VFs on PF bdf: 0000:c0:00.0 eth0 succ
0000:c0:00.1 'Ethernet Controller E810-C for QSFP 1592' if=eth1 drv=ice unused=vfio-pci 
Bind 0000:c0:11.0(eth4) to vfio-pci success
Bind 0000:c0:11.1(eth5) to vfio-pci success
Bind 0000:c0:11.2(eth6) to vfio-pci success
Bind 0000:c0:11.3(eth7) to vfio-pci success
Bind 0000:c0:11.4(eth8) to vfio-pci success
Bind 0000:c0:11.5(eth9) to vfio-pci success
Create 6 VFs on PF bdf: 0000:c0:00.1 eth1 succ
```

As we found previously, `eth1` is the proper port, so values for Virtual Functions derived from `eth1` should be used in commands, here:
- `0000:c0:11.0` (eth4)
- `0000:c0:11.1` (eth5)
- `0000:c0:11.2` (eth6)
- `0000:c0:11.3` (eth7)
- `0000:c0:11.4` (eth8)
- `0000:c0:11.5` (eth9)

e.g. `-p_port 0000:c0:11.0`.

### TCP/UDP port
> **Note:** In order to avoid port collisions, try to assign a separate and divergent set of ports for each container running on a server.

Use any unused set of ports out of the Well-known ports pool - `1024..65535`, e.g. `--expose=20000-20010`.

Out of those specified ports, choose a single `-udp_port`, e.g. `-udp_port 20000`.

### Final destination IP address
The `-p_rx_ip` parameter is used for the destination IP, e.g. `-p_rx_ip 192.168.2.1`.



## GPU-related settings
This chapter explains how to gather the platform information regarding graphics card/s and use them in commands.

### Gathering information about render device
Following commands can be used to determine the rendering device's location:
```shell
devices_path="/dev/dri/"
readarray -t rendering_devices <<< $(ls -l "$devices_path" | awk '//{if($1 ~ "^c" && $4 == "render"){print $NF}}')
rendering_device=${devices_path}${rendering_devices[0]}
echo "Found: ${rendering_devices[@]}"
echo "Selected: ${rendering_device}"
```
Variable `rendering_device` holds the location value and can be used in further command building.

### Using gathered system info in commands
This chapter explains how to use information from the previous chapter in the Intel® Tiber™ Broadcast Suite's parameters.

## CPU-related settings
This chapter explains how to gather the platform information regarding CPU and use them in commands.

### CPUset CPUs
For parameter `--cpuset-cpus` use information gathered in [Find proper NUMA node and CPUs](#find-proper-numa-node-and-cpus-for-the-nic) - `36-71,108-143`.

Dedicate a set of CPUs for specific container, e.g. `--cpuset-cpus=36-56`.

### CPUs for MTL
Are choosen automatically and checked with mtl_manger.
