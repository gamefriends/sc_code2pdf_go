# sc_code2pdf_go
软著代码自动生成PDF报告GO语言版本

## 功能
自动扫描代码目录，并输出PDF，，每页50行。如果超出600行，自动截取前后各300行代码。

## 使用方法

将代码目录生成PDF文件，每页50行代码，支持中文。

```shell
options:
  -h, --help           show this help message and exit
  -i, --input INPUT    源代码目录路径
  -o, --output OUTPUT  保存PDF文件的路径
```

> 示例:

```shell
# 内嵌字体版本
./code2pdf_embed -i path/code_dir -p path/output.pdf

# 非嵌入字体版本，需要有 fonts/SimSun.ttf 文件
./code2pdf_no_embed -i path/code_dir -p path/output.pdf
```

## 编译方法

```shell
# 编译内嵌字体版本
go build -tags embedfonts -ldflags "-X main.EmbededFontsStr=true" -o dist/code2pdf_embed

# 常规编译(不嵌入字体)
go build -tags noembedfonts -ldflags "-X main.EmbededFontsStr=false" -o dist/code2pdf_no_embed
```
