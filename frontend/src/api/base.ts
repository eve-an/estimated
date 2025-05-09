
export const API_URL = "";

export interface ServerResponse<T> {
  status: string;
  data: T | null;
  error: string | null;
};
