import { FilesResponseSchema } from '$lib/schemas/search.js';

export const load = async ({ fetch, params }) => {
	const searchData: Object = {
		query: 'Diagramm',
		file_types: ['answer', 'other', 'exam'],
		subject_ids: [],
		years: [2020, 2021]
	};

	let searchResp = await fetch('http://localhost:3000/search', {
		method: 'POST',
		headers: {
			Accept: 'application/json',
			'Content-Type': 'application/json'
		},

		body: JSON.stringify(searchData)
	});
	const searchRespData: FilesResponseSchema = await searchResp.json();

	const files = FilesResponseSchema.parse(searchRespData);

	return { files };
};
