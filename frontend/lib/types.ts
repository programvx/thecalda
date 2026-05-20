/** Mirrors the backend `model.User` JSON shape (GET /api/me). */
export type User = {
  uid: string;
  authUserId: string;
  email: string;
  fullName: string;
  createdAt: string;
  updatedAt: string;
};
