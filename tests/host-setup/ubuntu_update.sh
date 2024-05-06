#! /usr/bin/bash -e

echo "/************************************************************************************/"
echo "ubuntu install steps from public intel wiki page."
echo "for more details see:"
echo "    https://dgpu-docs.intel.com/driver/installation.html#ubuntu-install-steps"
echo "/************************************************************************************/"



echo "/************ Step: 1.3.3.1. Ubuntu Package Repository ************/"

export DEBIAN_FRONTEND=noninteractive
sudo apt update
sudo apt install -y gpg-agent wget

source /etc/os-release
if [[ ! " jammy " =~ " ${VERSION_CODENAME} " ]]; then
  echo "/************ Only ubuntu version jammy (22.04) is supported ************/"
  exit 1
else
  wget -qO - https://repositories.intel.com/gpu/intel-graphics.key | \
    sudo gpg --dearmor --output /usr/share/keyrings/intel-graphics.gpg
  echo "deb [arch=amd64 signed-by=/usr/share/keyrings/intel-graphics.gpg] https://repositories.intel.com/gpu/ubuntu ${VERSION_CODENAME}/lts/2350 unified" | \
    sudo tee /etc/apt/sources.list.d/intel-gpu-${VERSION_CODENAME}.list
  sudo apt update
fi



echo "/************ Step: 1.3.3.2. Ubuntu Package Installation ************/"

sudo apt install -y \
  libc6 \
  libstdc++6 \
  libigc1 \
  linux-headers-$(uname -r) \
  linux-modules-extra-$(uname -r) \
  flex bison \
  intel-fw-gpu intel-i915-dkms xpu-smi \
  intel-opencl-icd intel-level-zero-gpu level-zero \
  intel-media-va-driver-non-free libmfx1 libmfxgen1 libvpl2 \
  libegl-mesa0 libegl1-mesa libegl1-mesa-dev libgbm1 libgl1-mesa-dev libgl1-mesa-dri \
  libglapi-mesa libgles2-mesa-dev libglx-mesa0 libigdgmm12 libxatracker2 mesa-va-drivers \
  mesa-vdpau-drivers mesa-vulkan-drivers va-driver-all vainfo hwinfo clinfo \
  libigc-dev intel-igc-cm libigdfcl-dev libigfxcmrt-dev level-zero-dev

echo "/************ Step: 1.3.4. Configuring Render Group Membership ************/"

stat -c "%G" /dev/dri/render*
groups ${USER}
sudo gpasswd -a ${USER} render
newgrp render


echo "/************ All done please reboot machine to apply changes ************/"
