"use client"

import { useState, FormEvent } from "react"

export default function Search() {
    const [gameName, setGameName] = useState("");

    const onSubmit = (e: FormEvent<HTMLFormElement>) => {
        e.preventDefault();
    }

    return (
        <div>
            <form onSubmit={onSubmit}>
                <input type="text" lang="en" placeholder="Name of game" value={gameName} onChange={e => setGameName(e.target.value)}/>
                <button>Submit</button>
            </form>
        </div>
    )
}
