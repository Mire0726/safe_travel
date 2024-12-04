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

	// テーブルデータのリストを準備
	tables := []TemplateData{
		{
			Package:               "sql",
			CamelTableName:        "user",
			PascalTableName:       "User",
			PluralCamelTableName:  "users",
			PluralPascalTableName: "Users",
			TableName:             "users",
		},
		{
			Package:               "sql",
			CamelTableName:        "event",
			PascalTableName:       "Event",
			PluralCamelTableName:  "events",
			PluralPascalTableName: "Events",
			TableName:             "events",
		},
		{
			Package:               "sql",
			CamelTableName:        "transport",
			PascalTableName:       "Transport",
			PluralCamelTableName:  "transports",
			PluralPascalTableName: "Transports",
			TableName:             "transports",
		},
	}

	// 各テーブルデータに対してテンプレートを適用し、ファイルを生成
	for _, data := range tables {
		var buf bytes.Buffer
		err = tmpl.Execute(&buf, data)
		if err != nil {
			log.Fatalf("Error executing template for table %s: %v", data.CamelTableName, err)
		}

		// 出力ファイル名を生成
		outputFileName := fmt.Sprintf("../../backend/api/infrastructure/datastore/datastoresql/%s/%s_generated_sql.go", data.CamelTableName,data.CamelTableName)

		// 出力ファイルに書き込む（既存の場合は上書き）
		err = os.WriteFile(outputFileName, buf.Bytes(), 0o644)
		if err != nil {
			log.Fatalf("Error writing file for table %s: %v", data.CamelTableName, err)
		}

		log.Printf("Code generated successfully for table %s!", data.CamelTableName)
	}
}
