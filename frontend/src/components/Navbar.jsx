import React, { useState } from "react";
import { signOut } from "firebase/auth";
import { auth } from "../firebaseAuth";
import { useNavigate, Link } from "react-router-dom";

function NavBar() {
  const navigate = useNavigate();
  const [menuOpen, setMenuOpen] = useState(false);

  const handleLogout = async () => {
    try {
      await signOut(auth);
      navigate("/LoginPage");
      alert("Déconnexion réussie !");
    } catch (error) {
      alert(error.message);
    }
  };

  const links = [
    { path: "/news", label: "News" },
    { path: "/discussions", label: "Discussions" },
    { path: "/challenges", label: "Challenges" },
  ];

  return (
    <nav className="bg-white shadow-md px-6 py-4 flex justify-between items-center">
      <Link
        to="/"
        className="text-xl font-bold text-gray-800 hover:text-blue-600 transition-colors"
      >
        ProjetPC3R
      </Link>

      <ul className="flex space-x-6">
        {links.map(({ path, label }) => (
          <li key={path}>
            <Link
              to={path}
              className="text-gray-700 hover:text-blue-600 transition-colors"
            >
              {label}
            </Link>
          </li>
        ))}
      </ul>

      <div className="flex items-center space-x-4">
        <div className="relative">
          <button
            onClick={() => setMenuOpen(!menuOpen)}
            className="w-8 h-8 rounded-full bg-gray-300 flex items-center justify-center text-sm font-semibold text-white"
            title="Menu utilisateur"
          >
            U
          </button>
          {menuOpen && (
            <div className="absolute right-0 mt-2 min-w-[150px] bg-white border rounded shadow-md z-10 py-1">
              <button
                onClick={handleLogout}
                className="w-full text-left px-4 py-2 text-sm text-red-600 hover:bg-red-50 flex items-center gap-2"
              >
                🔓 Se déconnecter
              </button>
            </div>
          )}
        </div>
      </div>
    </nav>
  );
}

export default NavBar;
