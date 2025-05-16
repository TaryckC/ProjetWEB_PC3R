import React, { useEffect, useState } from "react";
import { useParams, useLocation } from "react-router-dom";
import MonacoEditor from "@monaco-editor/react";

const templates = {
  python: {
    code: `def twoSum(nums, target):
    return []`,
    funcName: "twoSum",
  },
  javascript: {
    code: `function twoSum(nums, target) {
  return [];
}`,
    funcName: "twoSum",
  },
};

function extractExamplesFromHTML(html) {
  const preBlocks = Array.from(html.matchAll(/<pre[^>]*>(.*?)<\/pre>/gs));
  const examples = [];

  for (const [, block] of preBlocks) {
    const lines = block.trim().split("\n");
    let inputLine = "",
      outputLine = "";

    for (let line of lines) {
      if (line.trim().startsWith("Input:")) {
        inputLine = line.trim().replace("Input:", "").trim();
      }
      if (line.trim().startsWith("Output:")) {
        outputLine = line.trim().replace("Output:", "").trim();
      }
    }

    if (inputLine && outputLine) {
      const inputObj = {};
      const parts = inputLine.split(/,\s*/);
      for (let part of parts) {
        const [key, value] = part.split("=");
        if (key && value) {
          try {
            inputObj[key.trim()] = JSON.parse(value.trim());
          } catch {
            inputObj[key.trim()] = value.trim();
          }
        }
      }
      examples.push({ input: inputObj, expected: outputLine });
    }
  }

  return examples;
}

export default function IdePage() {
  const { id } = useParams();
  const { state } = useLocation();
  const challenge = state?.challenge;

  const [language, setLanguage] = useState("python");
  const [code, setCode] = useState(templates["python"].code);
  const [output, setOutput] = useState("");
  const [loading, setLoading] = useState(false);
  const [examples, setExamples] = useState([]);

  useEffect(() => {
    setCode(templates[language].code);
  }, [language]);

  useEffect(() => {
    const html = challenge.description || challenge.question?.description;
    const autoExamples = extractExamplesFromHTML(html);
    setExamples(autoExamples);
  }, [challenge]);

  const languageIdMap = {
    python: 71,
    javascript: 63,
  };

  function buildTestCode(lang, funcName, inputObj) {
    const varNames = Object.keys(inputObj);
    const values = Object.values(inputObj);
    const declarations = varNames.map(
      (name, i) => `${name} = ${JSON.stringify(values[i])}`
    );
    const call = `${funcName}(${varNames.join(", ")})`;

    if (lang === "python") {
      return declarations.join("\n") + `\nprint(${call})`;
    } else if (lang === "javascript") {
      return declarations.join("\n") + `\nconsole.log(${call});`;
    }
    return "// unsupported language";
  }

  async function runExample(example) {
    const injected = buildTestCode(
      language,
      templates[language].funcName,
      example.input
    );
    const fullCode = code + "\n" + injected;

    setLoading(true);
    try {
      const response = await fetch("/compile", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          source_code: fullCode,
          language_id: languageIdMap[language],
          stdin: "",
        }),
      });

      const result = await response.json();
      const actual = result.stdout?.trim() || result.stderr?.trim();
      return { ...example, actual, pass: actual === example.expected };
    } catch (err) {
      return { ...example, actual: "Erreur", pass: false };
    } finally {
      setLoading(false);
    }
  }

  const runAllExamples = async () => {
    setOutput("Chargement...");
    const results = await Promise.all(examples.map(runExample));
    const report = results.map(
      (ex, i) =>
        `Cas ${i + 1} :\nInput: ${JSON.stringify(ex.input)}\nExpected: ${
          ex.expected
        }\nActual: ${ex.actual}\nRésultat: ${
          ex.pass ? "✅ Réussi" : "❌ Échoué"
        }\n`
    );
    setOutput(report.join("\n\n"));
  };

  if (!challenge) return <p className="p-6 text-red-500">Challenge manquant</p>;

  return (
    <div className="flex h-screen w-full">
      {/* Description */}
      <aside className="w-1/2 p-6 border-r overflow-y-auto">
        <h1 className="text-2xl font-bold mb-4">{challenge.title}</h1>
        <div
          className="prose max-w-none text-gray-800 overflow-x-auto break-words"
          dangerouslySetInnerHTML={{
            __html: challenge.description || challenge.question?.description,
          }}
        />
        <button
          onClick={runAllExamples}
          className="mt-6 bg-green-600 text-white px-4 py-2 rounded hover:bg-green-700"
        >
          Tester les cas d'exemples
        </button>
      </aside>

      {/* Éditeur + Résultat */}
      <main className="flex-1 p-6 flex flex-col bg-gray-50">
        <div className="flex items-center gap-4 mb-4">
          <label>Langage :</label>
          <select
            value={language}
            onChange={(e) => setLanguage(e.target.value)}
            className="border px-2 py-1 rounded"
          >
            {Object.keys(templates).map((lang) => (
              <option key={lang} value={lang}>
                {lang}
              </option>
            ))}
          </select>
        </div>

        <div className="flex-1 border rounded overflow-hidden mb-4">
          <MonacoEditor
            height="100%"
            language={language}
            theme="vs-dark"
            value={code}
            onChange={(value) => setCode(value)}
            options={{
              fontSize: 14,
              minimap: { enabled: false },
              automaticLayout: true,
            }}
          />
        </div>

        <div className="bg-black text-green-400 p-4 rounded text-sm font-mono overflow-x-auto">
          <p className="mb-1 font-semibold">Résultat :</p>
          <pre>{output}</pre>
        </div>
      </main>
    </div>
  );
}
