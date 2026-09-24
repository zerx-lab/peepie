package locale

func init() {
	registerCatalog(LanguageChineseSimplified, map[*string]string{
		// Monitoring Screen
		&MonitoringTitle:       "监控配置",
		&MonitoringDescription: "配置监控与可观测性平台，全面洞察系统运行状况",
		&MonitoringName:        "监控",
		&MonitoringOverview: `为生产级部署提供全面的监控与可观测性。

为什么需要监控：
• 追踪性能瓶颈：识别缓慢的 LLM 调用、数据库查询和系统资源占用
• 更快排查问题：详细的链路追踪有助于诊断分布式组件中的问题
• 优化成本：监控 Token 用量模式，优化开销较大的 LLM 交互
• 生产就绪：在关键环境中可靠运行的必要条件

平台选项：
Langfuse：专注于 LLM 的可观测性，提供对话追踪、提示词工程洞察和成本分析
可观测性：全栈监控，提供指标、链路追踪、日志和告警，覆盖基础设施与应用健康状况

快速设置：
• 开发环境：仅启用 Langfuse 以获取 LLM 洞察
• 生产环境：同时启用两个平台以实现全面监控
• 控制成本：使用内置模式，避免外部服务费用`,

		// Langfuse Integration
		&MonitoringLangfuseFormTitle:       "Langfuse 配置",
		&MonitoringLangfuseFormDescription: "配置用于 LLM 监控的 Langfuse 集成",
		&MonitoringLangfuseFormName:        "Langfuse",
		&MonitoringLangfuseFormOverview: `Langfuse 提供：
• 完整的对话追踪
• 模型性能指标
• 成本监控与优化
• 用户行为分析
• AI 交互的调试链路

可选择内置实例或连接外部服务。`,

		&MonitoringLangfuseEmbedded: "内置服务器",
		&MonitoringLangfuseExternal: "外部服务器",
		&MonitoringLangfuseDisabled: "已禁用",

		&MonitoringLangfuseDeploymentType:     "部署类型",
		&MonitoringLangfuseDeploymentTypeDesc: "选择 Langfuse 的部署类型",
		&MonitoringLangfuseBaseURL:            "服务器 URL",
		&MonitoringLangfuseBaseURLDesc:        "Langfuse 服务器地址（例如 https://cloud.langfuse.com）",
		&MonitoringLangfuseProjectID:          "项目 ID",
		&MonitoringLangfuseProjectIDDesc:      "Langfuse 中的项目标识符",
		&MonitoringLangfusePublicKey:          "公钥",
		&MonitoringLangfusePublicKeyDesc:      "用于访问项目的公开 API 密钥",
		&MonitoringLangfuseSecretKey:          "私钥",
		&MonitoringLangfuseSecretKeyDesc:      "用于访问项目的私密 API 密钥",
		&MonitoringLangfuseListenIP:           "监听 IP",
		&MonitoringLangfuseListenIPDesc:       "Docker 端口映射使用的绑定地址（例如 0.0.0.0 表示在所有网络接口上暴露）",
		&MonitoringLangfuseListenPort:         "监听端口",
		&MonitoringLangfuseListenPortDesc:     "Docker 为 Langfuse Web 界面暴露的外部 TCP 端口",

		&MonitoringLangfuseAdminEmail:        "管理员邮箱",
		&MonitoringLangfuseAdminEmailDesc:    "用于登录 Langfuse 管理面板的邮箱",
		&MonitoringLangfuseAdminPassword:     "管理员密码",
		&MonitoringLangfuseAdminPasswordDesc: "用于登录 Langfuse 管理面板的密码",
		&MonitoringLangfuseAdminName:         "管理员用户名",
		&MonitoringLangfuseAdminNameDesc:     "Langfuse 中的管理员用户名",
		&MonitoringLangfuseLicenseKey:        "企业版许可证密钥",
		&MonitoringLangfuseLicenseKeyDesc:    "Langfuse 企业版许可证密钥（可选）",

		&MonitoringLangfuseModeGuide: "选择部署方式：内置（本地掌控）、外部（云端/已有服务）、禁用（不做分析）",
		&MonitoringLangfuseEmbeddedHelp: `内置模式将部署完整的 Langfuse 服务栈：
• PostgreSQL + ClickHouse 数据库
• MinIO S3 存储 + Redis 缓存
• 完整的 LLM 对话追踪
• 成本分析与性能指标
• 私有数据保留在你的服务器上

资源要求：
• 至少约 2GB 内存、5GB 磁盘空间
• 对话日志需要额外存储空间
• 自动完成安装与维护

适用于：注重数据隐私、需要自定义配置或不希望依赖外部服务的团队。所有分析数据均存储在本地，并拥有完整的管理权限。

默认管理员访问：
• Web 界面：http://localhost:4000
• 登录名：admin@pentagi.com
• 密码：password（需修改）`,
		&MonitoringLangfuseExternalHelp: `外部模式将连接到 cloud.langfuse.com 或你已有的 Langfuse 服务器：

• 无需本地基础设施
• 托管式更新与维护
• 跨团队共享分析数据
• 可使用企业版功能
• 数据存储在外部提供商处

设置要求：
• Langfuse 账号和 API 密钥
• 需要互联网连接
• 项目 ID 和认证密钥

适用于：使用云服务、希望采用托管基础设施，或需要与组织内已有 Langfuse 部署集成的团队。`,
		&MonitoringLangfuseDisabledHelp: `Langfuse 已禁用。没有 LLM 可观测性，你将无法获得：

• 对话历史追踪
• Token 用量与成本分析
• 模型性能指标
• AI 交互的调试链路
• 用户行为分析
• 提示词工程洞察

建议在生产环境中启用，
以监控 AI 智能体性能
并有效优化成本。`,

		// Graphiti Integration
		&MonitoringGraphitiFormTitle:       "Graphiti 配置（beta）",
		&MonitoringGraphitiFormDescription: "配置 Graphiti 知识图谱集成",
		&MonitoringGraphitiFormName:        "Graphiti (beta)",
		&MonitoringGraphitiFormOverview: `⚠️  BETA 功能：此功能仍在积极开发中。请关注后续更新以获得改进和稳定性修复。

Graphiti 提供时序知识图谱能力：
• 实体与关系抽取
• AI 智能体的语义记忆
• 时序上下文追踪
• 按流程隔离的上下文搜索

Graphiti 复用在“LLM 提供商”和“嵌入”界面中配置的凭据。LiteLLM 凭据仍保存在 .env 中。在此选择提供商预设；在 ./graphiti/<provider>.yaml 中配置其模型。
支持 openai、gemini、litellm、custom。

可选择内置实例或连接外部服务。`,

		&MonitoringGraphitiEmbedded: "内置服务栈",
		&MonitoringGraphitiExternal: "外部服务",
		&MonitoringGraphitiDisabled: "已禁用",

		&MonitoringGraphitiDeploymentType:           "部署类型",
		&MonitoringGraphitiDeploymentTypeDesc:       "选择 Graphiti 的部署类型",
		&MonitoringGraphitiURL:                      "Graphiti 服务器 URL",
		&MonitoringGraphitiURLDesc:                  "Graphiti API 服务器地址",
		&MonitoringGraphitiTimeout:                  "请求超时",
		&MonitoringGraphitiTimeoutDesc:              "Graphiti 操作的超时时间（秒）",
		&MonitoringGraphitiLLMClientType:            "LLM 提供商预设",
		&MonitoringGraphitiLLMClientTypeDesc:        "提供商预设；凭据来自“LLM 提供商”（openai、gemini、litellm、custom），模型来自 ./graphiti/<provider>.yaml",
		&MonitoringGraphitiSeparateEmbedding:        "使用独立的嵌入端点",
		&MonitoringGraphitiSeparateEmbeddingDesc:    "使用共享的嵌入配置，而不是所选 LLM 的凭据",
		&MonitoringGraphitiSemaphoreLimit:           "协程上限",
		&MonitoringGraphitiSemaphoreLimitDesc:       "Graphiti 辅助协程的最大并发数",
		&MonitoringGraphitiLogLevel:                 "日志级别",
		&MonitoringGraphitiLogLevelDesc:             "Graphiti 日志详细程度",
		&MonitoringGraphitiSearchScope:              "搜索范围",
		&MonitoringGraphitiSearchScopeDesc:          "使用 flowid 实现隔离，或使用 all 进行可信的全局搜索测试",
		&MonitoringGraphitiIngestPolicyRules:        "摄取策略规则",
		&MonitoringGraphitiIngestPolicyRulesDesc:    "将消息模式映射到 REJECT、SKIP_LLM 或 PROCESS 的 JSON",
		&MonitoringGraphitiIngestPolicyField:        "摄取策略匹配字段",
		&MonitoringGraphitiIngestPolicyFieldDesc:    "用于策略匹配的消息字段",
		&MonitoringGraphitiIngestDefaultAction:      "默认摄取动作",
		&MonitoringGraphitiIngestDefaultActionDesc:  "没有摄取策略规则匹配时采用的动作",
		&MonitoringGraphitiIngestWorkerCount:        "摄取工作线程数",
		&MonitoringGraphitiIngestWorkerCountDesc:    "所有流程中重量级摄取任务的最大并发数",
		&MonitoringGraphitiIngestQueueMaxSize:       "摄取队列上限",
		&MonitoringGraphitiIngestQueueMaxSizeDesc:   "最大排队消息数；0 表示队列不限长度",
		&MonitoringGraphitiTaxonomyLayerProfile:     "分类层配置",
		&MonitoringGraphitiTaxonomyLayerProfileDesc: "启用的关系类别，或 full/minimal 别名",
		&MonitoringGraphitiCPUs:                     "Graphiti CPU 限制",
		&MonitoringGraphitiCPUsDesc:                 "Graphiti 容器的 CPU 限制",
		&MonitoringGraphitiMemory:                   "Graphiti 内存限制",
		&MonitoringGraphitiMemoryDesc:               "Graphiti 容器的内存限制",
		&MonitoringGraphitiNeo4jUser:                "Neo4j 用户名",
		&MonitoringGraphitiNeo4jUserDesc:            "访问 Neo4j 数据库的用户名",
		&MonitoringGraphitiNeo4jPassword:            "Neo4j 密码",
		&MonitoringGraphitiNeo4jPasswordDesc:        "访问 Neo4j 数据库的密码",
		&MonitoringGraphitiNeo4jDatabase:            "Neo4j 数据库",
		&MonitoringGraphitiNeo4jDatabaseDesc:        "Neo4j 数据库名称",
		&MonitoringGraphitiNeo4jCPUs:                "Neo4j CPU 限制",
		&MonitoringGraphitiNeo4jCPUsDesc:            "Neo4j 容器的 CPU 限制",
		&MonitoringGraphitiNeo4jMemory:              "Neo4j 内存限制",
		&MonitoringGraphitiNeo4jMemoryDesc:          "Neo4j 容器的内存限制",
		&MonitoringGraphitiNeo4jShmSize:             "Neo4j 共享内存限制",
		&MonitoringGraphitiNeo4jShmSizeDesc:         "/dev/shm 限制；实际使用量计入 Neo4j 内存",

		&MonitoringGraphitiModeGuide: "选择部署方式：内置（本地 Neo4j）、外部（已有 Graphiti）、禁用（无知识图谱）",
		&MonitoringGraphitiEmbeddedHelp: `⚠️  BETA：此功能仍在积极开发中。请关注后续更新以获得改进。

内置模式将部署完整的 Graphiti 服务栈：
• Neo4j 图数据库
• Graphiti API 服务
• 从智能体交互中自动抽取实体
• 时序关系追踪
• 私有知识图谱保留在你的服务器上

前提条件：
• 在“LLM 提供商”中配置提供商凭据（LiteLLM 则在 .env 中配置）
• 启用独立嵌入前，先配置共享嵌入
• 在 ./graphiti/<provider>.yaml 中编辑模型和调用参数

资源要求：
• 默认为 Graphiti 预留最多 2GB、为 Neo4j 预留最多 4GB 内存
• Neo4j 界面：http://localhost:7474
• Graphiti API：http://localhost:8000
• 自动完成安装与维护

保持 flowid 搜索范围以实现隔离。在内存受限的主机上使用有界的摄取队列。信号量/工作线程控制吞吐量；Graphiti 和 Neo4j 的 CPU/内存限制需随并发流程数增加而提高。

重试、合并抽取、锚点、遥测和诊断等精细调优仍可在 .env 中配置。

适用于：希望获得知识图谱能力，同时完全掌控数据与隐私的团队。`,
		&MonitoringGraphitiExternalHelp: `⚠️  BETA：此功能仍在积极开发中。请关注后续更新以获得改进。

外部模式将连接到你已有的 Graphiti 服务器：

• 无需本地基础设施
• 托管式更新与维护
• 跨团队共享知识图谱
• 数据存储在外部提供商处

设置要求：
• Graphiti 服务器 URL 及访问权限
• 需要网络连接
• 提供商、模型、嵌入、抽取和图谱设置均在外部服务器上配置

适用于：使用已有 Graphiti 部署或云服务的团队。`,
		&MonitoringGraphitiDisabledHelp: `Graphiti 已禁用。你将无法获得：

• 时序知识图谱
• 实体与关系抽取
• AI 智能体的语义记忆
• 按流程隔离的图谱召回
• 高级上下文搜索

禁用 Graphiti 时，Peepie 将继续使用其主要的向量记忆。`,

		// Observability Integration
		&MonitoringObservabilityFormTitle:       "可观测性配置",
		&MonitoringObservabilityFormDescription: "配置监控与可观测性服务栈",
		&MonitoringObservabilityFormName:        "可观测性",
		&MonitoringObservabilityFormOverview: `可观测性服务栈包括：
• Grafana 可视化仪表盘
• VictoriaMetrics 时序数据
• Jaeger 分布式链路追踪
• Loki 日志聚合
• OpenTelemetry 数据采集

监控 Peepie 性能与系统健康状况。`,

		&MonitoringObservabilityEmbedded: "内置服务栈",
		&MonitoringObservabilityExternal: "外部采集器",
		&MonitoringObservabilityDisabled: "已禁用",

		&MonitoringObservabilityDeploymentType:     "部署类型",
		&MonitoringObservabilityDeploymentTypeDesc: "选择监控的部署类型",
		&MonitoringObservabilityOTelHost:           "OpenTelemetry 主机",
		&MonitoringObservabilityOTelHostDesc:       "外部 OpenTelemetry 采集器的地址",

		&MonitoringObservabilityGrafanaListenIP:        "Grafana 监听 IP",
		&MonitoringObservabilityGrafanaListenIPDesc:    "Docker 端口映射使用的绑定地址（例如 0.0.0.0 表示在所有网络接口上暴露）",
		&MonitoringObservabilityGrafanaListenPort:      "Grafana 监听端口",
		&MonitoringObservabilityGrafanaListenPortDesc:  "Docker 为 Grafana Web 界面暴露的外部 TCP 端口",
		&MonitoringObservabilityOTelGrpcListenIP:       "OTel gRPC 监听 IP",
		&MonitoringObservabilityOTelGrpcListenIPDesc:   "Docker 端口映射使用的绑定地址（例如 0.0.0.0 表示在所有网络接口上暴露）",
		&MonitoringObservabilityOTelGrpcListenPort:     "OTel gRPC 监听端口",
		&MonitoringObservabilityOTelGrpcListenPortDesc: "Docker 为 OTel gRPC 接收器暴露的外部 TCP 端口",
		&MonitoringObservabilityOTelHttpListenIP:       "OTel HTTP 监听 IP",
		&MonitoringObservabilityOTelHttpListenIPDesc:   "Docker 端口映射使用的绑定地址（例如 0.0.0.0 表示在所有网络接口上暴露）",
		&MonitoringObservabilityOTelHttpListenPort:     "OTel HTTP 监听端口",
		&MonitoringObservabilityOTelHttpListenPortDesc: "Docker 为 OTel HTTP 接收器暴露的外部 TCP 端口",

		&MonitoringObservabilityModeGuide: "选择监控方式：内置（全栈）、外部（已有基础设施）、禁用（不监控）",
		&MonitoringObservabilityEmbeddedHelp: `内置模式将部署完整的监控服务：
• Grafana 仪表盘与告警
• VictoriaMetrics 时序数据库
• Jaeger 分布式链路追踪界面
• Loki 日志聚合系统
• ClickHouse 分析型数据库
• Node Exporter + cAdvisor 指标
• OpenTelemetry 数据采集

组件自动埋点，并预置了
系统健康、性能分析
和调试用的仪表盘。

资源要求：
• 至少约 1.5GB 内存、3GB 磁盘空间
• Grafana 界面：http://localhost:3000
• 性能剖析：http://localhost:7777

适用于：全面的系统可见性、
故障排查和性能调优。`,
		&MonitoringObservabilityExternalHelp: `外部模式将遥测数据发送到你已有的监控基础设施：

• 基于 HTTP/2 的 OTLP 协议（无 TLS）
• 你的采集器必须支持：
  - OTLP HTTP 接收器（端口 4318）
  - OTLP gRPC 接收器（端口 8148）
  - tls: insecure: true 设置
• 发送指标、链路追踪和日志
• 兼容企业级平台：
  Datadog、New Relic、Splunk 等

OTEL_HOST 示例：
your-collector:4318

采集器配置要求：
tls: insecure: true

适用于：已有监控基础设施
或集中式可观测性平台
的组织。`,
		&MonitoringObservabilityDisabledHelp: `可观测性已禁用。你将无法获得：

• 系统性能监控
• 分布式请求链路追踪
• 结构化日志聚合
• 资源使用分析
• 错误追踪与告警
• 性能瓶颈分析

建议在生产环境中启用，
以监控系统健康状况、排查问题
并有效优化性能。`,
	})
}
