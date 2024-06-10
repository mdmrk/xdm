import instance from "@/axios";

export const sendChat = async (id: string, body: string) => {
	try {
		const response = await instance.post(`/users/${id}/chat`, {
			content: body,
			recipient_id: id,
		});
		return response;
	} catch (error) {
		console.error(error.response.data);
	}
};

export const retrieveChat = async (id: string) => {
	try {
		const response = await instance.get(`/users/${id}/chat`);
		return response;
	} catch (error) {
		console.error(error.response.data);
	}
};
