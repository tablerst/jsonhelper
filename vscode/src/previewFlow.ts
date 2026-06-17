import { ParseLensBackend, ParseMode } from './backend';

export type ParsePreviewInput = {
  selection: string;
  document: string;
};

export type ParsePreviewServices = {
  getInput(): ParsePreviewInput | undefined;
  openJsonPreview(content: string): Promise<void>;
  showError(message: string): void;
  logInfo?(message: string): void;
  logError(title: string, message: string): void;
};

export async function parsePreview(
  mode: ParseMode,
  backend: ParseLensBackend,
  services: ParsePreviewServices
): Promise<void> {
  const input = services.getInput();
  if (input === undefined) {
    services.logInfo?.(`${mode} parse preview skipped: no active editor`);
    services.showError('ParseLens: open a document before parsing.');
    return;
  }

  const source = getSelectedTextOrDocument(input);
  if (source.length === 0) {
    services.logInfo?.(`${mode} parse preview skipped: empty source`);
    services.showError('ParseLens: selected text or document is empty.');
    return;
  }

  services.logInfo?.(`${mode} parse preview started (${source.length} source chars)`);

  try {
    const output = await backend.parsePreview(source, mode);
    services.logInfo?.(`${mode} parse preview parsed (${output.length} output chars)`);

    const content = ensureTrailingNewline(output);
    await services.openJsonPreview(content);
    services.logInfo?.(`${mode} parse preview opened`);
  } catch (error) {
    const message = errorToLogMessage(error);
    services.logError(`${mode} parse preview failed`, message);
    services.showError('ParseLens: parse preview failed. See output for details.');
  }
}

export function getSelectedTextOrDocument(input: ParsePreviewInput): string {
  if (input.selection.length > 0) {
    return input.selection;
  }
  return input.document;
}

export function ensureTrailingNewline(value: string): string {
  return value.endsWith('\n') ? value : `${value}\n`;
}

function errorToLogMessage(error: unknown): string {
  if (error instanceof Error) {
    return error.stack ?? error.message;
  }
  return String(error);
}
