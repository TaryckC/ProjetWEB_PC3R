function parseInputString(inputStr) {
  const obj = {};

  // Match "key = value" même si value contient des [], {}, "", () ou des virgules
  const regex =
    /([a-zA-Z_][a-zA-Z0-9_]*)\s*=\s*("[^"]*"|\[[^\]]*\]|\{[^\}]*\}|\([^\)]*\)|[^,]+)/g;

  let match;
  while ((match = regex.exec(inputStr)) !== null) {
    const key = match[1].trim();
    let value = match[2].trim();

    try {
      // Tuples → listes
      if (value.startsWith("(") && value.endsWith(")")) {
        value = "[" + value.slice(1, -1) + "]";
      }
      obj[key] = JSON.parse(value);
    } catch {
      obj[key] = value;
    }
  }

  return obj;
}

self.onmessage = function (e) {
  const html = e.data;

  const preBlocks = Array.from(html.matchAll(/<pre[^>]*>([\s\S]*?)<\/pre>/gi));
  const result = [];

  for (const [, content] of preBlocks) {
    const text = content.replace(/<[^>]+>/g, ""); // nettoie HTML
    const match = text.match(/Input:\s*([\s\S]*?)\nOutput:\s*([^\n]+)/);
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
