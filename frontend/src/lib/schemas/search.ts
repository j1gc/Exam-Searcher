import { z } from 'zod';

const ExamFileSchema = z.object({
	// exam
	exam_id: z.number(),
	subject_id: z.number(),
	exam_type: z.string(),
	difficulty: z.string(),
	task_label: z.string(),
	work_time_in_minutes: z.number(),
	// file
	file_id: z.number(),
	year: z.number(),
	file_path: z.string(),
	embedding_id: z.number()
});

const AnswerFileSchema = z.object({
	// answer
	answer_id: z.number(),
	subject_id: z.number(),
	// file
	file_id: z.number(),
	year: z.number(),
	file_path: z.string(),
	embedding_id: z.number()
});

const OtherFileSchema = z.object({
	// other
	other_id: z.number(),
	subject_id: z.number(),
	// file
	file_id: z.number(),
	year: z.number(),
	file_path: z.string(),
	embedding_id: z.number()
});

const FileTypeSchema = z.union([
	z.object({ Exam: ExamFileSchema }),
	z.object({ Answer: AnswerFileSchema }),
	z.object({ Other: OtherFileSchema })
]);

export const FileResponseSchema = z.object({
	file_type: z.string(),
	file: FileTypeSchema,
	similarity: z.float32()
});

export const SearchRequestSchema = z.object({
	query: z.string(),
	file_types: z.array(z.string()),
	subject_ids: z.array(z.number()),
	years: z.array(z.number())
});

export const FilesResponseSchema = z.array(FileResponseSchema);

export type FilesResponseSchema = z.infer<typeof FilesResponseSchema>;

export type SearchRequestSchema = z.infer<typeof SearchRequestSchema>;
