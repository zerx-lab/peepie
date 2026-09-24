package locale

func init() {
	registerCatalog(LanguageChineseSimplified, map[*string]string{
		// LLM Providers screen
		&LLMProvidersTitle:       "LLM 提供商配置",
		&LLMProvidersDescription: "为 AI 智能体配置大语言模型提供商",
		&LLMProvidersName:        "LLM 提供商",
		&LLMProvidersOverview: `Peepie 使用专门的 AI 智能体（研究员、开发者、执行者、渗透测试员），它们需要不同的 LLM 能力才能获得最佳的渗透测试效果。

为什么需要多个提供商：
• 智能体专业化：不同智能体受益于针对推理、编码或分析优化的模型
• 成本效率：复杂任务使用昂贵的推理模型（o3、grok-4、claude-sonnet-4、gemini-2.5-pro），简单操作使用更便宜的模型
• 性能优化：各提供商擅长的领域不同 - OpenAI 适合中等任务，Anthropic 适合复杂任务，Gemini 适合节省成本

提供商选择指南：
• 云端生产环境：OpenAI + Anthropic + Gemini，性能和可靠性业界领先
• 企业/合规：AWS Bedrock，满足 SOC2、HIPAA 要求，并可使用多个模型系列
• 隐私/本地部署：Ollama 或 vLLM 搭配 Llama 3.1、Qwen3 或其他开放模型，完全掌控数据

容器内的 /opt/pentagi/conf/ 目录提供了 OpenRouter、DeepInfra、vLLM、Ollama 等提供商的开箱即用配置`,

		// Provider titles and descriptions
		&LLMProviderOpenAI:        "OpenAI",
		&LLMProviderAnthropic:     "Anthropic",
		&LLMProviderGemini:        "Google Gemini",
		&LLMProviderBedrock:       "AWS Bedrock",
		&LLMProviderOllama:        "Ollama",
		&LLMProviderDeepSeek:      "DeepSeek",
		&LLMProviderGLM:           "GLM 智谱 AI",
		&LLMProviderKimi:          "Kimi 月之暗面",
		&LLMProviderQwen:          "Qwen 阿里云通义千问",
		&LLMProviderMiniMax:       "MiniMax",
		&LLMProviderCustom:        "自定义",
		&LLMProviderOpenAIDesc:    "业界领先的 GPT 模型，综合性能出色",
		&LLMProviderAnthropicDesc: "Claude 模型，推理能力和安全特性出众",
		&LLMProviderGeminiDesc:    "Google 先进的多模态模型，知识面广",
		&LLMProviderBedrockDesc:   "面向企业的 AWS 服务，可访问多家基础模型提供商",
		&LLMProviderOllamaDesc:    "本地和云端开源模型，兼顾隐私与灵活性",
		&LLMProviderDeepSeekDesc:  "先进的中国 AI 模型，推理和多语言能力强",
		&LLMProviderGLMDesc:       "智谱 AI 的 GLM 模型，适用于中英文任务",
		&LLMProviderKimiDesc:      "月之暗面的长上下文模型，适用于文档分析",
		&LLMProviderQwenDesc:      "阿里云的通义千问模型，适用于多语言任务",
		&LLMProviderMiniMaxDesc:   "MiniMax 的 M 系列模型，适用于智能体推理和长上下文任务",
		&LLMProviderCustomDesc:    "自定义 OpenAI 兼容端点，灵活性最高",

		// Provider-specific help text
		&LLMFormOpenAIHelp: `OpenAI 提供业界领先的模型，具备前沿的推理能力，非常适合复杂的渗透测试。

Peepie 默认模型：
• o1、o4-mini：高级推理模型，用于复杂漏洞分析和策略规划
• GPT-4.1、GPT-4.1-mini：旗舰模型，针对漏洞利用开发和代码生成优化
• 根据智能体类型和任务复杂度自动选择模型

主要优势：
• 最先进的推理能力，支持逐步分析（o 系列模型）
• 出色的编码能力，可用于定制漏洞利用开发和载荷生成
• 性能可靠，可用性稳定，API 文档详尽
• 在安全研究和渗透测试场景中久经验证

适用于：需要前沿 AI 能力的生产环境、重视性能胜过成本的团队
成本：高端定价，但优化后的配置可兼顾成本与质量

设置：在 https://platform.openai.com/api-keys 获取 API 密钥`,

		&LLMFormAnthropicHelp: `Anthropic Claude 模型推理和分析能力出众，擅长注重安全性的渗透测试。

Peepie 默认模型：
• Claude Sonnet-4：高端推理模型，用于复杂安全分析和策略性漏洞评估
• Claude 3.5 Haiku：高速模型，针对快速信息收集和简单解析任务优化
• 在所有安全测试场景中实现均衡的性价比

主要优势：
• 高度关注安全与伦理 - 在保持安全测试效果的同时减少有害输出
• 推理能力出众，适合有条理的漏洞分析和系统化的渗透测试方法
• 上下文窗口大，适合分析大型代码库和配置文件
• 擅长理解复杂的安全场景和法规合规要求

适用于：重视负责任测试实践的安全团队、注重合规的环境、详细分析
成本：中等定价，对推理密集型安全工作流来说性价比极高

设置：在 https://console.anthropic.com/ 获取 API 密钥`,

		&LLMFormGeminiHelp: `Google Gemini 结合了多模态能力和高级推理，非常适合全面的安全评估。

Peepie 默认模型：
• Gemini 2.5 Pro：高级推理模型，用于深度漏洞分析和复杂漏洞利用开发
• Gemini 2.5 Flash：高性能模型，兼顾速度与智能，适合大多数安全测试任务
• Gemini 2.0 Flash Lite：高性价比模型，适合快速扫描和信息收集操作
• 具备逐步分析的推理能力，可进行全面的渗透测试

主要优势：
• 多模态支持，可分析截图、网络拓扑图和安全文档
• 价格有竞争力，速率限制宽松，适合开发和测试环境
• 上下文窗口大（最高 2M Token），可分析超大型代码库和系统配置
• 在多种编程语言的代码分析和漏洞识别方面表现出色

适用于：预算有限的团队、开发环境、需要图像/文档分析的场景
成本：主流云提供商中最具性价比的选择，性能价格比极佳

设置：在 https://aistudio.google.com/app/apikey 获取 API 密钥`,

		&LLMFormBedrockHelp: `AWS Bedrock 提供对 20 多个基础模型的企业级访问，支持多种认证方式并具备增强的安全性。

Peepie 默认模型：
• Claude Sonnet-4.5（通过 Bedrock）：高端推理模型，具备 AWS 企业级安全性和扩展思考能力
• OpenAI GPT OSS 120B：强大的推理模型，适合科学分析和复杂安全任务
• Claude Haiku-4.5、DeepSeek V3.2、Qwen3-32B：高效模型，适用于特定智能体角色和成本优化
• 通过统一接口访问 Amazon Nova（多模态）、Mistral、Moonshot 等模型

认证方式（按优先级）：
1. 默认 AWS 认证（BEDROCK_DEFAULT_AUTH=true）：使用 AWS SDK 凭证链 - 推荐用于 EC2/ECS/Lambda
2. Bearer 令牌（BEDROCK_BEARER_TOKEN）：基于令牌的认证，适用于自定义认证场景
3. 静态凭证（ACCESS_KEY + SECRET_KEY）：传统 IAM 凭证，适用于开发和测试

主要优势：
• 企业合规：SOC2、HIPAA、FedRAMP 认证，具备数据驻留和治理控制
• 多提供商访问：来自 Anthropic、Amazon、OpenAI、Qwen、DeepSeek、Cohere、Mistral、Moonshot 的 20 多个模型
• 灵活认证：三种方式，适应不同的部署场景和安全要求
• 增强安全性：VPC 集成、CloudTrail 日志、IAM 控制、私有端点，实现完全隔离
• 区域部署：部署在首选 AWS 区域，优化延迟并满足数据主权要求

适用于：企业环境、受监管行业、需要合规控制和灵活认证的团队
成本：价格有竞争力，提供预置吞吐量选项，但新账户的速率限制较严格（2-20 次请求/分钟）
重要：对于生产环境的渗透测试工作流，请通过 AWS Service Quotas 控制台申请提高配额

设置：选择认证方式并配置凭证。在 https://docs.aws.amazon.com/bedrock/ 查看速率限制`,

		&LLMFormOllamaHelp: `Ollama 支持两种部署场景，灵活性极高。

场景 1：本地 Ollama 服务器（自托管）
• 在自己的硬件上运行 Ollama（建议 8GB+ 内存，GPU 可选但有帮助）
• 完全的数据隐私 - 所有处理均在本地完成
• 无持续费用 - 仅需基础设施
• 无需 API 密钥 - 通过网络访问控制认证
• 设置：从 https://ollama.ai/ 安装，并配置 OLLAMA_SERVER_URL=http://ollama-server:11434

场景 2：Ollama Cloud（托管服务）
• 云端托管模型，无需本地基础设施
• 无需硬件 - 模型在 Ollama 的基础设施上运行
• 按用量付费，并提供免费套餐
• 需要 API 密钥 - 在 https://ollama.com/settings/keys 生成
• 设置：在 https://ollama.com 注册，配置 OLLAMA_SERVER_URL=https://ollama.com + OLLAMA_SERVER_API_KEY=your_key

Peepie 默认模型：
• Llama 3.1:8b、Qwen3:32b 及其他开放模型
• 可自定义 - 在 100 多个可用模型之间切换
• 提供模型自动下载和加载选项，使用方便

主要优势：
• 两种部署选项：在隐私（本地）和便利（云端）之间选择
• 成本灵活：本地无持续费用，云端按用量付费
• 丰富的模型库：可使用最新的开源模型（Llama、Qwen、Mistral、Gemma 等）
• 支持离线隔离环境：本地部署可在隔离网络中运行

适用于：注重隐私的团队（本地）、预算有限的部署（云端）、有数据主权要求的组织
设置选项：从 https://10.10.10.10:11434 本地安装，或在 https://ollama.com 注册云服务`,

		&LLMFormDeepSeekHelp: `DeepSeek 提供推理能力强、支持多语言的先进 AI 模型。

Peepie 默认模型：
• deepseek-v4-flash：高性价比通用模型，适用于对话、代码生成和工具调用
• deepseek-v4-pro：更高阶的推理模型，适用于复杂逻辑、数学推理和安全分析
• 价格实惠，性能可与领先模型媲美

主要优势：
• 强大的编码和推理能力，适合安全分析和漏洞利用开发
• 多语言支持（中文和英文），适合国际化渗透测试场景
• 价格有竞争力，性价比出色
• OpenAI 兼容 API，集成无缝

LiteLLM 集成：
• 使用 LiteLLM 代理时，将提供商名称设置为 'deepseek'
• 启用模型前缀（例如 deepseek/deepseek-v4-flash），无需修改 config.yml
• 直接使用 DeepSeek API 时可选

适用于：需要多语言支持的团队、注重成本的部署、中文安全测试
成本：价格极具竞争力，性能表现强劲

设置：在 https://platform.deepseek.com/ 获取 API 密钥`,

		&LLMFormGLMHelp: `智谱 AI 的 GLM 提供先进的语言模型，具备强大的 NLP 和推理能力，由清华大学研发。

Peepie 默认模型：
• GLM-4-Air：高性能通用对话模型，针对常规任务和工具调用优化
• GLM-4-Plus：旗舰模型，具备强大的推理和代码生成能力
• GLM-Z1-Plus：高级推理模型，具备深度分析能力，适合安全研究

主要优势：
• 出色的中英文 NLP 能力
• 在多语言安全测试和分析场景中表现强劲
• GLM-4 和 GLM-Z1 模型系列，推理和编码能力增强
• OpenAI 兼容 API，易于集成

备用 API 端点：
• 国际：https://api.z.ai/api/paas/v4（默认）
• 中国：https://open.bigmodel.cn/api/paas/v4
• 编码专用：https://api.z.ai/api/coding/paas/v4

LiteLLM 集成：
• 使用 LiteLLM 代理时，将提供商名称设置为 'zai'
• 启用模型前缀（例如 zai/glm-4），无需修改 config.yml
• 直接使用 GLM API 时可选

适用于：中英文多语言渗透测试、在亚洲市场运营的团队
成本：价格有竞争力，多语言任务性能良好

设置：在 https://open.bigmodel.cn/ 获取 API 密钥`,

		&LLMFormKimiHelp: `月之暗面的 Kimi 提供超长上下文模型，非常适合分析大型代码库和文档。

Peepie 默认模型：
• Moonshot-v1-8k：长上下文模型，最多支持 8K Token，适用于通用对话
• Kimi-k2.5：高级模型，具备强大的推理和文档理解能力
• 针对处理大量文本和代码进行了优化

主要优势：
• 超长上下文窗口（最高 1M Token），可全面分析代码库
• 强大的中英文支持，适合多语言渗透测试
• 对于文档密集型安全评估和威胁情报分析性价比高
• 擅长理解复杂的系统架构和长篇技术文档

备用 API 端点：
• 国际：https://api.moonshot.ai/v1（默认）
• 中国：https://api.moonshot.cn/v1

LiteLLM 集成：
• 使用 LiteLLM 代理时，将提供商名称设置为 'moonshot'
• 启用模型前缀（例如 moonshot/kimi-k2.5），无需修改 config.yml
• 直接使用 Kimi API 时可选

适用于：大型代码库分析、文档密集型评估、安全研究中需要长上下文的团队
成本：价格有竞争力，长上下文场景性价比极高

设置：在 https://platform.moonshot.ai/ 获取 API 密钥`,

		&LLMFormQwenHelp: `阿里云百炼（DashScope）的通义千问提供强大的多语言模型，并具备多模态能力。

Peepie 默认模型：
• Qwen-Turbo：速度最快的轻量模型，适用于高频任务和实时响应场景
• Qwen-Plus：性能均衡的模型，适用于通用对话、代码生成和工具调用
• Qwen-Max：旗舰推理模型，指令遵循能力强，擅长处理复杂任务
• QwQ-Plus：深度推理模型，具备扩展思维链，适合复杂逻辑分析

主要优势：
• 强大的多语言支持（中文、英文及多种其他语言）
• 通过 Qwen-VL 提供多模态能力，可进行视觉安全分析
• 与阿里云集成，适合企业部署
• DashScope 生态提供更多 AI 服务和工具
• Qwen2.5、Qwen3 和 QwQ 模型系列，提供多种规模和专长

备用 API 端点：
• 美国：https://dashscope-us.aliyuncs.com/compatible-mode/v1（默认）
• 新加坡：https://dashscope-intl.aliyuncs.com/compatible-mode/v1
• 中国：https://dashscope.aliyuncs.com/compatible-mode/v1

LiteLLM 集成：
• 使用 LiteLLM 代理时，将提供商名称设置为 'dashscope'
• 启用模型前缀（例如 dashscope/qwen-plus），无需修改 config.yml
• 直接使用 Qwen API 时可选

适用于：在亚洲市场运营的团队、多语言安全测试、使用 Qwen-VL 进行视觉分析、阿里云生态集成
成本：价格有竞争力，针对不同用例提供灵活的档位

设置：在 https://dashscope.console.aliyun.com/ 获取 API 密钥`,

		&LLMFormMiniMaxHelp: `MiniMax 提供具备超大上下文窗口的 M 系列智能体模型，可通过 OpenAI 兼容 API 访问。

Peepie 默认模型：
• MiniMax-M3：最新旗舰模型（约 1M Token 上下文），适用于智能体推理、工具使用、编码和长上下文任务
• MiniMax-M2.7：上一代模型，具备强大的推理和编码能力
• MiniMax-M2.7-highspeed：M2.7 的低延迟版本，适用于快速响应场景

主要优势：
• M3 上下文窗口超大（约 1M Token），适合长代码库和报告
• 强大的智能体工具调用、JSON 输出和流式输出支持
• 提供 OpenAI 和 Anthropic 兼容端点，可直接替换集成

API 端点：
• 全球：https://api.minimax.io/v1（默认）

LiteLLM 集成：
• 使用 LiteLLM 代理时，将提供商名称设置为 'minimax'
• 启用模型前缀（例如 minimax/MiniMax-M3），无需修改 config.yml
• 直接使用 MiniMax API 时可选

适用于：受益于超大上下文和可靠工具调用的智能体安全工作流
成本：M 系列按 Token 计费，价格有竞争力

设置：在 https://platform.minimax.io/ 获取 API 密钥`,

		&LLMFormCustomHelp: `配置任意 OpenAI 兼容 API 端点，获得最大的灵活性，并与现有基础设施集成。

开箱即用的配置：
• vLLM 部署：高吞吐量的本地推理，GPU 利用率最优
• OpenRouter：通过单一 API 访问多家提供商的 200 多个模型，价格有竞争力
• DeepInfra：热门开放模型的无服务器推理，按用量付费
• Together AI、Groq、Fireworks：其他云提供商，各有专门的性能优化
• LiteLLM 代理：连接 100 多个提供商的通用网关，支持负载均衡和统一接口（使用 LLM_SERVER_PROVIDER 设置模型前缀）
• 部分推理模型和 LLM 提供商在使用工具调用时可能需要保留推理内容（LLM_SERVER_PRESERVE_REASONING=true）

常用本地部署选项：
• vLLM：生产级推理服务，支持 Qwen、Llama、Mistral 模型，具备批处理和 GPU 优化
• LocalAI：OpenAI 兼容的 API 封装，支持多种本地模型和嵌入服务
• Text Generation WebUI：社区热门界面，模型支持广泛并具备微调能力
• Hugging Face TGI：企业级文本生成推理，支持自动扩缩容和监控

主要优势：
• 无限灵活：可使用任意 OpenAI 兼容端点或服务
• 成本优化：选择价格有竞争力的提供商，或在自己的基础设施上部署模型
• 供应商独立：避免锁定，可在提供商和模型之间无缝切换
• 定制微调：部署基于你的安全测试场景训练的专用模型

适用于：有特定模型需求、成本优化需求或已有 LLM 基础设施的团队
LiteLLM 集成：将 LLM_SERVER_PROVIDER 设置为你的提供商名称（例如 "openrouter"、"moonshot"），即可在直接 API 访问和 LiteLLM 代理中使用相同的配置文件
可用示例：容器内的 /opt/pentagi/conf/ 目录提供了主流提供商的预配置`,

		// LLM Provider Form field labels and descriptions
		&LLMFormFieldBaseURL:           "基础 URL",
		&LLMFormFieldAPIKey:            "API 密钥",
		&LLMFormFieldDefaultAuth:       "使用默认 AWS 认证",
		&LLMFormFieldBearerToken:       "Bearer 令牌",
		&LLMFormFieldAccessKey:         "访问密钥 ID",
		&LLMFormFieldSecretKey:         "秘密访问密钥",
		&LLMFormFieldSessionToken:      "会话令牌",
		&LLMFormFieldRegion:            "区域",
		&LLMFormFieldModel:             "模型",
		&LLMFormFieldConfigPath:        "配置路径",
		&LLMFormFieldLegacyReasoning:   "旧版推理",
		&LLMFormFieldPreserveReasoning: "保留推理内容",
		&LLMFormFieldProviderName:      "提供商名称",
		&LLMFormFieldPullTimeout:       "模型拉取超时",
		&LLMFormFieldPullEnabled:       "自动拉取模型",
		&LLMFormFieldLoadModelsEnabled: "从服务器加载模型",
		&LLMFormBaseURLDesc:            "提供商的 API 端点 URL",
		&LLMFormAPIKeyDesc:             "用于认证的 API 密钥",
		&LLMFormDefaultAuthDesc:        "使用 AWS SDK 默认凭证链（环境变量、EC2 角色、~/.aws/credentials）- 优先级最高",
		&LLMFormBearerTokenDesc:        "用于认证的 Bearer 令牌 - 优先于静态凭证",
		&LLMFormAccessKeyDesc:          "用于静态凭证认证的 AWS 访问密钥 ID",
		&LLMFormSecretKeyDesc:          "用于静态凭证认证的 AWS 秘密访问密钥",
		&LLMFormSessionTokenDesc:       "用于临时凭证的 AWS 会话令牌（可选，与静态凭证一起使用）",
		&LLMFormRegionDesc:             "Bedrock 服务所在的 AWS 区域",
		&LLMFormModelDesc:              "此提供商使用的默认模型",
		&LLMFormConfigPathDesc:         "配置文件路径（可选）",
		&LLMFormLegacyReasoningDesc:    "启用旧版推理模式（true/false）",
		&LLMFormPreserveReasoningDesc:  "在多轮对话中保留推理内容（部分提供商需要）",
		&LLMFormProviderNameDesc:       "模型名称的提供商前缀（适用于 LiteLLM 代理）",
		&LLMFormPullTimeoutDesc:        "下载模型的超时时间，单位为秒（默认：600）",
		&LLMFormPullEnabledDesc:        "启动时自动下载所需模型",
		&LLMFormLoadModelsEnabledDesc:  "从 Ollama 服务器加载可用模型列表",
		&LLMFormOllamaAPIKeyDesc:       "Ollama Cloud API 密钥（可选，使用本地 Ollama 服务器时留空）",

		// LLM Provider Form status messages
		&LLMProviderFormTitle:       "LLM 提供商 %s 配置",
		&LLMProviderFormDescription: "配置大语言模型提供商设置",
		&LLMProviderFormName:        "LLM 提供商 %s",
		&LLMProviderFormOverview: `智能体角色分配：
• 主智能体和渗透测试员：使用推理模型（o3、grok-4、claude-sonnet-4、gemini-2.5-pro）进行复杂漏洞分析
• 助手和顾问：使用高级模型（o4-mini、claude-sonnet-4）进行策略规划和提供建议
• 编码员和安装员：使用精确模型（gpt-4.1、claude-sonnet-4）进行漏洞利用开发和系统配置
• 搜索员和信息增强员：使用快速模型（gpt-4.1-mini、claude-3.5-haiku、gemini-2.0-flash-lite）收集信息
• 简单任务：使用轻量模型进行 JSON 解析和基本操作

性能考量：
• 推理模型提供逐步分析，但速度较慢、成本较高
• 标准模型响应更快，适合高频的智能体交互
• 每种智能体类型都使用针对安全测试工作流优化的提供商专属模型配置

你的配置将决定每个智能体在不同渗透测试场景中使用哪些模型。`,
	})
}
