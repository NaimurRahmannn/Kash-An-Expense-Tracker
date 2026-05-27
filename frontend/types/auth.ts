export type LoginInput = {
  email: string;
  password: string;
};

export type RegisterInput = {
  name: string;
  email: string;
  password: string;
};

export type LoginResponseData = {
  user_id: number;
  name: string;
  email: string;
};
