import React from "react";
import { signOut } from "firebase/auth";
import { auth } from "../firebaseAuth";
import { useNavigate } from "react-router-dom";

function NavBar() {
    const navigate = useNavigate();

    const handleLogout = async () => {
        try {
            await signOut(auth);
            navigate("/LoginPage");
            alert("Déconnexion réussie !");
        } catch (error) {
            alert(error.message);
        }
    };

    return (
        <nav className="bg-white shadow-md px-6 py-4 flex justify-between items-center">
            <a href="/" className="text-xl font-bold text-gray-800 hover:text-blue-600 transition-colors">ProjetPC3R</a>
            <ul className="flex space-x-6">
                <li>
                    <a href="../pages/news" className="text-gray-700 hover:text-blue-600 transition-colors">News</a>
                </li>
                <li>
                    <a href="../pages/discussions" className="text-gray-700 hover:text-blue-600 transition-colors">Discussions</a>
                </li>
                <li>
                    <a href="../pages/challenges" className="text-gray-700 hover:text-blue-600 transition-colors">Challenges</a>
                </li>
            </ul>
            <div className="flex items-center space-x-4">
                <button
                    className="text-sm text-red-600 bg-red-300 hover:bg-red-600 hover:text-white transition focus:outline-none rounded-xl px-3 py-2"
                    onClick={handleLogout}
                >
                    Déconnexion
                </button>
                <div className="w-8 h-8 rounded-full bg-gray-300 flex items-center justify-center text-sm font-semibold text-white">
                    U
                </div>
            </div>
        </nav>
    );
}

export default NavBar;
