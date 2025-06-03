import Link from "next/link";

function page() {
  return (
    <div className="h-screen flex flex-col justify-center items-center">
      <div>Home</div>
      <Link href="/post">Shorten url</Link>
    </div>
  );
}

export default page;
