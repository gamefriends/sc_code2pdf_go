package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jung-kurt/gofpdf"
)

// 检查文件是否存在
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func main() {
	// 解析命令行参数
	inputPath := ""
	outputPath := ""
	for i := 1; i < len(os.Args); i++ {
		if os.Args[i] == "-i" && i+1 < len(os.Args) {
			inputPath = os.Args[i+1]
		} else if os.Args[i] == "-o" && i+1 < len(os.Args) {
			outputPath = os.Args[i+1]
		}
	}

	if inputPath == "" {
		fmt.Println("请提供源代码目录路径，使用 -i 选项")
		os.Exit(1)
	}

	info, err := os.Stat(inputPath)
	if err != nil || !info.IsDir() {
		fmt.Printf("输入的源代码目录路径无效: %s\n", inputPath)
		os.Exit(1)
	}

	codeName := filepath.Base(inputPath)
	if codeName == "." {
		// 如果输出的是当前目录相对路径： ./ ，则使用当前目录名称作为代码名称
		currentDir, err := os.Getwd()
		if err != nil {
			fmt.Printf("获取当前目录名称时出错: %s\n", err)
			os.Exit(1)
		}
		codeName = filepath.Base(currentDir)
	}
	fmt.Println(codeName)

	if outputPath == "" {
		outputPath = fmt.Sprintf("%s_%s.pdf", codeName, time.Now().Format("20060102"))
	}

	// 检查中文字体文件是否存在
	fontPath := "fonts/SimSun.ttf"
	if !fileExists(fontPath) {
		fmt.Printf("未找到中文字体文件: %s\n", fontPath)
		os.Exit(1)
	}

	// 创建 PDF 对象
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddUTF8Font("SimSun", "", fontPath)
	pdf.SetFont("SimSun", "", 10)

	allLines := []string{}
	// 读取所有代码文件
	filepath.Walk(inputPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && isCodeFile(info.Name()) {
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			lines := strings.Split(string(content), "\n")
			for _, line := range lines {
				if strings.TrimSpace(line) != "" {
					allLines = append(allLines, line)
				}
			}
		}
		return nil
	})

	linesPerPage := 50
	totalPages := (len(allLines) + linesPerPage - 1) / linesPerPage

	fmt.Printf("总共有 %d 行代码，共 %d 页。\n", len(allLines), totalPages)

	// 处理超过 60 页的情况，将多余的代码添加到多个 PDF 文件中
	if totalPages > 60 {
		for pageIndex := 0; pageIndex < totalPages; pageIndex += 60 {
			curOutputPath := fmt.Sprintf("%s_page%d.pdf", outputPath[:len(outputPath)-4], pageIndex/60)
			curPdf := gofpdf.New("P", "mm", "A4", "")
			curPdf.AddUTF8Font("SimSun", "", fontPath)
			curPdf.SetFont("SimSun", "", 10)
			start := pageIndex * linesPerPage
			end := min((pageIndex+60)*linesPerPage, len(allLines))
			for i := start; i < end; i += linesPerPage {
				if i != start {
					curPdf.AddPage()
					curPdf.SetFont("SimSun", "", 10)
				}
				for j, line := range allLines[i:min(i+linesPerPage, end)] {
					curPdf.Text(10, 280-3*float64(j), line)
				}
			}
			err = curPdf.OutputFileAndClose(curOutputPath)
			if err != nil {
				fmt.Println("生成 PDF 时出错:", err)
				os.Exit(1)
			}
		}
	} else {
		// 生成 PDF 内容
		for i := 0; i < len(allLines); i += linesPerPage {
			if i != 0 {
				pdf.AddPage()
				pdf.SetFont("SimSun", "", 10)
			}
			for j, line := range allLines[i:min(i+linesPerPage, len(allLines))] {
				pdf.Text(10, 280-3*float64(j), line)
			}
		}

		// 保存 PDF
		err = pdf.OutputFileAndClose(outputPath)
		if err != nil {
			fmt.Println("生成 PDF 时出错:", err)
			os.Exit(1)
		}
	}

	actualPages := (len(allLines) + linesPerPage - 1) / linesPerPage
	fmt.Printf("实际生成 %d 页。\n", actualPages)
	fmt.Printf("PDF [%s] 生成成功！\n", codeName)
}

func isCodeFile(filename string) bool {
	ext := filepath.Ext(filename)
	codeExtensions := []string{
		".java", ".py", ".ts", ".js", ".html", ".css", ".xml", ".sql", ".sh", ".properties", ".yml", ".yaml", ".json",
		".go", ".php", ".cpp", ".c", ".h", ".hpp", ".cs", ".rb", ".pl", ".lua", ".swift", ".kt", ".scala", ".groovy", ".gradle",
	}
	for _, codeExt := range codeExtensions {
		if ext == codeExt {
			return true
		}
	}
	return false
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
