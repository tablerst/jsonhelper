import assert from 'node:assert/strict';

import {
  ParsePreviewInput,
  ParsePreviewServices,
  ensureTrailingNewline,
  getSelectedTextOrDocument,
  parsePreview
} from '../src/previewFlow';
import { Json5helperBackend, ParseMode } from '../src/backend';

class RecordingBackend implements Json5helperBackend {
  public calls: Array<{ input: string; mode: ParseMode }> = [];
  public result = '{"ok":true}';
  public error: Error | undefined;

  public async parsePreview(input: string, mode: ParseMode): Promise<string> {
    this.calls.push({ input, mode });
    if (this.error !== undefined) {
      throw this.error;
    }
    return this.result;
  }
}

class RecordingServices implements ParsePreviewServices {
  public input: ParsePreviewInput | undefined = {
    selection: '',
    document: '{unquoted: true}'
  };
  public previews: string[] = [];
  public errors: string[] = [];
  public logs: Array<{ title: string; message: string }> = [];

  public getInput(): ParsePreviewInput | undefined {
    return this.input;
  }

  public async openJsonPreview(content: string): Promise<void> {
    this.previews.push(content);
  }

  public showError(message: string): void {
    this.errors.push(message);
  }

  public logError(title: string, message: string): void {
    this.logs.push({ title, message });
  }
}

async function run(): Promise<void> {
  assert.equal(ensureTrailingNewline('{"ok":true}'), '{"ok":true}\n');
  assert.equal(ensureTrailingNewline('{"ok":true}\n'), '{"ok":true}\n');

  assert.equal(
    getSelectedTextOrDocument({ selection: '{selected: true}', document: '{document: true}' }),
    '{selected: true}'
  );
  assert.equal(
    getSelectedTextOrDocument({ selection: '', document: '{document: true}' }),
    '{document: true}'
  );

  const backend = new RecordingBackend();
  const services = new RecordingServices();
  services.input = { selection: '{selected: true}', document: '{document: true}' };
  await parsePreview('json5', backend, services);

  assert.deepEqual(backend.calls, [{ input: '{selected: true}', mode: 'json5' }]);
  assert.deepEqual(services.previews, ['{"ok":true}\n']);
  assert.deepEqual(services.errors, []);

  const wholeDocumentBackend = new RecordingBackend();
  const wholeDocumentServices = new RecordingServices();
  wholeDocumentServices.input = { selection: '', document: '{document: true}' };
  await parsePreview('jsonc', wholeDocumentBackend, wholeDocumentServices);

  assert.deepEqual(wholeDocumentBackend.calls, [{ input: '{document: true}', mode: 'jsonc' }]);

  const failingBackend = new RecordingBackend();
  failingBackend.error = new Error('bad syntax');
  const failingServices = new RecordingServices();
  await parsePreview('json5', failingBackend, failingServices);

  assert.deepEqual(failingServices.previews, []);
  assert.equal(failingServices.errors[0], 'Json5helper: parse preview failed. See output for details.');
  assert.equal(failingServices.logs[0].message, 'bad syntax');

  const emptyBackend = new RecordingBackend();
  const emptyServices = new RecordingServices();
  emptyServices.input = { selection: '', document: '' };
  await parsePreview('repr', emptyBackend, emptyServices);

  assert.deepEqual(emptyBackend.calls, []);
  assert.equal(emptyServices.errors[0], 'Json5helper: selected text or document is empty.');
}

run().catch((error: unknown) => {
  console.error(error);
  process.exitCode = 1;
});
