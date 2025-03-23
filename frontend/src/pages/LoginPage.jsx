import { useState } from "react";
import { signInWithEmailAndPassword, signOut, createUserWithEmailAndPassword } from "firebase/auth";
import { auth } from "../firebaseAuth";
import { useNavigate } from "react-router-dom";
import "./LoginPage.css"


export default function Login() {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const navigate = useNavigate();

  const handleLogin = async (e) => {
    e.preventDefault();
    try {
      await signInWithEmailAndPassword(auth, email, password).then((currentUser)=>{
        alert("Connexion réussie !");
        navigate("/HomePage");
      });
      
    } catch (error) {
      alert(error.message);
    }
  };

  const handleRegister = async (e) => {
    e.preventDefault();
    try {
      await createUserWithEmailAndPassword(auth, email, password);
      alert("Inscription réussie !");
    } catch (error) {
      alert(error.message);
    }
  };



  return (
    <div className="container">
      <h2>Connexion</h2>
      <form className="form" onSubmit={handleLogin}>
        <input
          type="email"
          className="input"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          placeholder="Email"
        />
        <input
          type="password"
          className="input"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          placeholder="Mot de passe"
        />
        <button type="submit" className="button">
          Se connecter
        </button>
        <button type="button" className="button-outline" onClick={handleRegister}>
          S’inscrire
        </button>
      </form>
    </div>
  );
}
