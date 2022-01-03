# Use imutable image tags rather than mutable tags (like ubuntu:18.04)
FROM ubuntu:bionic-20200807

ADD .docker/sources.list /etc/apt/sources.list
ADD .docker/pip.conf /root/.pip/pip.conf

# RUN apt-get update -y \
#     && apt-get install -y \
#     libssl-dev python3-dev sshpass apt-transport-https jq moreutils vim \
#     ca-certificates curl gnupg2 software-properties-common python3-pip unzip rsync git \
#     && rm -rf /var/lib/apt/lists/*

RUN apt-get update -y \
    && apt-get install -y \
    libssl-dev sshpass apt-transport-https moreutils \
    ca-certificates curl gnupg2 python3-pip unzip rsync \
    && rm -rf /var/lib/apt/lists/*


# RUN curl -fsSL https://mirrors.tuna.tsinghua.edu.cn/docker-ce/linux/ubuntu/gpg | apt-key add - \
#     && add-apt-repository \
#     "deb [arch=amd64] https://mirrors.tuna.tsinghua.edu.cn/docker-ce/linux/ubuntu \
#     $(lsb_release -cs) \
#     stable" \
#     && apt-get update -y && apt-get install --no-install-recommends -y docker-ce \
#     && rm -rf /var/lib/apt/lists/*

RUN curl -o docker-ce-cli.deb https://mirrors.tuna.tsinghua.edu.cn/docker-ce/linux/ubuntu/dists/bionic/pool/stable/amd64/docker-ce-cli_20.10.9~3-0~ubuntu-bionic_amd64.deb \
    && dpkg -i docker-ce-cli.deb \
    && rm -rf docker-ce-cli.deb

# Some tools like yamllint need this
# Pip needs this as well at the moment to install ansible
# (and potentially other packages)
# See: https://github.com/pypa/pip/issues/10219
ENV LANG=C.UTF-8

WORKDIR /kuboard-spray
COPY ./requirements.txt ./requirements.txt
RUN /usr/bin/python3 -m pip install --no-cache-dir pip -U \
    && python3 -m pip install --no-cache-dir -r requirements.txt \
    && update-alternatives --install /usr/bin/python python /usr/bin/python3 1

ENV KUBOARD_SPRAY_WEB_DIR="/kuboard-spray/ui"
ENV KUBOARD_SPRAY_PORT=":80"
ENV GIN_MODE=release

COPY ./server/ansible-script ansible-script
COPY ./server/kuboard-spray kuboard-spray
COPY ./web/dist /kuboard-spray/ui
COPY ./server/pull-resource-package.sh pull-resource-package.sh

CMD [ "./kuboard-spray" ]
