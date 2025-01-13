package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	// 获取命令行参数
	args := os.Args[1:]

	// 检查是否请求帮助
	if len(args) > 0 && args[0] == "--help" {
		fmt.Println("使用说明:")
		fmt.Println("  - 使用 -old 和 -new 参数指定旧和新扩展名。")
		fmt.Println("  - 如果只提供 -new 参数：将所有文件的扩展名更改为该参数指定的新扩展名。")
		fmt.Println("  - 使用 -clear 参数清除文件扩展名中的非字母字符。")
		fmt.Println("  - 扩展名参数可以不带 '.'，程序会自动补充。")
		fmt.Println("示例:")
		fmt.Println("  go run main.go -old .txt -new .md  # 将所有 .txt 文件改为 .md")
		fmt.Println("  go run main.go -new .md            # 将所有文件改为 .md 扩展名")
		fmt.Println("  go run main.go -clear              # 清除所有文件扩展名中的非字母字符")
		return
	}

	var oldExt, newExt string
	clear := false

	// 解析参数
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-old":
			if i+1 < len(args) {
				oldExt = args[i+1]
				if !strings.HasPrefix(oldExt, ".") {
					oldExt = "." + oldExt
				}
				i++
			} else {
				log.Fatal("-old 参数后需要指定扩展名")
			}
		case "-new":
			if i+1 < len(args) {
				newExt = args[i+1]
				if !strings.HasPrefix(newExt, ".") {
					newExt = "." + newExt
				}
				i++
			} else {
				log.Fatal("-new 参数后需要指定扩展名")
			}
		case "-clear":
			clear = true
		default:
			log.Fatalf("未知参数: %s", args[i])
		}
	}

	if newExt == "" && !clear {
		log.Fatal("请提供 -new 参数以指定新的文件扩展名或使用 -clear 参数")
	}

	// 获取当前工作目录
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	// 使用 os.ReadDir 读取目录中的文件
	files, err := os.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	// 修改文件扩展名
	for _, file := range files {
		if file.IsDir() {
			continue // 跳过文件夹
		}
		oldName := file.Name()
		// 排除 .go 和 .mod 文件
		if strings.HasSuffix(oldName, ".go") || strings.HasSuffix(oldName, ".mod") {
			continue
		}
		if clear {
			// 清除扩展名中的非字母字符
			ext := filepath.Ext(oldName)
			cleanExt := "." + clearNonAlpha(ext[1:])
			newName := strings.TrimSuffix(oldName, ext) + cleanExt
			oldPath := filepath.Join(dir, oldName)
			newPath := filepath.Join(dir, newName)
			err := os.Rename(oldPath, newPath)
			if err != nil {
				log.Printf("无法重命名文件 %s: %v\n", oldName, err)
			} else {
				fmt.Printf("文件 %s 已重命名为 %s\n", oldName, newName)
			}
		} else if oldExt == "" {
			// 去除旧的后缀并添加新的后缀
			newName := strings.TrimSuffix(oldName, filepath.Ext(oldName)) + newExt
			oldPath := filepath.Join(dir, oldName)
			newPath := filepath.Join(dir, newName)
			err := os.Rename(oldPath, newPath)
			if err != nil {
				log.Printf("无法重命名文件 %s: %v\n", oldName, err)
			} else {
				fmt.Printf("文件 %s 已重命名为 %s\n", oldName, newName)
			}
		} else if strings.HasSuffix(oldName, oldExt) {
			// 如果文件名以旧的后缀结尾，则替换为新的后缀
			newName := strings.TrimSuffix(oldName, oldExt) + newExt
			oldPath := filepath.Join(dir, oldName)
			newPath := filepath.Join(dir, newName)
			err := os.Rename(oldPath, newPath)
			if err != nil {
				log.Printf("无法重命名文件 %s: %v\n", oldName, err)
			} else {
				fmt.Printf("文件 %s 已重命名为 %s\n", oldName, newName)
			}
		}
	}
}

// 辅助函数：清除字符串中的非字母字符
func clearNonAlpha(s string) string {
	var result strings.Builder
	for _, c := range s {
		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') {
			result.WriteRune(c)
		}
	}
	return result.String()
}
