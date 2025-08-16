import { DurableObject } from "cloudflare:workers";

export class TestDurableObject extends DurableObject<Env> {
  constructor(ctx: DurableObjectState, env: Env) {
    // Required, as we're extending the base class.
    super(ctx, env)
  }

  async sayHello(): Promise<string> {
    const result = this.ctx.storage.sql
      .exec("SELECT 'Hello, World!' as greeting")
      .one();

    return result.greeting.toString();
  }
}
