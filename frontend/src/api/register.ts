import { createSignal } from "solid-js";

export const [name, setName] = createSignal<string>("");

export async function register() {
  const response = await fetch("http://localhost:8080/api/v1/register", {
    method: "POST",
    credentials: "include",
  });

  const data = await response.json();
  setName(data.data);

  console.log(data);
}
