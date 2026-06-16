import { Json5helperBackend, ParseMode } from './backend';

export type ParsePreviewInput = {
  selection: string;
  document: string;
};

export type ParsePreviewServices = {
  chooseMode(): Promise<ParseMode | undefined>;
  getInput(): ParsePreviewInput | undefined;
  openJsonPreview(content: string): Promise<void>;
  showError(message: string): void;
  logError(title: string, message: string): void;
};

export async function parsePreview(
  backend: Json5helperBackend,
  services: ParsePreviewServices
): Promise<void> {
  const input = services.getInput();
  if (input === undefined) {
    services.showError('Json5helper: open a document before parsing.');
    return;
  }

  const mode = await services.chooseMode();
  if (mode === undefined) {
    return;
  }

  const source = getSelectedTextOrDocument(input);
  if (source.length === 0) {
    services.showError('Json5helper: selected text or document is empty.');
    return;
  }

  try {
    const output = await backend.parsePreview(source, mode);
    await services.openJsonPreview(ensureTrailingNewline(output));
  } catch (error) {
    const message = error instanceof Error ? error.message : String(error);
    services.logError(`${mode} parse preview failed`, message);
    services.showError('Json5helper: parse preview failed. See output for details.');
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
