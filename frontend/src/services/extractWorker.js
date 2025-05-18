function parseInputString(inputStr) {
  const obj = {};
  let i = 0;

  while (i < inputStr.length) {
    // 1. Lire la clé
    while (i < inputStr.length && /\s/.test(inputStr[i])) i++;
    const startKey = i;
    while (i < inputStr.length && /[a-zA-Z0-9_]/.test(inputStr[i])) i++;
    const key = inputStr.slice(startKey, i);

    // 2. Sauter l’égalité et les espaces
    while (i < inputStr.length && /\s|=/.test(inputStr[i])) i++;

    // 3. Lire la valeur
    let value = "";
    if ('[({"'.includes(inputStr[i])) {
      // On démarre une structure
      const opener = inputStr[i];
      const closer =
        opener === "["
          ? "]"
          : opener === "{"
            ? "}"
            : opener === "("
              ? ")"
              : opener;
      let depth = opener === '"' ? 0 : 1;
      value += opener;
      i++;
      while (i < inputStr.length) {
        const ch = inputStr[i];
        value += ch;
        if (ch === opener && opener !== '"') {
          depth++;
        } else if (ch === closer) {
          if (opener === '"' || --depth === 0) {
            i++;
            break;
          }
        } else if (ch === "\\" && opener === '"') {
          // échappement dans les chaînes
          i++;
          if (i < inputStr.length) {
            value += inputStr[i];
          }
        }
        i++;
      }
    } else {
      // simple token jusqu’à la virgule ou fin
      while (i < inputStr.length && inputStr[i] !== ",") {
        value += inputStr[i++];
      }
    }

    // 4. Nettoyage et conversion
    value = value.trim();
    try {
      // Convertir les tuples (…)
      if (value.startsWith("(") && value.endsWith(")")) {
        const inner = value.slice(1, -1).trim();
        value = inner === "" ? "[]" : `[${inner}]`;
      }
      obj[key] = JSON.parse(value);
    } catch (e) {
      // fallback
      obj[key] = value;
    }

    // 5. Sauter la virgule
    while (i < inputStr.length && /[,\s]/.test(inputStr[i])) i++;
  }

  return obj;
}

self.onmessage = function (e) {
  const html = e.data;

  const preBlocks = Array.from(html.matchAll(/<pre[^>]*>([\s\S]*?)<\/pre>/gi));
  const result = [];

  for (const [, content] of preBlocks) {
    const text = content.replace(/<[^>]+>/g, "").trim();
    const match = text.match(/Input:\s*([\s\S]*?)\nOutput:\s*([^\n]+)/);
    if (match) {
      const rawInput = match[1].trim();
      try {
        // Tentative de parsing direct comme JSON si la chaîne commence par { ou [
        if (rawInput.startsWith("{") || rawInput.startsWith("[")) {
          result.push({
            input: JSON.parse(rawInput),
            expected: match[2].trim(),
          });
        } else {
          const parsedInput = parseInputString(rawInput);
          result.push({
            input: parsedInput,
            expected: match[2].trim(),
          });
        }
      } catch (error) {
        console.error("Error parsing input:", rawInput, error);
        // Fallback: traiter toute la chaîne d'input comme une valeur
        result.push({
          input: { value: rawInput },
          expected: match[2].trim(),
        });
      }
    }
  }

  self.postMessage(result);
};
