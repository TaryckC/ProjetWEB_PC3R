import React from "react"
import { auth } from "../firebaseAuth";
import { useNavigate } from "react-router-dom";

import { signOut } from "firebase/auth";

export default function Home(){
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
        <div>
            <h1>Bienvenue a la page d'accueil!</h1>
            <button type="button" onClick={handleLogout}>Déconnexion</button>
        </div>

    )
}