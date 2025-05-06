import { createMemo, createSignal, onMount, type Component } from 'solid-js';

import { Chart } from './components/Chart';
import { addVote, Vote } from './api/vote';
import { eventVotes, initSSE } from './api/sse';
import { register, name } from './api/register';

function getFormattedTimestamp(): string {
  const date = new Date();
  return date.toISOString()
}

const App: Component = () => {
  const [votes, setVotes] = createSignal<Vote[]>([]);

  const mappedVotes = createMemo(() => {
    const event = eventVotes().votes;
    const n = name();
    const local = votes();

    return {
      ...event,
      [n]: [...(event[n] ?? []), ...local]
    };
  });

  const handleClick = async (e: MouseEvent) => {
    const buttonId = (e.target as HTMLButtonElement).id;
    const parsedValue = parseInt(buttonId);

    const newVote = { value: parsedValue, timestamp: getFormattedTimestamp() };
    setVotes((prev) => [...prev, newVote]);

    try {
      await addVote(newVote);
    } catch (error) {
      console.warn("could not send vote data to server", error, newVote)
    }
  };

  onMount(async () => {
    try {
      await register();
      initSSE();
    } catch (error) {
      console.warn("registration failed", error)
    }
  });

  return (
    <div class="flex flex-col items-center min-h-screen w-full bg-background text-foreground p-6">
      <div class="w-full max-w-4xl rounded-2xl text-foreground shadow-lg border border-border p-6 mb-8">
        <Chart data={mappedVotes()} name={name()} />
      </div>
      <div class="flex flex-row space-x-4 justify-center">
        {[1, 2, 3, 5, 8].map(num => (
          <button
            id={num.toString()}
            onClick={handleClick}
            class="rounded-2xl px-10 py-6 text-xl font-medium transition-transform duration-200 hover:-translate-y-2 shadow-md border border-border"
          >
            {num}
          </button>
        ))}
      </div>
    </div>
  );
};

export default App;
