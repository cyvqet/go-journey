# 定义伪目标（不是真实的文件名，总是执行命令而不是检查文件是否存在）
.PHONY: docker docker-clean docker-deploy all

# docker 构建镜像 - 构建Go程序并创建Docker镜像
docker:
	@echo "1. 删除旧的webook文件..."
	@rm -f webook || true  # -f: 强制删除，不提示错误；|| true: 即使文件不存在也不报错
    
	@echo "2. 整理Go模块依赖..."
	@go mod tidy  # 清理和更新go.mod文件中的依赖
    
	@echo "3. 交叉编译Go程序 (Linux ARM)..."
	@GOOS=linux GOARCH=arm go build -tags=k8s -o webook .
	# GOOS=linux: 目标操作系统为Linux
	# GOARCH=arm: 目标CPU架构为ARM
	# -tags=k8s: 使用k8s构建标签（条件编译特性）
	# -o webook: 输出文件名为webook
    
	@echo "4. 删除旧Docker镜像..."
	@docker rmi -f cyvqet/webook:v1.0 2>/dev/null || true
	# rmi -f: 强制删除镜像
	# 2>/dev/null: 将错误输出重定向到/dev/null（静默错误）
    
	@echo "5. 构建新Docker镜像..."
	@docker build -t cyvqet/webook:v1.0 .
	# -t: 给镜像打标签（格式：用户名/镜像名:版本）
	# . : 使用当前目录作为构建上下文（需要包含Dockerfile）

# 清理Kubernetes资源 - 删除已部署的服务和部署
docker-clean:
	@echo "清理Kubernetes资源..."
	@kubectl delete service webook-record 2>/dev/null || true
	# 删除名为webook-record的Service
	@kubectl delete deployment webook-record-service 2>/dev/null || true
	# 删除名为webook-record-service的Deployment
	@sleep 2  # 等待2秒让清理完成
	@echo "清理完成"

# 部署到Kubernetes - 应用配置文件到Kubernetes集群
docker-deploy:
	@echo "部署到Kubernetes..."
	@kubectl apply -f webook-deployment.yaml
	# 应用Deployment配置，创建Pod和副本集
	@kubectl apply -f webook-service.yaml
	# 应用Service配置，创建网络服务
	@echo "等待Pod启动..."
	@sleep 5  # 等待5秒让Pod启动
	@kubectl get pods -l app=webook-record
	# 查看标签为app=webook-record的所有Pod状态

# 一键构建、清理、部署 - 完整的CI/CD流程
all: docker docker-clean docker-deploy
	# 依赖顺序：先执行docker，然后docker-clean，最后docker-deploy
	@echo "✅ 所有步骤完成"
	@echo "访问方式："
	@echo "  1. kubectl port-forward service/webook-record 8080:98"
	# 端口转发：将本地8080端口转发到Service的98端口
	@echo "  2. 然后在浏览器访问 http://localhost:8080"