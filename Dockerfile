# Use imutable image tags rather than mutable tags (like ubuntu:18.04)
FROM ubuntu:focal-20220105

ADD .docker/pip.conf /root/.pip/pip.conf

# RUN apt-get update -y \
#     && apt-get install -y \
#     libssl-dev python3-dev sshpass apt-transport-https jq moreutils vim \
#     ca-certificates curl gnupg2 software-properties-common python3-pip unzip rsync git \
#     && rm -rf /var/lib/apt/lists/*

ENV LANG=C.UTF-8
ENV TZ Asia/Shanghai

RUN apt-get update -y \
    && DEBIAN_FRONTEND=noninteractive apt-get install -y \
    libssl-dev sshpass apt-transport-https moreutils \
    ca-certificates curl gnupg2 python3-pip unzip rsync tzdata \
    && rm -rf /var/lib/apt/lists/*


# RUN curl -fsSL https://mirrors.tuna.tsinghua.edu.cn/docker-ce/linux/ubuntu/gpg | apt-key add - \
#     && add-apt-repository \
#     "deb [arch=amd64] https://mirrors.tuna.tsinghua.edu.cn/docker-ce/linux/ubuntu \
#     $(lsb_release -cs) \
#     stable" \
#     && apt-get update -y && apt-get install --no-install-recommends -y docker-ce \
#     && rm -rf /var/lib/apt/lists/*

ARG arch
RUN curl -o docker-ce-cli.deb https://mirrors.tuna.tsinghua.edu.cn/docker-ce/linux/ubuntu/dists/focal/pool/stable/$arch/docker-ce-cli_20.10.12~3-0~ubuntu-focal_$arch.deb \
    && dpkg -i docker-ce-cli.deb \
    && rm -rf docker-ce-cli.deb

# RUN curl -o nerdctl.tar.gz https://github.com/containerd/nerdctl/releases/download/v0.15.0/nerdctl-0.15.0-linux-amd64.tar.gz \
#     && tar -xvf nerdctl.tar.gz \
#     && mv nerdctl /usr/local/bin/nerdctl \
#     && rm -f nerdctl.tar.gz containerd-rootless*.sh

# Some tools like yamllint need this
# Pip needs this as well at the moment to install ansible
# (and potentially other packages)
# See: https://github.com/pypa/pip/issues/10219

WORKDIR /kuboard-spray
COPY ./requirements.txt ./requirements.txt
RUN /usr/bin/python3 -m pip install --no-cache-dir pip -U \
    && python3 -m pip install --no-cache-dir -r requirements.txt \
    && update-alternatives --install /usr/bin/python python /usr/bin/python3 1

COPY .docker/ansible-patch/config/base.yml /usr/local/lib/python3.8/dist-packages/ansible/config/base.yml
COPY .docker/ansible-patch/plugins_callback/default.py /usr/local/lib/python3.8/dist-packages/ansible/plugins/callback/default.py
COPY .docker/ansible-patch/plugins_action/raw.py /usr/local/lib/python3.8/dist-packages/ansible/plugins/action/raw.py

ENV KUBOARD_SPRAY_WEB_DIR="/kuboard-spray/ui"
ENV KUBOARD_SPRAY_PORT=":80"
ENV GIN_MODE=release
ENV KUBOARD_SPRAY_LOGRUS_LEVEL="info"
ENV KUBOARD_SPRAY_ADMIN_LOGRUS_LEVEL = "info"

COPY ./admin/kuboard-spray-admin kuboard-spray-admin
COPY ./server/ansible-script ansible-script
COPY ./server/ansible-rpc ansible-rpc
COPY ./server/kuboard-spray kuboard-spray
COPY ./web/dist /kuboard-spray/ui
COPY ./server/pull-resource-package.sh pull-resource-package.sh
RUN chmod +x kuboard-spray kuboard-spray-admin

CMD [ "./kuboard-spray" ]
