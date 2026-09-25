package locale

func init() {
	registerCatalog(LanguageChineseSimplified, map[*string]string{
		// Search Engines screen strings
		&ToolsSearchEnginesFormTitle:       "搜索引擎配置",
		&ToolsSearchEnginesFormDescription: "配置 AI 智能体在测试期间用于收集情报的搜索引擎",
		&ToolsSearchEnginesFormName:        "搜索引擎",
		&ToolsSearchEnginesFormOverview: `可用的搜索引擎：
• DuckDuckGo - 免费搜索引擎（无需 API 密钥）
• Sploitus - 安全漏洞利用与漏洞数据库（无需 API 密钥）
• Perplexity - 具备推理能力的 AI 驱动搜索
• Tavily - 面向 AI 应用的搜索 API
• Firecrawl - 带整页抓取的网页搜索（支持自托管）
• Traversaal - 网页抓取与搜索
• Google Search - 需要 API 密钥和自定义搜索引擎 ID
• Searxng - 互联网元搜索引擎
• 内部分析引擎 - 可选的基于浏览器的后备方案（抓取 + 摘要，默认关闭，无需 API 密钥）

API 密钥获取地址：
• Perplexity: https://www.perplexity.ai/
• Tavily: https://tavily.com/
• Firecrawl: https://www.firecrawl.dev/
• Traversaal: https://traversaal.ai/
• Google: https://developers.google.com/custom-search/v1/introduction`,

		&ToolsSearchEnginesDuckDuckGo:               "DuckDuckGo 搜索",
		&ToolsSearchEnginesDuckDuckGoDesc:           "启用 DuckDuckGo 搜索（无需 API 密钥）",
		&ToolsSearchEnginesDuckDuckGoRegion:         "DuckDuckGo 地区",
		&ToolsSearchEnginesDuckDuckGoRegionDesc:     "DuckDuckGo 地区代码（例如 us-en、uk-en、cn-zh）",
		&ToolsSearchEnginesDuckDuckGoSafeSearch:     "DuckDuckGo 安全搜索",
		&ToolsSearchEnginesDuckDuckGoSafeSearchDesc: "DuckDuckGo 安全搜索（strict、moderate、off）",
		&ToolsSearchEnginesDuckDuckGoTimeRange:      "DuckDuckGo 时间范围",
		&ToolsSearchEnginesDuckDuckGoTimeRangeDesc:  "DuckDuckGo 时间范围（d：天，w：周，m：月，y：年）",
		&ToolsSearchEnginesSploitus:                 "Sploitus 搜索",
		&ToolsSearchEnginesSploitusDesc:             "启用 Sploitus 搜索漏洞利用和漏洞（无需 API 密钥）",
		&ToolsSearchEnginesPerplexityKey:            "Perplexity API 密钥",
		&ToolsSearchEnginesPerplexityKeyDesc:        "Perplexity AI 搜索的 API 密钥",
		&ToolsSearchEnginesTavilyKey:                "Tavily API 密钥",
		&ToolsSearchEnginesTavilyKeyDesc:            "Tavily 搜索服务的 API 密钥",
		&ToolsSearchEnginesBraveKey:                 "Brave 搜索 API 密钥",
		&ToolsSearchEnginesBraveKeyDesc:             "Brave 搜索引擎 API 密钥（https://brave.com/search/api/）",
		&ToolsSearchEnginesFirecrawlKey:             "Firecrawl API 密钥",
		&ToolsSearchEnginesFirecrawlKeyDesc:         "Firecrawl 搜索服务的 API 密钥",
		&ToolsSearchEnginesFirecrawlURL:             "Firecrawl API URL",
		&ToolsSearchEnginesFirecrawlURLDesc:         "Firecrawl API 基础 URL（使用云服务时留空；自托管时填写）",
		&ToolsSearchEnginesTraversaalKey:            "Traversaal API 密钥",
		&ToolsSearchEnginesTraversaalKeyDesc:        "Traversaal 网页抓取的 API 密钥",
		&ToolsSearchEnginesGoogleKey:                "Google Search API 密钥",
		&ToolsSearchEnginesGoogleKeyDesc:            "Google Custom Search API 密钥",
		&ToolsSearchEnginesGoogleCX:                 "Google 搜索引擎 ID",
		&ToolsSearchEnginesGoogleCXDesc:             "Google 自定义搜索引擎 ID",
		&ToolsSearchEnginesGoogleLR:                 "Google 语言限制",
		&ToolsSearchEnginesGoogleLRDesc:             "Google 搜索引擎语言限制（例如 lang_en、lang_zh-CN 等）",
		&ToolsSearchEnginesSearxngURL:               "Searxng 搜索 URL",
		&ToolsSearchEnginesSearxngURLDesc:           "Searxng 搜索引擎 URL",
		&ToolsSearchEnginesSearxngCategories:        "Searxng 搜索类别",
		&ToolsSearchEnginesSearxngCategoriesDesc:    "Searxng 搜索引擎类别（例如 general、it、web、news、technology、science、health、other）",
		&ToolsSearchEnginesSearxngLanguage:          "Searxng 搜索语言",
		&ToolsSearchEnginesSearxngLanguageDesc:      "Searxng 搜索引擎语言（en、fr、de、it、es、pt、ru、zh，留空表示所有语言）",
		&ToolsSearchEnginesSearxngSafeSearch:        "Searxng 安全搜索",
		&ToolsSearchEnginesSearxngSafeSearchDesc:    "Searxng 搜索引擎安全搜索（0：关闭，1：中等，2：严格）",
		&ToolsSearchEnginesSearxngTimeRange:         "Searxng 时间范围",
		&ToolsSearchEnginesSearxngTimeRangeDesc:     "Searxng 搜索引擎时间范围（day、month、year）",
		&ToolsSearchEnginesSearxngTimeout:           "Searxng 超时",
		&ToolsSearchEnginesSearxngTimeoutDesc:       "Searxng 请求超时时间（秒）",
		&ToolsSearchEnginesInternalEnabled:          "内部分析引擎",
		&ToolsSearchEnginesInternalEnabledDesc:      "为回答/研究类查询启用内置的浏览器分析后备方案（无需 API 密钥；它会抓取并摘要网页，因此需要已配置的抓取器以及至少一个可用的链接引擎，例如 DuckDuckGo 或 Google）",
		&ToolsSearchEnginesInternalMaxSites:         "内部引擎最大站点数",
		&ToolsSearchEnginesInternalMaxSitesDesc:     "每次查询抓取并摘要的最大页面数",
		&ToolsSearchEnginesInternalMaxSiteBytes:     "内部引擎单站点最大字节数",
		&ToolsSearchEnginesInternalMaxSiteBytesDesc: "截断前从每个页面读取的最大 Markdown 字节数",

		// Scraper screen strings
		&ToolsScraperFormTitle:       "抓取器配置",
		&ToolsScraperFormDescription: "配置网页抓取服务",
		&ToolsScraperFormName:        "抓取器",
		&ToolsScraperFormOverview: `使用 vxcontrol/scraper Docker 镜像进行内容提取与分析的网页抓取服务。

模式：
• 内置 - 运行本地抓取器容器（推荐）
• 外部 - 使用外部抓取服务
• 禁用 - 不提供网页抓取功能

Docker 镜像：https://hub.docker.com/r/vxcontrol/scraper

抓取器支持：
• 通过公共 URL 访问外部链接
• 通过私有 URL 访问内部/本地链接
• 内容提取与分析
• 多种输出格式`,

		&ToolsScraperModeTitle:                 "抓取器模式",
		&ToolsScraperModeDesc:                  "选择抓取服务的运行方式",
		&ToolsScraperEmbedded:                  "内置容器",
		&ToolsScraperExternal:                  "外部服务",
		&ToolsScraperDisabled:                  "禁用",
		&ToolsScraperPublicURL:                 "公共抓取器 URL",
		&ToolsScraperPublicURLDesc:             "用于抓取公共/外部网站的 URL。留空则使用与私有 URL 相同的值。",
		&ToolsScraperPublicURLEmbeddedDesc:     "内置抓取器的 URL（可选覆盖）。留空则使用与私有 URL 相同的值。",
		&ToolsScraperPrivateURL:                "私有抓取器 URL",
		&ToolsScraperPrivateURLDesc:            "用于抓取私有/内部网站的 URL",
		&ToolsScraperPublicUsername:            "公共 URL 用户名",
		&ToolsScraperPublicUsernameDesc:        "访问公共抓取器的用户名",
		&ToolsScraperPublicPassword:            "公共 URL 密码",
		&ToolsScraperPublicPasswordDesc:        "访问公共抓取器的密码",
		&ToolsScraperPrivateUsername:           "私有 URL 用户名",
		&ToolsScraperPrivateUsernameDesc:       "访问私有抓取器的用户名",
		&ToolsScraperPrivatePassword:           "私有 URL 密码",
		&ToolsScraperPrivatePasswordDesc:       "访问私有抓取器的密码",
		&ToolsScraperLocalUsername:             "本地 URL 用户名",
		&ToolsScraperLocalUsernameDesc:         "内置抓取服务的用户名",
		&ToolsScraperLocalPassword:             "本地 URL 密码",
		&ToolsScraperLocalPasswordDesc:         "内置抓取服务的密码",
		&ToolsScraperMaxConcurrentSessions:     "最大并发会话数",
		&ToolsScraperMaxConcurrentSessionsDesc: "并发抓取会话的最大数量",
		&ToolsScraperEmbeddedHelp:              "内置模式会运行一个本地抓取器容器，可同时访问公共和私有资源。默认配置使用 https://someuser:somepass@scraper/。",
		&ToolsScraperExternalHelp:              "外部模式使用独立的抓取服务。可根据需要为公共访问和私有访问配置不同的 URL。",
		&ToolsScraperDisabledHelp:              "抓取器已禁用。网页内容提取与分析功能将不可用。",

		// Docker Environment screen strings
		&ToolsDockerFormTitle:       "Docker 环境配置",
		&ToolsDockerFormDescription: "为工作节点容器配置 Docker 环境",
		&ToolsDockerFormName:        "Docker 环境",
		&ToolsDockerFormOverview: `• 工作节点隔离 - 容器为任务提供安全边界
• 网络能力 - 为渗透测试启用特权网络操作
• 容器管理 - 控制工作节点如何访问 Docker 守护进程
• 存储配置 - 定义工作区和产物存储
• 镜像选择 - 为不同类型的任务设置默认镜像

对于需要网络扫描、自定义工具和安全任务隔离的渗透测试工作流至关重要。`,

		&ToolsDockerGeneralHelp: `每个 AI 智能体任务都在隔离的 Docker 容器中运行，每个流程会自动分配两个端口（28000-32000 范围）。工作节点容器会按需从默认镜像或智能体选择的镜像创建。

基本设置需要启用相关能力：Docker 访问允许为专用工具启动额外的容器，网络管理则授予 nmap 等扫描工具所必需的底层网络权限。

默认通过 Docker 卷提供存储；指定工作目录时则使用主机目录。连接设置控制 Docker 守护进程的位置——标准部署使用本地套接字，分布式环境使用带 TLS 的远程 TCP。

默认镜像作为后备：一般任务使用标准镜像，而安全测试默认使用面向渗透测试的容器。公共 IP 为工作节点提供可供目标回连的地址，从而支持反向 Shell 攻击。通常它是运行工作节点容器所用 Docker 守护进程的主机的本地接口地址。

根据场景组合配置：启用两项能力以进行完整的渗透测试，使用工作目录持久保存产物，或为隔离的 Docker 环境配置远程连接。`,

		&ToolsDockerInside:       "Docker 访问",
		&ToolsDockerInsideDesc:   "允许工作节点管理 Docker 容器",
		&ToolsDockerNetAdmin:     "网络管理",
		&ToolsDockerNetAdminDesc: "为 nmap 等网络扫描工具授予 NET_ADMIN 能力",

		&ToolsDockerSocket:       "Docker 套接字",
		&ToolsDockerSocketDesc:   "主机文件系统上 Docker 套接字的路径",
		&ToolsDockerNetwork:      "Docker 网络",
		&ToolsDockerNetworkDesc:  "工作节点容器的自定义网络名称，或填写 'host' 直接访问主机网络",
		&ToolsDockerPublicIP:     "公共 IP 地址",
		&ToolsDockerPublicIPDesc: "用于 OOB 攻击中反向连接的公共 IP",

		&ToolsDockerWorkDir:     "工作目录",
		&ToolsDockerWorkDirDesc: "工作节点文件系统所用的主机目录（默认：Docker 卷）",

		&ToolsDockerDefaultImage:               "默认镜像",
		&ToolsDockerDefaultImageDesc:           "一般任务使用的默认 Docker 镜像",
		&ToolsDockerDefaultImageForPentest:     "渗透测试镜像",
		&ToolsDockerDefaultImageForPentestDesc: "安全测试任务使用的默认 Docker 镜像",

		&ToolsDockerHost:          "Docker 主机",
		&ToolsDockerHostDesc:      "Docker 守护进程连接（unix:// 或 tcp://）",
		&ToolsDockerTLSVerify:     "TLS 验证",
		&ToolsDockerTLSVerifyDesc: "为 Docker 连接启用 TLS 验证",
		&ToolsDockerCertPath:      "TLS 证书",
		&ToolsDockerCertPathDesc:  "包含 ca.pem、cert.pem、key.pem 文件的目录",

		&ToolsDockerInsideHelp: `Docker 访问使工作节点能够为专用工具和环境启动额外的容器。当任务需要默认镜像中没有的自定义软件时必须启用。

启用后，工作节点可以拉取并运行任意 Docker 镜像，为复杂的测试场景提供最大的灵活性。`,

		&ToolsDockerNetAdminHelp: `网络管理能力允许工作节点执行渗透测试所必需的底层网络操作。

以下场景需要此能力：
• 使用 nmap、masscan 进行网络扫描
• 自定义构造数据包
• 操作网络接口
• 原始套接字操作

对于全面的安全评估至关重要。`,

		&ToolsDockerSocketHelp: `Docker 套接字路径决定工作节点如何访问 Docker 守护进程。只能填写套接字文件的路径。与 Docker 访问配合使用以启用容器管理。

为增强安全性，建议使用 docker-in-docker（DinD），而不是将主 Docker 守护进程直接暴露给工作节点。
使用 DinD 时，请填写已绑定到主机文件系统的 DinD 容器 Docker 套接字文件路径。

示例：/var/run/docker.sock`,

		&ToolsDockerNetworkHelp: `Docker 网络控制工作节点容器的网络隔离模式：

桥接模式（自定义网络名称）：
• 容器之间隔离通信
• 从容器到主机的端口转发
• 更强的安全边界
• 基于网络的监控和过滤
• 推荐用于大多数场景

主机模式（值：'host'）：
• 直接访问主机网络接口
• 无端口转发——端口直接绑定到主机
• 原始数据包操作所必需
• 高级网络测试能力
• 隔离性较低——请谨慎使用

示例：
• 'pentagi-network' - 创建隔离的桥接网络
• 'host' - 启用直接访问主机网络

安全提示：主机网络模式会降低容器隔离性。仅在需要直接访问网络栈的高级渗透测试任务中使用。`,

		&ToolsDockerPublicIPHelp: `公共 IP 地址为工作节点提供可供反向连接的地址，从而支持带外（OOB）攻击技术。

工作节点会自动获得两个映射到该 IP 的随机端口（28000-32000 范围），用于接收被利用目标的回连。

默认情况下，智能体会尝试通过 api.ipify.org、ipinfo.io/ip 或 ifconfig.me 服务获取公网地址。`,

		&ToolsDockerWorkDirHelp: `工作目录指定工作节点存储在主机文件系统上的位置。设置后，将用主机目录挂载替代默认的 Docker 卷。

优点：
• 重启后存储仍然保留
• 直接访问文件系统
• 更方便管理产物
• 可自定义备份策略

默认为每个工作节点容器使用专用的 Docker 卷。

示例：/path/to/workdir/`,

		&ToolsDockerDefaultImageHelp: `当任务需求未指定特定容器镜像时，默认镜像作为工作节点的后备镜像。

应包含用于通用任务的基础实用程序和工具。默认值：debian:latest`,

		&ToolsDockerDefaultImageForPentestHelp: `渗透测试镜像是安全测试任务的默认镜像，应包含全面的安全工具和实用程序。

推荐的镜像包括 Kali Linux、Parrot Security 或自定义的安全专用容器。默认值：vxcontrol/kali-linux`,

		&ToolsDockerHostHelp: `Docker 主机用于启动主要的工作节点容器，并覆盖默认的 Docker 守护进程连接。支持 Unix 套接字和 TCP 连接。

示例：
• unix:///var/run/docker.sock（本地）
• tcp://docker-host:2376（远程）

远程连接请启用 TLS。`,

		&ToolsDockerTLSVerifyHelp: `TLS 验证可保护基于 TCP 的 Docker 守护进程连接。强烈建议远程 Docker 主机启用。

需要在指定的证书目录中提供有效证书。`,

		&ToolsDockerCertPathHelp: `TLS 证书目录必须包含：
• ca.pem - 证书颁发机构
• cert.pem - 客户端证书
• key.pem - 私钥

使用 TLS 管理工作节点容器时，安全的远程 Docker 连接需要这些文件。

示例：/path/to/certs`,

		// Embedder form strings
		&EmbedderFormTitle:       "嵌入配置",
		&EmbedderFormDescription: "为语义搜索和知识存储配置文本向量化",
		&EmbedderFormName:        "嵌入器",
		&EmbedderFormOverview: `文本嵌入将文档转换为向量，用于语义搜索和知识存储。
不同的提供商提供能力和价格各异的模型。

请谨慎选择，更换提供商需要重建所有已存储数据的索引。`,

		&EmbedderFormProvider:     "嵌入提供商",
		&EmbedderFormProviderDesc: "选择用于文本向量化的提供商。嵌入用于语义搜索和知识存储。",

		&EmbedderFormURL:     "API 端点 URL",
		&EmbedderFormURLDesc: "自定义 API 端点（留空则使用默认值）",

		&EmbedderFormAPIKey:     "API 密钥",
		&EmbedderFormAPIKeyDesc: "提供商的认证密钥（Ollama 无需填写）",

		&EmbedderFormModel:     "模型名称",
		&EmbedderFormModelDesc: "要使用的具体嵌入模型（留空则使用提供商默认模型）",

		&EmbedderFormBatchSize:     "批大小",
		&EmbedderFormBatchSizeDesc: "单批处理的文档数量（1-1000）",

		&EmbedderFormStripNewLines:     "去除换行符",
		&EmbedderFormStripNewLinesDesc: "嵌入前移除文本中的换行符（true/false）",

		&EmbedderFormMaxTextBytes:     "最大文本字节数",
		&EmbedderFormMaxTextBytesDesc: "发送到嵌入 API 的每个文本块的最大字节数（例如 8192）",

		&EmbedderFormHelpTitle: "嵌入配置",
		&EmbedderFormHelpContent: `为语义搜索和知识存储配置文本向量化。

如果未配置特定的嵌入设置，系统将使用 OpenAI 嵌入，并使用 LLM 提供商中的 API 密钥。

请谨慎更换提供商——不同的嵌入器生成的向量互不兼容，需要重建数据库索引。`,

		&EmbedderFormHelpOpenAI:      "OpenAI：最可靠的选项，质量出色。如果此处未设置，则需要 LLM 提供商中的 API 密钥。",
		&EmbedderFormHelpOllama:      "Ollama：本地嵌入，无需 API 密钥。需要运行 Ollama 服务器。",
		&EmbedderFormHelpHuggingFace: "HuggingFace：开源模型，需要 API 密钥。",
		&EmbedderFormHelpGoogleAI:    "Google AI：高质量嵌入，需要 API 密钥。",

		&EmbedderProviderDefault:         "默认（OpenAI）",
		&EmbedderProviderDefaultDesc:     "使用 OpenAI 嵌入，API 密钥取自 LLM 提供商配置",
		&EmbedderProviderOpenAI:          "OpenAI",
		&EmbedderProviderOpenAIDesc:      "OpenAI 文本嵌入 API（text-embedding-3-small、ada-002）",
		&EmbedderProviderOllama:          "Ollama",
		&EmbedderProviderOllamaDesc:      "用于开源嵌入模型的本地 Ollama 服务器",
		&EmbedderProviderMistral:         "Mistral",
		&EmbedderProviderMistralDesc:     "Mistral AI 嵌入模型",
		&EmbedderProviderJina:            "Jina",
		&EmbedderProviderJinaDesc:        "Jina AI 嵌入 API",
		&EmbedderProviderHuggingFace:     "HuggingFace",
		&EmbedderProviderHuggingFaceDesc: "用于嵌入模型的 HuggingFace 推理 API",
		&EmbedderProviderGoogleAI:        "Google AI",
		&EmbedderProviderGoogleAIDesc:    "Google AI 嵌入模型（embedding-001）",
		&EmbedderProviderVoyageAI:        "VoyageAI",
		&EmbedderProviderVoyageAIDesc:    "VoyageAI 嵌入 API",
		&EmbedderProviderDisabled:        "禁用",
		&EmbedderProviderDisabledDesc:    "完全禁用嵌入功能",

		&EmbedderURLPlaceholderOpenAI:      "https://api.openai.com/v1",
		&EmbedderURLPlaceholderOllama:      "http://localhost:11434",
		&EmbedderURLPlaceholderMistral:     "https://api.mistral.ai/v1",
		&EmbedderURLPlaceholderJina:        "https://api.jina.ai/v1",
		&EmbedderURLPlaceholderHuggingFace: "https://api-inference.huggingface.co",
		&EmbedderURLPlaceholderGoogleAI:    "不支持 - 使用默认端点",
		&EmbedderURLPlaceholderVoyageAI:    "不支持 - 使用默认端点",

		&EmbedderAPIKeyPlaceholderOllama:      "本地模型无需填写",
		&EmbedderAPIKeyPlaceholderMistral:     "Mistral API 密钥",
		&EmbedderAPIKeyPlaceholderJina:        "Jina API 密钥",
		&EmbedderAPIKeyPlaceholderHuggingFace: "HuggingFace API 密钥",
		&EmbedderAPIKeyPlaceholderGoogleAI:    "Google AI API 密钥",
		&EmbedderAPIKeyPlaceholderVoyageAI:    "VoyageAI API 密钥",
		&EmbedderAPIKeyPlaceholderDefault:     "提供商的 API 密钥",

		&EmbedderModelPlaceholderOpenAI:      "text-embedding-3-small",
		&EmbedderModelPlaceholderOllama:      "nomic-embed-text",
		&EmbedderModelPlaceholderMistral:     "mistral-embed",
		&EmbedderModelPlaceholderJina:        "jina-embeddings-v2-base-en",
		&EmbedderModelPlaceholderHuggingFace: "sentence-transformers/all-MiniLM-L6-v2",
		&EmbedderModelPlaceholderGoogleAI:    "gemini-embedding-001",
		&EmbedderModelPlaceholderVoyageAI:    "voyage-2",
		&EmbedderModelPlaceholderDefault:     "模型名称",

		&EmbedderHelpGeneral: `嵌入将文本转换为向量，用于语义搜索和知识存储。这使 Peepie 能够理解含义而不仅仅是关键词，让搜索结果更相关、更智能。

主要优势：
• 按含义而非精确词语查找文档
• 基于渗透测试结果构建智能知识库
• 让 AI 智能体快速定位相关信息
• 借助上下文数据支持高级推理

选择 Ollama 可实现完全本地处理——你的数据永远不会离开你的基础设施。其他提供商提供基于云的处理，模型能力和价格各不相同。

请谨慎配置，更换提供商需要重建整个知识库。`,

		&EmbedderHelpAttentionPrefix: "重要：",
		&EmbedderHelpAttention: `不同的嵌入提供商生成的向量互不兼容。更换提供商或模型会使现有的语义搜索失效。

你必须使用 etester 工具清空或重建整个知识库的索引：
• 运行 'etester flush' 清除旧的嵌入
• 运行 'etester reindex' 使用新提供商重建
• 对于大型数据集，此过程可能需要相当长的时间`,

		&EmbedderHelpAttentionSuffix: `除非绝对必要，否则不要更换提供商。`,

		&EmbedderHelpDefault: `默认模式使用 OpenAI 嵌入，并使用 LLM 提供商中配置的 API 密钥。

这是推荐大多数用户使用的选项，如果你已经配置了 OpenAI，则无需任何额外配置。`,

		&EmbedderHelpOpenAI: `直接访问 OpenAI API 生成嵌入。

API 密钥获取地址：
https://platform.openai.com/api-keys

推荐模型：
• text-embedding-3-small（性价比高，1536 维）
• text-embedding-3-large（质量最高，3072 维）
• text-embedding-ada-002（旧版，仍受支持）`,

		&EmbedderHelpOllama: `用于开源嵌入模型的本地 Ollama 服务器。

常用嵌入模型：
• nomic-embed-text（推荐，768 维）
• mxbai-embed-large（大模型，1024 维）
• snowflake-arctic-embed（支持多语言）

Ollama 安装地址：
https://ollama.com/

开始使用：ollama pull nomic-embed-text`,

		&EmbedderHelpMistral: `通过 API 使用 Mistral AI 嵌入模型。

API 密钥获取地址：
https://console.mistral.ai/

使用配置固定的 Mistral 嵌入模型。
无需选择模型——使用默认嵌入模型。`,

		&EmbedderHelpJina: `提供专用模型的 Jina AI 嵌入 API。

API 密钥获取地址：
https://jina.ai/

推荐模型：
• jina-embeddings-v2-base-en（通用，768 维）
• jina-embeddings-v2-small-en（轻量，512 维）
• jina-embeddings-v2-base-code（代码专用嵌入）`,

		&EmbedderHelpHuggingFace: `用于开源模型的 HuggingFace Inference API。

API 密钥获取地址：
https://huggingface.co/settings/tokens

常用模型：
• sentence-transformers/all-MiniLM-L6-v2（384 维）
• sentence-transformers/all-mpnet-base-v2（768 维）
• intfloat/e5-large-v2（1024 维）`,

		&EmbedderHelpGoogleAI: `Google AI 嵌入模型（Gemini）。

API 密钥获取地址：
https://aistudio.google.com/app/apikey

可用模型：
• gemini-embedding-001（最新模型，768 维）
• text-embedding-004（旧版 Vertex AI 模型）

使用 Google 的固定端点——不支持配置 URL。`,

		&EmbedderHelpVoyageAI: `针对检索优化的 VoyageAI 嵌入 API。

API 密钥获取地址：
https://www.voyageai.com/

推荐模型：
• voyage-2（通用，1024 维）
• voyage-large-2（质量最高，1536 维）
• voyage-code-2（代码嵌入，1536 维）`,

		&EmbedderHelpDisabled: `禁用所有嵌入功能。

这将会：
• 禁用语义搜索功能
• 关闭知识存储向量化
• 降低内存和计算需求

仅在你的使用场景不需要嵌入时推荐。`,
	})
}
