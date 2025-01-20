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
		printHelp()
		return
	}

	var oldExt, newExt string
	clear := false
	clearChars := "" // 用于存储要清除的字符集

	// 解析参数
	parseArgs(args, &oldExt, &newExt, &clear, &clearChars)

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
			clearExtension(dir, oldName, clearChars)
		} else if oldExt == "" {
			replaceExtension(dir, oldName, filepath.Ext(oldName), newExt)
		} else if strings.HasSuffix(oldName, oldExt) {
			replaceExtension(dir, oldName, oldExt, newExt)
		}
	}
}

// 打印帮助信息
func printHelp() {
	fmt.Println("使用说明:")
	fmt.Println("  - 使用 -old 和 -new 参数指定旧和新扩展名。")
	fmt.Println("  - 如果只提供 -new 参数：将所有文件的扩展名更改为该参数指定的新扩展名。")
	fmt.Println("  - 使用 -clear 参数清除文件扩展名中的非字母字符，或指定要清除的字符集。")
	fmt.Println("  - 扩展名参数可以不带 '.'，程序会自动补充。")
	fmt.Println("示例:")
	fmt.Println("  rename-ext -old .txt -new .md  # 将所有 .txt 文件改为 .md")
	fmt.Println("  rename-ext -new .md            # 将所有文件改为 .md 扩展名")
	fmt.Println("  rename-ext -clear              # 清除所有文件扩展名中的非字母字符")
	fmt.Println("  rename-ext -clear 123          # 清除所有文件扩展名中的 '1', '2', '3' 字符")
}

// 解析命令行参数
func parseArgs(args []string, oldExt, newExt *string, clear *bool, clearChars *string) {
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-old":
			if i+1 < len(args) {
				*oldExt = args[i+1]
				if !strings.HasPrefix(*oldExt, ".") {
					*oldExt = "." + *oldExt
				}
				i++
			} else {
				log.Fatal("-old 参数后需要指定扩展名")
			}
		case "-new":
			if i+1 < len(args) {
				*newExt = args[i+1]
				if !strings.HasPrefix(*newExt, ".") {
					*newExt = "." + *newExt
				}
				i++
			} else {
				log.Fatal("-new 参数后需要指定扩展名")
			}
		case "-clear":
			*clear = true
			if i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
				*clearChars = args[i+1]
				i++
			}
		default:
			log.Fatalf("未知参数: %s", args[i])
		}
	}

	if *newExt == "" && !*clear {
		log.Fatal("请提供 -new 参数以指定新的文件扩展名或使用 -clear 参数")
	}
}

// 清除扩展名中的指定字符
func clearExtension(dir, oldName, clearChars string) {
	ext := filepath.Ext(oldName)
	if len(ext) > 1 { // 确保扩展名长度大于1
		cleanExt := "." + clearSpecified(ext[1:], clearChars)
		newName := strings.TrimSuffix(oldName, ext) + cleanExt
		renameFile(dir, oldName, newName)
	} else {
		// 如果没有有效的扩展名，保持原文件名
		renameFile(dir, oldName, oldName)
	}
}

// 根据指定字符集清除字符串中的字符
func clearSpecified(s, clearChars string) string {
	var result strings.Builder
	for _, c := range s {
		if clearChars == "" {
			// 默认清除非字母字符
			if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') {
				result.WriteRune(c)
			}
		} else {
			// 清除指定字符集中的字符
			if !strings.ContainsRune(clearChars, c) {
				result.WriteRune(c)
			}
		}
	}
	return result.String()
}

// 替换文件扩展名
func replaceExtension(dir, oldName, oldExt, newExt string) {
	newName := strings.TrimSuffix(oldName, oldExt) + newExt
	renameFile(dir, oldName, newName)
}

// 重命名文件
func renameFile(dir, oldName, newName string) {
	// 去掉新文件名的最后一个字符如果是 '.'
	newName = strings.TrimSuffix(newName, ".")

	oldPath := filepath.Join(dir, oldName)
	newPath := filepath.Join(dir, newName)
	err := os.Rename(oldPath, newPath)
	if err != nil {
		log.Printf("无法重命名文件 %s: %v\n", oldName, err)
	} else {
		fmt.Printf("文件 %s 已重命名为 %s\n", oldName, newName)
	}
}
