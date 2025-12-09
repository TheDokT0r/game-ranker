Bun.serve({
  fetch(req: Request) {
    return new Response("Success!");
  },
});
