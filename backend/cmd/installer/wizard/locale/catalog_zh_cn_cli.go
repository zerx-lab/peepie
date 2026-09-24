package locale

func init() {
	registerCatalog(LanguageChineseSimplified, map[*string]string{
		&CLIFlagVersion:         "显示版本信息",
		&CLIFlagEnvFile:         "环境变量文件路径",
		&CLIFlagLanguage:        "界面语言（%s）；默认：通过 Ctrl+L 选择的语言，其次为系统语言",
		&CLIUsageTitle:          "Peepie 安装程序 v%s",
		&CLIUsageLine:           "用法：%s [选项]",
		&CLIUsageOptions:        "选项：",
		&CLIUsageExamples:       "示例：",
		&CLIExampleDefault:      "使用默认 .env 文件",
		&CLIExampleEnvFile:      "使用自定义环境变量文件",
		&CLIExampleLanguage:     "使用简体中文界面",
		&CLIExampleVersion:      "显示版本",
		&CLIUnsupportedLanguage: "不支持的语言 %q（支持：%s）",
		&CLIError:               "错误：%v",
		&CLIFailedInitState:     "初始化状态失败：%v",
		&CLIFailedMigrate:       "迁移设置失败：%v",
		&CLIFailedSyncNetwork:   "同步网络设置失败：%v",
		&CLIFailedGatherFacts:   "收集系统信息失败：%v",
		&CLIFailedHardening:     "执行安全加固失败：%v",
		&CLIApplicationError:    "应用程序错误：%v",
		&CLIStartupTitle:        "Peepie 安装程序 v%s",
		&CLIStartupEnvFile:      "环境变量文件：%s",
		&CLISystemNotReady:      "⚠️  系统尚未就绪，无法继续。请先解决上述问题。",
		&CLISystemReady:         "✅ 系统已就绪，可以继续。",
		&CLIPendingChanges:      "你有未应用的更改。",
		&CLIPendingChangesHint:  "请再次运行安装程序以继续，或提交你的更改。",
	})
}
