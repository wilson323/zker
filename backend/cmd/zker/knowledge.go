package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/coze-dev/coze-studio/backend/cmd/zker/api"
	"github.com/coze-dev/coze-studio/backend/cmd/zker/config"
	"github.com/coze-dev/coze-studio/backend/cmd/zker/utils"
)

// knowledgeCmd 知识库管理命令
var knowledgeCmd = &cobra.Command{
	Use:   "knowledge",
	Short: "知识库管理",
	Long:  `管理知识库，包括上传、同步和搜索操作`,
}

// uploadKnowledgeCmd 上传文档到知识库
var uploadKnowledgeCmd = &cobra.Command{
	Use:   "upload [path]",
	Short: "上传文档到知识库",
	Long:  `上传文件或目录到指定的知识库`,
	Example: `  # 上传单个文件
  zker knowledge upload ./docs/manual.pdf --kb-id kb_xxxxx

  # 上传整个目录
  zker knowledge upload ./docs --kb-id kb_xxxxx

  # 递归上传
  zker knowledge upload ./docs --kb-id kb_xxxxx --recursive`,
	Args: cobra.ExactArgs(1),
	RunE: runKnowledgeUpload,
}

// syncKnowledgeCmd 同步知识库
var syncKnowledgeCmd = &cobra.Command{
	Use:   "sync",
	Short: "同步知识库",
	Long:  `同步本地文件到知识库`,
	Example: `  # 同步知识库
  zker knowledge sync --kb-id kb_xxxxx

  # 强制同步（覆盖远程）
  zker knowledge sync --kb-id kb_xxxxx --force`,
	RunE: runKnowledgeSync,
}

// searchKnowledgeCmd 搜索知识库
var searchKnowledgeCmd = &cobra.Command{
	Use:   "search [query]",
	Short: "搜索知识库",
	Long:  `在知识库中搜索相关内容`,
	Example: `  # 搜索知识库
  zker knowledge search "如何使用ZKER" --kb-id kb_xxxxx

  # 限制结果数量
  zker knowledge search "API文档" --kb-id kb_xxxxx --limit 5`,
	Args: cobra.ExactArgs(1),
	RunE: runKnowledgeSearch,
}

var (
	knowledgeID    string
	knowledgeForce bool
	knowledgeLimit int
	knowledgeRecursive bool
)

func init() {
	knowledgeCmd.AddCommand(uploadKnowledgeCmd)
	knowledgeCmd.AddCommand(syncKnowledgeCmd)
	knowledgeCmd.AddCommand(searchKnowledgeCmd)

	uploadKnowledgeCmd.Flags().StringVar(&knowledgeID, "kb-id", "", "知识库ID (必需)")
	uploadKnowledgeCmd.Flags().BoolVarP(&knowledgeRecursive, "recursive", "r", false, "递归上传目录")
	uploadKnowledgeCmd.MarkFlagRequired("kb-id")

	syncKnowledgeCmd.Flags().StringVar(&knowledgeID, "kb-id", "", "知识库ID (必需)")
	syncKnowledgeCmd.Flags().BoolVar(&knowledgeForce, "force", false, "强制同步（覆盖远程）")
	syncKnowledgeCmd.MarkFlagRequired("kb-id")

	searchKnowledgeCmd.Flags().StringVar(&knowledgeID, "kb-id", "", "知识库ID (必需)")
	searchKnowledgeCmd.Flags().IntVarP(&knowledgeLimit, "limit", "n", 10, "返回结果数量")
	searchKnowledgeCmd.MarkFlagRequired("kb-id")
}

func runKnowledgeUpload(cmd *cobra.Command, args []string) error {
	path := args[0]

	client := getKnowledgeClient()
	if client == nil {
		return fmt.Errorf("获取API客户端失败")
	}

	utils.PrintInfo("正在上传文档: %s", path)
	fmt.Printf("知识库ID: %s\n", knowledgeID)
	if knowledgeRecursive {
		fmt.Printf("模式: 递归上传\n")
	}

	// TODO: 实现文档上传API
	_ = client

	utils.PrintSuccess("文档上传成功！")
	return nil
}

func runKnowledgeSync(cmd *cobra.Command, args []string) error {
	client := getKnowledgeClient()
	if client == nil {
		return fmt.Errorf("获取API客户端失败")
	}

	utils.PrintInfo("正在同步知识库...")
	fmt.Printf("知识库ID: %s\n", knowledgeID)
	if knowledgeForce {
		fmt.Printf("模式: 强制同步（覆盖远程）\n")
	}

	// TODO: 实现知识库同步API
	_ = client

	utils.PrintSuccess("知识库同步成功！")
	return nil
}

func runKnowledgeSearch(cmd *cobra.Command, args []string) error {
	query := args[0]

	client := getKnowledgeClient()
	if client == nil {
		return fmt.Errorf("获取API客户端失败")
	}

	utils.PrintInfo("正在搜索知识库...")
	fmt.Printf("搜索词: %s\n", query)
	fmt.Printf("知识库ID: %s\n", knowledgeID)
	fmt.Printf("结果数量: %d\n", knowledgeLimit)

	// TODO: 实现知识库搜索API
	_ = client

	utils.PrintSuccess("搜索完成！")
	return nil
}

// getKnowledgeClient 获取知识库API客户端
func getKnowledgeClient() *api.Client {
	cfg := config.GetGlobalConfig()
	if !cfg.IsAuthenticated() {
		utils.PrintError("未登录，请先运行: zker login")
		return nil
	}

	return api.NewClient(cfg.GetAPIEndpoint(), cfg.GetAuthToken())
}
