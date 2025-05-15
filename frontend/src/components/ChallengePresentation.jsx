import React, { useEffect, useState } from "react";
import Container from "./Container";

const CHALLENGE_TYPES = {
  DAILY: {
    key: "daily_challenge",
    id: ["0"]
  },
  CLASSIC: {
    key: "classic_challenges",
    ids: ["0", "1", "2", "3", "4", "5", "6"]
  }
};

async function fetchChallenges(type) {
  const endpoint = type === CHALLENGE_TYPES.DAILY.key
    ? "http://localhost:8080/daily-challenge"
    : "http://localhost:8080/classic-challenges";

  const response = await fetch(endpoint);
  if (!response.ok) throw new Error(`Erreur lors de la récupération des données : ${response.statusText}`);

  const data = await response.json();

  return Array.isArray(data) ? data : [data];
}

export function ChallengeBubble({ challenge, type }) {
  console.log("ChallengeBubble rendu :", { challenge, type });

  useEffect(() => {
    console.log("ChallengeBubble mounted/updated =>", { challenge, type });
  }, [challenge, type]);

  const data = challenge.question ?? challenge;
  if (!data.title) return null;

  if (type == CHALLENGE_TYPES.CLASSIC.key) {
    return (
      <div className="bg-blue-100 p-4 rounded shadow-md w-full max-w-xl hover:bg-blue-400 transition-colors duration-300 cursor-pointer">
        <h2 className="text-1xl font-semibold text-gray-800 mb-2">
          {data.title}
        </h2>
      </div>
    );
  }
  return (
<div className="bg-orange-200 p-4 rounded shadow-md w-full max-w-xl hover:bg-orange-600 transition-colors duration-300 cursor-pointer">
      <h5 className="text-1xl font-semibold text-gray-800 mb-2">
      {data.title}
    </h5>
      <p className="text-sm text-gray-500 mb-1">
        Difficulté : {data.difficulty || "N/A"}
      </p>
      <p className="text-sm text-gray-500 mb-4">
        Taux de réussite : {data.acRate?.toFixed(2) || "N/A"}%
      </p>
    </div>
  );
}

export default function ChallengePresentation() {
  const [dailyChallenges, setDailyChallenges] = useState([]);
  const [classicChallenges, setClassicChallenges] = useState([]);

  useEffect(() => {
    fetchChallenges(CHALLENGE_TYPES.DAILY.key)
      .then((res) => {
        console.log("daily =>", res);
        setDailyChallenges(res);
      })
      .catch(console.error);

    fetchChallenges(CHALLENGE_TYPES.CLASSIC.key)
      .then((res) => {
        console.log("classic =>", res);
        setClassicChallenges(res);
      })
      .catch(console.error);
  }, []);

  return (
    <div className="flex flex-row items-start gap-8 mt-3">

      <Container bgColor="bg-gray-200" className="flex-1">
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

      <Container bgColor="bg-gray-100" className="py-16 px-8">
        {dailyChallenges.map((c) => (
          <ChallengeBubble key={c.id} challenge={c} type={CHALLENGE_TYPES.DAILY.key} />
        ))}
        {classicChallenges.map((c) => (
          <ChallengeBubble key={c.id} challenge={c} type={CHALLENGE_TYPES.CLASSIC.key} />
        ))}
      </Container>

    </div>
  );
}
