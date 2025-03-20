import React from "react"
import "../styles.css"

function NavBar() {
    return (
        <nav className="nav">
            <a className="site-title">ProjetPC3R</a>
            <ul>
                <li>
                    <a href="../pages/news">News</a>
                </li>
                <li>
                    <a href="../pages/discussions">Discussions</a>
                </li>
                <li>
                    <a href="../pages/challenges">Challenges</a>
                </li>
            </ul>
        </nav>
    );
}

export default NavBar
