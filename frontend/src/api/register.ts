import { createSignal } from "solid-js";
import { API_URL } from "./base";

export const [name, setName] = createSignal<string>("");

export async function register() {
  const response = await fetch(`${API_URL}/api/v1/register`, {
    method: "POST",
    credentials: "include",
  });

  const data = await response.json();
  setName(data.data);

  console.log(data);
}
