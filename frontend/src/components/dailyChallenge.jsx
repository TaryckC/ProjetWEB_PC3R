import React, { useEffect, useState } from "react";
import { doc, getDoc } from "firebase/firestore";
import { db } from "../firebaseAuth";
import Container from "./Container";

async function fetchDailyChallenge() {
  const docRef = doc(db, "Challenges", "daily_challenge");
  const docSnap = await getDoc(docRef);
  if (docSnap.exists()) {
    const fullData = docSnap.data()
    console.log("challenge récupéré :", fullData);
    return fullData.data.activeDailyCodingChallengeQuestion.question;
  } else {
    return {};
  }
}

export default function DailyChallenge() {
  const [challenge, setChallenge] = useState({});

  useEffect(() => {
    fetchDailyChallenge()
      .then((data) => {
        console.log("challenge récupéré :", data);
        setChallenge(data);
      })
      .catch(console.error);
  }, []);

  return (
    <div className="flex flex-row justify-center items-start gap-6 px-4 mt-8">
      <Container bgColor="bg-gray-100">
        <div className="bg-white p-4 rounded shadow-md w-full max-w-xl">
          <h1 className="text-3xl font-semibold text-gray-800 mb-6 text-center">
            Daily Challenge !
          </h1>
          <h2 className="text-2xl font-semibold text-gray-800 mb-2">
            {challenge.title}
          </h2>
          <p className="text-sm text-gray-500 mb-1">
            Difficulté : {challenge.difficulty || "N/A"}
          </p>
          <p className="text-sm text-gray-500 mb-4">
            Taux de réussite : {challenge.acRate?.toFixed(2) || "N/A"}%
          </p>
          <button className="bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700 transition">
            Voir le challenge
          </button>
        </div>
      </Container>

      <Container bgColor="bg-gray-200">
        <button
          type="button"
          onClick={() => {
            import("../firebaseAuth").then(({ auth }) =>
              import("firebase/auth").then(({ signOut }) =>
                signOut(auth)
                  .then(() => {
                    window.location.href = "/LoginPage";
                    alert("Déconnexion réussie !");
                  })
                  .catch((err) => alert(err.message))
              )
            );
          }}
          className="bg-red-600 text-white px-6 py-2 rounded hover:bg-red-700 transition focus:outline-none focus:ring-2 focus:ring-red-400"
        >
          Déconnexion
        </button>
      </Container>
    </div>
  );
}
