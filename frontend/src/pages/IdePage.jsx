import React, { useEffect, useState, useRef } from "react";
import { useParams, useLocation } from "react-router-dom";
import MonacoEditor from "@monaco-editor/react";

function getMonacoLang(lang) {
  if (lang === "python3" || lang === "python2") return "python";
  return lang;
}

function extractFunctionName(code, lang) {
  if (!code) return "unknown";

  if (lang.toLowerCase().startsWith("python")) {
    const match = code.match(/def\s+([a-zA-Z_]\w*)\s*\(/);
    return match?.[1] || "unknown";
  }

  if (
    lang.toLowerCase() === "javascript" ||
    lang.toLowerCase() === "typescript"
  ) {
    const match = code.match(/function\s+([a-zA-Z_]\w*)\s*\(/);
    return match?.[1] || "unknown";
  }

  if (lang.toLowerCase() === "java") {
    const match = code.match(/public\s+\w+\s+([a-zA-Z_]\w*)\s*\(/);
    return match?.[1] || "unknown";
  }

  if (lang.toLowerCase() === "cpp" || lang.toLowerCase() === "c++") {
    const match = code.match(/\w+\s+([a-zA-Z_]\w*)\s*\(/);
    return match?.[1] || "unknown";
  }

  // Ajoute d'autres langages ici si nécessaire...

  return "unknown";
}

export default function IdePage() {
  const { id } = useParams();
  const { state } = useLocation();
  const [challenge, setChallenge] = useState(null);
  const templatesRef = useRef({});
  const [language, setLanguage] = useState("python");
  const [code, setCode] = useState("");
  const [output, setOutput] = useState("");
  const [loading, setLoading] = useState(false);
  const [examples, setExamples] = useState([]);

  const sendChallengeContent = async (titleSlug) => {
    try {
      const response = await fetch(
        `http://localhost:8080/challengeContent/${titleSlug}`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
        }
      );

      if (!response.ok) {
        throw new Error("Erreur côté serveur");
      }

      const result = await response.json();
      const question = result.data.question;

      templatesRef.current = Object.fromEntries(
        question.codeSnippets.map((snip) => {
          const funcName = extractFunctionName(snip.code, snip.lang);
          return [
            snip.lang.toLowerCase(),
            {
              code: snip.code,
              language: snip.lang,
              funcName: funcName, // dynamique
            },
          ];
        })
      );

      console.log("Template : ", templatesRef);

      setChallenge({
        title: question.title,
        titleSlug: question.titleSlug,
        description: question.content,
      });

      console.log("Succès:", question);
    } catch (error) {
      console.error("Erreur réseau:", error);
    }
  };

  useEffect(() => {
    const titleSlug = state?.challenge?.titleSlug || id;
    if (titleSlug) {
      sendChallengeContent(titleSlug);
    }
  }, [state, id]);

  useEffect(() => {
    const template = templatesRef.current[language];
    if (template && template.code) {
      setCode(template.code);
    }
  }, [challenge, language]);

  useEffect(() => {
    if (!challenge || !challenge.description) return;

    const worker = new Worker(
      new URL("../services/extractWorker.js", import.meta.url),
      {
        type: "module",
      }
    );

    worker.onmessage = function (e) {
      const extracted = e.data;
      console.log("Résultat extrait :", extracted);
      setExamples(extracted);
    };

    worker.postMessage(challenge.description);

    return () => worker.terminate(); // bonne pratique pour nettoyer
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

    if (lang === "python" || lang === "python3") {
      return (
        declarations.join("\n") +
        `\nprint(Solution().${funcName}(${varNames.join(", ")}))`
      );
    } else if (lang === "javascript") {
      return declarations.join("\n") + `\nconsole.log(${call});`;
    } else if (lang === "java") {
      // Java nécessite une méthode main
      const javaDeclarations = varNames.map(
        (name, i) => `var ${name} = ${JSON.stringify(values[i])};`
      );

      const javaCall = `System.out.println(new Solution().${funcName}(${varNames.join(
        ", "
      )}));`;

      return `
public class Main {
    public static void main(String[] args) {
        ${javaDeclarations.join("\n        ")}
        ${javaCall}
    }
}
`;
    }

    return "// unsupported language";
  }

  async function runExample(example) {
    const injected = buildTestCode(
      language,
      templatesRef.current[language].funcName,
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
    const results = [];
    for (const example of examples) {
      const result = await runExample(example);
      results.push(result);
      await new Promise((res) => setTimeout(res, 500));
    }
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
            {Object.keys(templatesRef.current).map((lang) => (
              <option key={lang} value={lang}>
                {lang}
              </option>
            ))}
          </select>
        </div>

        <div className="flex-1 border rounded overflow-hidden mb-4">
          <MonacoEditor
            height="100%"
            language={getMonacoLang(language)}
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
