import * as vscode from 'vscode';

import { Json5helperBackend, ParseMode, WasmJson5helperBackend } from './backend';
import { ParsePreviewInput, ParsePreviewServices, parsePreview } from './previewFlow';

type ParseChoice = {
  label: string;
  description: string;
  mode: ParseMode;
};

const parseChoices: ReadonlyArray<ParseChoice> = [
  {
    label: 'Parse as JSON',
    description: 'Preview strict JSON as canonical JSON',
    mode: 'json'
  },
  {
    label: 'Parse as JSONC',
    description: 'Preview JSON with comments as canonical JSON',
    mode: 'jsonc'
  },
  {
    label: 'Parse as JSON5',
    description: 'Preview JSON5 as canonical JSON',
    mode: 'json5'
  },
  {
    label: 'Parse as Python repr',
    description: 'Preview repr-like diagnostics as JSON',
    mode: 'repr'
  }
];

let outputChannel: vscode.OutputChannel | undefined;

export function activate(context: vscode.ExtensionContext): void {
  outputChannel = vscode.window.createOutputChannel('Json5helper');
  context.subscriptions.push(outputChannel);

  const backend = createBackend();
  context.subscriptions.push(
    vscode.commands.registerCommand('json5helper.parsePreview', () =>
      parsePreview(
        backend,
        createVsCodeParsePreviewServices(outputChannel!)
      )
    )
  );
}

export function deactivate(): void {
  outputChannel?.dispose();
  outputChannel = undefined;
}

function createBackend(): Json5helperBackend {
  return new WasmJson5helperBackend(async () => require('../wasm/json5helper_wasm.js'));
}

function createVsCodeParsePreviewServices(channel: vscode.OutputChannel): ParsePreviewServices {
  return {
    async chooseMode(): Promise<ParseMode | undefined> {
      const choice = await vscode.window.showQuickPick(parseChoices, {
        placeHolder: 'Choose how to parse the current selection or document'
      });
      return choice?.mode;
    },
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
      const document = await vscode.workspace.openTextDocument({
        content,
        language: 'json'
      });
      await vscode.window.showTextDocument(document, {
        preview: false,
        viewColumn: vscode.ViewColumn.Beside
      });
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
