import { readFileSync } from "fs";
import { dirname, resolve } from "path";
import { fileURLToPath } from "url";
import "./bin/wasm_exec";

globalThis.tryCatch = (fn) => {
    try {
        return { data: fn() };
    } catch (err) {
        if (!(err instanceof Error)) {
            if (err instanceof Object) {
                err = JSON.stringify(err);
            }

            err = new Error(err || "no error message");
        }

        return { error: err };
    }
};

let go = new Go();
let initiliazed = false;

export async function init() {
    if (go.exited) {
        initiliazed = false
    }

    if (!initiliazed) {
        const CURRENT_DIR = dirname(fileURLToPath(import.meta.url));
        const app = readFileSync(resolve(CURRENT_DIR, './bin/app.wasm'));
        let { instance } = await WebAssembly.instantiate(app, go.importObject);
        go.run(instance).finally(() => {
            initiliazed = false
        });
        initiliazed = true;
    }

}

