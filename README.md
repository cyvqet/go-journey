# go-journey

## 项目简介

**Go Journey** 是一个用于学习和实践 Go 工程化架构的仓库。  
从最初的单体应用，到模块化、服务化、再到微服务架构，每一步都包含了对应的代码示例与思考。

除此之外，还包含我在学习过程中编写的业务项目，用于将架构思想应用到真实业务中。


## 项目结构

```
go-journey/
│
├── cmd/                    # 顶层运行入口
│
├── framework/              # 框架演进主线
│   ├── 01-monolith/        # 单体架构
│   ├── 02-modular/         # 模块化架构
│   ├── 03-service-oriented/# 服务化架构
│   ├── 04-microservices/   # 微服务架构
│   └── 05-infra/           # 基础设施（日志、配置、监控）
│
├── business/               # 业务项目
│   
├── pkg/                    # 通用模块与工具库
│
├── docs/                   # 学习笔记与设计文档
│   ├── design/             # 架构图与设计说明
│   ├── notes/              # 学习笔记与技术总结
│   ├── roadmap.md          # 学习路线图
│   └── changelog.md        # 更新日志
│
├── examples/               # 实验与小型示例
│
├── scripts/                # 构建与部署脚本
│
├── Makefile
├── go.mod
└── README.md
```

## 框架演进路线

| 阶段 | 名称 | 关键主题 |
|------|------|-----------|
| ① | Monolith | 单体架构，所有模块集中于一处 |
| ② | Modular | 模块化拆分，内聚包结构与依赖注入 |
| ③ | Service-Oriented | 服务化：REST/gRPC 通信与解耦 |
| ④ | Microservices | 微服务架构：服务发现、配置中心、网关 |
| ⑤ | Infra | 监控、日志、容错、追踪体系 |

每个阶段都包含：
- 示例代码  
- README 讲解演进思路  
- 与上一阶段的差异说明  

## 学习目标

- 系统化理解 **Go 工程化架构演进路径**
- 通过实践掌握 **服务化与微服务设计**
- 记录并复盘个人在 Go 学习中的思考
- 搭建可复用的 **Go 项目骨架（Project Scaffold）**

## 作者

**cyvqet**

> 一名热爱 Go、系统架构与工程化的开发者。  
>  
> “架构不是炫技，而是演化出的理性结果。”


## License

本项目遵循 [Apache License](./LICENSE)
