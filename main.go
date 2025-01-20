package main

import (
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {

	var oldExt = flag.String("old", "", "旧的文件扩展名")
	var newExt = flag.String("new", "", "新的文件扩展名")
	var clear = flag.String("clear", "", "清除指定的文件后缀")
	var clearCN = flag.Bool("clear--zh_CN", false, "清除后缀中的中文字符")
	flag.Parse()
	if len(os.Args) == 1 {
		fmt.Println("Usage: ")
		fmt.Println("\trename-ext -clear=txt ")
		fmt.Println("\trename-ext -old=.txt -new=.md ")
		fmt.Println("\trename-ext -clear--zh_CN ")
		flag.Usage()
		fmt.Println()
		fmt.Println("注意: 文件最终重命名时,如果文件名最后一个字符为'.'时,会自动清除")
		return
	}
	// fmt.Println(*oldExt, *newExt, *clear, *clearCN)

	cur, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	files := listFiles(cur)
	// 根据不同的参数去各自处理
	if *clearCN {
		handleClearCN(files, cur)
	} else if *clear != "" {
		handleClear(files, cur, *clear)
	} else if *oldExt != "" || *newExt != "" {
		handleRepalce(files, cur, *oldExt, *newExt)
	}
}

func handleRepalce(files []fs.DirEntry, workspace, oldExt, newExt string) {
	if oldExt == newExt {
		return
	}
	if oldExt == "" {
		log.Fatal("-old 或 -new 不能为空")
	}
	if newExt == "" {
		log.Fatal("-old 或 -new 不能为空")
	}
	if newExt == "." {
		log.Fatal("-new 参数不能为'.'")
	}
	for _, item := range files {
		name := item.Name()
		dotIndex := strings.LastIndex(name, ".")
		if dotIndex != -1 {
			ext := name[dotIndex:]
			if strings.HasSuffix(ext, oldExt) {
				finalExt := strings.ReplaceAll(ext, oldExt, newExt)
				if finalExt[0] != '.' {
					finalExt = "." + finalExt
				}
				if finalExt == ext {
					return
				}
				oldname := filepath.Join(workspace, name)
				newname := filepath.Join(workspace, name[0:dotIndex]+finalExt)
				fmt.Println(oldname, newname)
			}
		}
	}
}

func handleClear(files []fs.DirEntry, workspace, clearChars string) {
	for _, item := range files {
		name := item.Name()
		dotIndex := strings.LastIndex(name, ".")
		if dotIndex != -1 {
			ext := name[dotIndex:]
			finalExt := strings.ReplaceAll(ext, clearChars, "")

			if finalExt == ext {
				return
			}

			// 构建旧/新路径
			oldname := filepath.Join(workspace, name)
			newname := filepath.Join(workspace, name[0:dotIndex]+finalExt)
			rename(oldname, newname)
		}
	}
}

func handleClearCN(files []fs.DirEntry, workspace string) {
	for _, item := range files {
		name := item.Name()
		dotIndex := strings.LastIndex(name, ".")
		if dotIndex != -1 {
			// 获取旧的文件扩展名
			ext := name[dotIndex:]
			var finalExt strings.Builder
			// 清除扩展名称中的非数字和字母
			for _, c := range ext {
				if (c == '.') || (c >= '0' && c <= '9') || (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') {
					finalExt.WriteRune(c)
				}
			}
			// 优化,如果前后扩展名称一致,不需要操作
			if finalExt.String() == ext {
				continue
			}
			// 构建旧/新名称
			oldname := filepath.Join(workspace, name)
			newname := filepath.Join(workspace, name[0:dotIndex]+finalExt.String())
			// 重名
			rename(oldname, newname)
		}
	}
}

// 文件重命名
func rename(oldname, newname string) {
	// 处理新文件名,如果最后一个为'.', 那么自动清除
	if newname[len(newname)-1] == '.' {
		newname = newname[0 : len(newname)-1]
	}
	os.Rename(oldname, newname)
	fmt.Println("done", oldname, "==>", newname)
}

// 获取当前文件夹下的文件,不包括文件夹
func listFiles(dirname string) []fs.DirEntry {
	entries, err := os.ReadDir(dirname)
	if err != nil {
		log.Fatal(err)
	}
	var result []fs.DirEntry
	for _, item := range entries {
		if !item.IsDir() {
			result = append(result, item)
		}
	}

	return result
}
