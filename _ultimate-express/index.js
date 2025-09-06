const { init } = require("./load-wasm");
const express = require("ultimate-express");

init();

const app = express();

app.all("*", async (r) => {
  console.log(r.header, r.headers);
  await init();
  return cf.fetch(r);
});

app.listen(5173, () => {
  console.log("Server is running on port 5173");
});
