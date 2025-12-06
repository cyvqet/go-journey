# 使用 Ubuntu 20.04 作为基础镜像
FROM ubuntu:20.04

# 将本地的 webook 可执行文件复制到镜像的 /app/webook 路径
COPY webook /app/webook

# 设置工作目录为 /app（后续命令默认在此目录下执行）
WORKDIR /app

# 指定容器启动时执行的命令：运行 webook 应用程序
CMD ["/app/webook"]