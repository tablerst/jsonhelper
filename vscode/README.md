# Json5helper

Preview JSON, JSONC, JSON5, and Python `repr`-like diagnostics as canonical JSON from the editor context menu.

## Usage

1. Select text, or leave the cursor in a document to parse the whole document.
2. Right-click in the editor.
3. Choose `Json5helper: Parse Preview`.
4. Pick `Parse as JSON`, `Parse as JSONC`, `Parse as JSON5`, or `Parse as Python repr`.
5. Review the result in the editable untitled JSON document opened beside the current editor.

The extension never modifies the source document. JSONC and JSON5 previews are lossy: comments and JSON5 source syntax are converted into canonical JSON.

## Local VSIX

From this directory:

```bash
npm install
npm run package
code --install-extension json5helper-vscode.vsix --force
```

Copy `json5helper-vscode.vsix` to another PC and install it from VS Code with `Extensions: Install from VSIX...`, or run the same `code --install-extension` command there.
