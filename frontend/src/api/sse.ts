import { createSignal, onCleanup } from "solid-js";
import { Vote } from "./vote";
import { API_URL } from "./base";

export interface Votes {
  votes: Record<string, Vote[]>;
};

export const [eventVotes, setEventVotes] = createSignal<Votes>({ "votes": {} });

export function initSSE() {
  const eventSource = new EventSource(`${API_URL}/api/v1/events`, { withCredentials: true });

  eventSource.onmessage = (event) => {
    try {
      console.log(event);
      const parsed = JSON.parse(event.data);
      const points: Votes = parsed.data;

      setEventVotes(points)
    } catch (err) {
      console.warn("Invalid SSE data:", event.data, err);
    }
  };

  eventSource.onerror = (err) => {
    console.error("EventSource failed:", err);
  };

  onCleanup(() => {
    eventSource.close();
  });
}
