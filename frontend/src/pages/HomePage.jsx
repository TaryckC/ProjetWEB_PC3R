import React from "react"
import { auth } from "../firebaseAuth";
import { useNavigate } from "react-router-dom";
import { signOut } from "firebase/auth";
import Container from "../components/Container";
import NavBar from "../components/Navbar";

export default function Home() {
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
        <>
            <NavBar />
            <div className="flex flex-row justify-center items-start gap-6 px-4 mt-8">
                <Container bgColor="bg-gray-100">
                    <h1 className="text-3xl font-semibold text-gray-800 mb-6 text-center">
                        Bienvenue à la page d'accueil !
                    </h1>
                </Container>

                <Container bgColor="bg-gray-200">
                    <button
                        type="button"
                        onClick={handleLogout}
                        className="bg-red-600 text-white px-6 py-2 rounded hover:bg-red-700 transition focus:outline-none focus:ring-2 focus:ring-red-400"
                    >
                        Déconnexion
                    </button>
                </Container>
            </div>
        </>
    );
}