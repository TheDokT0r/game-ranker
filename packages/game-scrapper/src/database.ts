import { Pool } from 'pg';
import { Game, getGameDataFromIgdb } from './api.js';

export async function initTable() {
    const pool = new Pool();

    const sqlStmt = `
      CREATE TABLE IF NOT EXISTS games (
		    id SERIAL PRIMARY KEY,
		    public_id UUID NOT NULL UNIQUE DEFAULT gen_random_uuid(),
		    name TEXT NOT NULL,
        publisher TEXT NOT NULL,
		    release_date DATE,
		    cover_url TEXT
	    );
    `;

    await pool.query("CREATE EXTENSION IF NOT EXISTS pgcrypto;");
    await pool.query(sqlStmt);
}

export async function addGameToDb(games: Game[]): Promise<void> {
  const pool = new Pool();
  
  games.forEach(async (game) => {
    const {name, releaseDate, publisher, coverUrl} = game;
    const sqlStmt = `
      INSERT INTO games (name, release_date, publisher, cover_url)
      VALUES ($1, $2, $3, $4)
    `;

    await pool.query(sqlStmt, [name, releaseDate, publisher, coverUrl]);
  });

}

export async function searchGames(name: string): Promise<Game[]> {
  const pool = new Pool();

  const result = await pool.query<Game>(`
      SELECT *
      FROM games
      WHERE name ILIKE $1
      LIMIT 10;
    `, [`%${name}`]);

    if (result.rowCount === 0) {
      const idgbGames = await getGameDataFromIgdb(name);
      if (idgbGames.length === 0) {
        return [];
      }

      await addGameToDb(idgbGames);
      return idgbGames
    }

    return result.rows;
}
