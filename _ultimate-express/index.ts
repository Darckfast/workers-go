import { Readable } from "node:stream";
import express from "ultimate-express";
import { init } from "./load-wasm.js";

init();

const app = express();

app.use(express.raw({ type: "*/*" }))
app.all("*", async (req, res) => {
  await init();
  // @ts-ignore
  req.body = ReadableStream.from(req.body)
  const rs: Response = await cf.fetch(req);

  res.writeHead(rs.status, rs.statusText, Object.fromEntries(rs.headers));
  // @ts-ignore
  Readable.fromWeb(rs.body).pipe(res);
});

app.listen(5173, () => {
  console.log("Server is running on port 5173");
});
