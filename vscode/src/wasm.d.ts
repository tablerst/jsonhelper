declare module '*parselens_wasm.js' {
  export enum WasmSyntax {
    Json = 0,
    Jsonc = 1,
    Json5 = 2
  }

  export function format_json(input: string, syntax: WasmSyntax, pretty: boolean): string;
  export function repr_json(input: string, pretty: boolean): string;
}
