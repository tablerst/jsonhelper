import * as path from 'path';
import * as vscode from 'vscode';

import { Json5helperBackend, ParseMode, WasmJson5helperBackend } from './backend';
import { ParsePreviewInput, ParsePreviewServices, parsePreview } from './previewFlow';

type ParseCommand = {
  command: string;
  mode: ParseMode;
  label: string;
};

const parseCommands: ReadonlyArray<ParseCommand> = [
  {
    command: 'json5helper.previewJson',
    label: 'Preview as JSON',
    mode: 'json'
  },
  {
    command: 'json5helper.previewJsonc',
    label: 'Preview as JSONC',
    mode: 'jsonc'
  },
  {
    command: 'json5helper.previewJson5',
    label: 'Preview as JSON5',
    mode: 'json5'
  },
  {
    command: 'json5helper.previewRepr',
    label: 'Preview as Python repr',
    mode: 'repr'
  }
];

let outputChannel: vscode.OutputChannel | undefined;

export function activate(context: vscode.ExtensionContext): void {
  outputChannel = vscode.window.createOutputChannel('Json5helper');
  context.subscriptions.push(outputChannel);
  outputChannel.appendLine(`[${new Date().toISOString()}] Json5helper extension activating`);

  const backend = createBackend(context);
  const previewPanel = new JsonPreviewPanel();

  for (const command of parseCommands) {
    context.subscriptions.push(
      vscode.commands.registerCommand(command.command, async () => {
        outputChannel?.appendLine(`[${new Date().toISOString()}] ${command.label} command invoked`);
        await parsePreview(
          command.mode,
          backend,
          createVsCodeParsePreviewServices(outputChannel!, previewPanel)
        );
      })
    );
  }

  outputChannel.appendLine(`[${new Date().toISOString()}] Json5helper extension activated`);
}

export function deactivate(): void {
  outputChannel?.dispose();
  outputChannel = undefined;
}

function createBackend(context: vscode.ExtensionContext): Json5helperBackend {
  return new WasmJson5helperBackend(async () =>
    require(path.join(context.extensionPath, 'wasm', 'json5helper_wasm.js'))
  );
}

function createVsCodeParsePreviewServices(
  channel: vscode.OutputChannel,
  previewPanel: JsonPreviewPanel
): ParsePreviewServices {
  return {
    getInput(): ParsePreviewInput | undefined {
      const editor = vscode.window.activeTextEditor;
      if (editor === undefined) {
        return undefined;
      }

      return {
        selection: editor.selection.isEmpty ? '' : editor.document.getText(editor.selection),
        document: editor.document.getText()
      };
    },
    async openJsonPreview(content: string): Promise<void> {
      previewPanel.show(content);
    },
    showError(message: string): void {
      void vscode.window.showErrorMessage(message);
    },
    logError(title: string, message: string): void {
      channel.appendLine(`[${new Date().toISOString()}] ${title}`);
      channel.appendLine(message);
      channel.appendLine('');
    }
  };
}

class JsonPreviewPanel {
  private panel: vscode.WebviewPanel | undefined;

  public show(content: string): void {
    if (this.panel === undefined) {
      this.panel = vscode.window.createWebviewPanel(
        'json5helperPreview',
        'Json5helper Preview',
        vscode.ViewColumn.Beside,
        {
          enableScripts: false,
          retainContextWhenHidden: true
        }
      );
      this.panel.onDidDispose(() => {
        this.panel = undefined;
      });
    } else {
      this.panel.reveal(vscode.ViewColumn.Beside);
    }

    this.panel.webview.html = renderPreviewHtml(content);
  }
}

function renderPreviewHtml(content: string): string {
  return `<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <style>
    body {
      margin: 0;
      padding: 16px;
      color: var(--vscode-editor-foreground);
      background: var(--vscode-editor-background);
      font-family: var(--vscode-editor-font-family);
      font-size: var(--vscode-editor-font-size);
    }
    pre {
      margin: 0;
      white-space: pre-wrap;
      word-break: break-word;
    }
  </style>
</head>
<body><pre>${escapeHtml(content)}</pre></body>
</html>`;
}

function escapeHtml(value: string): string {
  return value
    .replaceAll('&', '&amp;')
    .replaceAll('<', '&lt;')
    .replaceAll('>', '&gt;');
}
