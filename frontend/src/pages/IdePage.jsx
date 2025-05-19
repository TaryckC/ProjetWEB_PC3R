import React, { useEffect, useState, useRef } from "react";
import Navbar from "../components/Navbar";
const BACKEND_URL = "https://projetpc3r.alwaysdata.net";
import { useParams, useLocation } from "react-router-dom";
import MonacoEditor from "@monaco-editor/react";
import { auth } from "../firebaseAuth";
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
  const [forumMessages, setForumMessages] = useState([]);
  const [newMessage, setNewMessage] = useState("");

  const sendChallengeContent = async (titleSlug) => {
    try {
      const response = await fetch(
        `${BACKEND_URL}/challengeContent/${titleSlug}`,
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
              funcName: funcName,
            },
          ];
        })
      );

      setChallenge({
        title: question.title,
        titleSlug: question.titleSlug,
        description: question.content,
      });
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
    if (!challenge?.titleSlug) return;
    fetch(`${BACKEND_URL}/forum/challengeContent/${challenge?.titleSlug}`)
      .then((res) => res.json())
      .then((data) => {
        if (Array.isArray(data)) {
          setForumMessages(data);
        } else if (Array.isArray(data.messages)) {
          setForumMessages(data.messages);
        } else {
          setForumMessages([]);
        }
      })
      .catch((err) => console.error("Erreur forum:", err));
  }, [challenge]);

  useEffect(() => {
    if (!challenge || !challenge.description) return;

    const worker = new Worker(
      new URL("../services/extractWorker.js", import.meta.url),
      {
        type: "module",
      }
    );

    worker.onmessage = function (e) {
      setExamples(e.data);
    };

    worker.postMessage(challenge.description);

    return () => worker.terminate();
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
}`;
    }

    return "// unsupported language";
  }

  const runAllExamples = async () => {
    setOutput("Execution en cours...");
    setLoading(true);

    const submissionList = examples.map((example) => {
      const injected = buildTestCode(
        language,
        templatesRef.current[language].funcName,
        example.input
      );
      return {
        source_code: code + "\n" + injected,
        language_id: languageIdMap[language],
        stdin: "",
      };
    });

    try {
      const response = await fetch(
        "https://projetpc3r.alwaysdata.net/compile",
        {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify(submissionList),
        }
      );

      const judgeResults = await response.json();

      const results = examples.map((example, i) => {
        const actual =
          judgeResults[i]?.stdout?.trim() ||
          judgeResults[i]?.stderr?.trim() ||
          "";
        return {
          ...example,
          actual,
          pass: actual === example.expected,
        };
      });

      const report = results.map(
        (ex, i) =>
          `Cas ${i + 1} :\nInput: ${JSON.stringify(ex.input)}\nExpected: ${
            ex.expected
          }\nActual: ${ex.actual}\nRésultat: ${
            ex.pass ? "✅ Réussi" : "❌ Échoué"
          }\n`
      );
      setOutput(report.join("\n\n"));
    } catch (err) {
      console.error("Erreur d'exécution batch:", err);
      setOutput("❌ Erreur lors de l'exécution des tests");
    } finally {
      setLoading(false);
    }
  };

  function handlePostMessage() {
    const trimmedMessage = newMessage.trim();
    if (!trimmedMessage) return;

    const newEntry = {
      author: auth.currentUser?.displayName || "anonyme",
      content: trimmedMessage,
    };

    fetch(`${BACKEND_URL}/forum/challengeContent/${challenge.titleSlug}`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(newEntry),
    })
      .then((res) => res.text())
      .then(() => {
        setForumMessages((prev) => [...prev, newEntry]);
        setNewMessage("");
      })
      .catch((err) => console.error("❌ Erreur POST forum :", err));
  }

  if (!challenge)
    return (
      <div className="h-screen w-full flex justify-center items-center">
        <svg
          className="animate-spin h-16 w-16 text-gray-500"
          xmlns="http://www.w3.org/2000/svg"
          fill="none"
          viewBox="0 0 24 24"
        >
          <circle
            className="opacity-25"
            cx="12"
            cy="12"
            r="10"
            stroke="currentColor"
            strokeWidth="4"
          />
          <path
            className="opacity-75"
            fill="currentColor"
            d="M4 12a8 8 0 018-8v4a4 4 0 00-4 4H4z"
          />
        </svg>
      </div>
    );

  return (
    <div className="min-h-screen bg-gray-100">
      <Navbar />
      <div className="flex h-[calc(100vh-1rem)] w-full mt-4">
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

          <section className="mt-8">
            <h2 className="text-lg font-semibold mb-2">
              💬 Forum du challenge
            </h2>

            <div className="bg-gray-100 p-4 rounded max-h-60 overflow-y-auto mb-4">
              {forumMessages.length === 0 ? (
                <p className="text-gray-600 italic">
                  Aucun message pour l’instant.
                </p>
              ) : (
                <ul>
                  {forumMessages.map((msg, i) => (
                    <li key={i} className="mb-2">
                      <span className="font-bold">{msg.author}</span>:{" "}
                      {msg.content}
                    </li>
                  ))}
                </ul>
              )}
            </div>

            <div className="flex gap-2">
              <input
                type="text"
                className="flex-1 border rounded px-2 py-1"
                placeholder="Écris un message..."
                value={newMessage}
                onChange={(e) => setNewMessage(e.target.value)}
                onKeyDown={(e) => {
                  if (e.key === "Enter") {
                    e.preventDefault();
                    handlePostMessage();
                  }
                }}
              />
              <button
                onClick={handlePostMessage}
                className="bg-blue-600 text-white px-4 py-1 rounded hover:bg-blue-700"
              >
                Envoyer
              </button>
            </div>
          </section>
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
    </div>
  );
}
