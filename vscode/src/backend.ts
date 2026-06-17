export type ParseMode = 'json' | 'jsonc' | 'json5' | 'repr';

export interface ParseLensBackend {
  parsePreview(input: string, mode: ParseMode): Promise<string>;
}

type WasmModule = {
  WasmSyntax: {
    Json: unknown;
    Jsonc: unknown;
    Json5: unknown;
  };
  format_json(input: string, syntax: unknown, pretty: boolean): string;
  repr_json(input: string, pretty: boolean): string;
};

export class WasmParseLensBackend implements ParseLensBackend {
  private modulePromise: Promise<WasmModule> | undefined;

  public constructor(private readonly loadModule: () => Promise<WasmModule>) {
  }

  public async parsePreview(input: string, mode: ParseMode): Promise<string> {
    const wasm = await this.getModule();
    switch (mode) {
      case 'json':
        return wasm.format_json(input, wasm.WasmSyntax.Json, true);
      case 'jsonc':
        return wasm.format_json(input, wasm.WasmSyntax.Jsonc, true);
      case 'json5':
        return wasm.format_json(input, wasm.WasmSyntax.Json5, true);
      case 'repr':
        return wasm.repr_json(input, true);
    }
  }

  private getModule(): Promise<WasmModule> {
    this.modulePromise ??= this.loadModule();
    return this.modulePromise;
  }
}
