import { Client } from 'pg';

export async function initTable() {
    const client = new Client();
    await client.connect();

    const sqlStmt = `
      CREATE TABLE IF NOT EXISTS games (
		    id SERIAL PRIMARY KEY,
		    public_id UUID NOT NULL UNIQUE DEFAULT gen_random_uuid(),
		    name TEXT NOT NULL,
		    release_date DATE,
		    cover_art_url TEXT
	    );
    `;

    await client.query("CREATE EXTENSION IF NOT EXISTS pgcrypto;");
    await client.query(sqlStmt);

    await client.end();
}