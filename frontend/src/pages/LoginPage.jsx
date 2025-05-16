import { useState, useEffect } from "react";
import {
  signInWithEmailAndPassword,
  signOut,
  createUserWithEmailAndPassword,
} from "firebase/auth";
import { toast } from "react-toastify";
import { useNavigate } from "react-router-dom";
import {
  signInWithPopup,
  GoogleAuthProvider,
  setPersistence,
  browserLocalPersistence,
} from "firebase/auth";
import { provider, auth } from "../firebaseAuth";

export default function Login() {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const navigate = useNavigate();
  const [loading, setLoading] = useState(false);

  const handleLogin = async (e) => {
    e.preventDefault();
    setLoading(true);
    try {
      await setPersistence(auth, browserLocalPersistence);
      await signInWithEmailAndPassword(auth, email, password);
      toast.success("Connexion réussie.");
      navigate("/HomePage");
    } catch (error) {
      toast.error(error.message);
    } finally {
      setLoading(false);
    }
  };

  const handleRegister = () => {
    navigate("/SignUpPage");
  };

  const handleResetPassword = () => {
    navigate("/PasswordResetPage");
  };

  const handleGoogleSignIn = () => {
    signInWithPopup(auth, provider)
      .then((result) => {
        const credential = GoogleAuthProvider.credentialFromResult(result);
        const token = credential.accessToken;
        const user = result.user;

        toast("Connexion avec Google réussie !");
        navigate("/HomePage");
      })
      .catch((error) => {
        const errorCode = error.code;
        const errorMessage = error.message;
        const email = error.customData.email;
        const credential = GoogleAuthProvider.credentialFromError(error);
      });
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-100 px-4">
      <form
        onSubmit={handleLogin}
        className="bg-white w-full max-w-sm p-8 rounded-xl shadow-md"
        aria-label="Formulaire de connexion"
      >
        <h2 className="text-2xl font-semibold text-gray-800 mb-6 text-center">
          Connexion
        </h2>

        <div className="mb-4">
          <label
            htmlFor="email"
            className="block text-sm font-medium text-gray-700 mb-1"
          >
            Email
          </label>
          <input
            id="email"
            type="email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            className="w-full px-4 py-2 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-blue-500"
            placeholder="ex: moi@email.com"
            autoComplete="email"
            required
          />
        </div>

        <div className="mb-4">
          <label
            htmlFor="password"
            className="block text-sm font-medium text-gray-700 mb-1"
          >
            Mot de passe
          </label>
          <input
            id="password"
            type="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            className="w-full px-4 py-2 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-blue-500"
            placeholder="••••••••"
            autoComplete="current-password"
            required
          />
        </div>

        <button
          type="submit"
          disabled={loading}
          className={`w-full text-white py-2 px-4 rounded transition mb-3 
    ${
      loading
        ? "bg-blue-300 cursor-not-allowed"
        : "bg-blue-600 hover:bg-blue-700"
    }
    focus:outline-none focus:ring-2 focus:ring-blue-500`}
        >
          {loading ? "Connexion..." : "Se connecter"}
        </button>

        <button
          type="button"
          onClick={handleGoogleSignIn}
          className="w-full bg-white border border-gray-300 text-gray-700 py-2 px-4 rounded hover:bg-gray-100 transition mb-3 flex items-center justify-center gap-2"
        >
          <svg className="w-5 h-5" aria-hidden="true" viewBox="0 0 24 24">
            {/* icône Google simplifiée ici */}
          </svg>
          Se connecter avec Google
        </button>

        <button
          type="button"
          onClick={handleRegister}
          className="w-full text-blue-600 py-2 px-4 rounded hover:underline text-sm font-medium text-center"
        >
          S’inscrire
        </button>

        <p
          onClick={handleResetPassword}
          className="mt-3 text-sm text-center text-gray-500 hover:text-blue-600 hover:underline cursor-pointer"
        >
          Mot de passe oublié ?
        </p>
      </form>
    </div>
  );
}
