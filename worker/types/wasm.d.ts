declare module "*.wasm" {
  interface WasmExports { }
  const wasmModule: (
    imports?: WebAssembly.Imports,
  ) => Promise<{ instance: WebAssembly.Instance & { exports: WasmExports } }>;
  export default wasmModule;
}
