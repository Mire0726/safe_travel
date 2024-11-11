package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"os"
)

type TemplateData struct {
	Package               string
	CamelTableName        string
	PascalTableName       string
	PluralCamelTableName  string
	PluralPascalTableName string
	TableName             string
}

func main() {
	// テンプレートファイルを読み込む
	tmpl, err := template.ParseFiles("template.tpl")
	if err != nil {
		log.Fatalf("Error parsing template: %v", err)
	}

	// データを準備
	data := TemplateData{
		Package:               "repository",
		CamelTableName:        "user",
		PascalTableName:       "User",
		PluralCamelTableName:  "users",
		PluralPascalTableName: "Users",
		TableName:             "users",
	}

	// テンプレートにデータを適用
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		log.Fatalf("Error executing template: %v", err)
	}

	// 出力ファイル名を生成
	outputFileName := fmt.Sprintf("../../backend/api/infrastructure/sql/%s_generated_sql.go", data.CamelTableName)

	// 出力ファイルに書き込む
	err = os.WriteFile(outputFileName, buf.Bytes(), 0o644)
	if err != nil {
		log.Fatalf("Error writing file: %v", err)
	}

	log.Println("Code generated successfully!")
}
