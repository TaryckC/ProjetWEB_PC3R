import React from "react";
import {onAuthStateChanged, sendEmailVerification} from  "firebase/auth";
import { auth } from "../firebaseAuth";
import { useNavigate } from "react-router-dom";
import { useEffect } from "react";


export default function MailConf(){

    const navigate = useNavigate();
    useEffect(() => {
        const checkVerification = async () => {
          if (auth.currentUser) {
            await auth.currentUser.reload();
            if (auth.currentUser.emailVerified) {
              console.log("Passage à la page suivante");
              navigate("/HomePage");
            }
          }
        };
      
        checkVerification();
    }, []);
      


    
    
    const handleReSend = async ()=>{
        await sendEmailVerification(auth.currentUser);
    }

    const handleVerification = async () => {
        await auth.currentUser.reload();
        if (auth.currentUser.emailVerified)
            navigate("/HomePage");
        else 
            alert("L'adresse e-mail n’est pas encore vérifiée.");
            
    }


    return (
        <div style={{ textAlign: "center", padding: "2rem" }}>
          <h1>Veuillez confirmer votre mail pour pouvoir vous connecter</h1>
          <p>Un lien de vérification a été envoyé à votre adresse e-mail.</p>
    
          <button onClick={handleReSend} style={{ marginRight: "1rem" }}>
            Renvoyer
          </button>
          <button onClick={handleVerification}>
            J'ai vérifié mon mail
          </button>
        </div>
    );
}