export interface IUser {
	Id: string;
	Alias: string;
	Username: string;
	Seen: Date;
}

export interface AuthResponse {
	token: string;
	userId: string;
}

export interface IPost {
	Id: string;
	UserId: string;
	Body: string;
	Likes: number;
	CreatedAt: Date;
}

export interface IMessage {
	Id: string;
	UserId: string;
	Body: string;
	CreatedAt: Date;
}
