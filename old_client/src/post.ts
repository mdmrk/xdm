import instance from "@/axios";
import type { IPost } from "@/types";

export const getPosts = async (id: string) => {
  try {
    const response = await instance.get<IPost[]>(`/posts?id=${id}`);
    return response.data || [];
  } catch (error) {
    return [];
  }
};

export const likePost = async (id: string) => {
  try {
    const response = await instance.post(`/posts/${id}/like`);
    return response;
  } catch (error) {
    return;
  }
};

export const unlikePost = async (id: string) => {
  try {
    const response = await instance.delete(`/posts/${id}/like`);
    return response;
  } catch (error) {
    return;
  }
};

export const postPost = async (content: string) => {
  try {
    const response = await instance.post(`/posts`, {
      body: content
    });
    return response;
  } catch (error) {
    console.error(error)
  }
};
