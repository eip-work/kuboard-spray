> 这个文档记录了一些创建 multipass 实例时用到的命令。仅在开发过程中使用到。

```sh
multipass launch --name spray-01 -c 2 -d 40G -m 4G --network name="Default Switch",mode=manual,mac="52:54:00:4b:ab:cf" focal

multipass shell spray-01

sudo -i

# 配置 sshd
echo "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIIvLHTyRYx85hJMMVrg+3966N5z5ctPYRaJz/hev89DB shaohuanqing@chinamobile.com" >> /home/ubuntu/.ssh/authorized_keys
echo "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIIvLHTyRYx85hJMMVrg+3966N5z5ctPYRaJz/hev89DB shaohuanqing@chinamobile.com" >> /root/.ssh/authorized_keys

sed -i 's/#PermitRootLogin prohibit-password/PermitRootLogin yes/g' /etc/ssh/sshd_config
sed -i 's/#PubkeyAuthentication yes/PubkeyAuthentication yes/g' /etc/ssh/sshd_config
sed -i 's/#PasswordAuthentication yes/PasswordAuthentication yes/g' /etc/ssh/sshd_config
sed -i 's/#UseDNS no/UseDNS no/g' /etc/ssh/sshd_config

sed -i 's/#force_color_prompt=yes/force_color_prompt=yes/g' /root/.bashrc

systemctl restart sshd


# 配置静态地址
# cat << EOF > /etc/netplan/10-custom.yaml
# network:
#   version: 2
#   ethernets:
#     extra0:
#       dhcp4: no
#       match:
#         macaddress: "52:54:00:4b:ab:cd"
#       addresses: [172.23.154.251/20]
# EOF

# netplan apply

```
