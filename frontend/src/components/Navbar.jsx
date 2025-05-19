import React, { useState } from "react";
import { signOut } from "firebase/auth";
import { auth } from "../firebaseAuth";
import { useNavigate, Link } from "react-router-dom";
import { toast } from "react-toastify";

function NavBar() {
  const navigate = useNavigate();
  const [menuOpen, setMenuOpen] = useState(false);

  const handleLogout = async () => {
    try {
      await signOut(auth);
      navigate("/LoginPage");
      toast.success("Déconnexion réussie !");
    } catch (error) {
      toast.error(error.message);
    }
  };

  return (
    <nav className="bg-gray-400 shadow-md px-6 py-4 flex justify-between items-center">
      <div className="flex items-center space-x-8">
        <Link
          to="/HomePage"
          className="text-xl font-bold text-gray-800 hover:text-blue-600 transition-colors"
        >
          ProjetPC3R
        </Link>

        {/* Liens vers les pages */}
        <div className="flex space-x-4">
          <Link
            to="/news"
            className="text-gray-800 font-medium hover:text-blue-600 transition-colors"
          >
            News
          </Link>
          <Link
            to="/challenges"
            className="text-gray-800 font-medium hover:text-blue-600 transition-colors"
          >
            Challenges
          </Link>
          <Link
            to="/discussions"
            className="text-gray-800 font-medium hover:text-blue-600 transition-colors"
          >
            Discussions
          </Link>
        </div>
      </div>

      {/* Menu utilisateur (inchangé) */}
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
            <div className="absolute right-0 mt-2 w-max bg-white border rounded shadow-md z-10 py-1">
              <button
                onClick={handleLogout}
                className="w-full whitespace-nowrap text-left px-4 py-2 text-sm font-medium text-red-700 bg-red-100 rounded-md flex items-center gap-2 transition-colors duration-200 hover:bg-red-600 hover:text-white"
              >
                Déconnexion
              </button>
            </div>
          )}
        </div>
      </div>
    </nav>
  );
}

export default NavBar;
