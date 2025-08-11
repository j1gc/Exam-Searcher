import subjectMappingJson from '$lib/subject_mapping.json';
import { SubjectMappingSchema } from '$lib/schemas/subject_mapping.js';
import { getLocale } from './paraglide/runtime';

let subjectMappingData: SubjectMappingSchema = SubjectMappingSchema.parse(subjectMappingJson);

export let subjectMapping = new Map(
	subjectMappingData.map((subject) => [subject.id, getLocale() === 'de' ? subject.de : subject.en])
);

export function findSubjectIdByName(subjectName: string) {
	if (getLocale() === 'de') {
		return subjectMappingData.find((v) => v.de === subjectName)?.id;
	} else {
		return subjectMappingData.find((v) => v.en === subjectName)?.id;
	}
}
