import React, { useState } from "react";
import "../css/LoginPage.css";
import { sendPasswordResetEmail } from "firebase/auth";
import { auth } from "../firebaseAuth";
import { useNavigate } from "react-router-dom";


export default function ResetPassword() {
  const [email, setEmail] = useState("");

  const navigate = useNavigate();
  const handleReset = async (e) => {
    e.preventDefault();
    try {
      await sendPasswordResetEmail(auth, email);
      alert("Email de réinitialisation envoyé !");
      navigate("/HomePage");
    } catch (error) {
      alert(error.message);
    }
  };

  return (
    <div className="container">
      <h2>Réinitialisation du mot de passe</h2>
      <form className="form" onSubmit={handleReset}>
        <input
          type="email"
          className="input"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          placeholder="Votre adresse email"
        />
        <button type="submit" className="button">
          Envoyer l’e-mail de réinitialisation
        </button>
      </form>
    </div>
  );
}
