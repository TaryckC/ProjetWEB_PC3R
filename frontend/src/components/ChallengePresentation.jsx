import React, { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";

const CHALLENGE_TYPES = {
  DAILY: {
    key: "daily_challenge",
    id: ["0"],
  },
  CLASSIC: {
    key: "classic_challenges",
    ids: ["0", "1", "2", "3", "4", "5", "6"],
  },
};

async function fetchChallenges(type) {
  const endpoint =
    type === CHALLENGE_TYPES.DAILY.key
      ? "http://localhost:8080/daily-challenge"
      : "http://localhost:8080/classic-challenges";

  const response = await fetch(endpoint);
  if (!response.ok)
    throw new Error(
      `Erreur lors de la récupération des données : ${response.statusText}`
    );

  const data = await response.json();

  return Array.isArray(data) ? data : [data];
}

function ChallengeCard({ challenge, onClick }) {
  const data = challenge.question ?? challenge;

  return (
    <div
      onClick={() => onClick(data)}
      className="bg-white border border-gray-200 p-4 rounded-lg shadow-sm hover:shadow-md transition cursor-pointer mb-4"
    >
      <h3 className="text-lg font-semibold text-gray-800">{data.title}</h3>
      <p className="text-sm text-gray-500">
        Difficulté : {data.difficulty || "N/A"}
      </p>
      <p className="text-sm text-gray-500">
        Taux de réussite : {data.acRate?.toFixed(2) || "N/A"}%
      </p>
    </div>
  );
}

export default function ChallengePresentation() {
  const [dailyChallenges, setDailyChallenges] = useState([]);
  const [classicChallenges, setClassicChallenges] = useState([]);
  const [selectedChallenge, setSelectedChallenge] = useState(null);
  const navigate = useNavigate();

  useEffect(() => {
    fetchChallenges(CHALLENGE_TYPES.DAILY.key)
      .then(setDailyChallenges)
      .catch(console.error);

    fetchChallenges(CHALLENGE_TYPES.CLASSIC.key)
      .then(setClassicChallenges)
      .catch(console.error);
  }, []);

  return (
    <div className="flex flex-row w-full h-screen p-6 gap-6 bg-gray-50 overflow-hidden">
      {/* Colonne gauche : Liste des challenges */}
      <div className="w-1/3 overflow-y-auto h-full pr-4">
        <h2 className="text-xl font-bold text-gray-700 mb-4">
          🗓 Daily Challenge
        </h2>
        {dailyChallenges.map((c) => (
          <ChallengeCard
            key={c.id}
            challenge={c}
            onClick={setSelectedChallenge}
          />
        ))}

        <h2 className="text-xl font-bold text-gray-700 mt-8 mb-4">
          📘 Classic Challenges
        </h2>
        {classicChallenges.map((c) => (
          <ChallengeCard
            key={c.id}
            challenge={c}
            onClick={setSelectedChallenge}
          />
        ))}
      </div>

      {/* Colonne droite : Détail du challenge */}
      <div className="flex-1 overflow-y-auto h-full bg-white rounded-xl shadow p-6 border border-gray-200">
        {selectedChallenge ? (
          <>
            <h2 className="text-2xl font-bold text-gray-800 mb-2">
              {selectedChallenge.title}
            </h2>
            <p className="text-sm text-gray-600 mb-1">
              Difficulté : {selectedChallenge.difficulty || "N/A"}
            </p>
            <p className="text-sm text-gray-600 mb-4">
              Taux de réussite : {selectedChallenge.acRate?.toFixed(2) || "N/A"}
              %
            </p>
            <div
              className="text-gray-700 prose max-w-none"
              dangerouslySetInnerHTML={{
                __html:
                  selectedChallenge.description ||
                  selectedChallenge.question?.description ||
                  "<em>Aucune description disponible.</em>",
              }}
            />
            <button
              onClick={() =>
                navigate(
                  `/ide/${
                    selectedChallenge.titleSlug ||
                    selectedChallenge.question?.titleSlug ||
                    selectedChallenge.title
                  }`,
                  {
                    state: { challenge: selectedChallenge },
                  }
                )
              }
              className="mt-4 bg-blue-600 text-white px-5 py-2 rounded hover:bg-blue-700 transition"
            >
              Commencer
            </button>
          </>
        ) : (
          <p className="text-gray-500 italic">
            Sélectionne un challenge pour voir les détails.
          </p>
        )}
      </div>
    </div>
  );
}
