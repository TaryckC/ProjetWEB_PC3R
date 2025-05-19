import React from "react";
import NewsSection from "../components/NewsSection"; // à créer si pas encore fait
import NavBar from "../components/Navbar";

export default function NewsPage() {
  return (
    <>
      <NavBar />
      <div className="px-6 mt-8">
        <h1 className="text-2xl font-bold mb-4">📰 Actualités Tech</h1>
        <NewsSection topic="technology" />
      </div>
    </>
  );
}
