import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import { useState, useEffect } from "react";
import { onAuthStateChanged } from "firebase/auth";
import { auth } from "./firebaseAuth";
import MailConf from "./pages/MailConfirmationPage";
import Home from "./pages/HomePage";
import Login from "./pages/LoginPage";
import PrivateRoute from "./components/PrivateRoute";
import SignUP from "./pages/SignUpPage";
import ResetPassword from "./pages/PasswordResetPage";
import CodeEditor from "./pages/testBackend";
import IDE from "./pages/IdePage";
import { ToastContainer } from "react-toastify";
import "react-toastify/dist/ReactToastify.css";
import IdePage from "./pages/IdePage";
function App() {
  const [user, setUser] = useState(null);
  const [authChecked, setAuthChecked] = useState(false);

  useEffect(() => {
    const unsub = onAuthStateChanged(auth, (currentUser) => {
      setUser(currentUser);
      setAuthChecked(true); // ✅ maintenant on sait si on est connecté ou pas
    });

    return () => unsub();
  }, []);

  if (!authChecked) {
    return <h1>Chargement de la page d'acceuil...</h1>;
  }

  return (
    <Router basename="/ProjetWEB_PC3R">
      <ToastContainer position="bottom-center" autoClose={3000} />
      <Routes>
        <Route path="/" element={<Login />} />
        <Route
          path="/HomePage"
          element={
            <PrivateRoute user={user}>
              <Home />
            </PrivateRoute>
          }
        />
        <Route path="/MailConfirmationPage" element={<MailConf />} />
        <Route path="/SignUpPage" element={<SignUP />} />
        <Route path="/PasswordResetPage" element={<ResetPassword />} />
        <Route path="/LoginPage" element={<Login />} />
        <Route path="/ide/:id" element={<IdePage />} />
        <Route path="/news" element={<NewsPage />} />
      </Routes>
    </Router>
  );
}

export default App;
// code de route derivé de stackoverflow
