import ky from "ky";
import Fastify from "fastify";
import "@dotenvx/dotenvx/config";
import { initTable } from "./database.js";

async function getAccessToken() {
  const clientId = process.env.TWITCH_CLIENT_ID;
  const clientSecret = process.env.TWITCH_SECRET;

  if (!clientId || !clientSecret) {
    throw new Error("Missing Twitch credentials in environment variables");
  }

  const response = await ky.post(
    `https://id.twitch.tv/oauth2/token?client_id=${clientId}&client_secret=${clientSecret}&grant_type=client_credentials`
  );

  const data = await response.json<{ access_token: string }>();
  return data.access_token;
}

const token = await getAccessToken();
console.log("Got access token")

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
  const clientId = process.env.TWITCH_CLIENT_ID;

  const igdbResponse = await ky.post("https://api.igdb.com/v4/games", {
    headers: {
      "Client-ID": clientId!,
      Authorization: `Bearer ${token}`,
      "Content-Type": "text/plain",
    },
    body: `
      fields name, first_release_date, involved_companies.company.name, cover.url;
      search "${name}";
      limit 5;
    `,
  }).json();

  const games = (igdbResponse as any[]).map(game => ({
    name: game.name,
    releaseDate: game.first_release_date ? new Date(game.first_release_date * 1000).toDateString() : null,
    publisher: game.involved_companies?.[0]?.company?.name || null,
    coverUrl: game.cover?.url ? `https:${game.cover.url}` : null,
  }));

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

