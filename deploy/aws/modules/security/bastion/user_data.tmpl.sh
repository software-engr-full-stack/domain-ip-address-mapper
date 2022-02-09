#!/usr/bin/env bash

sed -i 's/^[ ]*#*[ ]*Port 22[ ]*$/Port ${ssh_ingress_port}/' /etc/ssh/sshd_config
service sshd restart || service ssh restart

sed -i 's/^[ ]*#[ ]*DNS[ ]*=[ ]*$/DNS=208.67.222.222 208.67.220.220/' /etc/systemd/resolved.conf
systemctl restart systemd-resolved.service || systemctl restart systemd-resolved.service
