package main

import (
	"fmt"
	"os"

	"github.com/sfuruya0612/snatch/cmd"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

func main() {
	outputDir := "./docs"

	root := cmd.GetRootCmd()

	updateCommandDocs(root)

	err := doc.GenMarkdownTree(root, outputDir)
	if err != nil {
		fmt.Printf("Error generating documentation: %v\n", err)
		os.Exit(1)
	}
}

// updateCommandDocs はコマンドのドキュメントを更新します
func updateCommandDocs(root *cobra.Command) {
	root.Example = `  # EC2インスタンスの一覧表示（タブ区切り）
  snatch ec2
  snatch ec2 ls

  # EC2インスタンスの一覧表示（CSV形式）
  snatch ec2 --output csv
  snatch ec2 -o csv

  # タグでフィルタリング
  snatch ec2 --tag Environment:Production

  # RDSインスタンスの一覧表示
  snatch rds
  snatch rds ls`

	// Long descriptionの更新
	root.Long = `Snatch is a CLI tool that helps you to manage AWS resources.

Features:
- EC2インスタンスの一覧表示
  - インスタンス名、ID、タイプ、状態などの表示
  - タグによるフィルタリング
  - タブ区切りとCSV形式の出力に対応

- RDSインスタンスの一覧表示
  - インスタンス名、クラス、エンジン情報などの表示
  - タブ区切りとCSV形式の出力に対応`
}
