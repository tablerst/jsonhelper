# json5helper

`json5helper` 正在从早期 Go 原型迁移为 Rust workspace，用于解析和格式化 JSON 系列配置格式。

当前 Rust 实现目标：

- 解析 JSON、JSONC、JSON5，并转换为 `serde_json::Value`
- 输出紧凑 JSON 或美化后的 JSON
- 提供 `json5helper` CLI 二进制命令
- 额外提供 Python `repr` 风格对象图到诊断 JSON 的有损转换
- 后续支持 WASI 构建，供编辑器集成使用

原 Go 代码仍保留在仓库中作为历史参考，主要位于 `internal/`、`pkg/`、`test/` 和 `jsonhelper.go`。

## Workspace 结构

```text
crates/
  json5helper-core/   # JSON、JSONC、JSON5 的解析与格式化 API
  json5helper-cli/    # CLI 二进制入口
  repr-json/          # Python repr-like 文本转诊断 JSON
```

## 常用命令

```bash
cargo test --workspace
cargo fmt --all
cargo clippy --workspace --all-targets
```

格式化 JSON5 文件：

```bash
json5helper fmt --syntax json5 example.json5
```

从 stdin 读取：

```bash
echo "{unquoted: 'value', trailing: [1, 2,]}" | json5helper fmt --syntax json5
```

转换 Python `repr` 风格诊断文本：

```bash
echo "AgentExecutor(verbose=True)" | json5helper repr-json
```

## 覆盖率目标

Rust 实现的目标是至少 80% 行覆盖率。安装 `cargo-llvm-cov` 后可运行：

```bash
cargo llvm-cov --workspace --fail-under-lines 80
```

## WASI

核心解析逻辑放在 `json5helper-core`，避免绑定文件系统，后续可以包装为 `wasm32-wasip2` 目标或供 VS Code 插件调用。VS Code 插件本身暂未在本仓库实现。
