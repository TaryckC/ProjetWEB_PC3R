import { useState, useEffect } from "react";
import { signInWithEmailAndPassword, signOut, createUserWithEmailAndPassword } from "firebase/auth";

import { useNavigate } from "react-router-dom";
import "../css/LoginPage.css"
import { signInWithPopup, GoogleAuthProvider } from "firebase/auth";
import { provider,auth } from "../firebaseAuth";

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

  const handleRegister = () => {
    navigate("/SignUpPage");
  };

  const handleResetPassword = () => {
    navigate("/PasswordResetPage")
  };

  const handleGoogleSignIn=()=> {
    signInWithPopup(auth, provider).then((result) => {
    // This gives you a Google Access Token. You can use it to access the Google API.
    const credential = GoogleAuthProvider.credentialFromResult(result);
    const token = credential.accessToken;
    // The signed-in user info.
    const user = result.user;
    
    alert("Connexion avec Google réussie !");
    navigate("/HomePage");
  }).catch((error) => {
    // Handle Errors here.
    const errorCode = error.code;
    const errorMessage = error.message;
    // The email of the user's account used.
    const email = error.customData.email;
    // The AuthCredential type that was used.
    const credential = GoogleAuthProvider.credentialFromError(error);
    // ...
  });
  }


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
        <button type="button" className="button" onClick={handleGoogleSignIn}>
          Se connecter avec google
        </button>
        <button type="submit" className="button">
          Se connecter
        </button>
        <button type="button" className="button-outline" onClick={handleRegister}>
          S’inscrire
        </button>
        <p onClick={handleResetPassword}>Mots de passe oublié?</p>
      </form>
    </div>
  );
}
