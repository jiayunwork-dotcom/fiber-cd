# fiber-cd：Go 阶跃光纤 NA 与色散 Web 服务（Sellmeier + Gloge + 前端控制台）

输入纤芯/包层折射率、芯径和波长，算出 NA、归一化频率 V、单模判定，以及材料色散与波导色散合成的总色散。V 必须用半径而不是直径；总零色散波长被负的 D_wg 从纯石英材料零点推移。还可以核算熔接重叠/错位（Marcuse 与 Petermann 光斑不得混用）以及近截止宏弯损耗。

## 构建 / 运行 / 测试

```text
go build ./...
./fiber-cd -http :8080
curl -s http://127.0.0.1:8080/api/example
go run . mode example/smf-1310.json
go test ./...
```

## 评测镜像

- `benzhi.Dockerfile`
- `build_benzhi_docker.sh`
- `BENZHI_README.md`（本文件）

```bash
chmod +x build_benzhi_docker.sh
./build_benzhi_docker.sh <image-name> linux/amd64
docker run -d -P --name fiber-cd-b14 <image-name>:latest
curl -s http://127.0.0.1:$(docker port fiber-cd-b14 8080 | cut -d: -f2)/api/example
docker rm -f fiber-cd-b14
```
