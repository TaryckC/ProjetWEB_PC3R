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
      ? "https://projetpc3r.alwaysdata.net/daily-challenge"
      : "https://projetpc3r.alwaysdata.net/classic-challenges";

  const response = await fetch(endpoint);
  if (!response.ok)
    throw new Error(
      `Erreur lors de la récupération des données : ${response.statusText}`
    );

  const data = await response.json();

  return Array.isArray(data) ? data : [data];
}

function ChallengeCard({ challenge, onClick }) {
  let data;

  if (challenge.Question) {
    // Daily challenge
    data = {
      ...challenge.Question,
      title: challenge.Question.Title,
      difficulty: challenge.Question.Difficulty,
      acRate: challenge.Question.ACRate,
    };
  } else if (challenge.question) {
    // Classic challenge
    data = {
      ...challenge.question,
      title: challenge.question.title,
      difficulty: challenge.question.difficulty,
      acRate: challenge.question.acRate,
    };
  } else {
    data = challenge;
  }

  return (
    <div
      onClick={() => onClick(data)}
      className="bg-white border border-gray-200 p-4 rounded-lg shadow-sm hover:shadow-md transition cursor-pointer mb-4"
    >
      <h3 className="text-lg font-semibold text-gray-800">{data.title}</h3>
      {!challenge.question && (
        <>
          <p className="text-sm text-gray-500">
            Difficulté : {data.difficulty || "N/A"}
          </p>
          <p className="text-sm text-gray-500">
            Taux de réussite : {typeof data.acRate === "number" ? data.acRate.toFixed(2) : "N/A"}%
          </p>
        </>
      )}
    </div>
  );
}

export default function ChallengePresentation() {
  const [dailyChallenges, setDailyChallenges] = useState([]);
  const [classicChallenges, setClassicChallenges] = useState([]);
  const [selectedChallenge, setSelectedChallenge] = useState(null);
  const navigate = useNavigate();

  // Chercher dynamiquement la description du daily challenge sélectionné
  useEffect(() => {
    if (
      selectedChallenge &&
      !selectedChallenge.description &&
      !selectedChallenge.question // ne pas fetch pour les classic
    ) {
      const slug = selectedChallenge.titleSlug || selectedChallenge.TitleSlug;
      if (!slug) return;

      console.log("📛 Slug utilisé :", selectedChallenge.titleSlug || selectedChallenge.TitleSlug);

      fetch(`https://projetpc3r.alwaysdata.net/challengeContent/${slug}`)
        .then((res) => res.json())
        .then((data) => {
          console.log("🧪 Description reçue :", data);
          setSelectedChallenge((prev) => ({
            ...prev,
            description: data.description || "<em>Aucune description disponible.</em>",
          }));
        })
        .catch((err) => {
          console.error("Erreur chargement description:", err);
        });
    }
  }, [selectedChallenge]);

  useEffect(() => {
    fetchChallenges(CHALLENGE_TYPES.DAILY.key)
      .then(data => {
        console.log("📅 Daily Challenges reçus :", data);
        setDailyChallenges(data);
      })
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
        {dailyChallenges.map((c, index) => (
          <ChallengeCard
            key={`daily-${c.id || c.question?.FrontendID || index}`}
            challenge={c}
            onClick={setSelectedChallenge}
          />
        ))}

        <h2 className="text-xl font-bold text-gray-700 mt-8 mb-4">
          📘 Classic Challenges
        </h2>
        {classicChallenges.map((c, index) => (
          <ChallengeCard
            key={`classic-${c.id || c.question?.FrontendID || index}`}
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
                  selectedChallenge.Description ||
                  "<em>Aucune description disponible.</em>",
              }}
            />
            <button
              onClick={() =>
                navigate(
                  `/ide/${selectedChallenge.titleSlug ||
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