import React from "react"
import { auth } from "../firebaseAuth";
import { useNavigate } from "react-router-dom";
import { signOut } from "firebase/auth";
import Container from "../components/Container";
import NavBar from "../components/Navbar";
import ChallengePresentation from "../components/ChallengePresentation";

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
                <ChallengePresentation />
            </div>
        </>
    );
}