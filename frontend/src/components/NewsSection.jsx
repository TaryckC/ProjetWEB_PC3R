import React, { useEffect, useState } from "react";

const TOPICS = [
  { label: "Technologie", value: "technology" },
  { label: "Programmation", value: "programming" },
  { label: "Cybersécurité", value: "cybersecurity" },
];

function NewsSection() {
  const [topic, setTopic] = useState("programming");
  const [articles, setArticles] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    setLoading(true); // relance le spinner à chaque changement
    fetch(`https://projetpc3r.alwaysdata.net/news?topic=${topic}`)
      .then((res) => res.json())
      .then((data) => {
        setArticles(data.data || []);
        setLoading(false);
      })
      .catch((err) => {
        console.error("Erreur news :", err);
        setLoading(false);
      });
  }, [topic]);

  return (
    <div className="space-y-6">
      {/* Sélecteur de topic */}
      <div>
        <label htmlFor="topic-select" className="block mb-2 font-medium">
          Choisir un thème :
        </label>
        <select
          id="topic-select"
          value={topic}
          onChange={(e) => setTopic(e.target.value)}
          className="px-4 py-2 border rounded bg-white text-gray-800"
        >
          {TOPICS.map((t) => (
            <option key={t.value} value={t.value}>
              {t.label}
            </option>
          ))}
        </select>
      </div>

      {/* Résultat des actualités */}
      {loading ? (
        <p>Chargement des actualités...</p>
      ) : articles.length === 0 ? (
        <p>Aucune actualité disponible.</p>
      ) : (
        <div className="space-y-4">
          {articles.slice(0, 5).map((article, index) => (
            <div key={index} className="border-b pb-2">
              <a
                href={article.url}
                target="_blank"
                rel="noreferrer"
                className="text-blue-600 font-semibold"
              >
                {article.title}
              </a>
              <p className="text-sm text-gray-600">{article.excerpt}</p>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}

export default NewsSection;
