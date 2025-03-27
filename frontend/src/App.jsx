import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import { useState, useEffect } from "react";
import { onAuthStateChanged } from "firebase/auth";
import { auth } from "./firebaseAuth";
import MailConf from "./pages/MailConfirmationPage";
import Home from './pages/HomePage';
import Login from './pages/LoginPage';
import PrivateRoute from "./components/PrivateRoute";
import SignUP from "./pages/SignUpPage";


function App() {
  const [user, setUser] = useState(null);

  useEffect(() => {
    const unsub = onAuthStateChanged(auth, (currentUser) => {
      setUser(currentUser);
    });

    return () => unsub();
  }, []);

  return (
    <Router>
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
      </Routes>
    </Router>
  );
}

export default App;
// code de route derivé de stackoverflow
