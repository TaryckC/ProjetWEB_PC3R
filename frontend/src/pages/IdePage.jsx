import React, { useState } from "react";
import Editor from "@monaco-editor/react";

// Les templates en été generés par chatgpt
const templates = {
  javascript: `function main() {
  console.log("Hello, world!");
}

main();`,

  python: `def main():
    print("Hello, world!")

if __name__ == "__main__":
    main()`,

  go: `package main

import "fmt"

func main() {
    fmt.Println("Hello, world!")
}`,

  java: `public class Main {
    public static void main(String[] args) {
        System.out.println("Hello, world!");
    }
}`,

  cpp: `#include <iostream>
using namespace std;

int main() {
    cout << "Hello, world!" << endl;
    return 0;
}`,

  html: `<!DOCTYPE html>
<html>
  <head>
    <title>Hello</title>
  </head>
  <body>
    <h1>Hello, world!</h1>
  </body>
</html>`,

  json: `{
  "message": "Hello, world!"
}`,
};

function IDE() {
  const [language, setLanguage] = useState("javascript");
  const [savedCodes, setSavedCodes] = useState({ javascript: templates.javascript });
  const [code, setCode] = useState(templates.javascript);

    const handleLanguageChange = (e) => {
    const newLang = e.target.value;

    // Sauvegarder le code actuel dans savedCodes
    setSavedCodes((prev) => ({
      ...prev,
      [language]: code,
    }));    

    // Charger le code existant pour le nouveau langage (ou le template si rien)
    setLanguage(newLang);
    setCode(savedCodes[newLang] || templates[newLang]);
  };

  const handleCodeChange = (newValue) => {
    setCode(newValue);
  };

  return (
    // Le formulaire ci dessous est dérivé de chatgpt
    <div>
      <label htmlFor="lang-select">Choisir le langage : </label>
      <select id="lang-select" value={language} onChange={handleLanguageChange}>
        {Object.keys(templates).map((lang) => (
          <option key={lang} value={lang}>
            {lang}
          </option>
        ))}
      </select>

      <Editor
        height="300px"
        language={language}
        value={code}
        onChange={handleCodeChange}
        theme="vs-dark"
      />
    </div>
  );
}

export default IDE;
