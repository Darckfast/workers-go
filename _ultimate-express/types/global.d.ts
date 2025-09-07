import type { Request } from "ultimate-express";

declare global {
  var cf: {
    fetch(r: Request<any>): Promise<Response>;
  };
}
