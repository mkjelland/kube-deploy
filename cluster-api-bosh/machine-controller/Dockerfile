FROM ubuntu:14.04

RUN apt-get update && apt-get -y install openssh-client bash
RUN echo "StrictHostKeyChecking no" >> /etc/ssh/ssh_config

RUN echo "UserKnownHostsFile /dev/null" >> /etc/ssh/ssh_config
ADD machine-controller /bin

ENTRYPOINT "/bin/machine-controller"