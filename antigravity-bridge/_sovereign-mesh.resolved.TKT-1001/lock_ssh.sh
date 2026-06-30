#!/bin/bash
# Remove any old AllowUsers lines
sed -i '/AllowUsers/d' /etc/ssh/sshd_config
echo "AllowUsers aellok@192.168.12.* aellok@127.0.0.1 aellok@::1 root@127.0.0.1 root@::1" >> /etc/ssh/sshd_config
systemctl restart ssh
