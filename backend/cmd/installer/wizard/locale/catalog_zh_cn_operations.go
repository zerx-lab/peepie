package locale

func init() {
	registerCatalog(LanguageChineseSimplified, map[*string]string{
		&MockScreenTitle:             "开发中的页面",
		&MockScreenDescription:       "此页面仍在开发中",
		&MockScreenUnderDevelopment:  "🚧 此页面仍在开发中",
		&MockScreenAvailableLater:    "此配置页面将在后续更新中提供。",
		&MockScreenGoBack:            "按 Enter 或 Esc 返回主菜单。",
		&MockScreenPendingMigration:  "⏳ 配置功能待迁移",
		&MockScreenNoticeTitle:       "开发提示",
		&MockScreenNoticeMigrating:   "此配置页面正在迁移到新界面。",
		&MockScreenExpectedFeatures:  "计划提供的功能：",
		&MockScreenFeatureForms:      "• 现代化表单界面",
		&MockScreenFeatureValidation: "• 改进的输入验证",
		&MockScreenFeatureUX:         "• 优化的用户体验",
		&MockScreenCheckBack:         "请关注后续更新。",

		&ApplyChangesFormTitle:       "应用配置更改",
		&ApplyChangesFormName:        "应用更改",
		&ApplyChangesFormDescription: "检查并应用配置更改",
		&ApplyChangesFormOverview: `此页面可供你检查所有待处理的配置更改，并将其应用到 Peepie 安装中。

应用更改后，系统将：
• 将所有修改后的环境变量保存到 .env 文件
• 使用新配置重启受影响的服务
• 根据需要安装其他组件`,
		&ApplyChangesNotStarted:               "配置更改已准备好应用",
		&ApplyChangesInProgress:               "正在应用配置更改…\n",
		&ApplyChangesCompleted:                "已成功应用配置更改\n",
		&ApplyChangesFailed:                   "应用配置更改失败",
		&ApplyChangesResetCompleted:           "已成功重置配置更改\n",
		&ApplyChangesTerminalIsNotInitialized: "终端尚未初始化",
		&ApplyChangesInstructions:             `按 Enter 开始应用配置更改。`,
		&ApplyChangesNoChanges:                "没有待处理的配置更改",
		&ApplyChangesInstallNotFound: `此系统尚未安装 Peepie。

将执行以下操作：
• 设置并验证 Docker 环境
• 创建 docker-compose.yml 文件
• 安装并启动 Peepie 核心服务`,
		&ApplyChangesInstallFoundLangfuse:      `• 安装 Langfuse 可观测性组件（docker-compose-langfuse.yml）`,
		&ApplyChangesInstallFoundObservability: `• 安装包含 Grafana、VictoriaMetrics 和 Jaeger 的完整可观测性组件（docker-compose-observability.yml）`,
		&ApplyChangesUpdateFound: `此系统已安装 Peepie。

将执行以下操作：
• 更新 .env 文件中的环境变量
• 重新创建并重启受影响的 Docker 容器
• 将新配置应用到运行中的服务`,
		&ApplyChangesWarningCritical:        "⚠️  检测到重要更改，将重启服务",
		&ApplyChangesWarningSecrets:         "🔒 检测到敏感信息，将安全存储",
		&ApplyChangesNoteBackup:             "💾 更改前将备份当前配置",
		&ApplyChangesNoteTime:               "⏱️  此过程耗时可能不到一分钟，具体取决于所选组件",
		&ApplyChangesStageValidation:        "正在验证环境和依赖项…",
		&ApplyChangesStageBackup:            "正在备份配置…",
		&ApplyChangesStageEnvFile:           "正在更新环境变量文件…",
		&ApplyChangesStageCompose:           "正在生成 Docker Compose 文件…",
		&ApplyChangesStageDocker:            "正在管理 Docker 容器…",
		&ApplyChangesStageServices:          "正在启动服务…",
		&ApplyChangesStageComplete:          "已成功应用配置更改",
		&ApplyChangesChangesTitle:           "待处理的配置更改",
		&ApplyChangesChangesCount:           "更改总数：%d",
		&ApplyChangesChangesMasked:          "（出于安全考虑已隐藏）",
		&ApplyChangesChangesEmpty:           "没有需要应用的更改",
		&ApplyChangesHelpTitle:              "应用配置更改",
		&ApplyChangesHelpContent:            `应用更改前，请确认当前配置。`,
		&ApplyChangesIntegrityPromptTitle:   "文件完整性检查",
		&ApplyChangesIntegrityPromptMessage: "检测到过期文件。\n是否要将它们更新到最新版本？",
		&ApplyChangesIntegrityOutdatedList:  "过期文件：\n%s\n确认更新？（y/n）",
		&ApplyChangesIntegrityChecking:      "正在收集文件完整性信息…",
		&ApplyChangesIntegrityNoOutdated:    "未发现过期文件。继续应用更改。",

		&MaintenanceTitle:       "系统维护",
		&MaintenanceDescription: "管理 Peepie 服务并执行维护操作",
		&MaintenanceName:        "维护",
		&MaintenanceOverview: `执行 Peepie 系统维护操作。

可用操作取决于当前系统状态；仅显示适用的操作。

操作包括：
• 服务生命周期管理（启动／停止／重启）
• 组件更新与下载
• 系统重置与清理
• 容器和镜像管理

每项操作都会实时更新状态，并在需要时请求确认。`,
		&MaintenanceStartPentagi:            "启动 Peepie",
		&MaintenanceStartPentagiDesc:        "启动所有已配置的 Peepie 服务",
		&MaintenanceStopPentagi:             "停止 Peepie",
		&MaintenanceStopPentagiDesc:         "停止所有运行中的 Peepie 服务",
		&MaintenanceRestartPentagi:          "重启 Peepie",
		&MaintenanceRestartPentagiDesc:      "重启所有 Peepie 服务",
		&MaintenanceDownloadWorkerImage:     "下载工作节点镜像",
		&MaintenanceDownloadWorkerImageDesc: "下载工作节点任务所需的渗透测试容器镜像",
		&MaintenanceUpdateWorkerImage:       "更新工作节点镜像",
		&MaintenanceUpdateWorkerImageDesc:   "将渗透测试容器镜像更新到最新版本",
		&MaintenanceUpdatePentagi:           "更新 Peepie",
		&MaintenanceUpdatePentagiDesc:       "将 Peepie 更新到最新版本",
		&MaintenanceUpdateInstaller:         "更新安装程序",
		&MaintenanceUpdateInstallerDesc:     "将此安装程序更新到最新版本",
		&MaintenanceFactoryReset:            "恢复出厂设置",
		&MaintenanceFactoryResetDesc:        "将 Peepie 重置为出厂默认设置",
		&MaintenanceRemovePentagi:           "移除 Peepie",
		&MaintenanceRemovePentagiDesc:       "移除 Peepie 容器，但保留数据",
		&MaintenancePurgePentagi:            "彻底清除 Peepie",
		&MaintenancePurgePentagiDesc:        "彻底移除 Peepie，包括所有数据",
		&MaintenanceResetPassword:           "重置管理员密码",
		&MaintenanceResetPasswordDesc:       "重置 Peepie 管理员密码",

		&ResetPasswordFormTitle:       "重置管理员密码",
		&ResetPasswordFormDescription: "重置 Peepie 管理员密码",
		&ResetPasswordFormName:        "重置密码",
		&ResetPasswordFormOverview: `重置默认管理员账户（admin@peepie.com）的密码。

此操作要求 Peepie 正在运行，并会更新 PostgreSQL 数据库中的密码。

输入两次新密码进行确认，然后按 Enter 应用更改。

密码要求：
• 至少 5 个字符
• 两次输入的密码必须一致`,
		&ResetPasswordNewPassword:         "新密码",
		&ResetPasswordNewPasswordDesc:     "输入新的管理员密码",
		&ResetPasswordConfirmPassword:     "确认密码",
		&ResetPasswordConfirmPasswordDesc: "再次输入新密码以确认",
		&ResetPasswordNotAvailable:        "只有 Peepie 正在运行时才能重置密码",
		&ResetPasswordAvailable:           "可以重置密码",
		&ResetPasswordInProgress:          "正在重置密码…",
		&ResetPasswordSuccess:             "已成功重置密码",
		&ResetPasswordErrorPrefix:         "错误： ",
		&ResetPasswordErrorEmptyPassword:  "密码不能为空",
		&ResetPasswordErrorShortPassword:  "密码至少需要 5 个字符",
		&ResetPasswordErrorMismatch:       "两次输入的密码不一致",
		&ResetPasswordHelpContent: `重置用于访问 Peepie 的管理员密码。

此操作：
• 更新 admin@peepie.com 账户的密码
• 将用户状态设为 'active'
• 要求 Peepie 数据库可访问
• 不影响其他用户账户

操作成功后，新密码立即生效。

在两个字段中输入相同的密码，然后按 Enter 确认更改。`,

		&ProcessorOperationFormTitle:                 "%s",
		&ProcessorOperationFormDescription:           "执行 %s 操作",
		&ProcessorOperationFormName:                  "%s",
		&ProcessorOperationNotStarted:                "已准备好执行 %s 操作",
		&ProcessorOperationInProgress:                "正在执行 %s 操作…\n",
		&ProcessorOperationCompleted:                 "%s 操作已成功完成\n",
		&ProcessorOperationFailed:                    "执行 %s 操作失败",
		&ProcessorOperationConfirmation:              "确定要执行 %s 吗？",
		&ProcessorOperationPressEnter:                "按 Enter 执行 %s",
		&ProcessorOperationPressYN:                   "按 Y 确认，按 N 取消",
		&ProcessorOperationRequiresConfirmationShort: "此操作需要确认",
		&ProcessorOperationCancelled:                 "操作已取消",
		&ProcessorOperationUnknown:                   "未知操作：%s",
		&ProcessorOperationStarting:                  "正在启动服务…",
		&ProcessorOperationStopping:                  "正在停止服务…",
		&ProcessorOperationRestarting:                "正在重启服务…",
		&ProcessorOperationDownloading:               "正在下载镜像…",
		&ProcessorOperationUpdating:                  "正在更新组件…",
		&ProcessorOperationResetting:                 "正在恢复出厂设置…",
		&ProcessorOperationRemoving:                  "正在移除容器…",
		&ProcessorOperationPurging:                   "正在清除所有数据…",
		&ProcessorOperationInstalling:                "正在安装 Peepie 服务…",
		&ProcessorOperationHelpTitle:                 "%s 操作",
		&ProcessorOperationHelpContent:               "此操作将执行 %s。",
		&ProcessorOperationHelpContentDownload:       "此操作将下载 %s 组件。",
		&ProcessorOperationHelpContentUpdate:         "此操作将更新 %s 组件。",
		&OperationTitleInstallPentagi:                "安装 Peepie",
		&OperationDescInstallPentagi:                 "安装并配置 Peepie 服务",
		&OperationTitleDownload:                      "下载 %s",
		&OperationDescDownloadComponents:             "下载 %s 组件",
		&OperationTitleUpdate:                        "更新 %s",
		&OperationDescUpdateToLatest:                 "将 %s 更新到最新版本",
		&OperationTitleExecute:                       "执行 %s",
		&OperationDescExecuteOn:                      "执行 %s，目标：%s",
		&OperationProgressExecuting:                  "正在执行 %s…",
		&ProcessorOperationTerminalNotInitialized:    "终端尚未初始化",

		&ProcessorHelpInstallPentagi: `此操作将：
• 为所选服务部署 Docker 容器
• 配置网络和数据卷
• 启动所有已启用的服务
• 根据配置设置监控

安装将使用你当前的配置。`,
		&ProcessorHelpStartPentagi: `此操作将启动：
• Peepie 核心 API 和 Web 界面
• 已配置的 Langfuse 分析服务（如果已启用）
• 可观测性组件（如果已启用）

服务将按正确的依赖顺序启动。`,
		&ProcessorHelpStopPentagi: `此操作将：
• 正常停止容器
• 保留所有数据和配置
• 关闭网络连接

之后可重启服务，不会丢失数据。`,
		&ProcessorHelpRestartPentagi: `此操作将：
• 停止运行中的容器
• 应用配置更改
• 以新状态启动服务

适用于配置更新后，或用于解决问题。`,
		&ProcessorHelpDownloadWorkerImage: `此大型镜像（6GB+）包含：
• Kali Linux 工具和实用程序
• 安全测试框架
• 网络分析软件

渗透测试操作需要此镜像。`,
		&ProcessorHelpUpdateWorkerImage: `此操作将：
• 拉取最新的渗透测试镜像
• 更新安全工具和框架
• 保留现有的工作节点容器

注意：下载量较大（6GB+）。`,
		&ProcessorHelpUpdatePentagi: `此操作将：
• 下载最新的容器镜像
• 对服务执行滚动更新
• 保留所有数据和配置

更新期间服务将短暂不可用。`,
		&ProcessorHelpUpdateInstaller: `此操作将：
• 下载最新的安装程序二进制文件
• 替换当前安装程序
• 退出，以便手动重启

更新后需要重启安装程序。`,
		&ProcessorHelpFactoryReset: `⚠️  警告：此操作将：
• 移除所有容器和网络
• 删除所有配置文件
• 清除存储的数据和数据卷
• 恢复默认设置

此操作无法撤销！`,
		&ProcessorHelpRemovePentagi: `此操作将：
• 停止并移除所有容器
• 移除 Docker 网络
• 保留数据卷和数据
• 保留配置文件

之后可重新安装，不会丢失数据。`,
		&ProcessorHelpPurgePentagi: `⚠️  警告：此操作将永久删除：
• 所有容器和镜像
• 所有数据卷
• 所有配置文件
• 所有存储的结果

此操作无法撤销！`,

		&ProcessorSectionCurrentState:                "当前状态",
		&ProcessorSectionPlanned:                     "计划执行的操作",
		&ProcessorSectionEffects:                     "影响",
		&ProcessorComponentPentagi:                   "Peepie",
		&ProcessorComponentLangfuse:                  "Langfuse",
		&ProcessorComponentObservability:             "可观测性",
		&ProcessorComponentWorkerImage:               "工作节点镜像",
		&ProcessorComponentComposeStacks:             "Compose 服务栈",
		&ProcessorComponentDefaultFiles:              "默认文件",
		&ProcessorItemComposeFiles:                   "Compose 文件",
		&ProcessorItemComposeStacksImagesVolumes:     "Compose 服务栈、镜像、数据卷",
		&ProcessorStateInstalled:                     "已安装",
		&ProcessorStateMissing:                       "未安装",
		&ProcessorStateRunning:                       "运行中",
		&ProcessorStateStopped:                       "已停止",
		&ProcessorStateEmbedded:                      "内置",
		&ProcessorStateExternal:                      "外部",
		&ProcessorStateConnected:                     "已连接",
		&ProcessorStateDisabled:                      "已禁用",
		&ProcessorStateUnknown:                       "未知",
		&PlannedWillStart:                            "将启动：",
		&PlannedWillStop:                             "将停止：",
		&PlannedWillRestart:                          "将重启：",
		&PlannedWillUpdate:                           "将更新：",
		&PlannedWillSkip:                             "将跳过：",
		&PlannedWillRemove:                           "将移除：",
		&PlannedWillPurge:                            "将清除：",
		&PlannedWillDownload:                         "将下载：",
		&PlannedWillRestore:                          "将恢复：",
		&EffectsStart:                                "Peepie Web 界面将可用。后台服务将按所需顺序启动。",
		&EffectsStop:                                 "Web 界面将不可用。正在进行的流程会安全暂停。再次启动 Peepie 后，流程将自动恢复，但当前智能体步骤的一小部分进度可能丢失。",
		&EffectsRestart:                              "服务将停止并以干净状态重新启动。预计会短暂中断，之后流程将自动恢复。",
		&EffectsUpdateAll:                            "拉取镜像，并按需重新创建服务。外部或已禁用的组件将跳过。预计会短暂中断。",
		&EffectsDownloadWorker:                       "不会影响运行中的工作节点容器。新流程将使用下载的镜像。若要让现有流程切换到新镜像，请先结束该流程，然后开始新任务或创建新助手。",
		&EffectsUpdateWorker:                         "拉取最新工作节点镜像。运行中的工作节点容器继续使用旧镜像；新容器将使用更新后的镜像。",
		&EffectsUpdateInstaller:                      "安装程序二进制文件将更新，随后应用将退出。重新启动安装程序以继续。",
		&EffectsFactoryReset:                         "移除容器、数据卷和网络，恢复默认 .env 文件及内置文件，得到全新初始状态。此操作无法撤销。",
		&EffectsRemove:                               "停止并移除容器，但保留数据卷和镜像。数据不会丢失。重新启动之前，Web 界面将不可用。",
		&EffectsPurge:                                "彻底清理：删除容器、镜像、数据卷和配置文件。此操作无法撤销。",
		&EffectsInstall:                              "创建必需的文件并启动服务。检测到外部组件时将跳过。",
		&LLMProviderUnknownName:                      "未知",
		&LLMProviderCustomName:                       "自定义",
		&LLMFormDefaultAuthEnabled:                   "已启用",
		&ServerSettingsExternalSSLInsecureEnabled:    "已启用（⚠ 不安全）",
		&ToolsAIAgentsSettingNotSet:                  "未设置",
		&ToolsSearchEnginesPerplexityModel:           "Perplexity 模型",
		&ToolsSearchEnginesPerplexityModelDesc:       "选择 Perplexity 模型",
		&ToolsSearchEnginesPerplexityContextSize:     "Perplexity 上下文大小",
		&ToolsSearchEnginesPerplexityContextSizeDesc: "选择 Perplexity 上下文大小",
		&ToolsSearchEnginesGoogleSearch:              "Google 搜索",
		&ToolsDockerSummaryCustomNetwork:             "自定义网络",
		&ToolsDockerSummaryPublicIP:                  "公网 IP",
		&ToolsDockerSummaryPentestImage:              "渗透测试镜像",
		&ToolsDockerSummaryTLSConnection:             "TLS 连接",
		&ToolsDockerSummaryRemoteConnection:          "远程连接",
		&ToolsDockerSummaryWorkerDaemon:              "工作节点 Docker 守护进程",
	})
}
