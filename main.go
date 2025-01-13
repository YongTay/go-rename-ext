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

	if len(args) == 0 {
		log.Fatal("请提供至少一个文件扩展名参数")
	}

	var oldExt, newExt string
	if len(args) == 1 {
		newExt = args[0]
		// 如果新扩展名没有以 '.' 开头，自动补充
		if !strings.HasPrefix(newExt, ".") {
			newExt = "." + newExt
		}
	} else {
		oldExt = args[0]
		newExt = args[1]
		// 如果旧或新扩展名没有以 '.' 开头，自动补充
		if !strings.HasPrefix(oldExt, ".") {
			oldExt = "." + oldExt
		}
		if !strings.HasPrefix(newExt, ".") {
			newExt = "." + newExt
		}
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
		if oldExt == "" {
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
