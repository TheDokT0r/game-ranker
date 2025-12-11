import ky from "ky";
import Fastify from "fastify";
import "@dotenvx/dotenvx/config";
import { initTable, searchGames } from "./database.js";

const fastify = Fastify({
  logger: true,
  routerOptions: {
    ignoreTrailingSlash: true,
  }
});

fastify.get("/", async function handler(reqest, reply) {
  return { hello: "world" };
});

fastify.get<{ Querystring: { name: string } }>("/game-data", async (request, reply) => {
  const { name } = request.query;
  const games = await searchGames(name);
  return { games };
});

console.log("generating db...");
await initTable();
console.log("finished!");

const port = Number(process.env.PORT);

try {
  await fastify.listen({ port });
} catch (err) {
  fastify.log.error(err);
  process.exit(1);
}

