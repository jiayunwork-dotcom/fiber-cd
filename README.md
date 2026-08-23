# fiber-cd

阶跃光纤数值孔径（NA）与色散核算命令行工具。用户给出纤芯/包层折射率、
芯径（直径）和波长，工具算出 NA、归一化频率 V、是否单模，以及由材料色散
D_mat 与波导色散 D_wg 合成的总色散 D_tot。材料色散来自熔融石英
（Malitson 1965）Sellmeier 系数对 n(λ) 的解析求导，波导色散采用阶跃近似
（Gloge 1971）的文献公式，覆盖 1310/1550 nm 通信窗口的快速手算复核。
本工具只做圆阶跃光纤波导，不含光模块 SKU、链路工单，也不是金属矩形波导
TE10 截止（那是另一个工具的事）。

## 用法

```text
go run . mode example/smf-1310.json
```

`mode` 子命令读取一份 JSON 配置，打印 NA、V、截止波长、单模状态与色散
三项。`example/smf-1310.json` 是通信级阶跃光纤（芯径 9.0 μm、Δ≈0.25%）
在 1310 nm 下的核算：判单模，且 D_tot 接近零（总零色散波长约 1311 nm）。

其余子命令：

```text
go run . spec example/smf-1310.json                    # 完整表征：光斑、模式数、斜率、群折射率
go run . design example/smf-1310.json                  # 单模设计限值：最大芯径 / 最大 NA
go run . probe example/smf-1310.json 1550.0            # 换波长核算单点状态
go run . boundary example/smf-1310.json                # 截止波长（V=2.405）处的色散
go run . sweep example/smf-1310.json 1260 1620 41      # 波长扫描表，含模式边界与零色散标记
go run . dump example/smf-1310.json                    # 回显规范化配置 JSON
go run . validate example/smf-1310.json                # 只做输入校验
```

任何配置路径写成 `-` 就从标准输入读 JSON，例如：

```bash
echo '{"n1":1.4656,"n2":1.4619,"core_diameter_um":9,"wavelength_nm":1310}' \
  | go run . mode -
```

配置字段：

```json
{
  "name": "SMF-1310 step-index",
  "n1": 1.4656,
  "n2": 1.4619,
  "core_diameter_um": 9.0,
  "wavelength_nm": 1310.0,
  "sweep_start_nm": 1200,
  "sweep_stop_nm": 1700,
  "sweep_steps": 51
}
```

`sweep_start_nm` / `sweep_stop_nm` / `sweep_steps` 可缺省（默认
1200–1700 nm、51 档），也可用 `sweep` 子命令的位置参数覆盖。

非法输入（n2 ≥ n1 即包层比纤芯还密、芯径 ≤ 0、波长 ≤ 0、文件不存在、
参数不合法）一律写入标准错误并返回非零退出码。

## 关键约定

- 数值孔径 `NA = √(n1² − n2²)`。
- 归一化频率 `V = 2π·a·NA/λ`，其中 `a` 是纤芯**半径**，由芯径除以 2
  得到；把芯径直接当半径会把 V 翻倍、截止波长差一倍、单模判定全错。
- 单模判定 `V < 2.405`（圆阶跃第二模式截止，J0 首根 2.4048…）。
- 材料色散 `D_mat = −(λ/c)·d²n/dλ²`，Sellmeier 系数固定为 Malitson
  1965 熔融石英；纯石英零色散波长约 1273 nm。
- 波导色散 `D_wg = −(n1·Δ/(λ·c))·V·d²(Vb)/dV²`，b(V) 用 Gloge 近似
  (1.1428 − 0.996/V)²，适用 1.5 < V < 2.5，域外饱和。
- 总色散 `D_tot = D_mat + D_wg`，二者在零色散波长附近符号翻转；
  总零色散波长由二分求解，区间内不跨零时明确报告未发现。
- 交叉规律：只把芯径加倍 → V 加倍，原单模波长可能变多模；只把 n1−n2
  缩小 → NA 与 V 下降；λ 加倍（仍在公式适用区）→ V 减半。

## 构建与测试

```text
go build ./...
go test ./...
```

纯标准库实现，无第三方依赖。`go run . mode example/smf-1310.json` 可复现
上文示例输出。

## 许可

MIT，见 [LICENSE](LICENSE)。
