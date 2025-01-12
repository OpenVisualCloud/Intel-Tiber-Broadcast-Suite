# Building commands

In order to properly build commands for Intel® Tiber™ Broadcast Suite's image, it is needed to gather information about the platform used for execution, and properly put those into the command.

> **Note:** This instruction bases the examples on a configuration of one of the machines, the responses may differ.


## 1. NIC-related settings

This section explains how to gather the platform information regarding network card/s and Virtual Functions, and use them in commands.

### 1.1. Container's IP address
Use the network defined for docker containers to determine a pool of addresses (this may also be taken from output of [`first_run.sh`](../first_run.sh)).
```bash
docker network ls
```

Response example:
```text
NETWORK ID     NAME           DRIVER    SCOPE
d25f79b8d31d   bridge         bridge    local
b334bae7126d   host           host      local
b49202de451f   my_net_801f0   bridge    local
c06348a2888f   none           null      local
```
`my_net_801f0` network is created by [`first_run.sh`](../first_run.sh).

Then determine the Subnet based on the docker's network ID or name:
```bash
docker network inspect my_net_801f0 | grep Subnet
```
or
```bash
docker network inspect b49202de451f | grep Subnet
```
Response example:
```text
                    "Subnet": "192.168.2.0/24",
```

IP address should be taken from the subnet `192.168.2.0/24`.

Chosen IP address must be used in `--ip=x.x.x.x` as part of `[docker_parameters]`, as well as `-p_sip x.x.x.x` part of `[broadcast_suite_parameters]`.

Example: `--ip=192.168.2.2` and matching `-p_sip 192.168.2.2`.

> **Note:** `[docker_parameters]` is used for every flag that stands between `docker run` and an image name (e.g. `tiber_broadcast_suite`). `[broadcast_suite_parameters]` determines every switch that is put after the image name.

### 1.2. Find a proper E810 Series NIC
Run `lshw` command for network interfaces and find all E810 Series cards.
```bash
sudo lshw -C network
```

Response example:
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

Check which card/s have `link=yes` in `configuration` part. Above, only `eth1` fulfills this requirement. A card with `link=no` should not be used.

Note the following information for all of the matching cards:
- bus info
- logical name
- serial
- configuration


### 1.3. Virtual Function's port address
Knowing the PCI address of the proper NIC (example: `eth1`), use the command below to determine addresses of Virtual Functions.

```bash
for vf in /sys/bus/pci/devices/<PHYSICAL_DEVICE_PCI_ADDRESS>/virtfn*
do
  basename $(readlink -f "$vf")
done
```

> **Note:** Physical address of device is indicated as `Slot` in `lspci`'s response, or `bus info` in `lshw`'s response.

Example:
```bash
for vf in /sys/bus/pci/devices/0000:c0:00.1/virtfn*
do
  basename $(readlink -f "$vf")
done
```

```text
0000:c0:11.0
0000:c0:11.1
0000:c0:11.2
0000:c0:11.3
0000:c0:11.4
0000:c0:11.5
```

Based on above response, Virtual Functions derived from `eth1`, that can be used in commands:
- `0000:c0:11.0`
- `0000:c0:11.1`
- `0000:c0:11.2`
- `0000:c0:11.3`
- `0000:c0:11.4`
- `0000:c0:11.5`

e.g. `-p_port 0000:c0:11.0`.

### 1.4. TCP/UDP port
> **Note:** In order to avoid port collisions, try to assign a separate and divergent set of ports for each container running on a server.

Use any unused set of ports out of the Well-known ports pool - `1024..65535`, e.g. `--expose=20000-20010`.

Out of those specified ports, choose a single `-udp_port`, e.g. `-udp_port 20000`.

### 1.5. Final destination IP address
The `-p_rx_ip` parameter is used for the destination IP, e.g. `-p_rx_ip 192.168.2.1`.



## 2. GPU-related settings
This section explains how to gather the platform information regarding graphics card/s and use them in commands.

### 2.1. Gathering information about render device
Following commands can be used to determine the rendering device's location:
```bash
devices_path="/dev/dri/"
readarray -t rendering_devices <<< $(ls -l "$devices_path" | awk '//{if($1 ~ "^c" && $4 == "render"){print $NF}}')
rendering_device=${devices_path}${rendering_devices[0]}
echo "Found: ${rendering_devices[@]}"
echo "Selected: ${rendering_device}"
```
Variable `rendering_device` holds the location value and can be used in further command building as a `-qsv_device ${rendering_device}`, e.g. `-qsv_device /dev/dri/renderD128` (as a part of `[broadcast_suite_parameters]`, `-hwaccel qsv` must be used as well in order to enable hardware acceleration).
