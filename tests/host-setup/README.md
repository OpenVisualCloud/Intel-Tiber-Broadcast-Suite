# Test machine setup 
Before running any tests, machines needs to be set up first.

## Prerequsities
* ubuntu 22.04 Jammy
* Flex 140 
* Flex 170
* Access to target machine

Ensure your machine has proper apt proxy setup, before execution. In `/etc/apt/apt.conf.d/proxy.conf` add 2 following lines:
```text
Acquire::http::Proxy "http://proxy-dmz.intel.com:911";
Acquire::https::Proxy "http://proxy-dmz.intel.com:912";
```

## Usage
0. Log into target machine and clone this repo.
1. Run `ubuntu_update.sh`.
2. Reboot machine to apply changes.
3. Download [ice driver](https://www.intel.com/content/www/us/en/download/19630/intel-network-adapter-driver-for-e810-series-devices-under-linux.html), e.g. `ice-1.13.7.tar.gz`, and put it in `tests/host-setup` directory.
4. Run `nic_driver_setup.sh` with tar name parameter  
` $ ./nic_driver_setup.sh <tar-file-name | default: ice-1.13.7.tar.gz>`.

Optional:

5. Download [nic driver update](https://www.intel.com/content/www/us/en/download/15084/816410/intel-ethernet-adapter-complete-driver-pack.html) package in version 29.0 and put it in `tests/host-setup` directory.
6. Update nic_driver firmware by executing
`nic_firmware_update.sh <tar-file-name | default: Release_29.0.zip> <eth_name | default: eth0>`.


## Troubleshooting
If in step #1 (1.3.4.) a `stat: cannot statx '/dev/dri/render': No such file or directory` error, or similar, is thrown, reboot your machine and retry.

If in step #4 (1.4) a `rmmod: ERROR: Module ice is in use by: irdma` error is thrown, execute `sudo modprobe -r idrma`, remove the `mtl` and `ice-x.xx.x` folders, and restart the script.
