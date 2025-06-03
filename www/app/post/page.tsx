"use client";

import { useRef } from "react";
import { toast } from "sonner";

function Page() {
  const inputRef = useRef<HTMLInputElement>(null);
  const submitForm = async () => {
    if (inputRef.current !== null) {
      if (inputRef.current?.value == "") {
        toast.warning("Please enter url");
        return;
      }
      const val = inputRef.current.value;
      const url = URL.parse(val);
      if (url === null) {
        toast.error("Enter valid url");
        return;
      }
      const response = await fetch("http://localhost:8080/create", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ Url: url.toString() }),
      });
      const res = await response.json();
      if (res.error) {
        toast.error(res.error);
        return;
      }
      toast.success(res.message);
      window.location.reload();
    }
  };
  return (
    <div className="h-screen flex flex-col justify-center items-center gap-5">
      <div className="flex gap-5 items-center p-10 border rounded-md">
        <label htmlFor="url_hash_input" className="text-xl">
          Enter url to hash
        </label>
        <input
          type="text"
          id="url_hash_input"
          className="border px-2 py-2"
          ref={inputRef}
        />
      </div>
      <button
        className="px-4 py-2 border rounded-md hover:cursor-pointer"
        onClick={submitForm}
      >
        Hash
      </button>
    </div>
  );
}

export default Page;
