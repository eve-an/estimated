import { API_URL, ServerResponse } from "./base";

export interface Vote {
  value: number;
  timestamp: string;
};

export async function addVote(vote: Vote) {
  const response = await fetch(`${API_URL}/api/v1/votes`, {
    method: "POST",
    credentials: "include",
    body: JSON.stringify(vote)
  });

  if (response.status !== 200) {
    console.log("error on add vote", response);
  }

  const serverResponse: ServerResponse<Vote> = await response.json();
  if (serverResponse.status !== "success") {
    console.log("unsuccessfull add", serverResponse);
  }
}

// async function getAllVotes(): Promise<Vote[]> {
//   const response = await fetch("http://localhost:8080/api/v1/votes", {
//     method: "GET",
//     credentials: "include",
//   });
//
//   if (response.status !== 200) {
//     console.log("error on add vote", response);
//   }
//
//   const serverResponse: ServerResponse<Vote[]> = await response.json();
//   if (serverResponse.status !== "success") {
//     console.log("unsuccessfull add", serverResponse);
//   }
//
//   return serverResponse.data as Vote[]
// }
//
