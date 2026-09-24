package locale

func init() {
	registerCatalog(LanguageChineseSimplified, map[*string]string{
		// Common status and UI strings
		&UIStatistics:       "统计",
		&UIStatus:           "状态：",
		&UIMode:             "模式：",
		&UINoConfigSelected: "未选择配置",
		&UILoading:          "加载中...",
		&UINotImplemented:   "尚未实现",
		&UIUnsavedChanges:   "有未保存的更改",
		&UIConfigSaved:      "配置已保存",

		// Status labels
		&StatusEnabled:       "已启用",
		&StatusDisabled:      "已禁用",
		&StatusConfigured:    "已配置",
		&StatusNotConfigured: "未配置",
		&StatusEmbedded:      "内置",
		&StatusExternal:      "外部",

		// Success/Warning messages
		&MessageSearchEnginesNone:       "⚠ 未配置搜索引擎",
		&MessageSearchEnginesConfigured: "✓ 已配置 %d 个搜索引擎",
		&MessageDockerConfigured:        "✓ Docker 环境已配置",
		&MessageDockerNotConfigured:     "⚠ Docker 环境未配置",

		// Legend
		&LegendConfigured:    "✓ 已配置",
		&LegendNotConfigured: "✗ 未配置",

		// Navigation
		&NavBack:       "Esc: 返回",
		&NavExit:       "Ctrl+Q: 退出",
		&NavUpDown:     "↑/↓: 滚动/选择",
		&NavLeftRight:  "←/→: 移动",
		&NavPgUpPgDown: "PgUp/PgDn: 翻页",
		&NavHomeEnd:    "Home/End: 开头/结尾",
		&NavEnter:      "Enter: 继续",
		&NavYn:         "Y/N: 接受/拒绝",
		&NavCtrlC:      "Ctrl+C: 取消",
		&NavCtrlS:      "Ctrl+S: 保存",
		&NavCtrlR:      "Ctrl+R: 重置",
		&NavCtrlH:      "Ctrl+H: 显示/隐藏",
		&NavTab:        "Tab: 补全",
		&NavCtrlL:      "Ctrl+L: 切换到 %s",
		&NavSeparator:  " • ",

		// Welcome screen
		&WelcomeFormTitle:       "欢迎使用 Peepie",
		&WelcomeFormDescription: "Peepie 是一个自主渗透测试平台，利用 AI 技术执行全面的安全评估。",
		&WelcomeFormName:        "欢迎",
		&WelcomeFormOverview: `系统检查将验证：
• 环境配置文件是否存在
• Docker API 是否可访问及版本是否兼容
• 工作节点环境是否就绪
• 系统资源（CPU、内存、磁盘空间）
• 外部依赖的网络连通性

所有检查通过后，按配置向导依次设置 LLM 提供商、监控和安全工具。

安装程序会引导你完成每个组件的设置，并针对不同的部署场景给出建议。`,

		&WelcomeConfigurationFailed: "⚠ 未通过的检查：%s",
		&WelcomeConfigurationPassed: "✓ 所有系统检查均已通过",

		&WelcomeWorkflowTitle: "安装流程：",
		&WelcomeWorkflowStep1: "1. 接受最终用户许可协议",
		&WelcomeWorkflowStep2: "2. 配置 LLM 提供商（OpenAI、Anthropic 等）",
		&WelcomeWorkflowStep3: "3. 设置集成（Langfuse、可观测性）",
		&WelcomeWorkflowStep4: "4. 配置安全设置",
		&WelcomeWorkflowStep5: "5. 部署并启动 Peepie 服务",
		&WelcomeSystemReady:   "✓ 系统已就绪 - 按 Enter 继续",

		// Troubleshooting
		&TroubleshootTitle: "未满足系统要求",

		&TroubleshootEnvFileTitle: "缺少环境配置",
		&TroubleshootEnvFileDesc:  "Peepie 配置需要 .env 文件，但未找到该文件或文件不可读。",
		&TroubleshootEnvFileFix: `解决方法：
1. 在安装目录中将 .env.example 复制为 .env
2. 编辑 .env，至少配置一个 LLM 提供商的 API 密钥
3. 确保文件具有读取权限（chmod 644 .env）

快速修复：
cp .env.example .env && chmod 644 .env`,

		&TroubleshootWritePermTitle: "需要写入权限",
		&TroubleshootWritePermDesc:  "安装程序需要配置目录的写入权限，以保存设置并部署服务。",
		&TroubleshootWritePermFix: `解决方法：
1. 检查目录权限：ls -la
2. 授予写入权限：chmod 755 .
3. 或在可写的位置运行安装程序
4. 确保有足够的可用磁盘空间`,

		&TroubleshootDockerNotInstalledTitle: "未安装 Docker",
		&TroubleshootDockerNotInstalledDesc:  "此系统未安装 Docker。Peepie 需要 Docker 来运行容器。",
		&TroubleshootDockerNotInstalledFix: `解决方法：
1. 安装 Docker Desktop：https://docs.docker.com/get-docker/
2. Linux：按照对应发行版的说明安装
3. 验证安装：docker --version
4. 确保 docker 命令位于 PATH 中`,

		&TroubleshootDockerNotRunningTitle: "Docker 守护进程未运行",
		&TroubleshootDockerNotRunningDesc:  "Docker 已安装，但守护进程未运行。Docker 服务必须处于运行状态。",
		&TroubleshootDockerNotRunningFix: `解决方法：
1. 启动 Docker Desktop（Windows/Mac）
2. Linux：sudo systemctl start docker
3. 检查状态：docker ps
4. 如果使用 DOCKER_HOST，请确认远程守护进程可访问`,

		&TroubleshootDockerPermissionTitle: "Docker 权限被拒绝",
		&TroubleshootDockerPermissionDesc:  "你的用户账户没有访问 Docker 的权限。这在 Linux 系统上很常见。",
		&TroubleshootDockerPermissionFix: `解决方法：
1. 将用户加入 docker 组：sudo usermod -aG docker $USER
2. 注销并重新登录以使更改生效
3. 或使用 sudo 运行（不建议用于生产环境）
4. 验证：docker ps（应无需 sudo 即可运行）`,

		&TroubleshootDockerAPITitle: "Docker API 连接失败",
		&TroubleshootDockerAPIDesc:  "无法连接到 Docker API，可能是配置或网络问题所致。",
		&TroubleshootDockerAPIFix: `解决方法：
1. 检查 DOCKER_HOST 环境变量
2. 确认 Docker 正在运行：docker version
3. 远程 Docker：确保网络连通
4. 如果使用 TCP 连接，请检查防火墙设置
5. 尝试：export DOCKER_HOST=unix:///var/run/docker.sock`,

		&TroubleshootDockerVersionTitle: "Docker 版本过旧",
		&TroubleshootDockerVersionDesc:  "你的 Docker 版本不兼容。Peepie 需要 Docker 20.0.0 或更高版本。",
		&TroubleshootDockerVersionFix: `解决方法：
1. 将 Docker 升级到 20.0.0 或更高版本
2. 访问 https://docs.docker.com/engine/install/

当前版本：%s
要求版本：20.0.0+`,

		&TroubleshootComposeTitle: "未找到 Docker Compose",
		&TroubleshootComposeDesc:  "需要 Docker Compose v2 插件，但 `docker compose` 不可用。",
		&TroubleshootComposeFix: `解决方法：
1. 安装或更新 Docker Desktop，或为 Docker Engine 安装 Docker Compose v2 插件
2. 确认插件可用：docker compose version
3. 如果只安装了旧版 docker-compose，请同时安装 Docker Compose v2 插件

Peepie 执行的是 "docker compose"，仅有旧版 "docker-compose" 是不够的。
文档：https://docs.docker.com/compose/install/`,

		&TroubleshootComposeVersionTitle: "Docker Compose 版本过旧",
		&TroubleshootComposeVersionDesc:  "你的 `docker compose` 版本不兼容。Peepie 需要 Docker Compose 1.25.0 或更高版本。",
		&TroubleshootComposeVersionFix: `当前版本：%s
要求版本：1.25.0+

解决方法：
1. 将 Docker Desktop 或 Docker Compose v2 插件更新到较新版本
2. 通过以下命令验证结果：docker compose version

文档：https://docs.docker.com/compose/install/`,

		&TroubleshootWorkerTitle: "无法访问工作节点 Docker 环境",
		&TroubleshootWorkerDesc:  "无法连接到工作节点容器所用的 Docker 环境，可能是远程或本地 Docker 设置有问题。",
		&TroubleshootWorkerFix: `解决方法：
1. 使用远程 Docker 时，在运行安装程序前设置环境变量：
   export DOCKER_HOST=tcp://remote:2376
   export DOCKER_CERT_PATH=/path/to/certs
   export DOCKER_TLS_VERIFY=1
2. 验证连接：docker -H $DOCKER_HOST ps
3. 使用本地 Docker 时，不要设置这些变量
4. 检查防火墙是否放行 Docker 端口（2375/2376）
5. 如果使用 TLS，请确保证书有效`,

		&TroubleshootCPUTitle: "CPU 核心数不足",
		&TroubleshootCPUDesc:  "Peepie 至少需要 2 个 CPU 核心才能正常运行。",
		&TroubleshootCPUFix: `你的系统有 %d 个 CPU 核心，但至少需要 2 个。

虚拟机：
1. 在虚拟机设置中增加 CPU 分配
2. 确保宿主机有足够的资源

Docker Desktop 用户：
Settings → Resources → CPUs：设置为 2 或更多`,

		&TroubleshootMemoryTitle: "内存不足",
		&TroubleshootMemoryDesc:  "可用内存不足以运行所选组件。",
		&TroubleshootMemoryFix: `内存要求：
• 基础系统：0.5 GB
• Peepie 核心：+0.5 GB
• Langfuse（如启用）：+1.5 GB
• 可观测性（如启用）：+1.5 GB

总需求：%.1f GB
可用：%.1f GB

解决方法：
1. 关闭不必要的应用程序
2. 提高 Docker 内存限制
3. 禁用可选组件（Langfuse/可观测性）`,

		&TroubleshootDiskTitle: "磁盘空间不足",
		&TroubleshootDiskDesc:  "可用磁盘空间不足以完成安装和运行。",
		&TroubleshootDiskFix: `磁盘要求：
• 基础安装：至少 5 GB
• 含组件：10 GB + 每个组件 2 GB
• 工作节点镜像：25 GB（包括 6GB+ 的 Kali 镜像）

需要：%.1f GB
可用：%.1f GB

解决方法：
1. 释放磁盘空间
2. 为 Docker 使用外部存储
3. 清理未使用的 Docker 资源：
   docker system prune -a`,

		&TroubleshootNetworkTitle: "网络连接失败",
		&TroubleshootNetworkDesc:  "无法访问所需的外部服务，因此无法下载 Docker 镜像和更新。",
		&TroubleshootNetworkFix: `未通过的检查：
%s

解决方法：
1. 验证互联网连接：ping docker.io
2. 检查 DNS 解析：nslookup docker.io
3. 如果位于代理之后，在运行安装程序前设置：
   export HTTP_PROXY=http://proxy:port
   export HTTPS_PROXY=http://proxy:port
4. 如需持久使用代理，添加到 .env：
   PROXY_URL=http://proxy:port
5. 检查防火墙是否允许出站 HTTPS（端口 443）
6. 如果 DNS 解析失败，尝试使用其他 DNS 服务器`,

		&TroubleshootFixHint: "\n解决上述问题后，请重新运行安装程序。",

		&NetworkFailureDNS:        "• docker.io 的 DNS 解析失败",
		&NetworkFailureHTTPS:      "• 无法通过 HTTPS 访问外部服务",
		&NetworkFailureDockerPull: "• 无法从镜像仓库拉取 Docker 镜像",

		// System checks
		&ChecksTitle:               "系统检查",
		&ChecksWarningFailed:       "⚠ 部分检查未通过",
		&CheckEnvironmentFile:      "环境文件",
		&CheckWritePermissions:     "写入权限",
		&CheckDockerAPI:            "Docker API",
		&CheckDockerVersion:        "Docker 版本",
		&CheckDockerCompose:        "Docker Compose",
		&CheckDockerComposeVersion: "Docker Compose 版本",
		&CheckWorkerEnvironment:    "工作节点环境",
		&CheckSystemResources:      "系统资源",
		&CheckNetworkConnectivity:  "网络连通性",

		// EULA screen
		&EULAFormDescription: "Peepie 使用的法律条款和条件",
		&EULAFormName:        "EULA",
		&EULAFormOverview: `请阅读并接受最终用户许可协议，以继续安装 Peepie。

EULA 包含：
• 软件许可条款和使用权利
• 责任限制和担保
• 数据收集和隐私政策
• 合规要求和限制
• 支持和维护条款

你必须滚动阅读整份文档并接受条款，才能继续安装。

使用方向键、Page Up/Down 或 Home/End 键浏览文档。`,

		&EULAErrorLoadingTitle:     "# 加载 EULA 出错\n\n加载 EULA 失败：%v",
		&EULAContentFallback:       "# EULA 内容\n\n%s\n\n---\n\n*注意：Markdown 渲染失败：%v*",
		&EULAConfigurationRead:     "✓ 已阅读 EULA",
		&EULAConfigurationAccepted: "✓ 已接受 EULA",
		&EULAConfigurationPending:  "⚠ 尚未阅读 EULA",
		&EULALoading:               "正在加载 EULA...",
		&EULAProgress:              "进度：%d%%",
		&EULAProgressComplete:      " • 已读完",

		// Main menu
		&MainMenuTitle:       "Peepie 配置",
		&MainMenuDescription: "配置 Peepie 的所有组件和设置",
		&MainMenuName:        "主菜单",
		&MainMenuOverview: `欢迎使用 Peepie 配置中心。

配置核心组件：
• LLM 提供商 - 用于自主测试的 AI 语言模型
• 监控 - 可观测性和分析平台
• 工具 - 增强测试的附加能力
• 系统设置 - 环境和部署选项

依次浏览各个部分以完成 Peepie 的设置。`,

		&MenuTitle:        "配置菜单",
		&MenuSystemStatus: "系统状态",

		&MainMenuStatusPentagiRunning:     "Peepie 已在运行",
		&MainMenuStatusPentagiNotRunning:  "已准备好启动 Peepie 服务",
		&MainMenuStatusUpToDate:           "Peepie 已是最新版本",
		&MainMenuStatusUpdatesAvailable:   "有可用更新",
		&MainMenuStatusReadyToStart:       "已准备好启动",
		&MainMenuStatusAllServicesRunning: "所有服务均在运行",
		&MainMenuStatusNoUpdatesAvailable: "没有可用更新",
	})
}
