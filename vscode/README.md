# Json5helper

Preview JSON, JSONC, JSON5, and Python `repr`-like diagnostics as canonical JSON from the editor context menu.

## Usage

1. Select text, or leave the cursor in a document to parse the whole document.
2. Right-click in the editor.
3. Choose `Json5helper: Parse Preview`.
4. Pick `Preview as JSON`, `Preview as JSONC`, `Preview as JSON5`, or `Preview as Python repr`.
5. Review the result in the read-only JSON preview editor opened beside the current editor.

The extension never modifies the source document and does not create dirty untitled files. The preview uses VS Code's normal JSON editor features, including syntax highlighting and folding. JSONC and JSON5 previews are lossy: comments and JSON5 source syntax are converted into canonical JSON.

## Local VSIX

From this directory:

```bash
npm install
npm run package
code --install-extension json5helper-vscode.vsix --force
```

Copy `json5helper-vscode.vsix` to another PC and install it from VS Code with `Extensions: Install from VSIX...`, or run the same `code --install-extension` command there.

For VS Code Insiders:

```bash
npm run package
npm run install:local:insiders
```

If the command only shows "Activating Extensions..." and no `Json5helper` output channel appears, close the old Extension Development Host or reload the target VS Code window. VS Code can keep running an older extension host after rebuilding the extension.
