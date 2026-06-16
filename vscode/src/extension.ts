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

let outputChannel: vscode.LogOutputChannel | undefined;
let extensionLogger: ExtensionLogger | undefined;

export function activate(context: vscode.ExtensionContext): void {
  outputChannel = vscode.window.createOutputChannel('Json5helper', { log: true }) as vscode.LogOutputChannel;
  extensionLogger = new ExtensionLogger(outputChannel, context.logUri);
  context.subscriptions.push(extensionLogger);

  extensionLogger.info('Json5helper extension activating');
  extensionLogger.info(`extensionPath=${context.extensionPath}`);
  extensionLogger.info(`extensionUri=${context.extensionUri.toString()}`);
  extensionLogger.info(`logUri=${context.logUri.toString()}`);
  extensionLogger.info(`wasmModulePath=${getWasmModulePath(context)}`);

  const backend = createBackend(context);
  const previewController = new JsonPreviewController(context, extensionLogger);
  context.subscriptions.push(previewController);

  for (const command of parseCommands) {
    context.subscriptions.push(
      vscode.commands.registerCommand(command.command, async () => {
        extensionLogger?.info(`${command.label} command invoked`);
        try {
          await parsePreview(
            command.mode,
            backend,
            createVsCodeParsePreviewServices(extensionLogger!, previewController)
          );
        } catch (error) {
          extensionLogger?.error(`${command.label} command failed unexpectedly`, error);
          extensionLogger?.show();
          void vscode.window.showErrorMessage('Json5helper: command failed unexpectedly. See output for details.');
        }
      })
    );
  }

  context.subscriptions.push(
    vscode.commands.registerCommand('json5helper.openDiagnosticLog', async () => {
      extensionLogger?.info('Open Diagnostic Log command invoked');
      await extensionLogger?.openLogDocument();
    })
  );

  extensionLogger.info('Json5helper extension activated');
}

export function deactivate(): void {
  outputChannel?.dispose();
  outputChannel = undefined;
  extensionLogger?.dispose();
  extensionLogger = undefined;
}

function createBackend(context: vscode.ExtensionContext): Json5helperBackend {
  return new WasmJson5helperBackend(async () => {
    const wasmModulePath = getWasmModulePath(context);
    extensionLogger?.info(`Loading WASM module from ${wasmModulePath}`);
    const wasm = require(wasmModulePath);
    extensionLogger?.info(`WASM module loaded with exports: ${Object.keys(wasm).join(', ')}`);
    return wasm;
  });
}

function getWasmModulePath(context: vscode.ExtensionContext): string {
  return path.join(context.extensionPath, 'wasm', 'json5helper_wasm.js');
}

function createVsCodeParsePreviewServices(
  logger: ExtensionLogger,
  previewController: JsonPreviewController
): ParsePreviewServices {
  return {
    getInput(): ParsePreviewInput | undefined {
      const editor = vscode.window.activeTextEditor;
      if (editor === undefined) {
        logger.info('No active text editor');
        return undefined;
      }

      logger.info(
        [
          `activeEditor.uri=${editor.document.uri.toString()}`,
          `language=${editor.document.languageId}`,
          `selectionEmpty=${editor.selection.isEmpty}`,
          `documentChars=${editor.document.getText().length}`
        ].join(' ')
      );

      return {
        selection: editor.selection.isEmpty ? '' : editor.document.getText(editor.selection),
        document: editor.document.getText()
      };
    },
    async openJsonPreview(content: string): Promise<void> {
      logger.info(`Opening JSON preview (${content.length} chars)`);
      await previewController.show(content);
      logger.info('JSON preview opened');
    },
    showError(message: string): void {
      logger.info(`Showing error message: ${message}`);
      void vscode.window.showErrorMessage(message);
    },
    logInfo(message: string): void {
      logger.info(message);
    },
    logError(title: string, message: string): void {
      logger.error(title, message);
      logger.show();
    }
  };
}

class JsonPreviewController implements vscode.Disposable {
  private provider: JsonPreviewDocumentProvider | undefined;
  private providerRegistration: vscode.Disposable | undefined;

  public constructor(
    private readonly context: vscode.ExtensionContext,
    private readonly logger: ExtensionLogger
  ) {
  }

  public async show(content: string): Promise<void> {
    if (this.provider === undefined) {
      this.logger.info('Registering JSON preview document provider');
      this.provider = new JsonPreviewDocumentProvider();
      this.providerRegistration = vscode.workspace.registerTextDocumentContentProvider(
        JsonPreviewDocumentProvider.scheme,
        this.provider
      );
      this.context.subscriptions.push(this.provider, this.providerRegistration);
    }

    await this.provider.show(content, this.logger);
  }

  public dispose(): void {
    this.providerRegistration?.dispose();
    this.provider?.dispose();
  }
}

class JsonPreviewDocumentProvider implements vscode.TextDocumentContentProvider, vscode.Disposable {
  public static readonly scheme = 'json5helper-preview';
  private static readonly uri = vscode.Uri.from({
    scheme: JsonPreviewDocumentProvider.scheme,
    path: '/Json5helper Preview.json'
  });

  private content = '';
  private readonly changeEmitter = new vscode.EventEmitter<vscode.Uri>();
  public readonly onDidChange = this.changeEmitter.event;

  public provideTextDocumentContent(): string {
    return this.content;
  }

  public async show(content: string, logger: ExtensionLogger): Promise<void> {
    this.content = content;
    logger.info(`Updating preview document ${JsonPreviewDocumentProvider.uri.toString()}`);
    this.changeEmitter.fire(JsonPreviewDocumentProvider.uri);

    logger.info('Opening preview text document');
    const document = await vscode.workspace.openTextDocument(JsonPreviewDocumentProvider.uri);
    logger.info(`Preview text document opened with language ${document.languageId}`);

    const jsonDocument = document.languageId === 'json'
      ? document
      : await vscode.languages.setTextDocumentLanguage(document, 'json');
    logger.info(`Preview text document language is ${jsonDocument.languageId}`);

    logger.info('Showing preview text document');
    await vscode.window.showTextDocument(jsonDocument, {
      preview: true,
      viewColumn: vscode.ViewColumn.Beside
    });
    logger.info('Preview text document shown');
  }

  public dispose(): void {
    this.changeEmitter.dispose();
  }
}

class ExtensionLogger implements vscode.Disposable {
  private readonly logFile: vscode.Uri;
  private writeQueue: Promise<void> = Promise.resolve();

  public constructor(
    private readonly channel: vscode.LogOutputChannel,
    logUri: vscode.Uri
  ) {
    this.logFile = vscode.Uri.joinPath(logUri, 'json5helper.log');
  }

  public info(message: string): void {
    this.write('info', message);
  }

  public error(title: string, error: unknown): void {
    this.write('error', `${title}\n${this.errorToLogMessage(error)}`);
  }

  public show(): void {
    this.channel.show(true);
  }

  public async openLogDocument(): Promise<void> {
    await this.flush();
    const document = await vscode.workspace.openTextDocument(this.logFile);
    await vscode.window.showTextDocument(document, {
      preview: false,
      viewColumn: vscode.ViewColumn.Beside
    });
  }

  public async flush(): Promise<void> {
    await this.writeQueue;
  }

  public dispose(): void {
    this.channel.dispose();
  }

  private write(level: 'info' | 'error', message: string): void {
    const timestamp = new Date().toISOString();
    const entry = `[${timestamp}] [${level}] ${message}`;
    if (level === 'error') {
      this.channel.error(message);
      console.error(`[Json5helper] ${message}`);
    } else {
      this.channel.info(message);
      console.log(`[Json5helper] ${message}`);
    }

    this.writeQueue = this.writeQueue
      .then(async () => {
        await vscode.workspace.fs.createDirectory(this.logFile.with({ path: path.posix.dirname(this.logFile.path) }));
        const previous = await this.readExistingLog();
        const next = Buffer.concat([previous, Buffer.from(`${entry}\n`, 'utf8')]);
        await vscode.workspace.fs.writeFile(this.logFile, next);
      })
      .catch((error: unknown) => {
        const fallback = this.errorToLogMessage(error);
        this.channel.error(`Failed to write diagnostic log: ${fallback}`);
        console.error(`[Json5helper] Failed to write diagnostic log: ${fallback}`);
      });
  }

  private async readExistingLog(): Promise<Uint8Array> {
    try {
      return await vscode.workspace.fs.readFile(this.logFile);
    } catch {
      return new Uint8Array();
    }
  }

  private errorToLogMessage(error: unknown): string {
    if (error instanceof Error) {
      return error.stack ?? error.message;
    }
    return String(error);
  }
}
