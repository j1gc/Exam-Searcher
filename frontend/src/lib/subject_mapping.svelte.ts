import subjectMappingJson from '$lib/subject_mapping.json';
import { SubjectMappingSchema } from '$lib/schemas/subject_mapping.js';

let subjectMappingData: SubjectMappingSchema = SubjectMappingSchema.parse(subjectMappingJson);
export let subjectMapping = new Map(
	subjectMappingData.map((subject) => [subject.id, subject.name])
);

export function findSubjectIdByName(subjectName: string) {
	return subjectMappingData.find((v) => v.name === subjectName)?.id;
}
