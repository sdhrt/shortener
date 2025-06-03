export function parseCSV(data: string): string[][] {
  const csvstring = data.trim();
  const csv: string[][] = [];
  const lines = csvstring.split("\n");
  lines.forEach((line) => {
    const tuple: string[] = [];
    line.split(",").forEach((word) => {
      tuple.push(word);
    });
    csv.push(tuple);
  });
  return csv;
}
