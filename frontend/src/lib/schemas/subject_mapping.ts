import z from 'zod';

export const SubjectMappingSchema = z.array(
	z.object({
		id: z.number(),
		name: z.string()
	})
);

export type SubjectMappingSchema = z.infer<typeof SubjectMappingSchema>;
