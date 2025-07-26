import { fail, superValidate } from 'sveltekit-superforms';
import type { PageServerLoad } from './$types';
import { zod4 } from 'sveltekit-superforms/adapters';
import type { Actions } from '@sveltejs/kit';
import { FilesResponseSchema } from '$lib/schemas/search';
import { SubjectMappingSchema } from '$lib/schemas/subject_mapping';
import subjectMappingJson from '../lib/subject_mapping.json';
