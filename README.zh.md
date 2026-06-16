# json5helper

`json5helper` 正在从早期 Go 原型迁移为 Rust workspace，用于解析和格式化 JSON 系列配置格式。

当前 Rust 实现覆盖：

- 将 JSON、JSONC、JSON5 解析为 `serde_json::Value`
- 输出紧凑 JSON 或 pretty JSON
- 提供 `json5helper` CLI 二进制命令
- 提供 VS Code 右键解析预览插件
- 将 Python `repr` 风格对象图转换为有损诊断 JSON
- 提供供编辑器集成使用的 WebAssembly 包装层

旧 Go 实现仍保留在仓库中作为历史参考，主要位于 `internal/`、`pkg/`、`test/` 和 `jsonhelper.go`。

## Workspace

```text
crates/
  json5helper-core/   # JSON、JSONC、JSON5 解析与格式化 API
  json5helper-cli/    # CLI 二进制入口
  json5helper-wasm/   # 编辑器集成用 WebAssembly 包装层
  repr-json/          # Python repr-like 文本转诊断 JSON
vscode/               # VS Code 插件包
```

## 常用命令

```bash
cargo test --workspace
cargo fmt --all
cargo clippy --workspace --all-targets
```

解析并格式化 JSON5：

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

## VS Code 插件

`vscode/` 包提供显式的右键解析预览工作流，而不是自动 formatter。插件不会注册 Format Document，不会在保存时运行，也不会修改源文档。

使用方式：

1. 选中文本，或不选择文本以使用整个当前文档。
2. 在编辑器里右键选择 `Json5helper: Parse Preview`。
3. 在子菜单中选择 `Preview as JSON`、`Preview as JSONC`、`Preview as JSON5` 或 `Preview as Python repr`。
4. 结果会在右侧 preview 面板中打开。

插件不会修改源文档，也不会创建需要保存/丢弃的 untitled 文件。

命令面板中也可以直接运行各个 `Json5helper` preview 命令，这适合用来确认插件是否已经正常激活。

JSONC 和 JSON5 预览是有损转换：注释、单引号、未加引号的 key、尾逗号等源码层语法会被转换为 canonical JSON。

插件开发命令：

```bash
cd vscode
npm install
npm run build
```

从仓库根目录可以直接使用 VS Code launch 配置 `Run Json5helper VS Code Extension`，它会先构建插件，再打开 Extension Development Host。默认配置保留普通 extension-host 环境，方便观察兼容性问题。如果其他已安装扩展污染日志或干扰调试，可以改用 `Run Json5helper VS Code Extension (Isolated)`；该配置会禁用其他扩展，并把临时 VS Code 状态放在 `.vscode-dev/` 下。

`npm run build` 会构建 Rust WebAssembly 包装层并编译 TypeScript 插件。必要时先安装一次性前置工具：

```bash
rustup target add wasm32-unknown-unknown
cargo install wasm-bindgen-cli --version 0.2.125 --locked
```

打包插件并同步到其他 PC：

```bash
cd vscode
npm run package
```

该命令会生成 `vscode/json5helper-vscode.vsix`。把这个文件复制到其他 PC 后，可以在 VS Code 中运行 `Extensions: Install from VSIX...` 安装，或使用命令：

```bash
code --install-extension json5helper-vscode.vsix --force
```

如果使用 VS Code Insiders：

```bash
code-insiders --install-extension json5helper-vscode.vsix --force
```

如果运行命令时只显示 "Activating Extensions..."，但没有出现 `Json5helper` output channel，先关闭旧的 Extension Development Host，或 reload 目标 VS Code 窗口。VS Code 在插件重新构建后仍可能继续运行旧的 extension host。

## 覆盖率目标

Rust 实现目标是至少 80% 行覆盖率。安装 `cargo-llvm-cov` 后可运行：

```bash
cargo llvm-cov --workspace --fail-under-lines 80
```

## WebAssembly

核心解析逻辑放在 `json5helper-core`，避免绑定文件系统。当前 VS Code 插件使用内置的 `wasm-bindgen` 后端，因此用户不需要安装 Rust、`json5helper` CLI 或 `wasmtime`。
