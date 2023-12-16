## Custom Dockerfile
FROM consol/rocky-xfce-vnc
ENV REFRESHED_AT 2022-10-12

# Switch to root user to install additional software
USER 0

## Install a gedit
RUN yum install -y gedit git \
    && yum clean all

RUN git clone https://github.com/radareorg/radare2 && radare2/sys/install.sh

## switch back to default user
USER 1000
