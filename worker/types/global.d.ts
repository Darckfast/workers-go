declare global {
  // interface Global {
  var cf: {
    fetch(r: Request, e: Env, ctx: ExecutionContext): Promise<Response>;
    email(
      m: ForwardableEmailMessage,
      e: Env,
      c: ExecutionContext,
    ): Promise<void>;
    scheduled(
      c: ScheduledController,
      e: Env,
      ctx: ExecutionContext,
    ): Promise<void>;
    queue(b: MessageBatch, e: Env, ctx: ExecutionContext): Promise<void>;
    tail(t: TraceItem[], e: Env, c: ExecutionContext): Promise<void>;
  };
  // }
}
export { };
