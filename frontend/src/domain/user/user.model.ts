export type UserId = string;

export type Role = 'admin' | 'user';

export interface User {
  readonly id: UserId;
  readonly name: string;
  readonly bio: string;
  readonly iconKey?: string;
  readonly role: Role;
}