import React, { useEffect, useState, useRef } from "react";
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

const CHALLENGE_KIND = {
  DAILY: 'daily',
  CLASSIC: 'classic',
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
  const isDaily = challenge._kind === CHALLENGE_KIND.DAILY;
  const title =
    challenge.Question?.Title ||
    challenge.question?.title ||
    challenge.title;

  const difficulty = isDaily
    ? challenge.Question?.Difficulty
    : challenge.question?.difficulty || challenge.difficulty || "N/A";

  const acRate = isDaily
    ? challenge.Question?.ACRate
    : challenge.question?.acRate ?? challenge.acRate;

  return (
    <div
      onClick={() => onClick(challenge)}
      className="bg-white border border-gray-200 p-4 rounded-lg shadow-sm hover:shadow-md transition cursor-pointer mb-4"
    >
      <h3 className="text-lg font-semibold text-gray-800">{title}</h3>
      {isDaily && difficulty && (
        <p className="text-sm text-gray-500">
          Difficulté : {difficulty}
        </p>
      )}
      {isDaily && typeof acRate === "number" && (
        <p className="text-sm text-gray-500">
          Taux de réussite : {acRate.toFixed(2)}%
        </p>
      )}
    </div>
  );
}

function DailyChallengeDetail({ challenge, navigate }) {
  return (
    <>
      <h2 className="text-2xl font-bold text-gray-800 mb-2">
        {challenge.Question.Title}
      </h2>
      <p className="text-sm text-gray-600 mb-1">
        Difficulté : {challenge.Question.Difficulty}
      </p>
      <p className="text-sm text-gray-600 mb-4">
        Taux de réussite : {challenge.Question.ACRate.toFixed(2)}%
      </p>
      <div
        className="text-gray-700 prose max-w-none"
        dangerouslySetInnerHTML={{
          __html: challenge.description || challenge.Question?.Description || "<em>Aucune description disponible.</em>",
        }}
      />
      <button
        onClick={() =>
          navigate(`/ide/${challenge.Question.TitleSlug}`, {
            state: { challenge },
          })
        }
        className="mt-4 bg-blue-600 text-white px-5 py-2 rounded hover:bg-blue-700 transition"
      >
        Commencer le Daily
      </button>
    </>
  );
}

function ClassicChallengeDetail({ challenge, navigate }) {
  return (
    <>
      <h2 className="text-2xl font-bold text-gray-800 mb-2">
        {challenge.question?.title || challenge.title}
      </h2>
      <div className="mb-4" />
      <div
        className="text-gray-700 prose max-w-none"
        dangerouslySetInnerHTML={{
          __html: challenge.question.description || "<em>Aucune description disponible.</em>",
        }}
      />
      <button
        onClick={() =>
          navigate(`/ide/${challenge.question?.titleSlug || challenge.titleSlug}`, {
            state: { challenge },
          })
        }
        className="mt-4 bg-green-600 text-white px-5 py-2 rounded hover:bg-green-700 transition"
      >
        Commencer le Classic
      </button>
    </>
  );
}

export default function ChallengePresentation() {
  const [dailyChallenges, setDailyChallenges] = useState([]);
  const [classicChallenges, setClassicChallenges] = useState([]);
  const [selectedChallenge, setSelectedChallenge] = useState(null);
  const navigate = useNavigate();
  const detailRef = useRef(null);

  useEffect(() => { }, [selectedChallenge]);

  // Daily description (pas de fetch, on utilise Question.Description)
  useEffect(() => {
    if (selectedChallenge?._kind === CHALLENGE_KIND.DAILY && !selectedChallenge.description) {
      setSelectedChallenge(prev => ({
        ...prev,
        description: prev.Question?.Description || "<em>Aucune description disponible.</em>",
      }));
    }
  }, [selectedChallenge?._kind]);

  // Classic description (fetch vers ton API)
  useEffect(() => {
    if (selectedChallenge?._kind === CHALLENGE_KIND.CLASSIC && !selectedChallenge.description) {
      const slug = selectedChallenge.titleSlug || selectedChallenge.TitleSlug;
      if (!slug) return;
      fetch(`https://projetpc3r.alwaysdata.net/challengeContent/${slug}`)
        .then(res => res.json())
        .then(data => {
          setSelectedChallenge(prev => ({
            ...prev,
            description: data.description || "<em>Aucune description disponible.</em>",
          }));
        })
        .catch(err => console.error("Erreur chargement description classique :", err));
    }
  }, [selectedChallenge]);

  useEffect(() => {
    fetchChallenges(CHALLENGE_TYPES.DAILY.key)
      .then(data => {
        setDailyChallenges(data);
      })
      .catch(console.error);

    fetchChallenges(CHALLENGE_TYPES.CLASSIC.key)
      .then(setClassicChallenges)
      .catch(console.error);
  }, []);

  useEffect(() => {
    if (detailRef.current) {
      detailRef.current.scrollTop = 0;
    }
  }, [selectedChallenge]);

  return (
    <div className="flex flex-row w-full h-screen p-6 gap-6 bg-gray-50 overflow-hidden">
      {/* Colonne gauche : Liste des challenges */}
      <div className="flex-none w-[300px] overflow-y-auto h-full pr-4">
        <h2 className="text-xl font-bold text-gray-700 mb-4">
          🗓 Daily Challenge
        </h2>
        {dailyChallenges.map((c, index) => (
          <ChallengeCard
            key={`daily-${c.id || c.question?.FrontendID || index}`}
            challenge={{ ...c, _kind: CHALLENGE_KIND.DAILY }}
            onClick={setSelectedChallenge}
          />
        ))}

        <h2 className="text-xl font-bold text-gray-700 mt-8 mb-4">
          📘 Classic Challenges
        </h2>
        {classicChallenges.map((c, index) => (
          <ChallengeCard
            key={`classic-${c.id || c.question?.FrontendID || index}`}
            challenge={{ ...c, _kind: CHALLENGE_KIND.CLASSIC }}
            onClick={setSelectedChallenge}
          />
        ))}
      </div>

      {/* Colonne droite : Détail du challenge */}
      <div
        ref={detailRef}
        className="flex-1 flex flex-col overflow-y-auto bg-white rounded-xl shadow p-6 border border-gray-200"
      >
        {!selectedChallenge ? (
          <p className="text-gray-500 italic">
            Sélectionne un challenge pour voir les détails.
          </p>
        ) : selectedChallenge._kind === CHALLENGE_KIND.DAILY ? (
          <DailyChallengeDetail challenge={selectedChallenge} navigate={navigate} />
        ) : (
          <ClassicChallengeDetail challenge={selectedChallenge} navigate={navigate} />
        )}
      </div>
    </div>
  );
}