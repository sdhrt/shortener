import { parseCSV } from "@/utils/csv";

export default async function ViewPage() {
  const response = (await fetch("http://localhost:8080/view")).text();
  const csv = parseCSV(await response);

  const headers = csv[0];
  const data = csv.slice(1);
  data.forEach((d) => (d[4] = new Date(d[4]).toLocaleString()));
  data.forEach((d) => (d[5] = new Date(d[5]).toLocaleString()));

  return (
    <>
      <div className="flex flex-col">
        <h1>Current hashed urls</h1>
        <table className="table-fixed">
          <thead>
            <tr>
              {headers.map((heading) => {
                return <th key={heading}>{heading}</th>;
              })}
            </tr>
          </thead>
          <tbody>
            {data.map((tuple, index) => {
              return (
                <tr key={index}>
                  {tuple.map((d, index) => (
                    <td key={index}>{d}</td>
                  ))}
                </tr>
              );
            })}
          </tbody>
        </table>
      </div>
    </>
  );
}
