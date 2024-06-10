import instance from "@/axios";
import type { IUser } from "@/types";

export const getUserById = async (id: string) => {
	try {
		const response = await instance.get<IUser>(`/users/${id}`);
		return response.data;
	} catch (error) {
		return [];
	}
};

export const getUserLikes = async (id: string) => {
	try {
		const response = await instance.get<string[]>(`/users/${id}/like`);
		return response.data;
	} catch (error) {
		return [];
	}
};

export const  followUser= async (id: string) => {
	try {
		const response = await instance.post(`/users/${id}/follow`);
		return response.data;
	} catch (error) {
		return [];
	}
};

export const  unfollowUser= async (id: string) => {
	try {
		const response = await instance.delete(`/users/${id}/follow`);
		return response.data;
	} catch (error) {
		return [];
	}
};
