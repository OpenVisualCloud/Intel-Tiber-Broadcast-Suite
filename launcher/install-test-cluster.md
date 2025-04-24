## Steps
1. Operating System: Ubuntu:22.04
2. Hardware Requirements: At least 2 CPUs, 2GB RAM, and 20GB of disk space.
3. Disable Swap: `sudo swapoff -a`
4. Install container runtime either docker or containerd. In this case, docker is used. Refer to this instruction in official manual: https://kubernetes.io/docs/setup/production-environment/container-runtimes/#docker and https://docs.docker.com/engine/install/#server and https://mirantis.github.io/cri-dockerd/usage/install/ (latest version)
> Hint: For CRI you can use: 
> ```bash
> wget https://github.com/Mirantis/cri-dockerd/releases/download/v0.3.17/cri-dockerd_0.3.17.3-0.ubuntu-jammy_amd64.deb && dpkg -i  cri-dockerd_0.3.17.3-0.ubuntu-jammy_amd64.deb
```
For Kubernetes v1.33:
5. Install essential tools
```bash
sudo apt-get update
sudo apt-get install -y apt-transport-https ca-certificates curl gpg
sudo mkdir -p -m 755 /etc/apt/keyrings
curl -fsSL https://pkgs.k8s.io/core:/stable:/v1.33/deb/Release.key | sudo gpg --dearmor -o /etc/apt/keyrings/kubernetes-apt-keyring.gpg
echo 'deb [signed-by=/etc/apt/keyrings/kubernetes-apt-keyring.gpg] https://pkgs.k8s.io/core:/stable:/v1.33/deb/ /' | sudo tee /etc/apt/sources.list.d/kubernetes.list
sudo apt-get update
sudo apt-get install -y kubelet kubeadm kubectl
sudo apt-mark hold kubelet kubeadm kubectl
sudo systemctl enable --now kubelet
```
6. Create cluster
```bash
sudo kubeadm init --pod-network-cidr=192.168.0.0/16 --cri-socket=unix:///var/run/cri-dockerd.sock
sudo systemctl daemon-reload
sudo systemctl restart kubelet
mkdir -p $HOME/.kube
sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
sudo chown $(id -u):$(id -g) $HOME/.kube/config
```
7. Install cni
`kubectl apply -f https://raw.githubusercontent.com/projectcalico/calico/v3.25.1/manifests/calico.yaml`
