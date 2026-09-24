package locale

func init() {
	registerCatalog(LanguageChineseSimplified, map[*string]string{
		// Server Settings screen strings
		&ServerSettingsFormTitle:       "服务器设置",
		&ServerSettingsFormDescription: "配置 Peepie 服务器的网络访问与公共路由",
		&ServerSettingsFormName:        "服务器设置",
		&ServerSettingsFormOverview: `• 网络绑定 - 控制 Peepie 监听的网络接口和端口
• 公共 URL - 外部访问地址及重定向中使用的可选基础路径
• CORS - 允许浏览器访问的来源
• 代理 - 访问 LLM/搜索提供商等出站流量使用的 HTTP/HTTPS 代理
• SSL 目录 - 包含 server.crt 和 server.key（PEM）的自定义证书目录
• 数据目录 - 智能体产物和流程工作区的持久化存储`,

		&ServerSettingsLicenseKey:     "许可证密钥",
		&ServerSettingsLicenseKeyDesc: "Peepie 许可证密钥，格式为 XXXX-XXXX-XXXX-XXXX",

		&ToolsDockerInsideHost:     "工作节点 Docker 守护进程地址",
		&ToolsDockerInsideHostDesc: "提供给工作节点容器的守护进程端点（例如 tcp://dind:2376）；留空则保持挂载套接字",

		&ToolsDockerInsideTLSVerify:     "工作节点 Docker TLS 验证",
		&ToolsDockerInsideTLSVerifyDesc: "为工作节点容器的 Docker 连接启用 TLS 验证",

		&ToolsDockerInsideCertPath:     "工作节点 Docker 证书路径",
		&ToolsDockerInsideCertPathDesc: "位于工作节点上的 TLS 证书目录，以只读方式挂载到工作节点容器中",

		&ToolsDockerInsideHostHelp: `启用 Docker 访问时提供给工作节点容器的 Docker 守护进程端点。

它会以 DOCKER_HOST 的形式注入每个沙箱，使容器内的 Docker CLI 与该守护进程通信。

设置此项后还会停止自动检测主机套接字并将其绑定挂载到沙箱中：智能体只能访问你指定的守护进程，而不是运行 Peepie 本身的守护进程。显式配置的 Docker 套接字仍然优先，并照常挂载。

留空则保持原有行为（自动检测并挂载套接字）。

示例：
• tcp://dind:2376
• tcp://10.0.0.5:2375
• unix:///var/run/docker.sock`,

		&ToolsDockerInsideTLSVerifyHelp: `为工作节点容器与 Docker 守护进程之间的连接启用 TLS 验证。

以 DOCKER_TLS_VERIFY 的形式注入每个沙箱。当守护进程要求双向 TLS 时，请与“工作节点 Docker 守护进程地址”和“工作节点 Docker 证书路径”配合使用。

对于普通（非 TLS）端点，请保持禁用。`,

		&ToolsDockerInsideCertPathHelp: `存放工作节点容器所用 TLS 客户端材料（ca.pem、cert.pem、key.pem）的目录。

以 DOCKER_CERT_PATH 的形式注入每个沙箱，并以只读方式绑定挂载到相同路径，因此该值在容器内解析结果一致。

重要：此路径在工作节点上解析——即由其 Docker 守护进程创建沙箱容器的机器——而不是运行本安装程序的机器。因此这里不会对其进行校验；请确保该路径在工作节点上存在。

示例：
• /etc/docker/certs
• /opt/pentagi/docker/ssl`,

		&ServerSettingsTenantID:     "租户 ID",
		&ServerSettingsTenantIDDesc: "可选的命名空间，用于在同一主机上部署多个实例（例如 acme、team_alpha）",

		&ServerSettingsPprofAddr:     "pprof 监听地址",
		&ServerSettingsPprofAddrDesc: "可选的 Go pprof 绑定地址；留空则禁用（例如 :7777、127.0.0.1:7778）",

		&ServerSettingsHost:     "服务器主机（监听 IP）",
		&ServerSettingsHostDesc: "Docker 端口映射使用的绑定地址（例如 0.0.0.0 表示在所有接口上暴露）",

		&ServerSettingsPort:     "服务器端口（监听端口）",
		&ServerSettingsPortDesc: "Docker 为 Peepie Web UI 暴露的外部 TCP 端口",

		&ServerSettingsPublicURL:     "公共 URL",
		&ServerSettingsPublicURLDesc: "用于重定向和链接的公共基础 URL（支持基础路径，例如 https://example.com/pentagi/）",

		&ServerSettingsCORSOrigins:     "CORS 来源",
		&ServerSettingsCORSOriginsDesc: "以逗号分隔的允许来源列表（例如 https://localhost:8443,https://localhost）",

		&ServerSettingsProxyURL:     "HTTP/HTTPS 代理",
		&ServerSettingsProxyURLDesc: "访问 LLM 和外部工具的出站请求所用代理（不用于 Docker API 访问）",

		&ServerSettingsProxyUsername:     "代理用户名",
		&ServerSettingsProxyUsernameDesc: "代理认证用户名（可选）",
		&ServerSettingsProxyPassword:     "代理密码",
		&ServerSettingsProxyPasswordDesc: "代理认证密码（可选）",

		&ServerSettingsHTTPClientTimeout:       "HTTP 客户端超时",
		&ServerSettingsHTTPClientTimeoutDesc:   "外部 API 调用（LLM 提供商、搜索引擎等）的超时时间（秒）",
		&ServerSettingsTerminalToolTimeout:     "终端工具超时",
		&ServerSettingsTerminalToolTimeoutDesc: "终端命令的默认超时时间（秒）（0 或负数 = 使用 3 小时上限）",

		&ServerSettingsExternalSSLCAPath:     "自定义 CA 证书路径",
		&ServerSettingsExternalSSLCAPathDesc: "容器内自定义根 CA 证书的路径（例如 /opt/pentagi/ssl/ca-bundle.pem）",

		&ServerSettingsExternalSSLInsecure:     "跳过 SSL 验证",
		&ServerSettingsExternalSSLInsecureDesc: "禁用 SSL/TLS 证书校验（仅用于使用自签名证书的测试）",

		&ServerSettingsSSLDir:     "SSL 目录",
		&ServerSettingsSSLDirDesc: "包含 PEM 格式 server.crt 和 server.key 的目录（server.crt 可包含完整证书链）",

		&ServerSettingsDataDir:     "数据目录",
		&ServerSettingsDataDirDesc: "存放所有智能体生成文件的目录；其中的 flow-N 子目录在工作节点容器中用作 /work",

		&ServerSettingsCookieSigningSalt:     "Cookie 签名盐值",
		&ServerSettingsCookieSigningSaltDesc: "用于签名 Cookie 的密钥（请妥善保密）",

		&ServerSettingsDatabaseExtensionsSchema:     "数据库扩展 Schema",
		&ServerSettingsDatabaseExtensionsSchemaDesc: "设置租户 ID 时存放共享扩展的 Schema（例如 public，Supabase 为 extensions）",

		&ServerSettingsDatabaseSearchPathViaOptions:     "通过 options 传递 search_path",
		&ServerSettingsDatabaseSearchPathViaOptionsDesc: "在 options 启动参数中发送租户 search_path（某些连接池需要）",

		// Hints for fields overview
		&ServerSettingsLicenseKeyHint:                   "许可证密钥",
		&ServerSettingsTenantIDHint:                     "租户 ID",
		&ServerSettingsPprofAddrHint:                    "pprof 地址",
		&ServerSettingsHostHint:                         "监听 IP",
		&ServerSettingsPortHint:                         "监听端口",
		&ServerSettingsPublicURLHint:                    "公共 URL",
		&ServerSettingsCORSOriginsHint:                  "CORS 来源",
		&ServerSettingsProxyURLHint:                     "代理 URL",
		&ServerSettingsProxyUsernameHint:                "代理用户名",
		&ServerSettingsProxyPasswordHint:                "代理密码",
		&ServerSettingsHTTPClientTimeoutHint:            "HTTP 超时",
		&ServerSettingsTerminalToolTimeoutHint:          "终端超时",
		&ServerSettingsExternalSSLCAPathHint:            "自定义 CA 路径",
		&ServerSettingsExternalSSLInsecureHint:          "跳过 SSL 验证",
		&ServerSettingsSSLDirHint:                       "SSL 目录",
		&ServerSettingsDataDirHint:                      "数据目录",
		&ServerSettingsDatabaseExtensionsSchemaHint:     "扩展 Schema",
		&ServerSettingsDatabaseSearchPathViaOptionsHint: "通过 options 传递 search_path",

		// Help texts per-field
		&ServerSettingsGeneralHelp: `Peepie 通过 Docker 暴露其 Web UI，主机和端口均可配置。

公共 URL 必须与用户访问服务器的方式一致。如果使用子路径（例如 /pentagi/），请在此处包含它。CORS 控制来自指定来源的浏览器访问。代理影响访问 LLM/搜索提供商以及工具所用其他外部服务的出站流量。

SSL 目录可用于提供自定义证书。设置后，服务器将使用该目录中的 server.crt 和 server.key。数据目录用于存放流程的产物和工作文件。`,

		&ServerSettingsLicenseKeyHelp: `Peepie 许可证密钥，格式为 XXXX-XXXX-XXXX-XXXX。用于与 PentAGI Cloud API 通信。`,

		&ServerSettingsTenantIDHelp: `可选标识符，当多个 Peepie 实例共享同一主机和相同的后端服务时，用于为 PostgreSQL Schema、数据目录、Docker 对象、Graphiti 组 ID、认证 Cookie 和遥测划分命名空间。

单实例部署请留空（默认）。设置后，其值必须匹配 ^[a-z][a-z0-9_]{0,31}$ ——以小写字母开头，其后为小写字母、数字或下划线，最多 32 个字符。不允许使用连字符。

示例：
• acme
• team_alpha
• staging01`,

		&ServerSettingsPprofAddrHelp: `Go pprof HTTP 端点的可选监听地址，用于 CPU/内存性能分析。

留空则保持 pprof 禁用（默认，生产环境推荐）。设置后，其值必须是 host:port 形式，例如 :7777 或 127.0.0.1:7778。

共享主机网络命名空间的实例需要使用不同的地址——端口不是字符串命名空间，无法从 TENANT_ID 推导。

示例：
• :7777
• 127.0.0.1:7778
• 0.0.0.0:7780`,

		&ServerSettingsHostHelp: `docker-compose 端口映射中发布端口的绑定地址。

示例：
• 127.0.0.1 — 仅本地访问
• 0.0.0.0 — 在所有接口上暴露`,

		&ServerSettingsPortHelp: `Peepie UI 的外部端口。该端口必须在主机上可用。示例：8443。`,

		&ServerSettingsPublicURLHelp: `设置用于重定向和链接的公共基础 URL。

示例：
• http://localhost:8443
• https://example.com/
• https://example.com/pentagi/（带基础路径）`,

		&ServerSettingsCORSOriginsHelp: `以逗号分隔的允许浏览器访问的来源。`,

		&ServerSettingsProxyURLHelp: `访问 LLM 提供商和外部工具的出站请求所用的 HTTP 或 HTTPS 代理。不用于 Docker API 通信。`,

		&ServerSettingsHTTPClientTimeoutHelp: `所有外部 HTTP/HTTPS API 调用的超时时间（秒），包括：
• LLM 提供商请求（OpenAI、Anthropic、Bedrock 等）
• 搜索引擎查询（Google、Tavily、Perplexity 等）
• 外部工具集成
• 嵌入生成请求

默认值：600 秒（10 分钟）
设置为 0 会禁用超时（不建议在生产环境中使用）
值过低可能导致正常的长时间请求失败。`,

		&ServerSettingsTerminalToolTimeoutHelp: `当智能体请求 timeout=0 或负数超时值时所使用的默认超时时间（秒）。

这会影响通过隔离终端容器执行的命令，包括扫描器和基于 CLI 的工具。

默认值：1200 秒（20 分钟）
允许范围：1–10800 秒（最长 3 小时）
小于等于 0 或大于 10800 的值会被限制为最大值（10800 秒 = 3 小时）；智能体永远不允许无限期运行。
工具调用显式提供的超时值在 1–10800 秒范围内时会覆盖此默认值。`,

		&ServerSettingsExternalSSLCAPathHelp: `容器内自定义 CA 证书文件（PEM 格式）的路径。

必须指向 /opt/pentagi/ssl/ 目录，该目录由主机上的 pentagi-ssl 卷挂载。

示例：
• /opt/pentagi/ssl/ca-bundle.pem
• /opt/pentagi/ssl/corporate-ca.pem

文件可以包含多个根证书和中间证书。`,

		&ServerSettingsExternalSSLInsecureHelp: `禁用连接 LLM 提供商和外部服务时的 SSL/TLS 证书校验。

⚠ 警告：仅用于使用自签名证书的测试。切勿在生产环境中启用。

启用后将绕过所有证书校验，使连接容易遭受中间人攻击。`,

		&ServerSettingsSSLDirHelp: `包含 PEM 格式 server.crt 和 server.key 的目录路径。server.crt 可包含完整证书链。设置后将覆盖默认生成证书的行为。`,

		&ServerSettingsDataDirHelp: `用于持久化数据的主机目录。Peepie 将智能体产物存放在 flow-N 子目录下，这些子目录映射为工作节点容器内的 /work。`,

		&ServerSettingsCookieSigningSaltHelp: `用于签名 Cookie 的密钥盐值。请妥善保密。`,

		&ServerSettingsDatabaseExtensionsSchemaHelp: `存放共享扩展（vector、pg_trgm）的 PostgreSQL Schema，每个租户都必须能通过其 search_path 访问这些扩展。

仅在设置租户 ID 时使用；留空则使用默认值 "public"，即标准 PostgreSQL 安装存放扩展的位置。当你的数据库遵循其他约定时请设置此项——Supabase 将扩展安装在 "extensions" 中；如果与此不匹配，启动会中止，并给出其找到的 Schema 名称。

示例：
• public
• extensions`,

		&ServerSettingsDatabaseSearchPathViaOptionsHelp: `以 options=--search_path=<value> 形式发送租户 search_path，而不是使用单独的 search_path 连接参数。

仅在设置租户 ID 时使用。直接连接 PostgreSQL 时请保持 false。当通过会转发 "options" 启动参数但会丢弃无法识别的单独 search_path 参数的连接池连接时启用此项——据报告某些版本的 Supabase Supavisor 存在这种情况。此方式不保证有效：如果该值未能到达后端，启动会失败并给出明确的 Schema 不匹配错误。

可选值：true、false`,

		// Human-in-the-loop screen strings
		&ToolsAIAgentsSettingsFormTitle:       "AI 智能体设置",
		&ToolsAIAgentsSettingsFormDescription: "配置 AI 智能体的全局行为",
		&ToolsAIAgentsSettingsFormName:        "AI 智能体设置",
		&ToolsAIAgentsSettingsFormOverview: `此部分配置 Peepie 中 AI 智能体的全局行为。

基本设置：
• 启用用户交互：允许智能体在需要时请求用户输入
• 使用多智能体模式：让助手协调多个专业智能体

执行监控（⚠️  BETA）：
• 启用执行监控：由导师自动监督并进行模式分析
• 相同工具调用阈值：触发导师审查前连续相同工具调用的次数
• 工具调用总数阈值：触发导师审查前的工具调用总次数

工具调用限制：
• 最大工具调用次数（通用智能体）：防止 Assistant、Primary Agent、Pentester、Coder、Installer 失控执行
• 最大工具调用次数（受限智能体）：防止 Searcher、Enricher、Memorist 等失控执行

任务规划（⚠️  BETA）：
• 启用任务规划：为专业智能体生成结构化执行计划

⚠️  BETA 功能仍在积极开发中，仅建议用于测试。`,

		&ToolsAIAgentsSettingHumanInTheLoop:          "启用用户交互",
		&ToolsAIAgentsSettingHumanInTheLoopDesc:      "允许智能体在需要时请求用户输入",
		&ToolsAIAgentsSettingUseAgents:               "使用多智能体模式",
		&ToolsAIAgentsSettingUseAgentsDesc:           "让助手协调多个专业智能体",
		&ToolsAIAgentsSettingExecutionMonitor:        "启用执行监控（beta）",
		&ToolsAIAgentsSettingExecutionMonitorDesc:    "自动调用导师进行执行模式分析",
		&ToolsAIAgentsSettingSameToolLimit:           "相同工具调用阈值",
		&ToolsAIAgentsSettingSameToolLimitDesc:       "触发导师审查前连续相同工具调用的次数",
		&ToolsAIAgentsSettingTotalToolLimit:          "工具调用总数阈值",
		&ToolsAIAgentsSettingTotalToolLimitDesc:      "触发导师审查前的工具调用总次数",
		&ToolsAIAgentsSettingMaxGeneralToolCalls:     "最大工具调用次数（通用智能体）",
		&ToolsAIAgentsSettingMaxGeneralToolCallsDesc: "Assistant、Primary Agent、Pentester、Coder、Installer 的最大工具调用次数",
		&ToolsAIAgentsSettingMaxLimitedToolCalls:     "最大工具调用次数（受限智能体）",
		&ToolsAIAgentsSettingMaxLimitedToolCallsDesc: "Searcher、Enricher、Memorist 等的最大工具调用次数",
		&ToolsAIAgentsSettingTaskPlanning:            "启用任务规划（beta）",
		&ToolsAIAgentsSettingTaskPlanningDesc:        "为专业智能体生成结构化执行计划",

		&ToolsAIAgentsSettingsHelp: `AI 智能体设置决定智能体如何协作、如何与用户交互以及如何控制执行。

基本设置：
• 启用用户交互：允许智能体在需要时请求用户输入
• 使用多智能体模式：让助手协调专业智能体处理复杂任务

执行监控（⚠️  BETA）：
自动调用顾问（导师）分析执行模式、检测循环、建议替代策略，并防止智能体执着于单一方法。阈值：连续相同调用次数（默认：5）和调用总次数（默认：10）。

任务规划（⚠️  BETA）：
在专业智能体开始工作前生成 3-7 步的执行计划。防止范围蔓延并提高成功率。当顾问使用增强配置（更强的模型或最高推理模式）时效果最佳。

工具调用限制（始终生效）：
硬性限制可防止无限循环：通用智能体默认 100 次，受限智能体默认 20 次。独立于 beta 功能工作。

小于 32B 的开源模型（Qwen3.5-27B、DeepSeek-V3、Llama-3.1-70B）：
✓ 启用这两项 beta 功能——对获得高质量结果至关重要
✓ 测试显示结果质量相比基线提升 2 倍
✓ 为顾问配置增强设置以获得最佳性能
✓ 非常适合使用本地 LLM 推理的离线隔离部署

性能：Token 用量/耗时增加 2-3 倍，小于 32B 的模型质量提升 2 倍。

⚠️  BETA 警告：功能仍在积极开发中。尽管处于 beta 阶段，仍推荐小于 32B 的开源模型使用。对于使用更大模型的云端 API，请保持禁用。

注意：更改需要重启服务后生效。`,
	})
}
