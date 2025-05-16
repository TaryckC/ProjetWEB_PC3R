self.onmessage = function (e) {
  const html = e.data;

  const preBlocks = Array.from(html.matchAll(/<pre[^>]*>([\s\S]*?)<\/pre>/gi));
  const result = [];

  for (const [, content] of preBlocks) {
    const text = content.replace(/<[^>]+>/g, ""); // supprime les balises HTML

    const match = text.match(/Input:\s*([^\n]+)[\s\S]*?Output:\s*([^\n]+)/);
    if (match) {
      const rawInput = match[1].trim();
      const parsedInput = parseInputString(rawInput);

      result.push({
        input: parsedInput,
        expected: match[2].trim(),
      });
    }
  }

  self.postMessage(result);
};

function parseInputString(inputStr) {
  const obj = {};
  const parts = inputStr.split(",").map((part) => part.trim());

  for (const part of parts) {
    const [key, rawValue] = part.split("=").map((s) => s.trim());
    let parsed;

    try {
      // Essaye de parser comme JSON
      if (rawValue.startsWith("(") && rawValue.endsWith(")")) {
        // Tuples → listes
        parsed = JSON.parse(rawValue.replace(/\(/g, "[").replace(/\)/g, "]"));
      } else {
        parsed = JSON.parse(rawValue);
      }
    } catch {
      // Si échec, garde en chaîne brute
      parsed = rawValue;
    }

    obj[key] = parsed;
  }

  return obj;
}
