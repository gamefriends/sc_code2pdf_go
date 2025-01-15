package main

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jung-kurt/gofpdf"
)

//go:embed fonts/*
var embeddedFonts embed.FS

// 检查文件是否存在
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// 每页宽度
var PAGE_WIDTH = 190.0

// 字体大小
var FONT_SIZE = 9.0

// 每页行数
var PAGE_LINES = 50

// 批量页数
var BATCH_PAGES = 30

// 行间距
var LINE_GAP = 5.3

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
	fontBytes, err := embeddedFonts.ReadFile("fonts/SimSun.ttf")
	if err != nil {
		fmt.Printf("读取嵌入字体文件时出错: %s\n", err)
		os.Exit(1)
	}
	pdf.AddUTF8FontFromBytes("SimSun", "", fontBytes)
	pdf.SetFont("SimSun", "", FONT_SIZE)

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

	totalLines := 0

	printLines := []string{}

	for _, line := range allLines {
		// 如果是空行，忽略
		if strings.TrimSpace(line) == "" {
			continue
		}
		lines := pdf.SplitText(line, PAGE_WIDTH)
		totalLines += len(lines)
		printLines = append(printLines, lines...)
	}
	totalPages := (totalLines + PAGE_LINES - 1) / PAGE_LINES

	fmt.Printf("总共有 %d 行代码，共 %d 页。\n", totalLines, totalPages)

	// 如果printLines为空，说明没有代码文件
	if len(printLines) == 0 {
		fmt.Println("没有找到任何代码文件。")
		os.Exit(1)
	}

	// 如果超出60页, 则 printLines = 前30页 + 后30页
	if totalPages > 60 {
		batch_lines := PAGE_LINES * BATCH_PAGES
		printLines = append(printLines[:batch_lines], printLines[totalLines-batch_lines:]...)
	}
	totalPages = (len(printLines) + PAGE_LINES - 1) / PAGE_LINES

	// 每页打印 PAGE_LINES 行代码
	pageNum := 1
	for i := 0; i < totalPages; i++ {
		start := i * PAGE_LINES
		end := min((i+1)*PAGE_LINES, totalLines)
		printPage(pdf, printLines[start:end])
		fmt.Printf("正在生成第 %d 页，共 %d 页\n", pageNum, totalPages)
		pageNum++
	}

	// 保存 PDF 文件
	err = pdf.OutputFileAndClose(outputPath)
	if err != nil {
		fmt.Printf("保存 PDF 文件时出错: %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("PDF 文件已保存到: %s\n", outputPath)

}

func printPage(pdf *gofpdf.Fpdf, lines []string) {
	pdf.AddPage()
	pdf.SetFont("SimSun", "", FONT_SIZE)
	pdf.SetXY(10, 10)
	for num, line := range lines {
		fmt.Printf("%d 行: %s\n", num+1, line)
		pdf.MultiCell(PAGE_WIDTH, float64(LINE_GAP), line, "", "L", false)
	}
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
