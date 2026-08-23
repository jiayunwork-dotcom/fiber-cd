# fiber-cd 是阶跃光纤数值孔径（NA）与色散核算命令行工具

fiber-cd 读入纤芯/包层折射率、芯径和波长，输出 NA、归一化频率 V、
单模状态，以及材料色散与波导色散合成的 D_tot，并可扫描波长定位单模
边界与零色散波长。

## 构建 / 运行 / 测试

```text
go build ./...
go run . mode example/smf-1310.json
go test ./...
```

## 评测镜像

本目录评测专用文件（勿覆盖项目自带 Dockerfile/README）：

- `benzhi.Dockerfile`
- `build_benzhi_docker.sh`
- `BENZHI_README.md`（本文件）

两种架构都要构建并进容器验证：

```bash
chmod +x build_benzhi_docker.sh
./build_benzhi_docker.sh <image-name> linux/arm64
./build_benzhi_docker.sh <image-name> linux/amd64
docker run -it <image-name>:latest
```
