import type { Buffer } from "node:buffer";
import { readFileSync } from "node:fs";
import { dirname, resolve } from "node:path";
import { Readable } from "node:stream";
import { ReadableStream } from "node:stream/web";
import { fileURLToPath } from "node:url";
import express from "ultimate-express";
import "./bin/wasm_exec.js";

// @ts-expect-error
globalThis.tryCatch = (fn: () => any) => {
  try {
    return { data: fn() };
  } catch (err) {
    if (!(err instanceof Error)) {
      if (err instanceof Object) {
        return { error: JSON.stringify(err) };
      }
    }

    return { error: err };
  }
};

const CURRENT_DIR = dirname(fileURLToPath(import.meta.url));
const binary = readFileSync(resolve(CURRENT_DIR, "./bin/app.wasm"));
const go = new Go();
let initiliazed = false;

export async function init() {
  if (go.exited) {
    initiliazed = false;
  }

  if (!initiliazed) {
    const { instance } = await WebAssembly.instantiate(binary, go.importObject);

    go.run(instance).finally(() => {
      initiliazed = false;
    });
    initiliazed = true;
  }
}

// warm up the init
init();

const app = express();

function bufferToUint8Array(buffer: Buffer) {
  const uintArray = new Uint8Array(buffer.length);
  for (let i = 0; i < buffer.length; ++i) {
    uintArray[i] = buffer[i];
  }

  return uintArray;
}

// the body must be raw, the parsing will happen within
// our Go code
app.use(express.raw({ type: "*/*" }));

app.all("*", async (req, res) => {
  // if the first request comes before the initialization
  // is done, or if the Go process exit, this will guarantee
  // it's up and running
  await init();

  // the express.raw gives us a buffer, but the lib expected a
  // ReadableStream of UInt8Array
  req.body = ReadableStream.from([bufferToUint8Array(req.body)]);

  const gores = await cf.fetch(req);

  res.writeHead(
    gores.status,
    gores.statusText,
    Object.fromEntries(gores.headers),
  );

  if (gores.body !== null) {
    Readable.fromWeb(gores.body as ReadableStream<any>).pipe(res);
  }
});

app.listen(5173, () => {
  console.log("Server is running on port 5173");
});
