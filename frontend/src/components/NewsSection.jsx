import React, { useEffect, useState } from "react";

function NewsSection({ topic = "technology" }) {
  const [articles, setArticles] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetch(`/news?topic=${topic}`)
      .then((res) => res.json())
      .then((data) => {
        setArticles(data.articles || []);
        setLoading(false);
      })
      .catch((err) => {
        console.error("Erreur news :", err);
        setLoading(false);
      });
  }, [topic]);

  if (loading) return <p>Chargement des news...</p>;

  return (
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
          <p className="text-sm text-gray-600">{article.description}</p>
        </div>
      ))}
    </div>
  );
}

export default NewsSection;
