declare global {
    class Go {
        importObject: WebAssembly.Imports;
        run(instance: WebAssembly.Instance): void;
    }
}

export default global;