package locale

func init() {
	registerCatalog(LanguageChineseSimplified, map[*string]string{
		// Summarizer Screen
		&SummarizerTitle:       "摘要器配置",
		&SummarizerDescription: "启用对话摘要，以降低 LLM 成本并改进上下文管理",
		&SummarizerName:        "摘要器",
		&SummarizerOverview: `优化上下文使用，降低 LLM 成本，并与模型能力相匹配。

何时需要调整摘要设置：
• Token 成本过高：缩小上下文（4K-12K，而非 22K+ Token）
• 出现“上下文过长”错误：按模型上限进行配置
• 对话连贯性差：提高上下文保留量以保证质量
• 模型类型不同：分别针对短上下文与长上下文模型调优

通用摘要：最大程度控制成本并精确调优，适合研究/分析类任务
助手摘要：智能管理上下文，为交互式会话提供最佳对话质量

快速见效：
• 降低成本：使用通用摘要，将“最近段落数”降至 1-2
• 上下文错误：按模型设置上限（8K/32K/128K）
• 质量优先：使用助手摘要并提高各项上限`,

		&SummarizerTypeGeneralName: "通用摘要",
		&SummarizerTypeGeneralDesc: "用于对话上下文管理的全局摘要设置",

		&SummarizerTypeGeneralInfo: `选择此项可最大程度控制成本，并兼容短上下文模型。

适用场景：
• 大幅降低成本：精细调整每个参数，使 Token 用量最小化
• 短上下文模型（8K-32K）：精确设定上限以避免溢出错误
• 研究/分析类任务：在不丢失关键数据的前提下可控压缩
• 自定义问答处理：完全控制问答对的处理方式

典型效果：
• 相比默认设置降低 40-70% 成本
• 4K-12K Token 上下文（助手模式为 22K+）
• 在 GPT-3.5、Claude Instant 及较小模型上表现更好
• 精确平衡对话记忆与新鲜上下文

最佳实践：
• 从 1-2 个“最近段落”开始，以最大程度节省成本
• 启用“大小管理”以自动防止溢出
• 仅在关键推理任务中禁用问答压缩`,

		&SummarizerTypeAssistantName: "助手摘要",
		&SummarizerTypeAssistantDesc: "针对 AI 助手上下文的专用摘要设置",

		&SummarizerTypeAssistantInfo: `选择此项可获得最佳对话质量和对话连贯性。

适用场景：
• 较长的推理链：为复杂的多步思考保留上下文
• 高质量对话：保持对话连贯性和助手风格
• 长上下文模型（64K+）：高效发挥模型全部能力
• 交互式会话：更好地记住用户偏好和对话历史

典型效果：
• 8K-40K Token 上下文，并进行智能压缩
• 对话连贯性优于手动设置
• 针对推理任务自动优化上下文
• 兼顾成本与质量（上下文约为通用模式的 3 倍）

最佳实践：
• 大多数场景使用默认设置即可——它们已预先优化
• 仅在非常复杂的任务中增加“最近段落数”
• 关注上下文用量——成本随 Token 数量增长
• 非常适合 GPT-4、Claude 及其他大上下文模型`,

		// Summarizer Form Screen
		&SummarizerFormGeneralTitle:   "通用摘要器配置",
		&SummarizerFormAssistantTitle: "助手摘要器配置",
		&SummarizerFormDescription:    "配置%s设置",

		&SummarizerFormPreserveLast:     "大小管理",
		&SummarizerFormPreserveLastDesc: "控制最后一个段落的压缩。启用：段落大小限制在 LastSecBytes 以内（上下文更小）。禁用：段落可自由增长（上下文更大）",

		&SummarizerFormUseQA:     "问答摘要",
		&SummarizerFormUseQADesc: "当问答内容总量超过 MaxQABytes 或 MaxQASections 上限时，启用问答对压缩",

		&SummarizerFormSumHumanInQA:     "压缩用户消息",
		&SummarizerFormSumHumanInQADesc: "在问答压缩中包含用户消息。禁用：保留用户原文（大多数情况下推荐）",

		&SummarizerFormLastSecBytes:     "段落大小上限",
		&SummarizerFormLastSecBytesDesc: "启用“大小管理”时每个最近段落的最大字节数。越大：每个段落细节越多，Token 用量越高",

		&SummarizerFormMaxBPBytes:     "响应大小上限",
		&SummarizerFormMaxBPBytesDesc: "单条 AI 响应在被压缩前的最大字节数。防止单条大型响应占据上下文",

		&SummarizerFormMaxQASections:     "问答段落上限",
		&SummarizerFormMaxQASectionsDesc: "触发问答压缩前的最大问答段落数。与 MaxQABytes 配合控制问答记忆总量",

		&SummarizerFormMaxQABytes:     "问答记忆总量",
		&SummarizerFormMaxQABytesDesc: "所有问答段落合计的最大字节数。超出时（同时超出 MaxQASections）触发问答压缩以满足上限",

		&SummarizerFormKeepQASections:     "最近段落数",
		&SummarizerFormKeepQASectionsDesc: "不经压缩直接保留的最近对话段落数量。影响上下文大小的首要参数",

		&SummarizerFormGeneralHelp: `上下文估算：4K-22K Token（典型值），最高 94K（最大设置）。

关键关系：
• 最近段落数：最关键——每 +1 约增加 1.5-9K Token
• 关闭大小管理：上下文增大 2-3 倍（压缩更少）
• 段落/响应上限：控制各组成部分的大小
• 问答记忆：超出上限时管理整体对话历史

参数相互作用：
• 仅当 MaxQABytes 与 MaxQASections 同时超出时才触发问答压缩
• 禁用大小管理 → 段落可增长到上限的 2 倍
• 响应上限可防止单条大型输出占据上下文
• 用户消息压缩（SummHumanInQA）可节省 5%，但会丢失原始措辞

针对较小模型可调低：
• 最近段落数：1-2（默认为 3+）
• 段落上限：25-35KB（而非 50KB+）
• 简单对话可禁用大小管理

常见错误：
• 最近段落数设置过高（上下文溢出的主要原因）
• 启用大小管理的同时将段落上限设得很低（过度压缩）
• 问答上限不匹配（字节数高 + 段落数低 = 无效）

当前算法会压缩较早的内容，同时保证最近上下文的质量。`,

		&SummarizerFormAssistantHelp: `针对需要上下文连贯性的交互式对话进行了优化。

默认调优（3 个最近段落，75KB 上限）：
• 典型范围：8K-40K Token
• 适用于：长对话、推理链、依赖上下文的任务
• 模型：适合 32K+ 上下文模型

按模型类型调整：
• 短上下文（≤16K）：最近段落数=1-2，段落上限=45KB
• 长上下文（128K+）：可将最近段落数提高到 5-7
• 高频聊天：将最近段落数降至 2 以加快响应

高级调优：
• 文档分析类对话可将问答记忆设为 200KB+
• 详细技术回答可将响应上限设为 24-32KB
• 保持用户消息不压缩（SummHumanInQA=false）以获得更好的上下文

性能优化：
• 助手模式下每个最近段落 ≈ 9-18KB
• 大小管理可将增长减少约 20%，但可能丢失细节
• 由于默认上限较大，问答压缩触发得较少

默认启用大小管理——在保持对话连贯的同时防止上下文溢出。
请关注实际 Token 用量，先调整最近段落数，再调整各项上限。`,

		&SummarizerContextEstimatedSize:    "预计上下文大小：%s\n%s",
		&SummarizerContextTokenRange:       "约 %s Token",
		&SummarizerContextTokenRangeMinMax: "约 %s-%s Token",
		&SummarizerContextRequires256K:     "需要 256K+ 上下文模型",
		&SummarizerContextRequires128K:     "需要 128K+ 上下文模型",
		&SummarizerContextRequires64K:      "需要 64K+ 上下文模型",
		&SummarizerContextRequires32K:      "需要 32K+ 上下文模型",
		&SummarizerContextRequires16K:      "需要 16K+ 上下文模型",
		&SummarizerContextFitsIn8K:         "适用于 8K+ 上下文模型",

		// Tools screen strings
		&ToolsTitle:       "工具配置",
		&ToolsDescription: "通过附加工具和选项增强智能体能力",
		&ToolsName:        "工具",
		&ToolsOverview: `为 AI 智能体配置附加工具和能力。
每个工具都可按需启用和配置。

可用设置：
• 人机协同 - 在测试过程中启用用户交互
• AI 智能体设置 - 配置 AI 智能体的全局行为
• 搜索引擎 - 配置外部搜索提供商
• 爬虫 - 网页内容抽取与分析
• Graphiti (beta) - 用于语义记忆的时序知识图谱
• Docker - 容器环境配置`,
	})
}
