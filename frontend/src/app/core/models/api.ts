export interface ApiUser {
  id: string;
  name: string;
  email: string;
  phone?: string | null;
  birthDate?: string | null;
  avatarUrl?: string | null;
  createdAt: string;
  updatedAt: string;
}

export interface TokenPair {
  accessToken: string;
  refreshToken: string;
  expiresIn: number;
}

export interface AuthPayload {
  user: ApiUser;
  tokens: TokenPair;
}

export interface ApiRoom {
  id: string;
  name: string;
  description: string;
  image: string;
  maxAdults: number;
  maxChildren: number;
  prices: Record<string, { from: string | number; to: string | number }>;
}
