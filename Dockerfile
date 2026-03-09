# 第一阶段：构建阶段
FROM --platform=$BUILDPLATFORM golang:1.24.2 AS builder

WORKDIR /app

# 复制 Go Workspace 文件
COPY go.work go.work.sum ./

# 复制所有模块的源代码
COPY server/ server/
COPY admin/ admin/
COPY go-yaml/ go-yaml/

# 下载依赖
RUN go mod download


# 设置构建参数
ARG TARGETOS
ARG TARGETARCH

# 编译二进制文件
RUN GOPROXY=https://goproxy.cn CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o pangee-cluster server/pangee-cluster.go
RUN GOPROXY=https://goproxy.cn CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o pangee-cluster-admin admin/pangee-cluster-admin.go

# 第二阶段：运行阶段
FROM --platform=$BUILDPLATFORM ubuntu:22.04

# 设置构建参数
ARG TARGETOS
ARG TARGETARCH

ADD .devcontainer/docker/pip.conf /root/.pip/pip.conf

ENV LANG=C.UTF-8
ENV TZ=Asia/Shanghai

# 使用清华大学镜像源
RUN sed -i 's/archive.ubuntu.com/mirrors.tuna.tsinghua.edu.cn/g' /etc/apt/sources.list \
    && sed -i 's/security.ubuntu.com/mirrors.tuna.tsinghua.edu.cn/g' /etc/apt/sources.list \
    && apt-get update -y \
    && DEBIAN_FRONTEND=noninteractive apt-get install -y \
    libssl-dev sshpass apt-transport-https moreutils git wget skopeo \
    ca-certificates curl gnupg2 python3-pip unzip rsync tzdata \
    && rm -rf /var/lib/apt/lists/*

RUN wget https://github.com/mikefarah/yq/releases/download/v4.47.1/yq_linux_$TARGETARCH -O /usr/local/bin/yq \
    && chmod +x /usr/local/bin/yq

RUN update-alternatives --install /usr/bin/python python /usr/bin/python3 1
RUN mkdir /pangee-cluster 

WORKDIR /pangee-cluster
COPY ./.devcontainer/docker/requirements.txt ./requirements.txt
RUN python3 --version
RUN /usr/bin/python3 -m pip install --no-cache-dir pip -U \
    && python3 -m pip install --no-cache-dir --root-user-action=ignore -r requirements.txt \
    && pip cache purge \
    && rm -rf $(pip cache dir)

COPY .devcontainer/docker/ansible-patch/config/base.yml /usr/local/lib/python3.10/dist-packages/ansible/config/base.yml
COPY .devcontainer/docker/ansible-patch/plugins_callback/default.py /usr/local/lib/python3.10/dist-packages/ansible/plugins/callback/default.py
COPY .devcontainer/docker/ansible-patch/plugins_callback/__init__.py /usr/local/lib/python3.10/dist-packages/ansible/plugins/callback/__init__.py
COPY .devcontainer/docker/ansible-patch/plugins_action/raw.py /usr/local/lib/python3.10/dist-packages/ansible/plugins/action/raw.py

COPY .devcontainer/docker/policy.json /root/.config/containers/policy.json

ENV PANGEE_CLUSTER_WEB_DIR="/pangee-cluster/ui"
ENV PANGEE_CLUSTER_PORT="9080"
ENV GIN_MODE=release
ENV PANGEE_CLUSTER_LOGRUS_LEVEL="info"
ENV PANGEE_CLUSTER_ADMIN_LOGRUS_LEVEL="info"
ENV PANGEE_CLUSTER_DEFAULT_ADMIN_PASSWORD="2432612431302464617343536e4858706f567278764f532f5274752f755a526c72694f30306d73456657674c7a39644b712f43436c6d64544e506375"
ENV TERM=xterm

# 从构建阶段复制编译好的二进制文件
COPY --from=builder /app/pangee-cluster pangee-cluster
COPY --from=builder /app/pangee-cluster-admin pangee-cluster-admin
RUN chmod +x pangee-cluster
RUN chmod +x pangee-cluster-admin

COPY ./server/ansible-script ansible-script
COPY ./server/ansible-rpc ansible-rpc
COPY ./web/dist /pangee-cluster/ui
COPY --chmod=+x ./server/pull-resource-package.sh pull-resource-package.sh

ENTRYPOINT ["./pangee-cluster"]
