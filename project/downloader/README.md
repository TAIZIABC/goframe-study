# downloader — 批量文件下载工具

从 URL 列表文件批量下载，支持并发、自动重试、进度条显示。

## 使用方式

```bash
# 从 URL 列表文件批量下载
go run main.go -file urls.txt -workers 5

# 直接传入 URL
go run main.go -out ./data https://example.com/a.zip https://example.com/b.zip

# 完整参数
go run main.go -file urls.txt -out ./data -workers 5 -retry 3 -timeout 120s
```

## 参数

| 参数 | 默认值 | 说明 |
|------|--------|------|
| `-file` | - | URL 列表文件（每行一个 URL，`#` 注释） |
| `-out` | `./downloads` | 下载输出目录 |
| `-workers` | `3` | 并发下载数 |
| `-retry` | `2` | 失败重试次数 |
| `-timeout` | `60s` | 单个文件下载超时 |

## URL 列表文件格式

```
# 注释行会跳过
https://example.com/file1.zip
https://example.com/file2.tar.gz
```

## 功能特点

- 并发下载（goroutine 池），可配置并发数
- 失败自动重试，递增等待间隔
- 实时进度条（已知大小显示百分比，未知显示已下载量）
- 彩色结果汇总（绿色成功/红色失败）
- 自动从 URL 提取文件名
- 失败时退出码为 1，方便脚本集成
