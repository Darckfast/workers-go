declare global {
  var cf: {
    fetch(r: Request): Promise<Response>;
  };
}
export { };
