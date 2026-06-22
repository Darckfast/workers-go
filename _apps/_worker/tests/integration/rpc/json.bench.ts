import { bench } from "vitest";

const payload = JSON.stringify(
  Array(500).fill({
    bool: true,
    number: 1,
    string: "my super string",
    list: [1, 2, 3, 4, 5, 6, 7],
  }),
);

bench(
  "http",
  async () => {
    const rs = await fetch("http://localhost:5173/echo", {
      body: payload,
      method: "POST",
      headers: {
        "content-type": "application/json",
      },
    });

    await rs.text();
  },
  {
    iterations: 100,
    time: 5000,
  },
);

bench(
  "rpc-uintarray",
  async () => {
    const rs = await fetch("http://localhost:5173/rpc", {
      body: payload,
      method: "POST",
      headers: {
        "content-type": "application/json",
      },
    });

    await rs.text();
  },
  {
    iterations: 100,
    time: 5000,
  },
);
