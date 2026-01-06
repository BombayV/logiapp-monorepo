import { error } from '@sveltejs/kit';

export const load = async () => {
	error(404, {
		message: 'No se encontro esta pagina.'
	});
};
