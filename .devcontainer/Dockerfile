FROM ubuntu:20.04

ARG HTTP_PROXY
ARG HTTPS_PROXY
ARG NO_PROXY

ENV HTTP_PROXY=${HTTP_PROXY}
ENV http_proxy=${HTTP_PROXY}
ENV HTTPS_PROXY=${HTTPS_PROXY}
ENV https_proxy=${HTTPS_PROXY}
ENV NO_PROXY=${NO_PROXY}
ENV no_proxy=${NO_PROXY}

RUN sed -i 's/#force_color_prompt=yes/force_color_prompt=yes/g' /root/.bashrc \
  && sed -i 's/ports\.ubuntu\.com/mirrors\.tuna\.tsinghua\.edu\.cn/g' /etc/apt/sources.list \
  && sed -i 's/archive\.ubuntu\.com/mirrors\.tuna\.tsinghua\.edu\.cn/g' /etc/apt/sources.list \
  && if [ -n "$HTTP_PROXY" ]; then echo "Acquire::http::Proxy \"$HTTP_PROXY\";" > /etc/apt/apt.conf.d/proxy.conf; fi \
  && if [ -n "$HTTPS_PROXY" ]; then echo "Acquire::https::Proxy \"$HTTPS_PROXY\";" >> /etc/apt/apt.conf.d/proxy.conf; fi \
  && apt update \
  && DEBIAN_FRONTEND=noninteractive apt install -y --no-install-recommends tzdata \
  && apt upgrade -y \
  && ln -fs /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
  && apt install curl git vim wget zip unzip rsync netcat lsof -y 

RUN wget https://download.docker.com/linux/ubuntu/dists/focal/pool/stable/$(dpkg --print-architecture)/docker-ce-cli_28.1.1-1~ubuntu.20.04~focal_$(dpkg --print-architecture).deb \
  && wget https://download.docker.com/linux/ubuntu/dists/focal/pool/stable/$(dpkg --print-architecture)/docker-buildx-plugin_0.23.0-1~ubuntu.20.04~focal_$(dpkg --print-architecture).deb \
  && dpkg -i docker-ce-cli_28.1.1-1~ubuntu.20.04~focal_$(dpkg --print-architecture).deb \
  && dpkg -i docker-buildx-plugin_0.23.0-1~ubuntu.20.04~focal_$(dpkg --print-architecture).deb \
  && rm docker-ce-cli_28.1.1-1~ubuntu.20.04~focal_$(dpkg --print-architecture).deb

RUN curl -fsSL https://deb.nodesource.com/setup_22.x -o nodesource_setup.sh \
  && bash nodesource_setup.sh \
  && apt-get install -y nodejs \
  && useradd -D

RUN wget https://go.dev/dl/go1.24.2.linux-$(dpkg --print-architecture).tar.gz \
  && rm -rf /usr/local/go && tar -C /usr/local -xzf go1.24.2.linux-$(dpkg --print-architecture).tar.gz

RUN apt install -y python3-pip \
  && apt install -y sshpass

RUN apt install sudo -y \
  && echo "ubuntu ALL=(ALL:ALL) NOPASSWD: ALL" >> /etc/sudoers

RUN wget -qO- https://get.pnpm.io/install.sh | ENV="/root/.bashrc" SHELL="/usr/bin/bash" bash -

RUN echo "export PATH=$PATH:/usr/local/go/bin:/root/.local/bin" >> /root/.profile \
  && echo "export GO111MODULE=on" >> /root/.profile \
  && echo "export GOPROXY=https://goproxy.cn" >> /root/.profile \
  && echo "source /root/.cmcc-venv/bin/activate" >> /root/.profile \
  && echo "export LC_ALL=C.UTF-8" >> /root/.profile

RUN apt install python3-venv -y

ADD docker/requirements.txt /root/requirements.txt

SHELL ["/bin/bash", "-c"] 

RUN python3 -m venv /root/.cmcc-venv \
  && source /root/.cmcc-venv/bin/activate \
  && pip install -U -r /root/requirements.txt

COPY docker/ansible-patch/config/base.yml /root/.cmcc-venv/lib/python3.8/site-packages/ansible/config/base.yml
COPY docker/ansible-patch/plugins_callback/default.py /root/.cmcc-venv/lib/python3.8/site-packages/ansible/plugins/callback/default.py
COPY docker/ansible-patch/plugins_callback/__init__.py /root/.cmcc-venv/lib/python3.8/site-packages/ansible/plugins/callback/__init__.py
COPY docker/ansible-patch/plugins_action/raw.py /root/.cmcc-venv/lib/python3.8/site-packages/ansible/plugins/action/raw.py

RUN wget https://github.com/mikefarah/yq/releases/download/v4.47.1/yq_linux_amd64 -O /usr/local/bin/yq \
    && chmod +x /usr/local/bin/yq

RUN source /root/.profile && go version && apt install -y libgpgme-dev libassuan-dev libbtrfs-dev pkg-config \
    && if [ -n "$HTTP_PROXY" ]; then git config --global http.proxy "$HTTP_PROXY"; fi \
    && if [ -n "$HTTPS_PROXY" ]; then git config --global https.proxy "$HTTPS_PROXY"; fi \
    && if [ -n "$NO_PROXY" ]; then git config --global http.noProxy "$NO_PROXY"; fi \
    && git clone https://github.com/containers/skopeo.git \
    && cd skopeo \
    && git checkout v1.21.0 \
    && make bin/skopeo \
    && mv bin/skopeo /usr/local/bin \
    && cd ../ \
    && rm -rf skopeo

COPY docker/policy.json /root/.config/containers/policy.json

RUN apt install ssh iproute2 -y && mkdir /run/sshd

ENV ROOT_PASSWORD=123
RUN echo "root:${ROOT_PASSWORD}" | chpasswd

RUN sed -i 's/#PermitRootLogin prohibit-password/PermitRootLogin yes/' /etc/ssh/sshd_config
RUN sed -i 's/#AuthorizedKeysFile/AuthorizedKeysFile/' /etc/ssh/sshd_config
RUN sed -i 's/#PubkeyAuthentication yes/PubkeyAuthentication yes/' /etc/ssh/sshd_config
RUN echo "export PANGEE_CLUSTER_LOGRUS_LEVEL=trace" >> /root/.profile
RUN echo "export PANGEE_CLUSTER_ADMIN_LOGRUS_LEVEL=trace" >> /root/.profile
RUN echo "export PANGEE_CLUSTER_DEFAULT_ADMIN_PASSWORD=2432612431302464617343536e4858706f567278764f532f5274752f755a526c72694f30306d73456657674c7a39644b712f43436c6d64544e506375" >> /root/.profile

ENTRYPOINT [ "/usr/sbin/sshd", "-D" ]
