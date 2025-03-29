import { useState } from "react";
import { signInWithEmailAndPassword, signOut, createUserWithEmailAndPassword, sendEmailVerification } from "firebase/auth";
import { auth } from "../firebaseAuth";
import { useNavigate } from "react-router-dom";

export default function SignUP() {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const navigate = useNavigate();


  const handleRegister = async (e) => {
    e.preventDefault();
    try {
      await createUserWithEmailAndPassword(auth, email, password).then((currentUser)=>{
         sendEmailVerification(currentUser.user);
         navigate("/MailConfirmationPage");
      });
      alert("Inscription réussie !");
    } catch (error) {
      alert(error.message);
    }
  };



  return (
    <div className="container">
      <h2>Inscription</h2>
      <form className="form" onSubmit={handleRegister}>
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
        <button type="button" className="button-outline" onClick={handleRegister}>
          S’inscrire
        </button>
      </form>
    </div>
  );
}
