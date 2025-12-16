import { error, fail } from '@sveltejs/kit';
import type { PageServerLoad, Actions } from './$types';
import { getFormByPublicId, submitForm } from '$lib/server/forms';

export const load: PageServerLoad = async (event) => {
	const { params } = event;
	const formId = params.form;
	if (!formId) {
		error(404, {
			message: 'No se encontro esa encuesta'
		});
	}

	const form = await getFormByPublicId(event, formId);
	if (!form) {
		error(404, {
			message: 'No se encontro esa encuesta'
		});
	}

	if (form.is_finished) {
		error(404, {
			message: 'Esta encuesta ya ha sido respondida'
		});
	}

	return {
		form
	};
};

export const actions: Actions = {
	user_form: async (event) => {
		const { request, params } = event;
		const formData = await request.formData();
		const driverRating = formData.get('driver_rating');
		const cargoCondition = formData.get('cargo_condition');
		const comments = formData.get('experience_comments');

		if (!driverRating || !cargoCondition) {
			return fail(400, {
				error: 'Faltan campos requeridos',
				missing: true
			});
		}

		const result = await submitForm(event, params.form, {
			driver_rating: parseInt(driverRating.toString()),
			cargo_condition: cargoCondition.toString(),
			comments: comments ? comments.toString() : ''
		});

		if (result.error) {
			return fail(400, {
				error: result.error
			});
		}

		return {
			success: true
		};
	}
};
